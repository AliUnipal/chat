package usersvc

import "github.com/AliUnipal/chat/pkg/secret"

func NewHasher() *passHasher {
	return &passHasher{}
}

type passHasher struct {
}

func (h *passHasher) Hash(s string) ([]byte, error) {
	return secret.Hash(s)
}

func (h *passHasher) Verify(s string, hash []byte) (bool, error) {
	return secret.Verify(s, hash)
}
