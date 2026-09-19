package a

import (
	"errors"
	"log"
)

// Generic methods (Go 1.27) are instantiated like generic functions: every
// instantiation is a wrapper that forwards to the single generic body, which
// is where the diagnostic is reported.

type genericMethod struct{}

func (genericMethod) explicit[T any](v T) error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func GenericMethodExplicit() {
	if (genericMethod{}).explicit[int](1) != nil {
		return
	}
}

func (genericMethod) inferred[T any](v T) error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func GenericMethodInferred() {
	if (genericMethod{}).inferred(1) != nil {
		return
	}
}

func (genericMethod) logged[T any](v T) error {
	return errors.New("")
}

func GenericMethodLogged() {
	log.Print((genericMethod{}).logged(1)) // OK: passed to external function.
}

func (genericMethod) twoInstantiations[T any](v T) error {
	return errors.New("")
}

func GenericMethodTwoInstantiations() {
	log.Print((genericMethod{}).twoInstantiations("")) // OK: passed to external function.
	if (genericMethod{}).twoInstantiations(1) != nil {
		return
	}
}

func (genericMethod) tuple[T any](v T) (T, error) { // want "error is only ever nil-checked; consider returning a bool instead"
	return v, errors.New("")
}

func GenericMethodTuple() {
	if _, err := (genericMethod{}).tuple(1); err != nil {
		return
	}
}

func (genericMethod) tupleOfErrors[T any]() (T, error) { // want "error is only ever nil-checked; consider returning a bool instead"
	var zero T
	return zero, errors.New("")
}

func GenericMethodTupleOfErrors() {
	v, err := (genericMethod{}).tupleOfErrors[error]()
	log.Print(v) // The generic result is not tracked, even when instantiated with error.
	if err != nil {
		return
	}
}

func (genericMethod) genericResult[T any]() T { // Not a candidate: generic result.
	var zero T
	return zero
}

func GenericMethodGenericResult() {
	if (genericMethod{}).genericResult[error]() != nil {
		return
	}
}

func (genericMethod) Exported[T any](v T) error { // Not a candidate: exported method name.
	return errors.New("")
}

func GenericMethodExported() {
	if (genericMethod{}).Exported(1) != nil {
		return
	}
}

func (genericMethod) methodValue[T any](v T) error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func GenericMethodValue() {
	f := genericMethod{}.methodValue[int]
	if f(1) != nil {
		return
	}
}

func (genericMethod) methodValueLogged[T any](v T) error {
	return errors.New("")
}

func GenericMethodValueLogged() {
	f := genericMethod{}.methodValueLogged[int]
	log.Print(f(1)) // OK: passed to external function via method value.
}

func (genericMethod) methodValueInferred[T any](v T) error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func GenericMethodValueInferred() {
	var f func(int) error = genericMethod{}.methodValueInferred
	if f(1) != nil {
		return
	}
}

func (genericMethod) methodExpr[T any](v T) error {
	return errors.New("")
}

var _genericMethodExpr = map[string]func(genericMethod, int) error{
	"a": genericMethod.methodExpr[int], // OK: method expression inserted into map.
}

func (genericMethod) methodExprNilChecked[T any](v T) error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func GenericMethodExprNilChecked() {
	if genericMethod.methodExprNilChecked[int](genericMethod{}, 1) != nil {
		return
	}
}

func (genericMethod) passedToLocal[T any](v T) error {
	return errors.New("")
}

func GenericMethodPassedToLocal() {
	runAndLogInt(genericMethod{}.passedToLocal[int]) // OK: the method's result is passed to an external function.
}

func (genericMethod) passedToNilCheckingLocal[T any](v T) error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func GenericMethodPassedToNilCheckingLocal() {
	if runAndNilCheckInt(genericMethod{}.passedToNilCheckingLocal[int]) {
		return
	}
}

func (genericMethod) onlyDeferred[T any]() error { // Not a candidate: never called for its result.
	return errors.New("")
}

func GenericMethodOnlyDeferred() {
	defer (genericMethod{}).onlyDeferred[int]()
	go (genericMethod{}).onlyDeferred[int]()
}

// Errors flowing through generic methods.

func genericMethodChain() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func (genericMethod) chain[T any]() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return genericMethodChain()
}

func GenericMethodChain() {
	if (genericMethod{}).chain[int]() != nil {
		return
	}
}

func genericMethodChainReturned() error {
	return errors.New("")
}

func (genericMethod) chainReturned[T any]() error {
	return genericMethodChainReturned()
}

func GenericMethodChainReturned() error {
	return (genericMethod{}).chainReturned[int]() // OK: returned from exported function.
}

func (genericMethod) nilCheck[T any](err error) bool {
	return err != nil
}

func GenericMethodNilCheckingParam() {
	err := func() error { return errors.New("") }() // want "error is only ever nil-checked; consider returning a bool instead"
	if (genericMethod{}).nilCheck[int](err) {
		return
	}
}

func (genericMethod) logErr[T any](err error) {
	log.Print(err)
}

func GenericMethodLoggingParam() {
	err := func() error { return errors.New("") }()
	(genericMethod{}).logErr[int](err) // OK: the parameter is passed to an external function.
}

func (genericMethod) closureInside[T any]() {
	err := func() error { return errors.New("") }() // want "error is only ever nil-checked; consider returning a bool instead"
	if err != nil {
		return
	}
}

func GenericMethodClosureInside() {
	(genericMethod{}).closureInside[int]()
}

func (genericMethod) closureReturned[T any]() func() error {
	return func() error { return errors.New("") } // want "error is only ever nil-checked; consider returning a bool instead"
}

func GenericMethodClosureReturned() {
	if (genericMethod{}).closureReturned[int]()() != nil {
		return
	}
}

func (b genericMethod) recursive[T any](n int) error { // want "error is only ever nil-checked; consider returning a bool instead"
	if n > 0 {
		return b.recursive[T](n - 1)
	}
	return errors.New("")
}

func GenericMethodRecursive() {
	if (genericMethod{}).recursive[int](2) != nil {
		return
	}
}

// Generic methods reached through generic functions and types.

func (genericMethod) calledFromGeneric[T any](v T) error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func callGenericMethod[T any](m genericMethod, v T) error { // want "error is only ever nil-checked; consider returning a bool instead"
	return m.calledFromGeneric[T](v)
}

func GenericMethodCalledFromGeneric() {
	if callGenericMethod(genericMethod{}, 1) != nil {
		return
	}
}

func (genericMethod) boundInGeneric[T any](v T) error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func bindGenericMethod[T any](m genericMethod) func(T) error {
	return m.boundInGeneric[T]
}

func GenericMethodBoundInGeneric() {
	if bindGenericMethod[int](genericMethod{})(1) != nil {
		return
	}
}

type genericMethodOuter struct{ genericMethod }

func (genericMethod) promoted[T any](v T) error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func GenericMethodPromoted() {
	if (genericMethodOuter{}).promoted(1) != nil {
		return
	}
}

type genericTypeMethod[U any] struct{}

func (genericTypeMethod[U]) both[T any](v T, u U) error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func GenericTypeMethodBoth() {
	if (genericTypeMethod[string]{}).both(1, "") != nil {
		return
	}
}

func (*genericTypeMethod[U]) viaInterface[T any](v T) error {
	return errors.New("")
}

func (g *genericTypeMethod[U]) plain() error {
	return g.viaInterface(1)
}

type genericTypeMethodIface interface{ plain() error }

func GenericTypeMethodViaInterface() {
	var i genericTypeMethodIface = &genericTypeMethod[string]{}
	log.Print(i.plain()) // OK: passed to external function via interface.
	if (&genericTypeMethod[int]{}).viaInterface(1) != nil {
		return
	}
}

// The helpers above take a func(int) error rather than reusing the func() error
// helpers of closures.go: as of x/tools v0.50.0 the SSA builder crashes when a
// method value of a generic method has the same parameters and results as a
// generic function instantiated later in the package (it canonicalizes the two
// signatures to one, ignoring the receiver), and generics.go instantiates
// several generic functions to func() error.

func runAndLogInt(fn func(int) error) {
	log.Print(fn(1))
}

func runAndNilCheckInt(fn func(int) error) bool {
	return fn(1) != nil
}
