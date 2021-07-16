package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/dp1140a/geoip/api/geo"
	"github.com/dp1140a/geoip/api/httplogging"
	"github.com/dp1140a/geoip/models"
	"net/http"

	"github.com/go-chi/jwtauth"
	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

var tokenAuth *jwtauth.JWTAuth

type Server struct {
	Config   *Config
	Logger   *log.Logger
	Context  context.Context
	Router   *chi.Mux
	Handlers []models.HandlerIFace
}

func init() {

	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:
}

func NewServer(ctx context.Context) (s *Server, err error) {
	config, err := InitConfig()
	s = &Server{
		Config:  config,
		Logger:  log.New(),
		Context: ctx,
		Router:  chi.NewRouter(),
	}
	tokenAuth = jwtauth.New("HS256", []byte(config.JWTSecret), nil)
	s.Logger.Formatter = &log.JSONFormatter{
		// disable, as we set our own
		DisableTimestamp: false,
	}

	//Add and init middleware
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.RedirectSlashes)
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.Compress(5, "gzip"))
	log.Info("Request Logger Enabled: ", config.LoggingConfig.Enabled)
	if config.LoggingConfig.Enabled == true {
		s.Router.Use(httplogging.NewStructuredLogger(s.Logger, s.Config.LoggingConfig))
	}
	//s.Router.Use(middleware.Logger)
	if config.EnableCORS {
		s.Router.Use(corsConfig().Handler)
	}

	//Append Services
	s.Handlers = append(s.Handlers, geo.NewGeoHandler(s.Context, tokenAuth)) // Geo Service

	//Append the base services Last
	s.Handlers = append(s.Handlers, NewBaseHandler(s.Logger, s.Router))

	s.Router.Mount("/v1", s.registerRoutes())
	return s, err
}

func (s *Server) registerRoutes() chi.Router {
	r := chi.NewRouter()
	for _, handler := range s.Handlers {
		log.Info("Opening Route Prefix ", handler.GetPrefix())
		r.Route(handler.GetPrefix(), func(r chi.Router) {
			for _, route := range handler.GetRoutes() {
				log.Info("Adding route ", route.Name)
				if route.Protected { //Protected Route
					r.Group(func(r chi.Router) {
						// Seek, verify and validate JWT tokens
						r.Use(jwtauth.Verifier(tokenAuth))
						r.Use(MyAuthenticator)
						r.Method(route.Method, route.Pattern, route.HandlerFunc)
					})
				} else { // Public Route
					r.Method(route.Method, route.Pattern, route.HandlerFunc)
				}
			}
		})
	}
	return r
}

func (s *Server) ServeAPI() {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%v", s.Config.Hostname, s.Config.Port),
		Handler: s.Router,
	}

	done := make(chan struct{})
	go func() {
		<-s.Context.Done()
		if err := srv.Shutdown(s.Context); err != nil {
			log.Fatal(err)
		}

		close(done)
	}()

	var cer tls.Certificate
	transport := "http"

	if s.Config.UseHttps && (s.Config.TLSCert != "" && s.Config.TLSKey != "") {
		var err error
		cer, err = tls.LoadX509KeyPair(s.Config.TLSCert, s.Config.TLSKey)
		if err != nil {
			log.Error("failed to load x509 key pair", err)
			log.Info("Stopping")
		}
		transport = "https"

		// Sensible default
		var tlsMinVersion uint16 = tls.VersionTLS12

		switch s.Config.TLSMinVersion {
		case "1.0":
			log.Warn("Setting the minimum version of TLS to 1.0 - this is discouraged. Please use 1.2 or 1.3")
			tlsMinVersion = tls.VersionTLS10
		case "1.1":
			log.Warn("Setting the minimum version of TLS to 1.1 - this is discouraged. Please use 1.2 or 1.3")
			tlsMinVersion = tls.VersionTLS11
		case "1.2":
			tlsMinVersion = tls.VersionTLS12
		case "1.3":
			tlsMinVersion = tls.VersionTLS13
		}

		strictCiphers := []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		}

		// nil uses the default cipher suite
		var cipherConfig []uint16 = nil

		// TLS 1.3 does not support configuring the Cipher suites
		if tlsMinVersion != tls.VersionTLS13 && s.Config.HttpTLSStrictCiphers {
			cipherConfig = strictCiphers
		}

		srv.TLSConfig = &tls.Config{
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			MinVersion:               tlsMinVersion,
			CipherSuites:             cipherConfig,
		}
	}

	log.Infof("Listening %s on %s:%v", transport, s.Config.Hostname, s.Config.Port)
	if cer.Certificate != nil {
		if err := srv.ListenAndServeTLS(s.Config.TLSCert, s.Config.TLSKey); err != http.ErrServerClosed {
			log.Error(err)
		}
	} else {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Error(err)
		}
	}
	<-done
}

func MyAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			log.Error("Token Error: ", err)
			http.Error(w, http.StatusText(401), 401)
			return
		}

		if token == nil || !token.Valid {
			log.Error("Token Invalid")
			http.Error(w, http.StatusText(401), 401)
			return
		}
		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

func corsConfig() *cors.Cors {
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	return cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           86400, // Maximum value not ignored by any of major browsers
	})
}
