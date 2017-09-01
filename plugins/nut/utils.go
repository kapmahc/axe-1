package nut

import (
	"math/rand"
	"time"
)

//Random randome string
func Random(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	buf := make([]rune, n)
	for i := range buf {
		buf[i] = letters[rd.Intn(len(letters))]
	}
	return string(buf)
}
