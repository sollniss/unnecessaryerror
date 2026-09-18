package a

import (
	"errors"
	"io"
	"os"
)

func methodCalled() error {
	return errors.New("")
}

func MethodCalled() string {
	err := methodCalled()
	return err.Error() // OK: method called.
}

func typeAssert() error {
	return errors.New("")
}

func TypeAssert() {
	err := typeAssert()
	if _, ok := err.(*os.SyscallError); ok { // OK: type asserted.
		return
	}
}

func typeSwitch() error {
	return errors.New("")
}

func TypeSwitch() {
	err := typeSwitch()
	switch err.(type) { // OK: type switch.
	case *os.PathError:
		return
	}
}

func comparedToSentinel() error {
	return errors.New("")
}

func ComparedToSentinel() {
	err := comparedToSentinel()
	if err == io.EOF { // OK: compared to a value other than nil.
		return
	}
}

var errSentinel = errors.New("")

func switchedOnSentinel() error {
	return errors.New("")
}

func SwitchedOnSentinel() {
	switch switchedOnSentinel() {
	case errSentinel: // OK: compared to a value other than nil.
		return
	}
}

func comparedToOtherErr() error {
	return errors.New("")
}

func ComparedToOtherErr(other error) {
	if err := comparedToOtherErr(); err != nil && err != other { // OK: compared to a value other than nil.
		return
	}
}

func errorsIs() error {
	return errors.New("")
}

func ErrorsIs() {
	if errors.Is(errorsIs(), io.EOF) { // OK: passed to external function.
		return
	}
}

func panicked() error {
	return errors.New("")
}

func Panicked() {
	if err := panicked(); err != nil {
		panic(err) // OK: panicked.
	}
}

func returnInStruct() error {
	return errors.New("")
}

type _returnInStruct struct {
	err error
}

func ReturnInStruct() _returnInStruct {
	err := returnInStruct()
	if err != nil {
		return _returnInStruct{err: err} // OK: returned from exported function.
	}
	return _returnInStruct{}
}

func structField() error {
	return errors.New("")
}

func StructField() {
	x := struct{ err error }{err: structField()} // OK: assigned to struct field.
	_ = x
}

func mapIndex() error {
	return errors.New("")
}

func MapIndex() {
	m := map[int]error{}
	m[0] = mapIndex() // OK: inserted into map.
}

func mapKey() error {
	return errors.New("")
}

func MapKey(m map[error]bool) {
	if m[mapKey()] { // OK: used as map key.
		return
	}
}

func sliceIndex() error {
	return errors.New("")
}

func SliceIndex() {
	s := make([]error, 1)
	s[0] = sliceIndex() // OK: inserted into slice.
}

func appended() error {
	return errors.New("")
}

func Appended(errs []error) []error {
	return append(errs, appended()) // OK: passed to builtin function.
}

func sentOnChannel() error {
	return errors.New("")
}

func SentOnChannel(ch chan error) {
	ch <- sentOnChannel() // OK: sent on channel.
}

func sentViaSelect() error {
	return errors.New("")
}

func SentViaSelect(ch chan error) {
	select {
	case ch <- sentViaSelect(): // OK: sent on channel.
	default:
	}
}

func storedViaPointer() error {
	return errors.New("")
}

func StoredViaPointer(dst *error) {
	*dst = storedViaPointer() // OK: stored through a pointer of unknown origin.
}

func assignedToExportedGlobalVar1() error {
	return errors.New("")
}

var AssignedToExportedGlobalVar1 = assignedToExportedGlobalVar1() // OK: assigned to global variable.

func assignedToExportedGlobalVar2() error {
	return errors.New("")
}

var AssignedToExportedGlobalVar2_ error

func AssignedToExportedGlobalVar2() {
	AssignedToExportedGlobalVar2_ = assignedToExportedGlobalVar2() // OK: assigned to global variable.
}

func assignedToUnexportedGlobalVar() error {
	return errors.New("")
}

var _assignedToUnexportedGlobalVar = assignedToUnexportedGlobalVar() // OK: assigned to global variable.

var shadowedGlobal error

func shadowedGlobalVar() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func ShadowedGlobalVar() {
	var shadowedGlobal error // Shadows the global of the same name.
	func() {
		shadowedGlobal = shadowedGlobalVar()
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

func nilCheckedInGlobalInit() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

var _nilCheckedInGlobalInit = nilCheckedInGlobalInit() == nil
