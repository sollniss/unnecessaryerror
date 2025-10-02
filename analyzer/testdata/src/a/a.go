package a

import (
	"errors"
	"log"
	"os"
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

func MultiFuncReturnErr_noret() error {
	err := multiFuncReturnErr()
	if err != nil {
		return nil
	}
	return nil
}

func MultiFuncReturnErr_ret() error {
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

type receiver struct{}

func (r *receiver) methodReturnsErr() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func (r *receiver) MethodReturnsErr() {
	err := r.methodReturnsErr()
	if err != nil {
		return
	}
}

func passToExternalFunc() error {
	return errors.New("")
}

func PassToExternalFunc() {
	err := passToExternalFunc()
	log.Print(err) // OK: passed to external function.
}

func passToBuiltinFunc() error {
	return errors.New("")
}

func PassToBuiltinFunc() {
	err := passToBuiltinFunc()
	println(err) // OK: passed to builtin function.
}

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

func assignedToVar() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func AssignedToVar() {
	err := assignedToVar()
	x := err
	_ = x
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

func sliceIndex() error {
	return errors.New("")
}

func SliceIndex() {
	s := make([]error, 1)
	s[0] = sliceIndex() // OK: inserted into slice.
}

func Closure() {
	err := func() error { return errors.New("") }() // want "error is only ever nil-checked; consider returning a bool instead"
	if err != nil {
		return
	}
	return
}

func assignedToExportedGlobalVar() error {
	return errors.New("")
}

var AssignedToExportedGlobalVar = assignedToExportedGlobalVar() // OK: assigned to global variable.

func assignedToUnexportedGlobalVar() error {
	return errors.New("")
}

var _assignedToUnexportedGlobalVar = assignedToUnexportedGlobalVar() // OK: assigned to global variable.

func genericFunc1[T any]() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func GenericFunc1A() {
	err := genericFunc1[int]()
	if err != nil {
		return
	}
}

func GenericFunc1B() {
	err := genericFunc1[string]()
	if err != nil {
		return
	}
}

func genericFunc2[T any]() error {
	return errors.New("")
}

func GenericFunc2A() {
	err := genericFunc2[int]()
	if err != nil {
		return
	}
}

func GenericFunc2B() error {
	err := genericFunc2[string]()
	if err != nil {
		return err // OK: returned from exported function.
	}
	return nil
}

func genericError[T any]() T {
	var zero T
	return zero
}

func GenericError() {
	if err := genericError[error](); err != nil { // OK: ignore generic return types.
		return
	}
	return
}

func genericError2[T error]() T {
	var zero T
	return zero
}

func GenericError2() {
	if err := genericError2[error](); err != nil { // OK: ignore generic return types.
		return
	}
	return
}

func callchain1() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func callchain1A(err error) error { // want "error is only ever nil-checked; consider returning a bool instead"
	return err
}

func callchain1B(err error) error { // want "error is only ever nil-checked; consider returning a bool instead"
	return callchain1A(err)
}

func Callchain1() error {
	err := callchain1A(callchain1B(callchain1()))
	if err != nil {
		return err
	}
	return nil
}

func callchain2() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func callchain2A() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return callchain2()
}

func callchain2B() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return callchain2A()
}
func Callchain2() {
	err := callchain2B()
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

func nilCheckInDefer() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func NilCheckInDefer() {
	err := nilCheckInDefer()
	defer func() {
		if err != nil {
			return
		}
	}()
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
