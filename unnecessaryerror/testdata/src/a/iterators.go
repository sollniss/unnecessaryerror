package a

import (
	"errors"
	"iter"
	"log"
)

// Range-over-func loops compile their bodies into synthetic closures.

func seq(yield func(int) bool) {
	yield(0)
}

func LoggedInRangeFunc() {
	for range seq {
		err := func() error { return errors.New("") }()
		log.Print(err) // OK: passed to external function.
	}
}

func NilCheckedInRangeFunc() {
	for range seq {
		err := func() error { return errors.New("") }() // want "error is only ever nil-checked; consider returning a bool instead"
		if err != nil {
			return
		}
	}
}

func yieldedErr() error {
	return errors.New("")
}

func errSeq() iter.Seq2[int, error] {
	return func(yield func(int, error) bool) {
		yield(0, yieldedErr()) // OK: passed to a function value.
	}
}

func ErrSeq() {
	for _, err := range errSeq() {
		log.Print(err)
	}
}
