package random

import (
	crand "crypto/rand"
	"encoding/hex"
	"math/rand"
	"time"
)

//根据最大值生成随机整数
func RandomN(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

//生成随机bytes
func GetRandomBytes(len int) []byte {
	s := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]byte, len)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < len; i++ {
		result[i] = s[r.Intn(62)]
	}

	return result
}

func GetRandomString(l int) string {
	return string(GetRandomBytes(l))
}

//l should be even
func GetRandomString2(l int) string {
	if l < 2 {
		l = 2
	}
	if l%2 != 0 {
		l = l + 1
	}

	b := make([]byte, l/2)
	_, err := crand.Read(b)
	if err != nil {
		return GetRandomString(l)
	}

	return hex.EncodeToString(b)
}
