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

// Errors passed as arguments to functions of the same package.

func logErr(err error) {
	log.Print(err)
}

func LoggedByLocal() {
	err := func() error { return errors.New("") }()
	logErr(err) // OK: the parameter is passed to an external function.
}

func LoggedByDeferredLocal() {
	err := func() error { return errors.New("") }()
	defer logErr(err) // OK: the parameter is passed to an external function.
}

func LoggedByGoLocal() {
	err := func() error { return errors.New("") }()
	go logErr(err) // OK: the parameter is passed to an external function.
}

func LoggedByDeferredExternal() {
	err := func() error { return errors.New("") }()
	defer log.Print(err) // OK: passed to external function.
}

func isErr(err error) bool {
	return err != nil
}

func NilCheckedByLocal() {
	err := func() error { return errors.New("") }() // want "error is only ever nil-checked; consider returning a bool instead"
	if isErr(err) {
		return
	}
}

func ignoreErr(error) error {
	return errors.New("")
}

func IgnoredByLocal() error {
	err := func() error { return errors.New("") }() // want "error is only ever nil-checked; consider returning a bool instead"
	// The parameter is never used, even though the result is returned.
	return ignoreErr(err)
}

type errHolder struct {
	err error
}

func (h *errHolder) set(err error) {
	h.err = err
}

func PassedToStoringLocal(h *errHolder) {
	err := func() error { return errors.New("") }()
	h.set(err) // OK: the parameter is stored in a struct field.
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
	err := func() error { return errors.New("") }()
	if nilCheckAll(err) { // OK: variadic arguments are stored in a slice.
		return
	}
}
