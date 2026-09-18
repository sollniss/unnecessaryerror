package a

import (
	"errors"
	"log"
)

func genericFunc1[T any]() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func GenericFunc1_1() {
	err := genericFunc1[int]()
	if err != nil {
		return
	}
}

func GenericFunc1_2() {
	err := genericFunc1[string]()
	if err != nil {
		return
	}
}

func genericFunc2[T any]() error {
	return errors.New("")
}

func GenericFunc2_1() {
	err := genericFunc2[int]()
	if err != nil {
		return
	}
}

func GenericFunc2_2() error {
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
}

func genericError2[T error]() T {
	var zero T
	return zero
}

func GenericError2() {
	if err := genericError2[error](); err != nil { // OK: ignore generic return types.
		return
	}
}

func calledFromGeneric() error {
	return errors.New("")
}

func genericCaller[T any]() error {
	return calledFromGeneric()
}

func GenericCaller() error {
	return genericCaller[int]() // OK: returned from exported function.
}

func nilCheckedFromGeneric() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func genericNilChecker[T any]() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return nilCheckedFromGeneric()
}

func GenericNilChecker() {
	if genericNilChecker[int]() != nil {
		return
	}
}

type genericType[T any] struct{}

func (g *genericType[T]) method() error {
	return errors.New("")
}

func GenericType() {
	log.Print((&genericType[int]{}).method()) // OK: passed to external function.
}

func (g *genericType[T]) nilCheckedMethod() error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

func GenericTypeNilChecked() {
	if (&genericType[int]{}).nilCheckedMethod() != nil {
		return
	}
	if (&genericType[string]{}).nilCheckedMethod() != nil {
		return
	}
}

// A method table built from generic adapters in a package-level initializer.

type handlerServer struct{}

type handlerSession struct{ server *handlerServer }

type handlerFunc[P any] func(ss *handlerSession, p P) (any, error)

func serverMethod[P any](f func(*handlerServer, *handlerSession, P) (any, error)) handlerFunc[P] {
	return func(ss *handlerSession, p P) (any, error) {
		return f(ss.server, ss, p)
	}
}

type handlerInfo struct {
	handle any
}

func newHandlerInfo[P any](h handlerFunc[P]) handlerInfo {
	return handlerInfo{handle: h}
}

func handledErr() error {
	return errors.New("")
}

func (s *handlerServer) handle(ss *handlerSession, params *int) (any, error) {
	return nil, handledErr()
}

var _handlerInfos = map[string]handlerInfo{
	"handle": newHandlerInfo(serverMethod((*handlerServer).handle)), // OK: stored in a struct field via generic adapters.
}
