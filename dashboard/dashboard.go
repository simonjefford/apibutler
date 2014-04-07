package dashboard

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"fourth.com/apibutler/apiproxyserver"
	"fourth.com/apibutler/metadata"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
)

var (
	apiStorage metadata.ApiStorage
)

func NewDashboardServer(path string, proxy apiproxyserver.APIProxyServer) http.Handler {
	m := martini.New()
	m.Use(martini.Logger())
	l := log.New(os.Stdout, "[dashboard server] ", 0)
	m.Map(l)
	m.MapTo(proxy, (*apiproxyserver.APIProxyServer)(nil))
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
	r.Post("/apis", apisPostHandler)
	r.Get("/apis", apisGetHandler)
	r.Put("/apis/:id", apisPutHandler)
	r.Get("/apps", appsGetHandler)
	r.Get("/apps/:id", appGetHandler)
	r.Put("/apps/:id", appPutHandler)
	m.Action(r.Handle)
}

type ApisPayload struct {
	Apis []*metadata.Api `json:"apis"`
}

type SingleApiPayload struct {
	Api metadata.Api `json:"api"`
}

type SingleAppPayload struct {
	App metadata.Application `json:"app"`
}

type ApplicationsPayload struct {
	Apps []*metadata.Application `json:"apps"`
}

func apisGetHandler(rdr render.Render) {
	apis, err := apiStorage.Apis()
	if err != nil {
		rdr.JSON(500, nil)
	}
	a := ApisPayload{apis}
	rdr.JSON(200, a)
}

func appsGetHandler(rdr render.Render) {
	a := ApplicationsPayload{metadata.GetApplicationsList()}
	rdr.JSON(200, a)
}

type statusResponse struct {
	Message string `json:message`
}

func appGetHandler(res http.ResponseWriter, req *http.Request, rdr render.Render, params martini.Params) {
	id, _ := strconv.Atoi(params["id"])
	a := SingleAppPayload{*metadata.GetSingleApplication(id)}
	rdr.JSON(200, a)
}

func appPutHandler(req *http.Request, rdr render.Render, params martini.Params, proxy apiproxyserver.APIProxyServer) {
	decoder := json.NewDecoder(req.Body)
	var a SingleAppPayload
	decoder.Decode(&a)
	a.App.ID, _ = strconv.Atoi(params["id"])
	metadata.ChangeApplication(a.App)
	proxy.UpdateApps(metadata.GetApplicationsTable())
	rdr.JSON(http.StatusAccepted, a)
}

func apisPutHandler(res http.ResponseWriter, req *http.Request, rdr render.Render, params martini.Params) {
	decoder := json.NewDecoder(req.Body)
	var a SingleApiPayload
	decoder.Decode(&a)
	id, _ := strconv.Atoi(params["id"])
	a.Api.ID = int64(id)
	log.Println(a)
	rdr.JSON(http.StatusCreated, a)
}

func apisPostHandler(res http.ResponseWriter, req *http.Request, rdr render.Render, proxy apiproxyserver.APIProxyServer) {
	decoder := json.NewDecoder(req.Body)
	var a SingleApiPayload
	err := decoder.Decode(&a)
	if err != nil {
		rdr.JSON(http.StatusBadRequest, statusResponse{err.Error()})
		return
	}
	log.Println(a)
	apiStorage.AddApi(&a.Api)
	apis, _ := apiStorage.Apis()
	proxy.UpdateApis(apis)
	rdr.JSON(http.StatusCreated, a)
}
