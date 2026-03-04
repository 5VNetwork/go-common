package random

import "math/rand"

// RandomString generates a random word of specified length
func RandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		// Generate random index
		n := rand.Intn(len(letters))
		result[i] = letters[n]
	}

	return string(result)
}
