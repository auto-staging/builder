package model

import (
	"crypto/sha256"
	"encoding/hex"
)

func getRandomValueForBranch(branch string) string {
	random := sha256.Sum256([]byte(branch))
	return hex.EncodeToString(random[5:10])
}
