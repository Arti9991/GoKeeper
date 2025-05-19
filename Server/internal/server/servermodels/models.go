package servermodels

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

type UserInfo struct {
	UserID   string
	Register bool
}

var ErrorBadToken = errors.New("Error in auth interceptor: Bad token in session info")

// функция хэширования для паролей (sha256)
func CodePassword(src string) string {
	hash := sha256.New()
	hash.Write([]byte(src))
	dst := hash.Sum(nil)
	return hex.EncodeToString(dst)
}

var ErrorNoSuchUser = errors.New("No user with this login")
