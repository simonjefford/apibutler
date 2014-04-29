package dashboard

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"fourth.com/apibutler/metadata"
)

type dummyApiServer struct {
	apis []*metadata.Api
}

func (s *dummyApiServer) UpdateApis(apis []*metadata.Api) {
	s.apis = apis
}

func (s *dummyApiServer) UpdateApps(apps metadata.ApplicationTable) {
}

func (s *dummyApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

type dummyApiStore struct {
}

func (s *dummyApiStore) AddApi(a *metadata.Api) {
}

func (s *dummyApiStore) Apis() ([]*metadata.Api, error) {
	return []*metadata.Api{
		&metadata.Api{Fragment: "/cool"},
	}, nil
}

func (s *dummyApiStore) Forget(path string) {
}

func TestRouter(t *testing.T) {
	apiserver := &dummyApiServer{}
	apistore := &dummyApiStore{}

	d := NewDashboardServer("/", apiserver, apistore)

	req, err := http.NewRequest("GET", "http://example.com/apis", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	d.ServeHTTP(w, req)
	body := w.Body.String()
	if !strings.Contains(body, `"fragment":"/cool"`) {
		t.Fatalf("unexpected response: %s", body)
	}
}
