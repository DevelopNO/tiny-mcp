package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Signer struct {
	secret []byte
}

func NewSigner(secret string) *Signer {
	return &Signer{secret: []byte(secret)}
}

func (s *Signer) Sign(peerIP, team string) (string, error) {
	claims := jwt.MapClaims{
		"peer_ip": peerIP,
		"team":    team,
		"exp":     time.Now().Add(5 * time.Minute).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(s.secret)
}

func (s *Signer) Verify(token string) (peerIP, team string, err error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return "", "", err
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok || !parsed.Valid {
		return "", "", errors.New("invalid token")
	}
	ip, ok1 := claims["peer_ip"].(string)
	teamVal, ok2 := claims["team"].(string)
	if !ok1 || !ok2 {
		return "", "", errors.New("missing claims")
	}
	return ip, teamVal, nil
}
