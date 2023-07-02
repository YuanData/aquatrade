package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "thequickbrownfoxjumpsoverthelazydog"

func init() {
	rand.NewSource(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomAccount() string {
	return RandomString(5)
}

func RandomMoney() int64 {
	return RandomInt(0, 15000)
}

func RandomCurrency() string {
	currencies := []string{"AUD", "JPY", "RMB", "CHF"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@mail.test", RandomString(5))
}
