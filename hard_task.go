package main

import "strings"

func del_duplicates(a []int) []int {
	m := make(map[int]bool)
	var ans []int

	for i := 0; i < len(a); i++ {
		_, ok := m[a[i]]

		if !ok {
			ans = append(ans, a[i])
			m[a[i]] = true
		}
	}

	return ans
}

func bubble_sort(a []int) []int {
	stop := true

	for stop {
		for i := 0; i < len(a) - 1; i++ {
			if a[i] > a[i + 1] {
				a[i], a[i + 1] = a[i + 1], a[i]
				stop = false
			}
		}
		stop = !stop
	}

	return a
}

func fibbonaci(n int) []int {
	fib := make([]int, n)
	fib[0] = 1
	fib[1] = 1
	for i := 2; i < n; i++ {
		fib[i] = fib[i - 1] + fib[i - 2];
	}

	return fib
}

func count(a []int, x int) int {
	ans := 0
	for i := 0; i < len(a); i++ {
		if a[i] == x {
			ans++
		}
	}

	return ans
}

func intersection(a []int, b []int) []int {
	a = del_duplicates(a)
	b = del_duplicates(b)

	m := make(map[int]bool)
	for _, elem := range a {
		m[elem] = true
	}

	var ans[]int
	for _, elem := range b {
		_, ok := m[elem]
		if ok {
			ans = append(ans, elem)
		}
	}

	return ans
}

func anargamm(a, b string) bool {
	a = strings.ToLower(a)
	b = strings.ToLower(b)

	am := make(map[rune]int)
	bm := make(map[rune]int)

	for _, x := range a {
		am[x]++
	}
	for _, x := range b {
		bm[x]++
	}
	
	for _, x := range a {
		if am[x] != bm[x] {
			return false
		}
	}

	for _, x := range b {
		if am[x] != bm[x] {
			return false
		}
	}

	return true
}

func merge(a, b []int) []int {
	i := 0
	j := 0
	var ans []int

	for i < len(a) && j < len(b) {
		if a[i] <= b[j] {
			ans = append(ans, a[i])
			i++
		} else {
			ans = append(ans, b[j])
			j++
		}
	}

	if i < len(a) {
		for i < len(a) {
			ans = append(ans, a[i])
			i++
		}
	}
	if j < len(b) {
		for j < len(b) {
			ans = append(ans, b[j])
			j++
		}
	}

	return ans
}