package analyzer

import (
	"go/token"
	"go/types"
	"maps"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

var Analyzer = &analysis.Analyzer{
	Name:     "unnecessaryerror",
	Doc:      "Finds errors that are only ever nil-checked and could be replaced with a bool.",
	Run:      run,
	Requires: []*analysis.Analyzer{buildssa.Analyzer},
}

// go run ./cmd/unnecessaryerror -- ./analyzer/testdata/src/a/...
func run(pass *analysis.Pass) (interface{}, error) {
	ssaInput := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)

	fnErrs := findReturnedErrors(ssaInput.SrcFuncs)
	checking := make(map[ssa.Instruction]struct{})
	seen := make(map[ssa.Instruction]bool)

	var isUsed func(instr ssa.Instruction) bool
	isUsed = func(instr ssa.Instruction) (result bool) {
		if v, ok := seen[instr]; ok {
			return v
		}
		if _, ok := checking[instr]; ok {
			return false
		}
		checking[instr] = struct{}{}
		defer func() {
			seen[instr] = result
			delete(checking, instr)
		}()

		switch v := instr.(type) {
		case *ssa.FieldAddr, *ssa.IndexAddr, *ssa.MapUpdate, *ssa.Send, *ssa.TypeAssert:
			return true
		case *ssa.Call:
			callee := v.Call.StaticCallee()
			if callee == nil {
				return true // Method call.
			}
			if callee.Pkg != nil && callee.Pkg != v.Parent().Pkg {
				return true // Call to external function.
			}
			refs := v.Referrers()
			if refs == nil {
				return false
			}
			for _, r := range *refs {
				if isUsed(r) {
					return true
				}
			}
		case *ssa.Return:
			containingFn := v.Parent()
			if !isUnexpFunc(containingFn) {
				return true
			}
			callers := findCaller(ssaInput.SrcFuncs, containingFn)
			for _, caller := range callers {
				if isUsed(caller) {
					return true
				}
			}
			// There are no callers in the current package.

			// Function containing the return might be called externally.
			// Find out if it is passed to an external function.
			refs := containingFn.Referrers()
			if refs != nil {
				for _, r := range *refs {
					if isUsed(r) {
						return true
					}
				}
			}
			// Function is not called (inside the current package).
			// See if it is assigned somewhere.
			var buf [10]*ssa.Value
			for _, f := range ssaInput.SrcFuncs {
				for _, b := range f.Blocks {
					for _, instr := range b.Instrs {
						for _, op := range instr.Operands(buf[:0]) {
							if *op == containingFn {
								if isUsed(instr) {
									return true
								}
							}
						}
					}
				}
			}

		case *ssa.Store:
			obj := pass.Pkg.Scope().Lookup(v.Addr.Name())
			if obj != nil {
				return true // Object that was assigned to is global.
			}
			if addr, ok := v.Addr.(ssa.Instruction); ok {
				if isUsed(addr) {
					return true
				}
			}
			refs := v.Referrers()
			if refs != nil {
				for _, r := range *refs {
					if isUsed(r) {
						return true
					}
				}
			}
			use := findUsingInstruction(v, v.Addr)
			if use != nil {
				return isUsed(use)
			}
		case interface {
			Referrers() *[]ssa.Instruction
			Name() string
		}:
			refs := v.Referrers()
			if refs == nil {
				return false
			}
			for _, r := range *refs {
				if isUsed(r) {
					return true
				}
			}
		}
		return false
	}
	for fn, errIdxs := range fnErrs {
		for i, errs := range errIdxs {
			var used bool
		errIdxLoop:
			for _, err := range errs {
				for _, ref := range *err.Referrers() {
					if isUsed(ref) {
						used = true
						break errIdxLoop
					}
				}
			}
			if !used {
				pass.Reportf(returnValPos(fn, i), "error is only ever nil-checked; consider returning a bool instead")
			}
		}

	}
	//reportStable(pass, fnErrs, isUsed)
	return nil, nil
}

// reportStable iterates fnErrs in a stable order and reports the findings.
func reportStable(
	pass *analysis.Pass,
	fnErrs map[*ssa.Function]map[int][]ssa.Value,
	isUsed func(instr ssa.Instruction) bool,
) {
	sortedFns := slices.SortedFunc(maps.Keys(fnErrs), func(a, b *ssa.Function) int {
		return int(a.Pos() - b.Pos())
	})
	for _, fn := range sortedFns {
		for _, i := range slices.Sorted(maps.Keys(fnErrs[fn])) {
			errs := fnErrs[fn][i]
			var used bool
		errIdxLoop:
			for _, err := range errs {
				for _, ref := range *err.Referrers() {
					if isUsed(ref) {
						used = true
						break errIdxLoop
					}
				}
			}
			if !used {
				pass.Reportf(returnValPos(fn, i), "error is only ever nil-checked; consider returning a bool instead")
			}
		}
	}
}

func findReturnedErrors(funcs []*ssa.Function) map[*ssa.Function]map[int][]ssa.Value {
	errs := make(map[*ssa.Function]map[int][]ssa.Value)
	for _, f := range funcs {
		for _, block := range f.Blocks {
			for _, instr := range block.Instrs {
				call, ok := instr.(*ssa.Call)
				if !ok {
					continue
				}
				fn := call.Call.StaticCallee()
				if !isUnexpFunc(fn) {
					continue
				}
				if hasGenericReturns(fn) {
					continue
				}
				rets := returnErrs(call)
				if len(rets) == 0 {
					continue
				}
				// Track only the origin of generic functions.
				// Doing the check here allows us to capture generic errors too.
				if orig := fn.Origin(); orig != nil {
					fn = orig
				}
				if _, ok := errs[fn]; !ok {
					errs[fn] = make(map[int][]ssa.Value)
				}
				for i, v := range rets {
					errs[fn][i] = append(errs[fn][i], v)
				}
			}
		}
	}
	return errs
}

var errType = types.Universe.Lookup("error").Type()

func returnErrs(call *ssa.Call) map[int]ssa.Value {
	if call.Type() == errType {
		return map[int]ssa.Value{0: call}
	}

	if _, ok := call.Type().(*types.Tuple); !ok {
		return nil
	}

	// If the return type is a tuple, find possible extractions of error values.
	errs := make(map[int]ssa.Value, 1) // Most functions will only return a single error.
	for _, r := range *call.Referrers() {
		if extr, ok := r.(*ssa.Extract); ok {
			if extr.Type() == errType {
				errs[extr.Index] = extr
			}
		}
	}
	return errs
}

// findCaller finds all calls to the given function in the SSA representation
func findCaller(funcs []*ssa.Function, fn *ssa.Function) []*ssa.Call {
	calls := make([]*ssa.Call, 0)
	for _, f := range funcs {
		for _, block := range f.Blocks {
			for _, instr := range block.Instrs {
				if call, ok := instr.(*ssa.Call); ok {
					if call.Call.StaticCallee() == fn {
						calls = append(calls, call)
					}
				}
			}
		}
	}

	return calls
}

func findUsingInstruction(after ssa.Instruction, val ssa.Value) ssa.Instruction {
	// Search inside current block first, then in subsequent blocks.
	var buf [10]*ssa.Value
	currBlock := after.Block()
	var search bool
	for _, instr := range currBlock.Instrs {
		if !search {
			if instr == after {
				search = true
				continue
			}
		}
		for _, op := range instr.Operands(buf[:0]) {
			if *op == val {
				return instr
			}
		}
	}
	blocks := currBlock.Parent().Blocks
	for i := currBlock.Index + 1; i < len(blocks); i++ {
		for _, instr := range blocks[i].Instrs {
			for _, op := range instr.Operands(buf[:0]) {
				if *op == val {
					return instr
				}
			}
		}
	}
	return nil
}

func isUnexpFunc(fn *ssa.Function) bool {
	if fn == nil {
		return false
	}
	name := fn.Name()
	if strings.Contains(name, "$") {
		return true // Treat closures the same as unexported functions.
	}
	return !token.IsExported(name)
}

func hasGenericReturns(fn *ssa.Function) bool {
	orig := fn.Origin()
	if orig == nil {
		return false
	}
	sig := orig.Signature
	results := sig.Results()
	for i := 0; i < results.Len(); i++ {
		res := results.At(i)
		if _, ok := res.Type().(*types.TypeParam); ok {
			return true
		}
	}
	return false
}

func returnValPos(fn *ssa.Function, retIdx int) token.Pos {
	if fn == nil {
		return token.NoPos
	}
	sig := fn.Signature
	results := sig.Results()
	res := results.At(retIdx)
	return res.Pos()
}
