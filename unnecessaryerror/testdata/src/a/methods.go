package a

import (
	"errors"
	"log"

	"a/external"
)

type Receiver1 struct{}

func (r *Receiver1) methodReturnsErr() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func (r *Receiver1) MethodReturnsErr() {
	err := r.methodReturnsErr()
	if err != nil {
		return
	}
}

type calledViaInterface interface {
	calledViaInterface() error
}

type _calledViaInterface struct{}

func (_calledViaInterface) calledViaInterface() error {
	return errors.New("")
}

func CalledViaInterface(c calledViaInterface) {
	if err := c.calledViaInterface(); err != nil {
		log.Print(err) // OK: passed to external function via interface call.
	}
}

func CalledViaInterface_1() {
	if err := (_calledViaInterface{}).calledViaInterface(); err != nil {
		return
	}
}

type nilCheckedViaInterface interface {
	nilCheckedViaInterface() error
}

type _nilCheckedViaInterface struct{}

func (*_nilCheckedViaInterface) nilCheckedViaInterface() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func NilCheckedViaInterface(c nilCheckedViaInterface) {
	if err := c.nilCheckedViaInterface(); err != nil {
		return
	}
}

func NilCheckedViaInterface_1() {
	if err := (&_nilCheckedViaInterface{}).nilCheckedViaInterface(); err != nil {
		return
	}
}

type sameNameInterface interface {
	sameName() error
}

type _sameNameMatching struct{}

func (_sameNameMatching) sameName() error {
	return errors.New("")
}

type _sameNameMismatching struct{}

func (_sameNameMismatching) sameName(int) error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func SameNameInterface(s sameNameInterface) {
	log.Print(s.sameName()) // OK for _sameNameMatching: may be the dynamic target.
	if (_sameNameMatching{}).sameName() != nil {
		return
	}
	if (_sameNameMismatching{}).sameName(0) != nil { // Different signature: cannot be the dynamic target.
		return
	}
}

type methodExpr struct{}

func methodExprErr() error {
	return errors.New("")
}

func (methodExpr) methodExpr() error {
	return methodExprErr()
}

func MethodExpr() error {
	return methodExpr.methodExpr(methodExpr{}) // OK: returned from exported function via method expression.
}

type methodValue struct{}

func methodValueErr() error {
	return errors.New("")
}

func (methodValue) methodValue() error {
	return methodValueErr()
}

func MethodValue() {
	f := methodValue{}.methodValue
	log.Print(f()) // OK: passed to external function via method value.
}

type promotedInner struct{}

func promotedErr() error {
	return errors.New("")
}

func (promotedInner) promoted() error {
	return promotedErr()
}

type promotedOuter struct{ promotedInner }

func Promoted() error {
	return promotedOuter{}.promoted() // OK: returned from exported function via promoted method.
}

type _methodAssignedToExternalStructField struct{}

func (s _methodAssignedToExternalStructField) methodAssignedToExternalStructField() error {
	return errors.New("")
}

func (s _methodAssignedToExternalStructField) methodAssignedToExternalStructField1() error {
	err := s.methodAssignedToExternalStructField()
	if err != nil {
		return err
	}
	return nil
}

func MethodAssignedToExternalStructField() {
	var x _methodAssignedToExternalStructField
	s := external.MethodAssignedToExternalStructField{
		Func: x.methodAssignedToExternalStructField1, // OK: method value assigned to struct field.
	}
	_ = s
}

type methodExprValue struct{}

func methodExprValueErr() error {
	return errors.New("")
}

func (methodExprValue) methodExprValue() error {
	return methodExprValueErr()
}

var _methodExprValue = map[string]func(methodExprValue) error{
	"a": methodExprValue.methodExprValue, // OK: method expression inserted into map.
}

type methodExprPassedToLocal struct{}

func methodExprPassedToLocalErr() error {
	return errors.New("")
}

func (methodExprPassedToLocal) methodExprPassedToLocal() error {
	return methodExprPassedToLocalErr()
}

func callAndLog(f func(methodExprPassedToLocal) error) {
	log.Print(f(methodExprPassedToLocal{}))
}

func MethodExprPassedToLocal() {
	callAndLog(methodExprPassedToLocal.methodExprPassedToLocal) // OK: the result is passed to an external function.
}

type methodExprNilChecked struct{}

func methodExprNilCheckedErr() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func (methodExprNilChecked) methodExprNilChecked() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return methodExprNilCheckedErr()
}

func callAndNilCheck(f func(methodExprNilChecked) error) bool {
	return f(methodExprNilChecked{}) != nil
}

func MethodExprNilChecked() {
	if callAndNilCheck(methodExprNilChecked.methodExprNilChecked) {
		return
	}
}
