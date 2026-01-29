package secret

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func Hash(pass string) ([]byte, error) {
	const (
		memory      = 64 * 1024
		iterations  = 3
		parallelism = 4
		keyLength   = 32
		saltLength  = 16
	)

	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	hash := argon2.IDKey(
		[]byte(pass),
		salt,
		iterations,
		memory,
		parallelism,
		keyLength,
	)

	encoded := fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		memory,
		iterations,
		parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)

	return []byte(encoded), nil
}

func Verify(pass string, hash []byte) (bool, error) {
	parts := strings.Split(string(hash), "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}

	// parts:
	// 0 = ""
	// 1 = argon2id
	// 2 = v=19
	// 3 = m=65536,t=3,p=4
	// 4 = salt
	// 5 = hash

	var memory, iterations uint32
	var parallelism uint8

	_, err := fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&memory,
		&iterations,
		&parallelism,
	)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	h := argon2.IDKey(
		[]byte(pass),
		salt,
		iterations,
		memory,
		parallelism,
		uint32(len(expectedHash)),
	)

	return subtle.ConstantTimeCompare(h, expectedHash) == 1, nil
}
