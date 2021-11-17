package userdao

import (
	"math/rand"
)

var (
	table = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func genToken(length int) string {
	res := ""
	for i := 0; i < length; i++ {
		idx := rand.Intn(62)
		res = res + table[idx:idx+1]
	}

	return res
}
