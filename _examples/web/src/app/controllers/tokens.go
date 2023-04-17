package controllers

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/68696c6c/goat"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/pkg/errors"

	"github.com/68696c6c/web/app/lib/auth"
	"github.com/68696c6c/web/app/repos"
	"github.com/68696c6c/web/utils"
)

type TokensController interface {
	AuthorizeForm(cx *gin.Context)
	Authorize(cx *gin.Context)
	Callback(cx *gin.Context)
	Exchange(cx *gin.Context)
	View(cx *gin.Context)
}

type accessTokens struct {
	repo repos.UsersRepo
	auth auth.Service
}

func NewTokensController(repo repos.UsersRepo, authService auth.Service) TokensController {
	return accessTokens{
		repo: repo,
		auth: authService,
	}
}

// AuthorizeForm renders a very basic login form example for use with the included Insomnia requests.  Since Goat is
// primarily intended for making CLI and headless REST applications, it does not provide support rendering HTML views.
// In a real application, this page should exist in the front-end instead of here.
func (c accessTokens) AuthorizeForm(cx *gin.Context) {
	responseType := cx.Query("response_type")
	if responseType == "" {
		goat.RespondBadRequest(cx, errors.New("missing response_type"))
		return
	}
	clientId := cx.Query("client_id")
	if clientId == "" {
		goat.RespondBadRequest(cx, errors.New("missing client_id"))
		return
	}
	redirectUri := cx.Query("redirect_uri")
	if redirectUri == "" {
		goat.RespondBadRequest(cx, errors.New("missing redirect_uri"))
		return
	}
	codeChallenge := cx.Query("code_challenge")
	if codeChallenge == "" {
		goat.RespondBadRequest(cx, errors.New("missing code_challenge"))
		return
	}

	content, err := parseAuthCodeFormTemplate(authCodeFormData{
		GrantType:     "authorization_code",
		SubmitUrl:     goat.GetUrl(auth.AuthorizeLinkKey).String(),
		ClientId:      clientId,
		RedirectUri:   goat.GetUrl(auth.AuthorizeCallbackLinkKey).String(),
		Username:      "",
		ResponseType:  responseType,
		CodeChallenge: codeChallenge,
	})
	if err != nil {
		goat.RespondServerError(cx, errors.Wrap(err, "failed to build template"))
		return
	}

	_, err = cx.Writer.Write([]byte(content))
	if err != nil {
		goat.RespondServerError(cx, errors.Wrap(err, "failed to render template"))
		return
	}

	cx.Writer.WriteHeader(http.StatusOK)
	return
}

func (c accessTokens) Authorize(cx *gin.Context) {
	err := c.auth.Authorize(cx.Writer, cx.Request)
	if err != nil {
		goat.RespondServerError(cx, err)
		return
	}
}

// Callback is an example Oauth redirect_uri for use with the included Insomnia requests.
// In a real application, this page should exist in the front-end instead of here.
func (c accessTokens) Callback(cx *gin.Context) {

	accessCode := cx.Query("code")
	if accessCode == "" {
		goat.RespondBadRequest(cx, errors.New("missing access code param"))
		return
	}

	goat.RespondOk(cx, map[string]string{
		"code": accessCode,
	})
}

func (c accessTokens) Exchange(cx *gin.Context) {
	err := c.auth.GetToken(cx.Writer, cx.Request)
	if err != nil {
		goat.RespondServerError(cx, err)
		return
	}
}

// View renders the JWT and associated user (subject) of the current request.
// This endpoint is just for validating the auth flow works as expected and isn't necessary in real-world applications.
func (c accessTokens) View(cx *gin.Context) {
	accessJwt := cx.GetString(utils.ContextKeyCurrentUserToken)
	if accessJwt == "" {
		goat.RespondServerError(cx, errors.New("jwt is not set"))
		return
	}

	claims, err := c.auth.GetTokenClaims(accessJwt)
	if err != nil {
		goat.RespondServerError(cx, errors.Wrap(err, "failed to parse jwt claims"))
		return
	}

	goat.RespondOk(cx, map[string]any{
		"token":  accessJwt,
		"claims": claims,
	})
}

func parseAuthCodeFormTemplate(data authCodeFormData) (string, error) {
	var tpl bytes.Buffer
	t := template.Must(template.New("auth code form").Parse(authCodeFormTemplate))
	err := t.Execute(&tpl, data)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse auth code form")
	}
	return tpl.String(), nil
}

type authCodeFormData struct {
	GrantType     oauth2.GrantType `json:"grant_type"`
	SubmitUrl     string           `json:"submit_url"`
	ClientId      string           `json:"client_id"`
	RedirectUri   string           `json:"redirect_uri"`
	Username      string           `json:"username"`
	ResponseType  string           `json:"response_type"`
	CodeChallenge string           `json:"code_challenge"`
}

const authCodeFormTemplate = `<html>
	<body>
		<form action="{{ .SubmitUrl }}" method="POST">
			<input type="hidden" name="grant_type" value="{{ .GrantType }}" />
			<input type="hidden" name="client_id" value="{{ .ClientId }}" />
			<input type="hidden" name="redirect_uri" value="{{ .RedirectUri }}" />
			<input type="hidden" name="response_type" value="{{ .ResponseType }}" />
			<input type="hidden" name="code_challenge" value="{{ .CodeChallenge }}" />
			<label for="username">Email</label>
			<input id="username" type="text" name="username" value="{{ .Username }}" />
			<button type="submit">Submit</button>
		</form>
	</body>
</html>`
