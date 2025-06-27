package validation

import "strconv"

// Luhn валидирует входящий аргумент согласно алгоритму Луна
func Luhn(num string) bool {
	number, err := strconv.Atoi(num)
	if err != nil {
		return false
	}

	checksum := func(number int) int {
		var luhn int

		for i := 0; number > 0; i++ {
			cur := number % 10

			if i%2 == 0 {
				cur = cur * 2
				if cur > 9 {
					cur = cur%10 + cur/10
				}
			}

			luhn += cur
			number = number / 10
		}
		return luhn % 10
	}

	return (number%10+checksum(number/10))%10 == 0
}
