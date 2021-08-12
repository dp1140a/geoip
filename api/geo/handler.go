package geo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dp1140a/geoip/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type GeoHandler struct {
	models.Handler
	tokenAuth *jwtauth.JWTAuth
}

func NewGeoHandler(ctx context.Context, tokenAuth *jwtauth.JWTAuth) models.HandlerIFace {
	gh := GeoHandler{
		Handler: models.Handler{
			Prefix:  "/geo",
			Service: NewGeoIPService(ctx),
		},
		tokenAuth: tokenAuth,
	}
	gh.Routes = []models.Route{
		{
			Name:        "getGeoIP",
			Method:      http.MethodGet,
			Pattern:     "/{ipaddress}",
			HandlerFunc: gh.locate,
		},
		{
			Name:        "getGeoIP",
			Method:      http.MethodPost,
			Pattern:     "/batch",
			HandlerFunc: gh.batchLogEntries,
		},
	}

	return gh
}

func (gh GeoHandler) GetRoutes() []models.Route {
	return gh.Routes
}

func (gh GeoHandler) GetService() models.Service {
	return gh.Service
}

func (gh GeoHandler) GetPrefix() string {
	return gh.Prefix
}

/**
GET /geo/{ipaddress}  (200, 500) -- Returns JSON of GeoLocation
*/
func (gh *GeoHandler) locate(w http.ResponseWriter, r *http.Request) {
	record, err := gh.Service.(GeoIPService).locate(chi.URLParam(r, "ipaddress"))
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Println(record)
	json.NewEncoder(w).Encode(record)
}

/**
POST /geo/lookup [logEntries] (201, 500) -- For a batch of line protocol from tail will perform a geoip lookup and send a new point to InfluxDB
*/
func (gh *GeoHandler) batchLogEntries(w http.ResponseWriter, r *http.Request) {
	// Receive a block of line protocol lines
	// tail,host=dyson,path=/home/dave/Desktop/logTest/filter.log data_length="0",offset="0",dest_ip="173.160.205.9",tracker="1521923716",reason="match",program="filterlog",pid="45720",ttl="113",iface="em0",flags="none",ip_ver="4",direction="in",src_port="14553",rule="102",action="block",tos="0x20",src_ip="205.185.117.79",dest_port="22",id="14012",timestamp="Jul  8 16:54:53",length="48" 1625793903302961350
	b, _ := io.ReadAll(r.Body)

	err := gh.Service.(GeoIPService).mapPointsWithLookup(b)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}

	w.WriteHeader(http.StatusCreated)
}
