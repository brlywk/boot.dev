package main

import (
	mymw "brlywk/bootdev/webserver/middleware"
	"brlywk/bootdev/webserver/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

// ----- Config ----------------------------------

var dbPath = "database.json"

var fileRoutes = map[string]string{
	"/app": "./app",
}

var metricsConfig = routes.MetricsConfig{
	FileserverHits: 0,
}

// ----- Main ------------------------------------

func main() {
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(mymw.Cors)

	port := 8080

	urlRoot := "/app"
	fsRoot := http.FileServer(http.Dir(fileRoutes[urlRoot]))
	fileHandler := http.StripPrefix(urlRoot, fsRoot)

	// File Routing
	// Route:	/app
	r.Handle(urlRoot, metricsConfig.MiddlewareMetricsInc(fileHandler))
	// Route:	/app/*
	r.Handle(fmt.Sprintf("%v/*", urlRoot), metricsConfig.MiddlewareMetricsInc(fileHandler))

	// Route:	/api/
	apiRouter := routes.CreateApiRouter(dbPath)
	// Route:	/admin/
	adminRouter := routes.CreateAdminRouter(metricsConfig)

	// Mount all the routerz!
	r.Mount("/api", apiRouter)
	r.Mount("/admin", adminRouter)

	// Run server
	log.Println("Server listening on port 8080")
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), r); err != nil {
		log.Fatal("Unable to run server")
	}
}
