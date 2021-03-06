package dashboard

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"fourth.com/apibutler/apiproxyserver"
	"fourth.com/apibutler/metadata"
	"fourth.com/apibutler/middleware"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
)

func NewDashboardServer(path string, proxy apiproxyserver.APIProxyServer, apiStore metadata.ApiStore, stackStore middleware.StackStore) http.Handler {
	m := martini.New()
	m.Use(martini.Logger())
	l := log.New(os.Stdout, "[dashboard server] ", 0)
	m.Map(l)
	m.MapTo(proxy, (*apiproxyserver.APIProxyServer)(nil))
	m.MapTo(apiStore, (*metadata.ApiStore)(nil))
	m.MapTo(stackStore, (*middleware.StackStore)(nil))
	m.Use(martini.Recovery())
	m.Use(martini.Static(path))
	m.Use(render.Renderer())
	setupRouter(m)

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
	r.Get("/middlewares", middlewaresGetHandler)
	r.Get("/stacks", stacksGetHandler)
	r.Get("/stacks/:id", stackGetHandler)
	r.Post("/stacks", stacksPostHandler)
	m.Action(r.Handle)
}

type ApisPayload struct {
	Apis []*metadata.Api `json:"apis"`
}

type SingleApiPayload struct {
	Api metadata.Api `json:"api"`
}

type SingleStackPayload struct {
	Stack *middleware.Stack `json:"stack"`
}

type SingleAppPayload struct {
	App *metadata.Application `json:"app"`
}

type ApplicationsPayload struct {
	Apps []*metadata.Application `json:"apps"`
}

type MiddlewaresPayload struct {
	Middlewares []*middleware.Definition `json:"middlewares"`
}

type StacksPayload struct {
	Stacks []*middleware.Stack `json:"stacks"`
}

type StackPayload struct {
	Stack *middleware.Stack `json:"stack"`
}

func middlewaresGetHandler(rdr render.Render) {
	mw := middleware.Definitions()

	rdr.JSON(http.StatusOK, &MiddlewaresPayload{mw})
}

func stacksGetHandler(rdr render.Render, stackStore middleware.StackStore) {
	// TODO err here
	st, _ := stackStore.Stacks()

	p := &StacksPayload{st}

	rdr.JSON(http.StatusOK, p)
}

func stackGetHandler(rdr render.Render, stackStore middleware.StackStore, params martini.Params) {
	id := params["id"]
	stack, _ := stackStore.Stack(id)

	rdr.JSON(http.StatusOK, &StackPayload{stack})
}

func stacksPostHandler(req *http.Request, rdr render.Render, stackStore middleware.StackStore) {
	decoder := json.NewDecoder(req.Body)
	var s SingleStackPayload
	err := decoder.Decode(&s)
	if err != nil {
		rdr.JSON(http.StatusBadRequest, statusResponse{err.Error()})
		return
	}
	log.Println(&s)
	stackStore.AddStack(s.Stack)
	rdr.JSON(http.StatusCreated, s)
}

func apisGetHandler(rdr render.Render, apiStore metadata.ApiStore) {
	apis, err := apiStore.Apis()
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

func appGetHandler(rdr render.Render, params martini.Params) {
	id, _ := strconv.Atoi(params["id"])
	a := SingleAppPayload{metadata.GetSingleApplication(id)}
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

func apisPutHandler(req *http.Request, rdr render.Render, params martini.Params) {
	decoder := json.NewDecoder(req.Body)
	var a SingleApiPayload
	decoder.Decode(&a)
	log.Println(a)
	rdr.JSON(http.StatusCreated, a)
}

func apisPostHandler(req *http.Request, rdr render.Render, proxy apiproxyserver.APIProxyServer, apiStorage metadata.ApiStore) {
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
