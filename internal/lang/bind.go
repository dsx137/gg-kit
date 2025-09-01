package lang

func Bind[T any](f func(x any)) *T {
	x := new(T)
	f(x)
	return x
}

func BindR[T any, R any](f func(x any) R) (*T, R) {
	x := new(T)
	return x, f(x)
}

func BindPtr[T any](f func(x any)) *T {
	var x *T
	f(&x)
	return x
}

func BindPtrR[T any, R any](f func(x any) R) (*T, R) {
	var x *T
	return x, f(&x)
}
