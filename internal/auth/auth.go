package auth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/amirdaraby/go-todo-list-api/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type AuthIdKey string

func NewToken(userId uint) (string, error) {

	config := config.Get()

	privateKey, err := parseECDSAPrivateKey(config.PrivateKey)

	if err != nil {
		return "", errors.New("private key cannot be parsed")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"id": userId,
	})

	tokenString, err := token.SignedString(privateKey)

	if err != nil {
		return "", err
	}

	return tokenString, err
}

// returns user id if token is valid
func ValidateToken(tokenString string) (uint, error) {

	config := config.Get()

	publicKey, err := parseECDSAPublicKey(config.PublicKey)

	if err != nil {
		return 0, err
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodECDSA); !ok {
			return 0, errors.New("invalid token")
		}

		return publicKey, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		id, ok := claims["id"].(float64)

		if !ok {
			return 0, errors.New("invalid token")
		}

		return uint(id), nil
	}

	return 0, errors.New("invalid token")
}

func parseECDSAPublicKey(pemKey []byte) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemKey))

	if block == nil {
		return nil, errors.New("invalid public key format")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := pub.(*ecdsa.PublicKey)

	if !ok {
		return nil, errors.New("key is not an ECDSA public key")
	}

	return publicKey, nil
}

func parseECDSAPrivateKey(pemKey []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemKey))

	if block == nil {
		return nil, errors.New("invalid private key format")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)

	return privateKey, err
}
