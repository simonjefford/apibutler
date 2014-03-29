package dashboard

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"fourth.com/apibutler/metadata"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
)

var (
	apiStorage metadata.ApiStorage
)

func NewDashboardServer(path string) http.Handler {
	m := martini.New()
	m.Use(martini.Logger())
	l := log.New(os.Stdout, "[dashboard server] ", 0)
	m.Map(l)
	m.Use(martini.Recovery())
	m.Use(martini.Static(path))
	m.Use(render.Renderer())
	setupRouter(m)

	a, err := metadata.GetApiStore()
	apiStorage = a

	if err != nil {
		panic(err)
	}

	return m
}

func setupRouter(m *martini.Martini) {
	r := martini.NewRouter()
	r.Post("/paths", pathsPostHandler)
	r.Get("/paths", pathsGetHandler)
	r.Get("/apps", appsGetHandler)
	r.Put("/paths/:id", pathsPutHandler)
	m.Action(r.Handle)
}

type ApiPayload struct {
	Apis []metadata.Api `json:"paths"`
}

type SingleApiPayload struct {
	Api metadata.Api `json:"path"`
}

func pathsGetHandler(rdr render.Render) {
	p := ApiPayload{apiStorage.Apis()}
	rdr.JSON(200, p)
}

func appsGetHandler(rdr render.Render) {
	rdr.JSON(200, metadata.GetApplicationsList())
}

type statusResponse struct {
	Message string `json:message`
}

func pathsPutHandler(res http.ResponseWriter, req *http.Request, rdr render.Render, params martini.Params) {
	decoder := json.NewDecoder(req.Body)
	var p SingleApiPayload
	decoder.Decode(&p)
	id, _ := strconv.Atoi(params["id"])
	p.Api.ID = int64(id)
	log.Println(p)
	rdr.JSON(http.StatusCreated, p)
}

func pathsPostHandler(res http.ResponseWriter, req *http.Request, rdr render.Render) {
	decoder := json.NewDecoder(req.Body)
	var p SingleApiPayload
	err := decoder.Decode(&p)
	if err != nil {
		rdr.JSON(http.StatusBadRequest, statusResponse{err.Error()})
		return
	}
	log.Println(p)
	apiStorage.AddApi(p.Api)
	rdr.JSON(http.StatusCreated, p)
}
