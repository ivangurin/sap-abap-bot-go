package utils

type Map[TKey comparable, TArg any] map[TKey]TArg

func (m Map[TKey, TArg]) Add(key TKey, value TArg) {
	m[key] = value
}

func (m Map[TKey, TArg]) Delete(key TKey) {
	delete(m, key)
}

func (m Map[TKey, TArg]) Keys() []TKey {
	res := make([]TKey, 0, len(m))
	for key := range m {
		res = append(res, key)
	}
	return res
}

func (m Map[TKey, TArg]) Has(key TKey) bool {
	_, exists := m[key]
	return exists
}

func ToMap[S ~[]TArg, TKey comparable, TValue, TArg any](collection S, selector func(arg TArg) (TKey, TValue)) Map[TKey, TValue] {
	result := make(Map[TKey, TValue], len(collection))
	for _, element := range collection {
		k, v := selector(element)
		result[k] = v
	}
	return result
}

func ToMapByField[S ~[]TArg, TKey comparable, TArg any](collection S, selector func(arg TArg) TKey) Map[TKey, TArg] {
	result := make(Map[TKey, TArg], len(collection))
	for _, item := range collection {
		result[selector(item)] = item
	}
	return result
}
