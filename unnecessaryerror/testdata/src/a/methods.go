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

type embeddedIface interface{ embeddedIface() error }

type outerIface interface {
	embeddedIface
}

type _embeddedIface struct{}

func (_embeddedIface) embeddedIface() error {
	return errors.New("")
}

func EmbeddedIface(o outerIface) {
	log.Print(o.embeddedIface()) // OK: passed to external function via embedded interface method.
	if (_embeddedIface{}).embeddedIface() != nil {
		return
	}
}

type _anonIface struct{}

func (_anonIface) anonIface() error {
	return errors.New("")
}

func AnonIface(v any) {
	log.Print(v.(interface{ anonIface() error }).anonIface()) // OK: passed to external function via anonymous interface.
	if (_anonIface{}).anonIface() != nil {
		return
	}
}

type twoImpls interface{ twoImpls() error }

type _twoImpls1 struct{}

func (_twoImpls1) twoImpls() error { // Conservative: either implementation may be the dynamic target.
	return errors.New("")
}

type _twoImpls2 struct{}

func (_twoImpls2) twoImpls() error {
	return errors.New("")
}

func TwoImpls(t twoImpls) {
	log.Print(t.twoImpls())
	if (_twoImpls1{}).twoImpls() != nil {
		return
	}
	if (_twoImpls2{}).twoImpls() != nil {
		return
	}
}

type genericGetter interface{ get() error }

type genericImpl[T any] struct{}

func (*genericImpl[T]) get() error {
	return errors.New("")
}

func GenericImpl() {
	var g genericGetter = &genericImpl[int]{}
	log.Print(g.get()) // OK: passed to external function via interface.
	if (&genericImpl[string]{}).get() != nil {
		return
	}
}

type unexportedType struct{}

func (unexportedType) Exported() error { // Not a candidate: exported method name.
	return errors.New("")
}

func UnexportedType() {
	if (unexportedType{}).Exported() != nil {
		return
	}
}
