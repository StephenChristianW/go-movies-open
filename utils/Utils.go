package utils

import "math/rand/v2"

func ListToMap(list []string) map[string]struct{} {
	m := make(map[string]struct{}, len(list))
	for _, item := range list {
		m[item] = struct{}{}
	}
	return m
}

// SliceRandN 从切片中随机获取 n 个不重复元素（若 n 超出长度则返回全部）
func SliceRandN[T any](src []T, n int) []T {
	length := len(src)
	if n > length {
		n = length
	}
	tmp := make([]T, len(src))
	copy(tmp, src)
	Shuffle(tmp)
	return tmp[:n]
}

func Shuffle[T any](slice []T) {
	n := len(slice)
	for i := n - 1; i > 0; i-- {
		j := rand.IntN(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}
