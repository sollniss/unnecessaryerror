package a

import (
	"errors"
	"fmt"
	"io"
	"log"
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

func WrappedWithErrorf() error {
	err := func() error { return errors.New("") }()
	return fmt.Errorf("wrap: %w", err) // OK: passed to external function.
}

func Phi(cond bool) {
	var err error
	if cond {
		err = func() error { return errors.New("") }()
	} else {
		err = func() error { return errors.New("") }()
	}
	log.Print(err) // OK: passed to external function.
}

func PhiNilChecked(cond bool) {
	var err error
	if cond {
		err = func() error { return errors.New("") }() // want "error is only ever nil-checked; consider returning a bool instead"
	} else {
		err = func() error { return errors.New("") }() // want "error is only ever nil-checked; consider returning a bool instead"
	}
	if err != nil {
		return
	}
}

func PhiInLoop() {
	var err error
	for i := 0; i < 3; i++ {
		err = func() error { return errors.New("") }()
		if err == nil {
			break
		}
	}
	log.Print(err) // OK: passed to external function.
}

func ForInit() {
	for err := func() error { return errors.New("") }(); err != nil; err = nil { // want "error is only ever nil-checked; consider returning a bool instead"
		return
	}
}

func LoadedThroughPointer() {
	err := func() error { return errors.New("") }()
	p := &err
	log.Print(*p) // OK: passed to external function.
}

func LoadedThroughPointerNilChecked() {
	err := func() error { return errors.New("") }() // want "error is only ever nil-checked; consider returning a bool instead"
	p := &err
	if *p != nil {
		return
	}
}

func AddressPassedToExternal() {
	err := func() error { return errors.New("") }()
	log.Print(&err) // OK: address passed to external function.
}

type errAlias = error

func aliasReturn() errAlias { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func AliasReturn() {
	if aliasReturn() != nil {
		return
	}
}

type namedErrIface interface{ Error() string }

func namedIfaceReturn() namedErrIface { // Not an error: not tracked.
	return errors.New("")
}

func NamedIfaceReturn() {
	if namedIfaceReturn() != nil {
		return
	}
}

type convertedErr error

func ConvertedToNamedType() {
	err := convertedErr(func() error { return errors.New("") }())
	log.Print(err) // OK: passed to external function.
}

func ConvertedToNamedTypeNilChecked() {
	err := convertedErr(func() error { return errors.New("") }()) // want "error is only ever nil-checked; consider returning a bool instead"
	if err != nil {
		return
	}
}

func SwitchNilAndSentinel() {
	switch func() error { return errors.New("") }() {
	case nil:
	case io.EOF: // OK: compared to a value other than nil.
	}
}

func SwitchNilOnly() {
	switch func() error { return errors.New("") }() { // want "error is only ever nil-checked; consider returning a bool instead"
	case nil:
	default:
	}
}

func MapKeyCommaOk(m map[error]bool) {
	err := func() error { return errors.New("") }()
	if _, ok := m[err]; ok { // OK: used as map key.
		return
	}
}

func MapUpdateKey(m map[error]bool) {
	m[func() error { return errors.New("") }()] = true // OK: used as map key.
}

func ArrayElement() {
	var arr [1]error
	arr[0] = func() error { return errors.New("") }() // OK: inserted into array.
}

func ArrayLiteral() {
	arr := [...]error{func() error { return errors.New("") }()} // OK: inserted into array.
	_ = arr
}

func StructInChannel(ch chan struct{ err error }) {
	ch <- struct{ err error }{err: func() error { return errors.New("") }()} // OK: assigned to struct field.
}

func MethodValueOfError() {
	err := func() error { return errors.New("") }()
	f := err.Error // OK: method value.
	_ = f
}

func MethodExprOnError() string {
	err := func() error { return errors.New("") }()
	return error.Error(err) // OK: method called via method expression.
}

func TypeAssertToAnonInterface() {
	err := func() error { return errors.New("") }()
	if _, ok := err.(interface{ Timeout() bool }); ok { // OK: type asserted.
		return
	}
}

func ComparedToEachOther() {
	a := func() error { return errors.New("") }()
	b := func() error { return errors.New("") }()
	if a == b { // OK: compared to a value other than nil.
		return
	}
}

func nilCheckedInInit() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func init() {
	if nilCheckedInInit() != nil {
		return
	}
}
