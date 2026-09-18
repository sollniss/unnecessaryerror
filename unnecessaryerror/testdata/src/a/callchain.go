package a

import (
	"errors"
	"log"
)

func callchain1() error {
	return errors.New("")
}

func callchain1_1(err error) error {
	return err
}

func callchain1_2(err error) error {
	return err
}

func Callchain1() error {
	err := callchain1_2(callchain1_1(callchain1()))
	if err != nil {
		return err // OK: returned from exported function.
	}
	return nil
}

func callchain2() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func callchain2_1(err error) error { // want "error is only ever nil-checked; consider returning a bool instead"
	return err
}

func callchain2_2(err error) error { // want "error is only ever nil-checked; consider returning a bool instead"
	return err
}

func Callchain2() error {
	err := callchain2_2(callchain2_1(callchain2()))
	if err != nil {
		return nil
	}
	return nil
}

func callchain3() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func callchain3_1() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return callchain3()
}

func callchain3_2() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return callchain3_1()
}

func Callchain3() {
	err := callchain3_2()
	if err != nil {
		return
	}
}

func callchain4() error {
	return errors.New("")
}

func callchain4_1() error {
	err := callchain4()
	if err != nil {
		return err
	}
	return nil
}

type callchain4Rec struct{}

func (r *callchain4Rec) Callchain4() error {
	err := callchain4_1()
	if err != nil {
		return err // OK: returned from exported method.
	}
	return nil
}

func loggedByLocal() error {
	return errors.New("")
}

func logErr(err error) {
	log.Print(err)
}

func LoggedByLocal() {
	logErr(loggedByLocal()) // OK: the parameter is passed to an external function.
}

func loggedByDeferredLocal() error {
	return errors.New("")
}

func LoggedByDeferredLocal() {
	defer logErr(loggedByDeferredLocal()) // OK: the parameter is passed to an external function.
}

func loggedByGoLocal() error {
	return errors.New("")
}

func LoggedByGoLocal() {
	go logErr(loggedByGoLocal()) // OK: the parameter is passed to an external function.
}

func loggedByDeferredExternal() error {
	return errors.New("")
}

func LoggedByDeferredExternal() {
	defer log.Print(loggedByDeferredExternal()) // OK: passed to external function.
}

func nilCheckedByLocal() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func isErr(err error) bool {
	return err != nil
}

func NilCheckedByLocal() {
	if isErr(nilCheckedByLocal()) {
		return
	}
}

func ignoredByLocal() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func ignoreErr(error) error {
	return errors.New("")
}

func IgnoredByLocal() error {
	return ignoredByLocal_1()
}

func ignoredByLocal_1() error {
	return ignoreErr(ignoredByLocal()) // The parameter is never used, even though the result is returned.
}

func passedToStoringLocal() error {
	return errors.New("")
}

type errHolder struct {
	err error
}

func (h *errHolder) set(err error) {
	h.err = err
}

func PassedToStoringLocal(h *errHolder) {
	h.set(passedToStoringLocal()) // OK: the parameter is stored in a struct field.
}

func passedToVariadicLocal() error {
	return errors.New("")
}

func nilCheckAll(errs ...error) bool {
	for _, err := range errs {
		if err != nil {
			return true
		}
	}
	return false
}

func PassedToVariadicLocal() {
	if nilCheckAll(passedToVariadicLocal()) { // OK: variadic arguments are stored in a slice.
		return
	}
}

func recursive() error { // want "error is only ever nil-checked; consider returning a bool instead"
	if recursive() != nil {
		return errors.New("")
	}
	return nil
}

func Recursive() {
	if recursive() != nil {
		return
	}
}
