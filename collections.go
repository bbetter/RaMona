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
