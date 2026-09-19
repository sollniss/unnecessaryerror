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

func closureInGeneric[T any]() {
	err := func() error { return errors.New("") }() // want "error is only ever nil-checked; consider returning a bool instead"
	if err != nil {
		return
	}
}

func ClosureInGeneric() {
	closureInGeneric[int]()
}

func closureInGenericLogged[T any]() {
	err := func() error { return errors.New("") }()
	log.Print(err) // OK: passed to external function.
}

func ClosureInGenericLogged() {
	closureInGenericLogged[int]()
}

// Interface calls dispatch to methods of generic types through the wrappers
// of their instantiations, whose concrete signatures decide whether they
// implement the interface.

type genericParam[T any] struct{}

func (*genericParam[T]) genericParam(v T) error { // OK: passed to external function via interface.
	return errors.New("")
}

type genericParamIface interface{ genericParam(int) error }

func GenericParam(g genericParamIface) {
	log.Print(g.genericParam(1))
	if (&genericParam[int]{}).genericParam(1) != nil {
		return
	}
}

type genericParamNilChecked[T any] struct{}

func (*genericParamNilChecked[T]) genericParamNilChecked(v T) error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

type genericParamNilCheckedIface interface{ genericParamNilChecked(int) error }

func GenericParamNilChecked() {
	var g genericParamNilCheckedIface = &genericParamNilChecked[int]{}
	if g.genericParamNilChecked(1) != nil {
		return
	}
}

type genericValueRecv[T any] struct{}

func (genericValueRecv[T]) genericValueRecv(v T) error { // OK: passed to external function via interface.
	return errors.New("")
}

type genericValueRecvIface interface{ genericValueRecv(int) error }

func GenericValueRecv() {
	var g genericValueRecvIface = genericValueRecv[int]{}
	log.Print(g.genericValueRecv(1))
	if (genericValueRecv[int]{}).genericValueRecv(1) != nil {
		return
	}
}

type genericTypeParamCall[T any] struct{}

func (*genericTypeParamCall[T]) genericTypeParamCall(v T) error { // OK: passed to external function via type parameter.
	return errors.New("")
}

type genericTypeParamCallIface interface{ genericTypeParamCall(int) error }

func callGenericTypeParam[T genericTypeParamCallIface](g T) {
	log.Print(g.genericTypeParamCall(1))
}

func GenericTypeParamCall() {
	// The instantiation is never converted to an interface.
	callGenericTypeParam(&genericTypeParamCall[int]{})
	if (&genericTypeParamCall[int]{}).genericTypeParamCall(1) != nil {
		return
	}
}

type genericTypeParamCallNilChecked[T any] struct{}

func (*genericTypeParamCallNilChecked[T]) genericTypeParamCallNilChecked(v T) error { // want "error is only ever nil-checked; consider returning a bool instead"
	return errors.New("")
}

type genericTypeParamCallNilCheckedIface interface {
	genericTypeParamCallNilChecked(int) error
}

func callGenericTypeParamNilChecked[T genericTypeParamCallNilCheckedIface](g T) bool {
	return g.genericTypeParamCallNilChecked(1) != nil
}

func GenericTypeParamCallNilChecked() {
	if callGenericTypeParamNilChecked(&genericTypeParamCallNilChecked[int]{}) {
		return
	}
}

type genericCtor[T any] struct{}

func (*genericCtor[T]) genericCtor(v T) error { // OK: passed to external function via interface.
	return errors.New("")
}

func newGenericCtor[T any]() *genericCtor[T] {
	return &genericCtor[T]{}
}

type genericCtorIface interface{ genericCtor(int) error }

func GenericCtor() {
	// The instantiation is only spelled out inside the constructor.
	var g genericCtorIface = newGenericCtor[int]()
	log.Print(g.genericCtor(1))
	if newGenericCtor[string]().genericCtor("") != nil { // A different instantiation, which does not implement the interface.
		return
	}
}

// Generic results are ignored, but error results next to them are tracked.

func genericTupleNilChecked[T any]() (T, error) { // want "error is only ever nil-checked; consider returning a bool instead"
	var zero T
	return zero, errors.New("")
}

func GenericTupleNilChecked() {
	v, err := genericTupleNilChecked[int]()
	if err != nil {
		return
	}
	_ = v
}

func genericTupleLogged[T any]() (T, error) {
	var zero T
	return zero, errors.New("")
}

func GenericTupleLogged() {
	_, err := genericTupleLogged[int]()
	log.Print(err) // OK: passed to external function.
}

func genericTupleOfErrors[T any]() (T, error) { // want "error is only ever nil-checked; consider returning a bool instead"
	var zero T
	return zero, errors.New("")
}

func GenericTupleOfErrors() {
	v, err := genericTupleOfErrors[error]()
	log.Print(v) // The generic result is not tracked, even when instantiated with error.
	if err != nil {
		return
	}
}

type genericTupleMethod[T any] struct{}

func (*genericTupleMethod[T]) genericTupleMethod(v T) (T, error) { // OK: passed to external function via interface.
	return v, errors.New("")
}

type genericTupleMethodIface interface {
	genericTupleMethod(string) (string, error)
}

func GenericTupleMethod() {
	var g genericTupleMethodIface = &genericTupleMethod[string]{}
	log.Print(g.genericTupleMethod(""))
}

func onlyDeferredGeneric[T any]() error { // Not a candidate: never called for its result.
	return errors.New("")
}

func OnlyDeferredGeneric() {
	defer onlyDeferredGeneric[int]()
	go onlyDeferredGeneric[int]()
}
