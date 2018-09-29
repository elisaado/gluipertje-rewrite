package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

func empty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func randomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(b)
}
