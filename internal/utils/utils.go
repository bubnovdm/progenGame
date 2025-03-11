package utils

// min возвращает минимальное из двух целых чисел
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max возвращает максимальное из двух целых чисел
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// abs возвращает абсолютное значение целого числа
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
