package oauth

import (
	"github.com/codegangsta/martini"
	"net/http"
	"strings"
)

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
