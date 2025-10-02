# unnecessaryerror
unnecessaryerror is a linter that checks whether an error could be replaced with a bool.

Example:

```go
func foo() error { // error is only ever nil-checked; consider returning a bool instead
	return errors.New("error")
}

func Bar() {
	err := foo()
	if err != nil {
		doA()
	}
	doB()
}
```

For more examples, refer to the [testdata](https://github.com/sollniss/unnecessaryerror/blob/main/analyzer/testdata/src/a/a.go) directory.

Run with
```sh
go run github.com/sollniss/unnecessaryerror/cmd/unnecessaryerror@latest ./...
```

