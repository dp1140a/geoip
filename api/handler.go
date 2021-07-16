package api

import (
	"encoding/json"
	"fmt"
	"github.com/dp1140a/geoip/models"
	"github.com/dp1140a/geoip/version"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type BaseHandler struct {
	models.Handler
	mux *chi.Mux
}

func NewBaseHandler(logger *log.Logger, m *chi.Mux) models.HandlerIFace {
	bh := BaseHandler{
		mux: m,
		Handler: models.Handler{
			Prefix: "/",
		},
	}

	bh.Routes = []models.Route{
		{
			Name:        "ping",
			Method:      http.MethodGet,
			Pattern:     "/ping",
			HandlerFunc: bh.pingHandler,
		},
		{
			Name:        "api",
			Method:      http.MethodGet,
			Pattern:     "/api",
			HandlerFunc: bh.getAPI,
		},
		{
			Name:        "version",
			Method:      http.MethodGet,
			Pattern:     "/version",
			HandlerFunc: bh.getVersion,
		},
	}
	return bh
}

func (bh BaseHandler) GetRoutes() []models.Route {
	return bh.Routes
}

func (bh BaseHandler) GetService() models.Service {
	return bh.Service
}

func (bh BaseHandler) GetPrefix() string {
	return bh.Prefix
}

/**
GET /ping  (200) -- Returns "OK"
*/
func (bh *BaseHandler) pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("x-powered-by", "bacon")
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
}

/**
GET /api  (200) -- Returns JSON of API
*/
func (bh *BaseHandler) getAPI(w http.ResponseWriter, _ *http.Request) {
	var routes []string
	chi.Walk(bh.mux, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		routes = append(routes, fmt.Sprintf("%s %s", method, route))
		return nil
	})
	w.Header().Set("x-powered-by", "bacon")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(routes)
}

/**
GET /version  (200) -- Returns JSON of the current version
*/
func (bh *BaseHandler) getVersion(w http.ResponseWriter, _ *http.Request) {
	versionInfo := version.NewVersionInfo()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(versionInfo)
}
