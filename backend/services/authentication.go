package authentication

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	clientAuth "backend/clients"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func GetIssuerClient() (*oauth2.Config, error) {
	ctx := context.Background()

	issuer, err := oidc.NewProvider(ctx, os.Getenv("SSO_REALMS"))
	if err != nil {
		return nil, err
	}

	config := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Endpoint:     issuer.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
	return config, nil
}

func GetAuthorizationUrl(e echo.Context) error {
	config, err := GetIssuerClient()
	if err != nil {
		e.Logger().Error(err)
		e.String(http.StatusInternalServerError, "Internal Server Error")
		return err
	}
	e.Redirect(http.StatusTemporaryRedirect, config.AuthCodeURL("state"))
	return nil
}

func GetTokenClient(e echo.Context) error {
	code := e.QueryParam("code")

	response, err := clientAuth.GetTokenSSO(code)
	if err != nil {
		return err
	}

	var jsonTokenResponse map[string]interface{}
	json.Unmarshal([]byte(string(response)), &jsonTokenResponse)

	return e.JSON(http.StatusOK, jsonTokenResponse)
}

func GetLogout(e echo.Context) error {
	token := e.QueryParam("token")
	fmt.Println(token)
	if token == "" || token == "undefined" {
		return e.String(http.StatusBadRequest, "Bad Request")
	}

	clientAuth.GetLogoutSSO(token)

	return nil
}
