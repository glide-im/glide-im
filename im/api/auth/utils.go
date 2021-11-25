package auth

import (
	"math/rand"
	"time"
)

var (
	table = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func genToken(length int) string {
	res := ""
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		idx := rnd.Intn(62)
		res = res + table[idx:idx+1]
	}
	return res
}
