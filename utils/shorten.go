package utils

import (
	"fmt"
)

var base62 []rune = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

// Base10 to Base64 key generator
func GetShortCode(n int64) string {
	fmt.Println("Shortening URL...")

	var r []rune = []rune{}
	for n != 0 {
		r = append(r, base62[n%62])
		n /= 62
	}

	// Fill out to satisfy 7 length
	for len(r) < 7 {
		r = append(r, '0')
	}

	// Reverse the rune slice
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	fmt.Println("Key Generated: ", r)

	return string(r)
}
