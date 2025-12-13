package lang

func Bind[T any](f func(x *T)) *T {
	x := new(T)
	f(x)
	return x
}

func BindR[T any, R any](f func(x *T) R) (*T, R) {
	x := new(T)
	return x, f(x)
}

func BindPtr[T any](f func(x **T)) *T {
	var x *T
	f(&x)
	return x
}

func BindPtrR[T any, R any](f func(x **T) R) (*T, R) {
	var x *T
	return x, f(&x)
}

// --------------------- EXPAND ---------------------

func UnmarshalTo[T any, D any](f func(D, any) error, data D) (*T, error) {
	return BindR(func(v *T) error { return f(data, v) })
}

func ShouldBindTo[T any](f func(x *T) error) (*T, error) {
	return BindR(f)
}
