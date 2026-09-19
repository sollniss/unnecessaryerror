// Package a is the testdata for the unnecessaryerror analyzer.
//
// Each case consists of an unexported function returning an error and
// one or more callers. A `want` comment marks functions that should be
// reported; `OK:` comments explain why a function is not reported.
//
// A report is a property of the whole function, so cases cannot share the
// function under test. Cases that only care about what happens to the
// returned value define it inline as an immediately called closure:
//
//	err := func() error { return errors.New("") }() // followed by the want comment
//
// Cases about named functions, call chains, methods, generics or functions
// used as values declare a named function instead.
//
// The cases are grouped by topic:
//
//   - basic.go: plain nil checks
//   - uses.go: values escaping into data structures, comparisons, panics, ...
//   - callchain.go: errors passed through other functions of the package
//   - closures.go: errors captured by or assigned inside closures, named results
//   - funcvalues.go: functions and closures used as values
//   - methods.go: methods, interfaces, method expressions and method values
//   - external.go: values and functions passed to other packages
//   - generics.go: generic functions and types
//   - iterators.go: range-over-func loops
package a

import (
	"errors"
	"log"
)

func ifNotNil() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func IfNotNil() error {
	err := ifNotNil()
	if err != nil {
		return errors.New("")
	}
	return nil
}

func ifNil() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func IfNil() error {
	err := ifNil()
	if err == nil {
		return errors.New("")
	}
	return nil
}

func multiRet() (int, error) { // want "error is only ever nil-checked; consider returning a bool instead"
	return 1, errors.New("")
}

func MultiRet() (int, error) {
	i, err := multiRet()
	if err != nil {
		return 0, errors.New("")
	}
	return i, nil
}

func switchNotNil() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func SwitchNotNil() error {
	err := switchNotNil()
	switch {
	case err != nil:
		return errors.New("")
	}
	return nil
}

func multiFuncReturnErr() error {
	return errors.New("")
}

func MultiFuncReturnErr_1() error {
	err := multiFuncReturnErr()
	if err != nil {
		return nil
	}
	return nil
}

func MultiFuncReturnErr_2() error {
	err := multiFuncReturnErr()
	if err != nil {
		return err // OK: returned from exported function.
	}
	return nil
}

func multiFuncReturnMultiErr() (error, error) { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New(""), errors.New("")
}

func MultiFuncReturnMultiErr() error {
	err1, err2 := multiFuncReturnMultiErr()
	if err1 != nil {
		return err2 // err2 OK: returned from exported function.
	}
	return nil
}

func namedReturn() (err error) { // want "error is only ever nil-checked; consider returning a bool instead"
	err = errors.New("")
	return
}

func NamedReturn() {
	err := namedReturn()
	if err != nil {
		return
	}
}

// functions starting with an underscore are considered unexported.
func _specialUnexported() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func SpecialUnexported() {
	err := _specialUnexported()
	if err != nil {
		return
	}
}

func nilCheckInForLoop() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func NilCheckInForLoop() {
	err := nilCheckInForLoop()
	for err != nil {
		return
	}
}

func assignedToVar() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func AssignedToVar() {
	err := assignedToVar()
	x := err
	_ = x
}

func nilCheckStoredInVar() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func NilCheckStoredInVar() bool {
	failed := nilCheckStoredInVar() != nil // The result of the nil check is used, but not the error itself.
	return failed
}

func discarded() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func Discarded() {
	_ = discarded()
}

func nilCheckViaAny() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func NilCheckViaAny() {
	var v any = nilCheckViaAny()
	if v != nil {
		return
	}
}

func tupleDiscarded() (int, error) { // want "error is only ever nil-checked; consider returning a bool instead"
	return 0, errors.New("")
}

func TupleDiscarded() int {
	x, _ := tupleDiscarded()
	return x
}

func tupleDiscardedOnce() (int, error) {
	return 0, errors.New("")
}

func TupleDiscardedOnce() int {
	x, _ := tupleDiscardedOnce()
	_, err := tupleDiscardedOnce()
	log.Print(err) // OK: passed to external function.
	return x
}

func errorFirst() (error, int) { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New(""), 0
}

func ErrorFirst() int {
	err, x := errorFirst()
	if err != nil {
		return 0
	}
	return x
}
