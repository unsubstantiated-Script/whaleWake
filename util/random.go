package util

import (
	"encoding/hex"
	"github.com/google/uuid"
	"math/rand"
	"strconv"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const hexAlphabet = "0123456789abcdef"

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

// RandomHexString returns a random hexadecimal string of length n
func RandomHexString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = hexAlphabet[rand.Intn(len(hexAlphabet))]
	}
	return string(b)
}

// RandomSymmetricKey returns a random 32-byte symmetric key as a hex string
func RandomSymmetricKey() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(key)
}

// RandomUserName  returns a random owner name
func RandomUserName() string {
	return RandomString(6)
}

// RandomUUID returns a random UUID string
func RandomUUID() uuid.UUID {
	return uuid.New()
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
