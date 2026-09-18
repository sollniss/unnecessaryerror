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

func NilCheckInDefer() {
	err := func() error { return errors.New("") }() // want "error is only ever nil-checked; consider returning a bool instead"
	defer func() {
		if err != nil {
			return
		}
	}()
}

func CapturedByDefer() {
	err := func() error { return errors.New("") }()
	defer func() {
		log.Print(err) // OK: passed to external function.
	}()
}

func CapturedByGo() {
	err := func() error { return errors.New("") }()
	go func() {
		log.Print(err) // OK: passed to external function.
	}()
}

func CapturedByCalledClosure() {
	err := func() error { return errors.New("") }()
	f := func() {
		log.Print(err) // OK: passed to external function.
	}
	f()
}

func CapturedByUncalledClosure() {
	err := func() error { return errors.New("") }() // want "error is only ever nil-checked; consider returning a bool instead"
	f := func() {
		if err != nil {
			return
		}
	}
	f()
}

func AssignedInClosure() {
	var err error
	func() {
		err = func() error { return errors.New("") }()
	}()
	log.Print(err) // OK: passed to external function.
}

func AssignedInNestedClosure() {
	var err error
	func() {
		func() {
			err = func() error { return errors.New("") }()
		}()
	}()
	log.Print(err) // OK: passed to external function.
}

func NilCheckedAfterClosureAssign() {
	var err error
	func() {
		err = func() error { return errors.New("") }() // want "error is only ever nil-checked; consider returning a bool instead"
	}()
	if err != nil {
		return
	}
}

// The cases below are about functions used as values, so they need named functions:
// a closure that is passed around instead of called is not a candidate.

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
