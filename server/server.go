package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"kando-backend/config"
	"kando-backend/handlers"
	"kando-backend/ioc"
	"kando-backend/log"
	"kando-backend/middlewares"
	"net/http"
	"os"
	"os/signal"
)

func ServeApi(dp *ioc.DependencyProvider) {
	log.Logger.Infof("Serving api on %s:%d",
		config.C.Server.Host,
		config.C.Server.Port)

	r := mux.NewRouter()

	r.Use(middlewares.AccessLogMiddleware)

	r.Use(middlewares.MaxReadBytesMiddleware)
	r.Use(middlewares.EnforceJsonMiddleware)

	r.Use(middlewares.ScopeMiddleware(dp))
	r.Use(middlewares.ErrorHandlingMiddleware)

	// unauthorized api endpoints
	r.HandleFunc("/api/health/", handlers.HealthHandler)
	r.HandleFunc("/api/users/", handlers.CreateUserHandler).Methods("POST")
	r.HandleFunc("/api/sessions/", handlers.CreateSessionHandler).Methods("POST")

	// authorized api endpoints
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middlewares.AuthenticationMiddleware)

	sessions := api.PathPrefix("/sessions").Subrouter()
	sessions.HandleFunc("/", handlers.DeleteSessionHandler).Methods("DELETE")

	manufacturers := api.PathPrefix("/manufacturers").Subrouter()
	manufacturers.HandleFunc("/", handlers.CreateManufacturerHandler).Methods("POST")
	manufacturers.HandleFunc("/", handlers.GetManufacturersHandler).Methods("GET")
	manufacturers.HandleFunc("/{manufacturerId}/", handlers.UpdateManufacturerHandler).Methods("PUT")

	assetTypes := api.PathPrefix("/asset-types").Subrouter()
	assetTypes.HandleFunc("/", handlers.CreateAssetTypeHandler).Methods("POST")
	assetTypes.HandleFunc("/", handlers.GetAssetTypesHandler).Methods("GET")
	assetTypes.HandleFunc("/{assetTypeId}/", handlers.UpdateAssetTypeHandler).Methods("PUT")

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%d", config.C.Server.Host, config.C.Server.Port),
		WriteTimeout: config.C.Server.WriteTimeout,
		ReadTimeout:  config.C.Server.ReadTimeout,
	}

	go serve(srv)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), config.C.Server.ShutdownWait)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		panic(err)
	}
}

func serve(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil {
		log.Logger.Fatalf("Failed to serve api: %v", err)
	}
}
