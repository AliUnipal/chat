package usersvc

import (
	"context"
	"crypto/rsa"
	"errors"
	"time"

	"github.com/AliUnipal/chat/internal/models/user"
	"github.com/AliUnipal/chat/pkg/errcodes"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func NewJwtManager(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *jwtManager {
	return &jwtManager{
		privateKey,
		publicKey,
	}
}

type jwtManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func (s jwtManager) CreateToken(_ context.Context, u user.User) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Issuer:    "chat",
		Subject:   u.ID.String(),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	})

	token, err := claims.SignedString(s.privateKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s jwtManager) VerifyToken(_ context.Context, token string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}

	_, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New(errcodes.InvalidInput)
		}

		return s.publicKey, nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	usrId, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}

	return usrId, nil
}
