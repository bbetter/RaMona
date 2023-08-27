package utils

func Filter[Item any](vs []Item, f func(Item) bool) []Item {
	vsf := make([]Item, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func Any[Item any](vs []Item, f func(Item) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

func Map[Item1 any, Item2 any](vs []Item1, f func(Item1) Item2) []Item2 {
	vsf := make([]Item2, 0)
	for _, v := range vs {
		vsf = append(vsf, f(v))
	}
	return vsf
}

func Distinct[Item comparable](vs []Item) []Item {
	unique := make(map[Item]bool)
	result := make([]Item, 0, len(vs))
	for _, val := range vs {
		if !unique[val] {
			unique[val] = true
			result = append(result, val)
		}
	}
	return result
}
