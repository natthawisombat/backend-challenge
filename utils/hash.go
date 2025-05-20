package utils

import (
	"crypto/sha512"
	"fmt"
)

func Hash(data string) string {
	hasher := sha512.New()
	hasher.Write([]byte(data))
	bs := hasher.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
