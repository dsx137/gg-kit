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
