package features

// luhn algorithm for creaditcard num checking
func luhnAlgorithm(number string) bool {
	sum := 0
	isSecond := false

	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')
		if isSecond {
			digit *= 2
		}
		sum += digit / 10
		sum += digit % 10

		isSecond = !isSecond
	}
	if sum%10 == 0 {
		return true
	}
	return false
}