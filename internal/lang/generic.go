package lang

func BindAlloc[T any](f func(x *T)) *T {
	x := new(T)
	f(x)
	return x
}

func Bind[T any](f func(x **T)) **T {
	var x *T
	f(&x)
	return &x
}

func BindAllocR[T any, R any](f func(x *T) R) (*T, R) {
	x := new(T)
	return x, f(x)
}

func BindR[T any, R any](f func(x **T) R) (*T, R) {
	var x *T
	return x, f(&x)
}
