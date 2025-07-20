package linq

func Pipe1[T any, R any](
	src Stream[T],
	sink func(Stream[T]) (R, error),
) (R, error) {
	return sink(src)
}

func Pipe2[T any, U any, R any](
	src Stream[T],
	f1 func(Stream[T]) Stream[U],
	sink func(Stream[U]) (R, error),
) (R, error) {
	return sink(f1(src))
}

func Pipe3[T any, U any, V any, R any](
	src Stream[T],
	f1 func(Stream[T]) Stream[U],
	f2 func(Stream[U]) Stream[V],
	sink func(Stream[V]) (R, error),
) (R, error) {
	return sink(f2(f1(src)))
}

func Pipe4[T any, U any, V any, W any, R any](
	src Stream[T],
	f1 func(Stream[T]) Stream[U],
	f2 func(Stream[U]) Stream[V],
	f3 func(Stream[V]) Stream[W],
	sink func(Stream[W]) (R, error),
) (R, error) {
	return sink(f3(f2(f1(src))))
}

func Pipe5[T any, U any, V any, W any, X any, R any](
	src Stream[T],
	f1 func(Stream[T]) Stream[U],
	f2 func(Stream[U]) Stream[V],
	f3 func(Stream[V]) Stream[W],
	f4 func(Stream[W]) Stream[X],
	sink func(Stream[X]) (R, error),
) (R, error) {
	return sink(f4(f3(f2(f1(src)))))
}

func Pipe6[T any, U any, V any, W any, X any, Y any, R any](
	src Stream[T],
	f1 func(Stream[T]) Stream[U],
	f2 func(Stream[U]) Stream[V],
	f3 func(Stream[V]) Stream[W],
	f4 func(Stream[W]) Stream[X],
	f5 func(Stream[X]) Stream[Y],
	sink func(Stream[Y]) (R, error),
) (R, error) {
	return sink(f5(f4(f3(f2(f1(src))))))
}

func Pipe7[T any, U any, V any, W any, X any, Y any, Z any, R any](
	src Stream[T],
	f1 func(Stream[T]) Stream[U],
	f2 func(Stream[U]) Stream[V],
	f3 func(Stream[V]) Stream[W],
	f4 func(Stream[W]) Stream[X],
	f5 func(Stream[X]) Stream[Y],
	f6 func(Stream[Y]) Stream[Z],
	sink func(Stream[Z]) (R, error),
) (R, error) {
	return sink(f6(f5(f4(f3(f2(f1(src)))))))
}
