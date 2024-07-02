package main

import (
	"fmt"
)

func hello_word() {
	fmt.Println("Добрый вечер")
}

func sum(x int, y int) int {
	return x + y
}

func is_even(x int) bool {
	if x % 2 == 0 {
		return true
	}
	return false
}

func find_max(a, b, c int) int {
	var ans int = a
	if ans < b {
		ans = b
	}
	if ans < c {
		ans = c
	}

	return ans
}

func factorial(n int) int {
	var ans int = 1
	for i := 2; i <= n; i++ {
		ans *= i
	}

	return ans
}

func is_vowel(a uint8) bool {
	switch a {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	default:
		return false
	}
}

func print_primes(n int) {
	p := make([]bool, n + 1)
	for i := 2; i <= n; i++ {
		p[i] = true
	}
	for i := 2; i <= n; i++ {
		if p[i] {
			fmt.Println(i)

			for j := i * i; j <= n; j += i {
				p[j] = false
			}
		}
	}
}

func reverse(s string) string {
	new_str := []rune(s)
	for i, j := 0, len(new_str) - 1; i < j; i, j = i + 1, j - 1 {
		new_str[i], new_str[j] = new_str[j], new_str[i]
	}

	return string(new_str)
}

func array_sum(arr []int) int {
	ans := 0
	for i := 0; i < len(arr); i++ {
		ans += arr[i]
	}

	return ans
}

type Rectangle struct {
	width int
	height int
}

func (rec Rectangle) area() int {
	return rec.width * rec.height
}


func main() {
	k := []int{4, 7, 10, 11}
	l := []int{1, 2, 8, 9}
	fmt.Println(merge(k, l))
}