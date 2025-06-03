package util

import (
	"fmt"
	"math/rand"
	"strings"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(letters)

	for i := 0; i < n; i++ {
		c := letters[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() string {
	m := RandomInt(0, 10000)
	c := RandomInt(0, 99)
	return fmt.Sprintf("%d.%02d", m, c)
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "BRL"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
