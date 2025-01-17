package util

import (
	"math/rand"
	"strconv"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// RandomInt returns a random integer between min and max
func RandomInt(min, max int64) int64 {
	return rand.Int63n(max-min+1) + min
}

// RandomString returns a random string of length n
func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}

// RandomUserName  returns a random owner name
func RandomUserName() string {
	return RandomString(6)
}

// RandomBusinessName returns a random business name
func RandomBusinessName() string {
	return RandomString(6) + " " + RandomString(6)
}

// RandomStreetAddress returns a random street address
func RandomStreetAddress() string {
	streetNumber := strconv.FormatInt(RandomInt(1, 999), 10)
	return streetNumber + " " + RandomString(6) + " " + RandomString(6)
}

// RandomCountryCodeOrState returns a random country code or state
func RandomCountryCodeOrState() string {
	return strings.ToUpper(RandomString(2))
}

// RandomPassword returns a random password
func RandomPassword() string {
	return RandomString(16)
}

// RandomEmail returns a random email
func RandomEmail() string {
	return RandomString(6) + "@" + RandomString(4) + ".com"
}
