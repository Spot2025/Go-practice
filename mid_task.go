package main

import "fmt"

func convertion(temp float64) float64 {
	f := (temp * 9 / 5) + 32
	return f
}

func countdown(n int) {
	for i := n; i > 0; i-- {
		fmt.Println(i)
	}
}

func str_len(s string) int {
	ans := 0
	for range s {
		ans++
	}
	return ans
}

func is_here() {
	var n int
	fmt.Println("Введите размер массива")
	fmt.Scanln(&n)

	arr := make([]int, n)
	fmt.Println("Введите массив")
	for i := 0; i < n; i++ {
		fmt.Scan(&arr[i])
	}

	var x int
	fmt.Println("Что нужно найти")
	fmt.Scanln(&x)

	for i := 0; i < n; i++ {
		if x == arr[i] {
			fmt.Println("Найдено")
			return
		}
	}

	fmt.Println("Не найдено")
}

func mult_table(n int) {
	fmt.Print(" ")
	for i := 1; i <= n; i++ {
		fmt.Printf("\t%v", i)
	}
	fmt.Print("\n")
	for i := 1; i <= n; i++ {
		fmt.Print(i)
		for j := 1; j <= n; j++ {
			fmt.Printf("\t%v", i * j)
		}

		print("\n")
	}
}

func is_palindrome(s string) bool {
	str := []rune(s)
	for i, j := 0, len(str) - 1; i < j; i, j = i + 1, j - 1 {
		if str[i] != str[j] {
			return false
		}
	}
	return true
}

func min_max(arr []int) (int, int) {
	min, max := arr[0], arr[0]

	for i := 1; i < len(arr); i++ {
		if arr[i] < min {
			min = arr[i]
		}
		if arr[i] > max {
			max = arr[i]
		}
	}

	return min, max
}

func remove(a []int, i int) []int {
	if i > len(a) || i < 0 {
		return a
	}

	return append(a[:i], a[i + 1:]...)
}

func find(a []int, x int) int {
	for i := 0; i < len(a); i++ {
		if a[i] == x {
			return i
		}
	}

	return -1
}