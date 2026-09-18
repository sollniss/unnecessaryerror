package a

import (
	"errors"
	"log"
)

func Closure() {
	err := func() error { return errors.New("") }() // want "error is only ever nil-checked; consider returning a bool instead"
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

func capturedByDefer() error {
	return errors.New("")
}

func CapturedByDefer() {
	err := capturedByDefer()
	defer func() {
		log.Print(err) // OK: passed to external function.
	}()
}

func capturedByGo() error {
	return errors.New("")
}

func CapturedByGo() {
	err := capturedByGo()
	go func() {
		log.Print(err) // OK: passed to external function.
	}()
}

func capturedByCalledClosure() error {
	return errors.New("")
}

func CapturedByCalledClosure() {
	err := capturedByCalledClosure()
	f := func() {
		log.Print(err) // OK: passed to external function.
	}
	f()
}

func capturedByUncalledClosure() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func CapturedByUncalledClosure() {
	err := capturedByUncalledClosure()
	f := func() {
		if err != nil {
			return
		}
	}
	f()
}

func assignedInClosure() error {
	return errors.New("")
}

func AssignedInClosure() {
	var err error
	func() {
		err = assignedInClosure()
	}()
	log.Print(err) // OK: passed to external function.
}

func assignedInNestedClosure() error {
	return errors.New("")
}

func AssignedInNestedClosure() {
	var err error
	func() {
		func() {
			err = assignedInNestedClosure()
		}()
	}()
	log.Print(err) // OK: passed to external function.
}

func nilCheckedAfterClosureAssign() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func NilCheckedAfterClosureAssign() {
	var err error
	func() {
		err = nilCheckedAfterClosureAssign()
	}()
	if err != nil {
		return
	}
}

func closurePassedToLocal() error {
	return errors.New("")
}

func runAndLog(fn func() error) {
	log.Print(fn())
}

func ClosurePassedToLocal() {
	runAndLog(func() error { // OK: the closure's result is passed to an external function.
		return closurePassedToLocal()
	})
}

func closurePassedToNilCheckingLocal() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func runAndNilCheck(fn func() error) bool {
	return fn() != nil
}

func ClosurePassedToNilCheckingLocal() {
	if runAndNilCheck(func() error {
		return closurePassedToNilCheckingLocal()
	}) {
		return
	}
}

func funcPassedToLocal() error {
	return errors.New("")
}

func FuncPassedToLocal() {
	runAndLog(funcPassedToLocal) // OK: the function's result is passed to an external function.
	if funcPassedToLocal() != nil {
		return
	}
}

func funcPassedToNilCheckingLocal() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func FuncPassedToNilCheckingLocal() {
	if runAndNilCheck(funcPassedToNilCheckingLocal) {
		return
	}
	if funcPassedToNilCheckingLocal() != nil {
		return
	}
}
