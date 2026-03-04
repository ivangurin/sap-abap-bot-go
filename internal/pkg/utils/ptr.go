package utils

//go:fix inline
func Ptr[T any](val T) *T {
	return new(val)
}
