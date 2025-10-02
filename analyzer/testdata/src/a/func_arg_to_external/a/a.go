// Helper package.
package a

func CallFunc(fn func() error) error {
	return fn()
}

type Exported struct{}

func (e *Exported) CallFuncInMethod(fn func() error) error {
	return fn()
}
