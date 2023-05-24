package random

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"minik8s/pkg/types"
)

const (
	letterBytes = "bcdfghjklmnpqrstvwxyz"
	digitBytes  = "2456789"
)

func GenerateRandomString(length int) string {
	charset := letterBytes + digitBytes
	chars := make([]byte, length)

	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		chars[i] = charset[randomIndex.Int64()]
	}

	return string(chars)
}

func GenerateUUID() types.UID {
	uuid := make([]byte, 16)

	_, err := rand.Read(uuid)
	if err != nil {
		return ""
	}

	// 设置 UUID 版本号为 4
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	// 设置 UUID 变体为 RFC 4122
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return types.UID(fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]))
}
