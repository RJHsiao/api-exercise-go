package utilities

import (
	"crypto/sha256"
	"fmt"
)

// GetSha256SumFromString generate sha256sum string by given string
func GetSha256SumFromString(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
