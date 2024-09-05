package typed

func MapFromList[ListT any, MapKeyT comparable](items []ListT, mapper func(el *ListT) MapKeyT) map[MapKeyT]ListT {
	result := make(map[MapKeyT]ListT)
	for _, item := range items {
		key := mapper(&item)
		result[key] = item
	}

	return result
}

func UniqueList[ListT any, MapKeyT comparable](items []ListT, mapper func(el *ListT) MapKeyT) []ListT {
	uniqueMap := MapFromList(items, mapper)

	newUniqueList := []ListT{}
	for _, it := range uniqueMap {
		newUniqueList = append(newUniqueList, it)
	}

	return newUniqueList
}

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		Isolate(func() {
			us[i] = f(ts[i])
		})
	}
	return us
}
