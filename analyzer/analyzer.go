package analyzer

import (
	"go/token"
	"go/types"
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
	seen := make(map[ssa.Instruction]struct{})

	var isUsed func(instr ssa.Instruction) bool
	isUsed = func(instr ssa.Instruction) bool {
		if _, ok := seen[instr]; ok {
			return false
		}
		seen[instr] = struct{}{}

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
			for _, r := range *v.Referrers() {
				isUsed(r)
			}
		case *ssa.Return:
			if !isUnexpFunc(v.Parent()) {
				return true
			}
			callers := findCaller(ssaInput.SrcFuncs, v.Parent())
			for _, caller := range callers {
				isUsed(caller)
			}
		case *ssa.Store:
			if addr, ok := v.Addr.(ssa.Instruction); ok {
				if isUsed(addr) {
					return true
				}
			}
		case interface {
			Referrers() *[]ssa.Instruction
			Name() string
		}:
			for _, r := range *v.Referrers() {
				if isUsed(r) {
					return true
				}
			}
		}
		return false
	}
	for fn, errIdxs := range fnErrs {
		for i, errs := range errIdxs {
			var ok bool
		errIdxLoop:
			for _, err := range errs {
				for _, ref := range *err.Referrers() {
					if isUsed(ref) {
						ok = true
						break errIdxLoop
					}
				}
			}
			if !ok {
				pass.Reportf(returnValPos(fn, i), "error is only ever nil-checked; consider returning a bool instead")
			}
		}

	}
	return nil, nil
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
