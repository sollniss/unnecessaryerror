package a

import (
	"errors"
	"log"

	"a/external"
)

func passedToExternal() error {
	return errors.New("")
}

func PassedToExternal() {
	err := passedToExternal()
	log.Print(err) // OK: passed to external function.
}

func passedToBuiltinFunc() error {
	return errors.New("")
}

func PassedToBuiltinFunc() {
	err := passedToBuiltinFunc()
	println(err) // OK: passed to builtin function.
}

func funcPassedToExternal_1() error {
	return errors.New("")
}

func funcPassedToExternal_2() {
	err := external.FuncPassedToExternal(funcPassedToExternal_1) // OK: passed to external function.
	log.Print(err)
}

func funcPassedToExternalMethod_1() error {
	return errors.New("")
}

func funcPassedToExternalMethod_2() {
	err := (&external.Exported{}).FuncPassedToExternalMethod(funcPassedToExternalMethod_1) // OK: passed to external function.
	log.Print(err)
}

func closurePassedToExternal() error {
	return errors.New("")
}

func ClosurePassedToExternal() {
	err := external.ClosurePassedToExternal(func() error { // OK: passed to external function.
		return closurePassedToExternal()
	})
	log.Print(err)
}

func assignedInClosureToExternal() error {
	return errors.New("")
}

func AssignedInClosureToExternal() {
	var err1 error
	err2 := external.AssignedInClosure(func() error {
		err1 = assignedInClosureToExternal()
		return errors.New("")
	})
	log.Print(err1, err2) // OK: passed to external function.
}

func assignedToExternalGlobalVar() error {
	return errors.New("")
}

func AssignedToExternalGlobalVar() {
	external.Var = assignedToExternalGlobalVar() // OK: assigned to global variable.
}

func assignedToExternalStructField1() error {
	return errors.New("")
}

func AssignedToExternalStructField1() {
	s := external.AssignedToExternalStructField1{
		Err: assignedToExternalStructField1(), // OK: Assigned to struct field.
	}
	_ = s
}

func assignedToExternalStructField2() error {
	return errors.New("")
}

func assignedToExternalStructField2_1() error {
	err := assignedToExternalStructField2()
	if err != nil {
		return err
	}
	return nil
}

func AssignedToExternalStructField2() {
	s := external.AssignedToExternalStructField2{
		Err: assignedToExternalStructField2_1(), // OK: Assigned to struct field.
	}
	_ = s
}

func funcAssignedToExternalStructField1() error {
	return errors.New("")
}

func FuncAssignedToExternalStructField1() {
	s := external.FuncAssignedToExternalStructField1{
		Func: funcAssignedToExternalStructField1, // OK: Func assigned to struct field.
	}
	_ = s
}

func funcAssignedToExternalStructField2() error {
	return errors.New("")
}

func funcAssignedToExternalStructField2_1() error {
	err := funcAssignedToExternalStructField2()
	if err != nil {
		return err
	}
	return nil
}

func FuncAssignedToExternalStructField2() {
	s := external.FuncAssignedToExternalStructField2{
		Func: funcAssignedToExternalStructField2_1, // OK: Caller assigned to struct field.
	}
	_ = s
}

func funcAssignedToGlobalStructField() error {
	return errors.New("")
}

func funcAssignedToGlobalStructField_1() error {
	err := funcAssignedToGlobalStructField()
	if err != nil {
		return err
	}
	return nil
}

// Package-level initializers run in the package initializer, not in any source function.
var _funcAssignedToGlobalStructField = external.FuncAssignedToExternalStructField1{
	Func: funcAssignedToGlobalStructField_1, // OK: caller assigned to struct field.
}
