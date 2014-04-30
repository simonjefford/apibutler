package testhelpers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type CheckableResponse struct {
	*httptest.ResponseRecorder
}

func (r *CheckableResponse) CheckStatus(expected int, t *testing.T) {
	if expected != r.Code {
		t.Errorf("Response was not %d, was %d", expected, r.Code)
	}
}

func (r *CheckableResponse) CheckBody(expected string, t *testing.T) {
	body := r.Body.String()
	if body != expected {
		t.Errorf("Unexpected response, \"%s\". Was the wrong endpoint hit?", body)
	}
}

func (r *CheckableResponse) CheckBodySubstring(substr string, t *testing.T) {
	body := r.Body.String()
	if !strings.Contains(body, substr) {
		t.Errorf("Response did not contain %s.", substr)
	}
}

func (r *CheckableResponse) CheckHeader(name, value string, t *testing.T) {
	actual := r.Header().Get(name)
	if actual != value {
		t.Errorf("Header %s was %s not %s", name, actual, value)
	}
}

func MakeTestableRequest(h http.Handler, r *http.Request) *CheckableResponse {
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, r)
	return &CheckableResponse{rw}
}
