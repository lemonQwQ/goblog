package util

import (
	"crypto/rand"
	"math/big"
)

// GetRandom 获取随机数
func GetRandom(k int) string {
	str := ""
	for i := 0; i < k; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(10))
		str += num.String()
	}
	return str
}
