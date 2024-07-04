package random

import (
	"math/rand"
)

func NewRandomString(length int) string {
	chars := []rune("abcdefghigkmnopqrstvuwxyz")
	alias := make([]rune, length)
	for i := 0; i < length; i++ {
		alias[i] = chars[rand.Intn(len(chars))]
	}
	return string(alias)
}
