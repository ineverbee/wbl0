package app

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ineverbee/wbl0/internal/store"
	"github.com/stretchr/testify/require"
)

func TestHandlers(t *testing.T) {
	router := mux.NewRouter()
	router.Handle("/", limit(errorHandler(GetHomePageHandler()))).Methods("GET", "POST")
	router.Handle("/data/{id}", limit(errorHandler(GetDataPageHandler()))).Methods("GET")

	app = &App{
		&http.Server{},
		&store.DBMock{},
		&store.CacheMock{},
	}

	tc := []struct {
		method, target string
		body           io.Reader
		code           int
	}{
		{"GET", "/notfound", nil, http.StatusNotFound},
		{"GET", "/", nil, http.StatusOK},
		{"GET", "/data/1", nil, http.StatusOK},
		{"GET", "/data/-10", nil, http.StatusBadRequest},
		{"GET", "/data/NaN", nil, http.StatusBadRequest},
	}
	for _, c := range tc {
		request(t, router, c.method, c.target, c.body, c.code)
	}

	req := httptest.NewRequest("POST", "/", nil)
	req.URL.RawQuery += "id=NaN"
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)

	req = httptest.NewRequest("POST", "/", nil)
	req.URL.RawQuery += "id=1"
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusFound, rr.Code)
}

func request(t *testing.T, handler http.Handler, method, target string, body io.Reader, code int) {
	req := httptest.NewRequest(method, target, body)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	require.Equal(t, code, rr.Code)
}
