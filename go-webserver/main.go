package main

import (
	"brlywk/bootdev/webserver/config"
	"brlywk/bootdev/webserver/db"
	mymw "brlywk/bootdev/webserver/middleware"
	"brlywk/bootdev/webserver/routes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
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
	godotenv.Load()

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(mymw.Cors)

	port := 8080

	db, err := db.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Database Error: %v\n", err)
	}

	cfg := config.ApiConfig{
		Db:        db,
		JwtSecret: []byte(os.Getenv("JWT_SECRET")),
		PolkaKey:  os.Getenv("POLKA_KEY"),
		TokenSettings: config.TokenSettings{
			AccessIssuer:     "chirpy-access",
			AccessExpiresIn:  time.Hour * 1,
			RefreshIssuer:    "chirpy-refresh",
			RefreshExpiresIn: time.Hour * 24 * 60,
		},
	}

	urlRoot := "/app"
	fsRoot := http.FileServer(http.Dir(fileRoutes[urlRoot]))
	fileHandler := http.StripPrefix(urlRoot, fsRoot)

	// File Routing
	// Route:	/app
	r.Handle(urlRoot, metricsConfig.MiddlewareMetricsInc(fileHandler))
	// Route:	/app/*
	r.Handle(fmt.Sprintf("%v/*", urlRoot), metricsConfig.MiddlewareMetricsInc(fileHandler))

	// Route:	/api/
	apiRouter := routes.CreateApiRouter(cfg)
	// Route:	/admin/
	adminRouter := routes.CreateAdminRouter(metricsConfig)

	// Mount all the routerz!
	r.Mount("/api", apiRouter)
	r.Mount("/admin", adminRouter)

	// Run server
	log.Println("Server listening on port 8080")
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), r); err != nil {
		log.Fatalf("Unable to run server: %v", err)
	}
}
