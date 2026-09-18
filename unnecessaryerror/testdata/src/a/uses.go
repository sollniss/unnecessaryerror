package a

import (
	"errors"
	"io"
	"os"
)

func MethodCalled() string {
	err := func() error { return errors.New("") }()
	return err.Error() // OK: method called.
}

func TypeAssert() {
	err := func() error { return errors.New("") }()
	if _, ok := err.(*os.SyscallError); ok { // OK: type asserted.
		return
	}
}

func TypeSwitch() {
	err := func() error { return errors.New("") }()
	switch err.(type) { // OK: type switch.
	case *os.PathError:
		return
	}
}

func ComparedToSentinel() {
	err := func() error { return errors.New("") }()
	if err == io.EOF { // OK: compared to a value other than nil.
		return
	}
}

var errSentinel = errors.New("")

func SwitchedOnSentinel() {
	switch func() error { return errors.New("") }() {
	case errSentinel: // OK: compared to a value other than nil.
		return
	}
}

func ComparedToOtherErr(other error) {
	err := func() error { return errors.New("") }()
	if err != nil && err != other { // OK: compared to a value other than nil.
		return
	}
}

func ErrorsIs() {
	err := func() error { return errors.New("") }()
	if errors.Is(err, io.EOF) { // OK: passed to external function.
		return
	}
}

func Panicked() {
	err := func() error { return errors.New("") }()
	if err != nil {
		panic(err) // OK: panicked.
	}
}

type _returnInStruct struct {
	err error
}

func ReturnInStruct() _returnInStruct {
	err := func() error { return errors.New("") }()
	if err != nil {
		return _returnInStruct{err: err} // OK: returned from exported function.
	}
	return _returnInStruct{}
}

func StructField() {
	err := func() error { return errors.New("") }()
	x := struct{ err error }{err: err} // OK: assigned to struct field.
	_ = x
}

func MapIndex() {
	m := map[int]error{}
	m[0] = func() error { return errors.New("") }() // OK: inserted into map.
}

func MapKey(m map[error]bool) {
	err := func() error { return errors.New("") }()
	if m[err] { // OK: used as map key.
		return
	}
}

func SliceIndex() {
	s := make([]error, 1)
	s[0] = func() error { return errors.New("") }() // OK: inserted into slice.
}

func Appended(errs []error) []error {
	err := func() error { return errors.New("") }()
	return append(errs, err) // OK: passed to builtin function.
}

func SentOnChannel(ch chan error) {
	ch <- func() error { return errors.New("") }() // OK: sent on channel.
}

func SentViaSelect(ch chan error) {
	err := func() error { return errors.New("") }()
	select {
	case ch <- err: // OK: sent on channel.
	default:
	}
}

func StoredViaPointer(dst *error) {
	*dst = func() error { return errors.New("") }() // OK: stored through a pointer of unknown origin.
}

var AssignedToExportedGlobalVar1 = func() error { return errors.New("") }() // OK: assigned to global variable.

var AssignedToExportedGlobalVar2_ error

func AssignedToExportedGlobalVar2() {
	AssignedToExportedGlobalVar2_ = func() error { return errors.New("") }() // OK: assigned to global variable.
}

var _assignedToUnexportedGlobalVar = func() error { return errors.New("") }() // OK: assigned to global variable.

var shadowedGlobal error

func ShadowedGlobalVar() {
	var shadowedGlobal error // Shadows the global of the same name.
	func() {
		shadowedGlobal = func() error { return errors.New("") }() // want "error is only ever nil-checked; consider returning a bool instead"
	}()
	if shadowedGlobal != nil {
		return
	}
}

func funcAssignedToGlobalMap() error {
	return errors.New("")
}

func funcAssignedToGlobalMap_1() error {
	err := funcAssignedToGlobalMap()
	if err != nil {
		return err
	}
	return nil
}

var _funcAssignedToGlobalMap = map[string]func() error{
	"a": funcAssignedToGlobalMap_1, // OK: caller inserted into map.
}

var _nilCheckedInGlobalInit = func() error { return errors.New("") }() == nil // want "error is only ever nil-checked; consider returning a bool instead"
