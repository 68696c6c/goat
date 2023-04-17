package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/manage"
	oauthModels "github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/pkg/errors"

	"github.com/68696c6c/web/app/models"
	"github.com/68696c6c/web/utils"
)

type Client oauthModels.Client

func (c Client) toOauthClient() *oauthModels.Client {
	var secret string
	if c.Public {
		secret = ""
	} else {
		secret = c.Secret
	}
	return &oauthModels.Client{
		ID:     c.ID,
		Secret: secret,
		Domain: c.Domain,
		Public: c.Public,
		UserID: c.UserID,
	}
}

type Config struct {
	SignatureKey string
	Clients      []Client
}

type Service interface {
	GenerateClientCredentials() (Client, error)
	Authorize(w http.ResponseWriter, r *http.Request) error
	GetToken(w http.ResponseWriter, r *http.Request) error
	ValidateToken(r *http.Request) (oauth2.TokenInfo, error)
	GetTokenClaims(token string) (*Claims, error)
	GetCurrentUserAndClaims(cx *gin.Context) (*models.User, *Claims, error)
}

type auth struct {
	store  *store.ClientStore
	server *server.Server
	tokens TokensService
}

type UserAuthHandler func(cx context.Context, email string) (*models.User, error)

func NewAuthService(c Config, userAuthHandler UserAuthHandler) (Service, error) {
	j := NewTokensService(c.SignatureKey, userAuthHandler)

	clientStore := store.NewClientStore()
	for _, client := range c.Clients {
		err := clientStore.Set(client.ID, client.toOauthClient())
		if err != nil {
			return nil, err
		}
	}

	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// Generate JWT tokens
	manager.MapAccessGenerate(&j)

	manager.MapClientStorage(clientStore)
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)

	srv := server.NewServer(server.NewConfig(), manager)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	srv.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (string, error) {
		email := r.PostForm.Get("username")
		user, err := userAuthHandler(context.Background(), email)
		if err != nil {
			return "", errors.Wrap(err, "failed to authenticate user")
		}
		return user.Email, nil
	})

	return &auth{
		store:  clientStore,
		server: srv,
		tokens: j,
	}, nil
}

func (a *auth) GenerateClientCredentials() (Client, error) {
	clientId, err := makeRandomHexString()
	if err != nil {
		return Client{}, err
	}

	clientSecret, err := makeRandomHexString()
	if err != nil {
		return Client{}, err
	}

	return Client{
		ID:     clientId,
		Secret: clientSecret,
	}, nil
}

func (a *auth) Authorize(w http.ResponseWriter, r *http.Request) error {
	return a.server.HandleAuthorizeRequest(w, r)
}

func (a *auth) GetToken(w http.ResponseWriter, r *http.Request) error {
	return a.server.HandleTokenRequest(w, r)
}

func (a *auth) ValidateToken(r *http.Request) (oauth2.TokenInfo, error) {
	return a.server.ValidationBearerToken(r)
}

func (a *auth) GetTokenClaims(token string) (*Claims, error) {
	return a.tokens.GetClaims(token)
}

func (a *auth) GetCurrentUserAndClaims(cx *gin.Context) (*models.User, *Claims, error) {
	user, ok := utils.GetCurrentUser(cx)
	if !ok || user == nil {
		return nil, nil, errors.New("invalid user")
	}

	token, ok := utils.GetCurrentUserToken(cx)
	if !ok || token == "" {
		return nil, nil, errors.New("invalid jwt")
	}

	claims, err := a.GetTokenClaims(token)
	if err != nil {
		return nil, nil, err
	}

	return user, claims, nil
}

func makeRandomHexString() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
