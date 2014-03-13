package dashboard

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"fourth.com/ratelimit/applications"
	"fourth.com/ratelimit/limiter"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
)

var (
	ratelimiter limiter.RateLimit
)

func NewDashboardServer(r limiter.RateLimit, path string) http.Handler {
	ratelimiter = r
	m := martini.New()
	m.Use(martini.Logger())
	l := log.New(os.Stdout, "[dashboard server] ", 0)
	m.Map(l)
	m.Use(martini.Recovery())
	m.Use(martini.Static(path))
	m.Use(render.Renderer())
	setupRouter(m)
	return m
}

func setupRouter(m *martini.Martini) {
	r := martini.NewRouter()
	r.Post("/paths", pathsPostHandler)
	r.Get("/paths", pathsGetHandler)
	r.Get("/apps", appsGetHandler)
	m.Action(r.Handle)
}

func pathsGetHandler(rdr render.Render) {
	rdr.JSON(200, ratelimiter.Paths())
}

func appsGetHandler(rdr render.Render) {
	rdr.JSON(200, applications.Get())
}

type statusResponse struct {
	Message string `json:message`
}

func pathsPostHandler(res http.ResponseWriter, req *http.Request, rdr render.Render) {
	decoder := json.NewDecoder(req.Body)
	var p limiter.Path
	err := decoder.Decode(&p)
	if err != nil {
		rdr.JSON(http.StatusBadRequest, statusResponse{err.Error()})
		return
	}
	log.Println(p)
	ratelimiter.AddPath(p)
	rdr.JSON(http.StatusCreated, statusResponse{"Created"})
}
