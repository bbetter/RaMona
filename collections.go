package main

func Filter[K any](vs []K, f func(K) bool) []K {
    vsf := make([]K, 0)
    for _, v := range vs {
        if f(v) {
            vsf = append(vsf, v)
        }
    }
    return vsf
}

func Any[K any](vs []K, f func(K) bool) bool {
    for _, v := range vs {
        if f(v) {
            return true
        }
    }
    return false
}

func Map[K any, T any](vs []K, f func(K) T) []T {
    vsf := make([]T, 0)
    for _, v := range vs {
        vsf = append(vsf, f(v))
    }
    return vsf
}

func Distinct[K comparable](vs []K) []K{
    unique := make(map[K]bool)
    result := make([]K, 0, len(vs))
    for _, val := range vs {
        if !unique[val] {
            unique[val] = true
            result = append(result, val)
        }
    }
    return result
}
