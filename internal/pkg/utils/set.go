package utils

type Set[T comparable] map[T]struct{}

func NewSet[T comparable]() Set[T] {
	return make(Set[T])
}

func ToSet[S ~[]TArg, TKey comparable, TArg any](collection S, selector func(arg TArg) TKey) Set[TKey] {
	result := make(Set[TKey], len(collection))
	for _, item := range collection {
		result[selector(item)] = struct{}{}
	}
	return result
}

func (s Set[T]) Add(elem T) {
	s[elem] = struct{}{}
}

func (s Set[T]) Has(elem T) bool {
	_, ok := s[elem]

	return ok
}

func (s Set[T]) Delete(elem T) {
	delete(s, elem)
}

func (s Set[T]) ToSlice() []T {
	if s == nil {
		return nil
	}

	ret := make([]T, 0, len(s))
	for val := range s {
		ret = append(ret, val)
	}

	return ret
}
