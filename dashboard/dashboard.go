package dashboard

import (
	"encoding/json"
	"fourth.com/ratelimit/limiter"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"log"
	"net/http"
	"os"
)

var (
	ratelimiter *limiter.RateLimit
)

func NewDashboardServer(r *limiter.RateLimit, path string) http.Handler {
	ratelimiter = r
	rtr := martini.NewRouter()
	m := martini.New()
	m.Use(martini.Logger())
	l := log.New(os.Stdout, "[dashboard server] ", 0)
	m.Map(l)
	m.Use(martini.Recovery())
	m.Use(martini.Static(path))
	m.Use(render.Renderer())
	rtr.Post("/paths", pathsPostHandler)
	rtr.Get("/paths", pathsGetHandler)
	m.Action(rtr.Handle)
	return &martini.ClassicMartini{m, rtr}
}

func pathsGetHandler(rdr render.Render) {
	rdr.JSON(200, ratelimiter.Paths())
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
