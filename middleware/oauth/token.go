package oauth

import (
	"log"
	"net/http"
	"strings"

	"fourth.com/apibutler/jsonconfig"
	"fourth.com/apibutler/middleware"
	"github.com/codegangsta/martini"
)

func init() {
	log.Println("Registering auth middleware")
	middleware.Register(
		"auth",
		middleware.NewDefinition("Authorization", authConstructor))
}

type AccessToken interface {
	AccessToken() string
	String() string
}

type accessToken struct {
	bearerToken string
}

func (a *accessToken) String() string {
	return a.AccessToken()
}

func (a *accessToken) AccessToken() string {
	return a.bearerToken
}

func authConstructor(cfg jsonconfig.Obj) (martini.Handler, error) {
	return GetIdFromRequest, nil
}

func GetIdFromRequest(req *http.Request, res http.ResponseWriter, ctx martini.Context) {
	token := req.Header.Get("Authorization")
	if token == "" {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	parts := strings.SplitAfter(token, " ")
	access := &accessToken{
		parts[1],
	}
	ctx.MapTo(access, (*AccessToken)(nil))
}
