package a

import (
	"errors"
	"log"
)

// Functions and closures used as values rather than called directly.

func calledThroughLocalVar() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func CalledThroughLocalVar() {
	f := calledThroughLocalVar
	if f() != nil {
		return
	}
}

func calledThroughPhi1() error {
	return errors.New("")
}

func calledThroughPhi2() error {
	return errors.New("")
}

func CalledThroughPhi(cond bool) {
	f := calledThroughPhi1
	if cond {
		f = calledThroughPhi2
	}
	log.Print(f()) // OK: result passed to external function.
}

func calledThroughPhiNilChecked1() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func calledThroughPhiNilChecked2() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func CalledThroughPhiNilChecked(cond bool) {
	f := calledThroughPhiNilChecked1
	if cond {
		f = calledThroughPhiNilChecked2
	}
	if f() != nil {
		return
	}
}

type errFunc func() error

func (f errFunc) run() error {
	return f()
}

func convertedToNamedFuncType() error {
	return errors.New("")
}

func ConvertedToNamedFuncType() {
	log.Print(errFunc(convertedToNamedFuncType).run()) // OK: result passed to external function.
}

func convertedToNamedFuncTypeNilChecked() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

type errFuncNilChecked func() error

func (f errFuncNilChecked) run() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return f()
}

func ConvertedToNamedFuncTypeNilChecked() {
	if errFuncNilChecked(convertedToNamedFuncTypeNilChecked).run() != nil {
		return
	}
}

type methodValueNilChecked struct{}

func (methodValueNilChecked) method() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func MethodValueNilChecked() {
	f := methodValueNilChecked{}.method
	if f() != nil {
		return
	}
}

type ptrMethodExpr struct{}

func (*ptrMethodExpr) method() error {
	return errors.New("")
}

var _ptrMethodExpr = map[string]func(*ptrMethodExpr) error{
	"a": (*ptrMethodExpr).method, // OK: method expression inserted into map.
}
