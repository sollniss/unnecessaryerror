package b

import (
	"a/func_arg_to_external/a"
	"errors"
	"log"
)

func causeErr1() error {
	return errors.New("")
}

func Exported1() {
	err := a.CallFunc(func() error { // OK: passed to external function.
		return causeErr1()
	})
	log.Print(err)
}

func causeErr2() error {
	return errors.New("")
}

func Exported2() {
	var err error
	err = a.CallFunc(func() error {
		err = causeErr2()
		return err
	})
	log.Print(err)
}
