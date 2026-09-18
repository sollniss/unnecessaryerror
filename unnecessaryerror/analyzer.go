package unnecessaryerror

import (
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

// Analyzer that finds errors that are only ever nil-checked and could be replaced with a bool.
var Analyzer = &analysis.Analyzer{
	Name:     "unnecessaryerror",
	Doc:      "Finds errors that are only ever nil-checked and could be replaced with a bool.",
	Run:      run,
	Requires: []*analysis.Analyzer{buildssa.Analyzer},
}

const msg = "error is only ever nil-checked; consider returning a bool instead"

func run(pass *analysis.Pass) (any, error) {
	ssaInput := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)

	funcs := ssaInput.SrcFuncs
	// Package-level variable initializers live in the synthetic package
	// initializer, which is not among the source functions.
	if init := ssaInput.Pkg.Func("init"); init != nil {
		funcs = append(slices.Clone(funcs), init)
	}
	idx := buildIndex(funcs)
	t := &tracker{
		pkg:  ssaInput.Pkg,
		idx:  idx,
		memo: make(map[node]bool),
	}

	var diags []token.Pos
	for fn, errIdxs := range findReturnedErrors(idx) {
		if t.fnUsed(fn) {
			continue // Passed around as a value; its results may be used anywhere.
		}
		for i, errs := range errIdxs {
			if !slices.ContainsFunc(errs, t.valueUsed) {
				diags = append(diags, retValPos(fn, i))
			}
		}
	}
	// Map iteration order is random; report in source order.
	slices.Sort(diags)
	for _, pos := range diags {
		pass.Reportf(pos, msg)
	}
	return nil, nil
}

// index holds package-wide lookup tables built once per package.
type index struct {
	// funcs contains all functions with a body that are reachable from the
	// source functions, including synthetic wrappers (thunks, bound method
	// wrappers, generic instantiation wrappers) that the SSA builder emits
	// for method expressions, method values and generic calls.
	funcs []*ssa.Function
	// callers maps a function to all static calls to it.
	callers map[*ssa.Function][]*ssa.Call
	// invokes maps a method name to all dynamic (interface) calls of that name.
	invokes map[string][]*ssa.Call
	// fnUses maps a function to its uses as a value, i.e. anywhere but the
	// callee position of a static call. Instantiations of generic functions
	// are keyed by their origin.
	fnUses map[*ssa.Function][]fnUse
	// closures maps a function to the instructions that create closures from it.
	closures map[*ssa.Function][]*ssa.MakeClosure
}

func buildIndex(src []*ssa.Function) *index {
	idx := &index{
		callers:  make(map[*ssa.Function][]*ssa.Call),
		invokes:  make(map[string][]*ssa.Call),
		fnUses:   make(map[*ssa.Function][]fnUse),
		closures: make(map[*ssa.Function][]*ssa.MakeClosure),
	}
	seen := make(map[*ssa.Function]bool)
	work := slices.Clone(src)
	var buf [10]*ssa.Value
	for len(work) > 0 {
		fn := work[len(work)-1]
		work = work[:len(work)-1]
		if seen[fn] || len(fn.Blocks) == 0 {
			continue
		}
		seen[fn] = true
		idx.funcs = append(idx.funcs, fn)

		for _, block := range fn.Blocks {
			for _, instr := range block.Instrs {
				var calleePos *ssa.Value
				switch v := instr.(type) {
				case ssa.CallInstruction:
					common := v.Common()
					calleePos = &common.Value
					call, isCall := instr.(*ssa.Call)
					if common.IsInvoke() {
						if isCall {
							name := common.Method.Name()
							idx.invokes[name] = append(idx.invokes[name], call)
						}
					} else if callee := common.StaticCallee(); callee != nil {
						if isCall {
							idx.callers[callee] = append(idx.callers[callee], call)
						}
						work = append(work, callee)
					}
				case *ssa.MakeClosure:
					closure := v.Fn.(*ssa.Function)
					idx.closures[closure] = append(idx.closures[closure], v)
					work = append(work, closure)
				}
				for _, op := range instr.Operands(buf[:0]) {
					if op == calleePos {
						continue
					}
					if f, ok := (*op).(*ssa.Function); ok {
						key := origin(f)
						idx.fnUses[key] = append(idx.fnUses[key], fnUse{f, instr})
						// Method expressions and instantiations used as values
						// are wrappers whose bodies call the actual function.
						work = append(work, f)
					}
				}
			}
		}
	}
	return idx
}

// fnUse is an instruction that uses fn as a value.
type fnUse struct {
	fn    *ssa.Function
	instr ssa.Instruction
}

// origin returns the generic function fn was instantiated from, or fn itself.
func origin(fn *ssa.Function) *ssa.Function {
	if orig := fn.Origin(); orig != nil {
		return orig
	}
	return fn
}

// node is a value together with one of its referrers.
type node struct {
	val   ssa.Value
	instr ssa.Instruction
}

// tracker follows values through the package to decide whether they are used
// in any way other than a nil check.
type tracker struct {
	pkg *ssa.Package
	idx *index
	// memo caches results that hold independently of the query they were found in.
	memo map[node]bool
	// visited holds the nodes reached by the current query.
	visited map[node]bool
}

// valueUsed reports whether val is used in some way other than a nil check.
func (t *tracker) valueUsed(val ssa.Value) bool {
	return t.query(func() bool { return t.refsUsed(val) })
}

// fnUsed reports whether fn is used as a value in a way that may use its results.
func (t *tracker) fnUsed(fn *ssa.Function) bool {
	return t.query(func() bool { return t.fnRefsUsed(fn) })
}

// fnRefsUsed reports whether any use of fn as a value is a use.
func (t *tracker) fnRefsUsed(fn *ssa.Function) bool {
	return slices.ContainsFunc(t.idx.fnUses[origin(fn)], func(use fnUse) bool {
		return t.used(use.fn, use.instr)
	})
}

// query runs search as a new query.
func (t *tracker) query(search func() bool) bool {
	t.visited = make(map[node]bool)
	if search() {
		return true
	}
	// Nothing reachable from the query is a use, so the same holds for every
	// node visited on the way.
	for n := range t.visited {
		t.memo[n] = false
	}
	return false
}

// refsUsed reports whether any referrer of val uses it.
func (t *tracker) refsUsed(val ssa.Value) bool {
	refs := val.Referrers()
	if refs == nil {
		return false
	}
	return slices.ContainsFunc(*refs, func(instr ssa.Instruction) bool {
		return t.used(val, instr)
	})
}

// used reports whether instr uses val.
func (t *tracker) used(val ssa.Value, instr ssa.Instruction) bool {
	n := node{val, instr}
	if r, ok := t.memo[n]; ok {
		return r
	}
	if t.visited[n] {
		return false // Cycle; the answer is determined by the node that started it.
	}
	t.visited[n] = true
	if t.usedBy(val, instr) {
		t.memo[n] = true // A use stays a use no matter where the query started.
		return true
	}
	return false
}

func (t *tracker) usedBy(val ssa.Value, instr ssa.Instruction) bool {
	switch v := instr.(type) {
	case *ssa.If, *ssa.DebugRef:
		return false
	case *ssa.BinOp:
		// Comparing against nil is the one thing we are looking for.
		// Anything else, e.g. comparing against a sentinel error, is a use.
		return !isNilCheck(v)
	case ssa.CallInstruction:
		return t.callUses(val, v)
	case *ssa.Return:
		return t.returnUses(val, v)
	case *ssa.Store:
		if v.Val == val {
			return t.addrUsed(v.Addr)
		}
		return false // val is the address being overwritten.
	case *ssa.MakeClosure:
		if v.Fn == val {
			return t.refsUsed(v)
		}
		closure := v.Fn.(*ssa.Function)
		for i, b := range v.Bindings {
			if b == val && t.refsUsed(closure.FreeVars[i]) {
				return true
			}
		}
		return false
	case *ssa.Lookup:
		if v.Index == val {
			return true // Used as map key.
		}
		return t.refsUsed(v)
	case *ssa.Phi, *ssa.Extract, *ssa.UnOp,
		*ssa.MakeInterface, *ssa.ChangeInterface, *ssa.ChangeType, *ssa.Convert, *ssa.MultiConvert:
		// val is passed through unchanged; follow the result.
		return t.refsUsed(v.(ssa.Value))
	default:
		// Anything else (stored into a struct, slice, map or channel, panicked,
		// type asserted, ...) is a use.
		return true
	}
}

// callUses reports whether the call instr uses val.
func (t *tracker) callUses(val ssa.Value, instr ssa.CallInstruction) bool {
	common := instr.Common()
	if common.IsInvoke() {
		// Either a method is called on val or val is passed to an unknown implementation.
		return true
	}
	if common.Value == val {
		// val itself is called; it is used if its result is.
		if call, ok := instr.(*ssa.Call); ok {
			return t.refsUsed(call)
		}
		return false // go and defer discard the result.
	}
	callee := common.StaticCallee()
	if callee == nil || len(callee.Blocks) == 0 || (callee.Pkg != nil && callee.Pkg != t.pkg) {
		return true // Passed to a builtin, a function value or an external function.
	}
	// Passed to a function of this package; follow the parameter inside it.
	for i, arg := range common.Args {
		if arg != val {
			continue
		}
		if i >= len(callee.Params) || t.refsUsed(callee.Params[i]) {
			return true
		}
	}
	return false
}

// returnUses reports whether returning val from instr's function counts as a use.
func (t *tracker) returnUses(val ssa.Value, instr *ssa.Return) bool {
	fn := instr.Parent()
	if !isUnexportedFunc(fn) {
		return true // Could be used by other packages.
	}
	results := fn.Signature.Results()
	for _, call := range t.idx.callers[fn] {
		if t.resultUsed(val, instr, call, results) {
			return true
		}
	}
	if recv := fn.Signature.Recv(); recv != nil {
		for _, call := range t.idx.invokes[fn.Name()] {
			if !implementedBy(call.Common(), fn) {
				continue
			}
			if t.resultUsed(val, instr, call, results) {
				return true
			}
		}
	}
	return t.fnRefsUsed(fn)
}

// resultUsed reports whether the result of call that corresponds to val in ret is used.
func (t *tracker) resultUsed(val ssa.Value, ret *ssa.Return, call *ssa.Call, results *types.Tuple) bool {
	if results.Len() == 1 {
		return t.refsUsed(call)
	}
	for i, res := range ret.Results {
		if res != val {
			continue
		}
		for _, ref := range *call.Referrers() {
			if extr, ok := ref.(*ssa.Extract); ok && extr.Index == i && t.refsUsed(extr) {
				return true
			}
		}
	}
	return false
}

// addrUsed reports whether a value stored at addr is used.
func (t *tracker) addrUsed(addr ssa.Value) bool {
	switch a := addr.(type) {
	case *ssa.Alloc:
		return t.refsUsed(a)
	case *ssa.FreeVar:
		if t.refsUsed(a) {
			return true
		}
		// The variable is shared with the enclosing function.
		closure := a.Parent()
		i := slices.Index(closure.FreeVars, a)
		for _, mc := range t.idx.closures[closure] {
			if t.addrUsed(mc.Bindings[i]) {
				return true
			}
		}
		return false
	default:
		// Globals, fields, elements and pointers of unknown origin.
		return true
	}
}

// implementedBy reports whether the dynamic call could dispatch to the method fn.
func implementedBy(call *ssa.CallCommon, fn *ssa.Function) bool {
	if !types.Identical(call.Method.Type(), fn.Signature) {
		return false
	}
	iface, ok := call.Value.Type().Underlying().(*types.Interface)
	if !ok {
		return false
	}
	recv := fn.Signature.Recv().Type()
	if types.Implements(recv, iface) {
		return true
	}
	if _, isPtr := recv.(*types.Pointer); !isPtr {
		return types.Implements(types.NewPointer(recv), iface)
	}
	return false
}

// isNilCheck reports whether op compares against nil.
func isNilCheck(op *ssa.BinOp) bool {
	if op.Op != token.EQL && op.Op != token.NEQ {
		return false
	}
	return isNil(op.X) || isNil(op.Y)
}

func isNil(v ssa.Value) bool {
	c, ok := v.(*ssa.Const)
	return ok && c.IsNil()
}

// findReturnedErrors finds all unexported functions that return an error,
// mapped to the returned errors by result index.
func findReturnedErrors(idx *index) map[*ssa.Function]map[int][]ssa.Value {
	errs := make(map[*ssa.Function]map[int][]ssa.Value)
	add := func(fn *ssa.Function, call *ssa.Call) {
		rets := returnErrs(call)
		if len(rets) == 0 {
			return
		}
		if _, ok := errs[fn]; !ok {
			errs[fn] = make(map[int][]ssa.Value)
		}
		for i, v := range rets {
			errs[fn][i] = append(errs[fn][i], v)
		}
	}
	for _, f := range idx.funcs {
		for _, block := range f.Blocks {
			for _, instr := range block.Instrs {
				call, ok := instr.(*ssa.Call)
				if !ok {
					continue
				}
				if fn := reportedFunc(call.Call.StaticCallee()); fn != nil {
					add(fn, call)
				}
			}
		}
	}
	// Methods may also be called through an interface.
	for fn := range errs {
		if fn.Signature.Recv() == nil {
			continue
		}
		for _, call := range idx.invokes[fn.Name()] {
			if implementedBy(call.Common(), fn) {
				add(fn, call)
			}
		}
	}
	return errs
}

// reportedFunc returns the source function that a diagnostic about fn
// should be attached to, or nil if fn is not a candidate.
func reportedFunc(fn *ssa.Function) *ssa.Function {
	if !isUnexportedFunc(fn) {
		return nil
	}
	// Track generic functions by their origin, so that all
	// instantiations count towards the same function.
	fn = origin(fn)
	if fn.Synthetic != "" {
		return nil // Wrapper; the wrapped function is reported instead.
	}
	if hasGenericReturns(fn) {
		return nil
	}
	return fn
}

var errType = types.Universe.Lookup("error").Type()

// returnErrs returns all errors returned from the [ssa.Call] mapped to their index.
//
// Does not check whether the underlying type is an error.
func returnErrs(call *ssa.Call) map[int]ssa.Value {
	if types.Identical(call.Type(), errType) {
		return map[int]ssa.Value{0: call}
	}

	if _, ok := call.Type().(*types.Tuple); !ok {
		return nil
	}

	// If the return type is a tuple, find possible extractions of error values.
	errs := make(map[int]ssa.Value, 1) // Most functions will only return a single error.
	for _, r := range *call.Referrers() {
		if extr, ok := r.(*ssa.Extract); ok {
			if types.Identical(extr.Type(), errType) {
				errs[extr.Index] = extr
			}
		}
	}
	return errs
}

// isUnexportedFunc returns true if fn is unexported or a closure.
func isUnexportedFunc(fn *ssa.Function) bool {
	if fn == nil {
		return false
	}
	name := fn.Name()
	if strings.Contains(name, "$") {
		return true // Treat closures the same as unexported functions.
	}
	return !ast.IsExported(name)
}

// hasGenericReturns returns true if fn has generic return values, i.e.
//
//	func[T any]() T
func hasGenericReturns(fn *ssa.Function) bool {
	for res := range fn.Signature.Results().Variables() {
		if _, ok := res.Type().(*types.TypeParam); ok {
			return true
		}
	}
	return false
}

// retValPos returns the position of the nth return value of fn.
func retValPos(fn *ssa.Function, n int) token.Pos {
	return fn.Signature.Results().At(n).Pos()
}
