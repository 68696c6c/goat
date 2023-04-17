package auth

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type TokensService struct {
	signatureKey    []byte
	userAuthHandler UserAuthHandler
}

func NewTokensService(signatureKey string, userAuthHandler UserAuthHandler) TokensService {
	return TokensService{
		signatureKey:    []byte(signatureKey),
		userAuthHandler: userAuthHandler,
	}
}

func (j TokensService) GetClaims(token string) (*Claims, error) {
	tokenInfo, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (any, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("parse error")
		}
		return j.signatureKey, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse jwt")
	}

	claims, ok := tokenInfo.Claims.(*Claims)
	if !ok || !tokenInfo.Valid {
		return nil, errors.Wrap(err, "failed to parse jwt claims")
	}

	return claims, nil
}

func (j TokensService) Token(cx context.Context, data *oauth2.GenerateBasic, isGenRefresh bool) (string, string, error) {
	u, err := j.userAuthHandler(cx, data.UserID)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to load user")
	}

	claims := &Claims{
		CustomClaims: GetCustomClaimsByUserLevel(u.Level),
		StandardClaims: jwt.StandardClaims{
			Audience:  data.Client.GetID(),
			Subject:   data.UserID,
			ExpiresAt: data.TokenInfo.GetAccessCreateAt().Add(data.TokenInfo.GetAccessExpiresIn()).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	access, err := token.SignedString(j.signatureKey)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to sign jwt")
	}

	var refresh string
	if isGenRefresh {
		t := uuid.NewSHA1(uuid.Must(uuid.NewRandom()), []byte(access)).String()
		refresh = base64.URLEncoding.EncodeToString([]byte(t))
		refresh = strings.ToUpper(strings.TrimRight(refresh, "="))
	}

	return access, refresh, nil
}
