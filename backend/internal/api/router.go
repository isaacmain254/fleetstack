package api

import (
	"database/sql"
	"log"

	"github.com/gorilla/mux"

	"github.com/isaacmain254/fleetstack/backend/internal/github"
	"github.com/isaacmain254/fleetstack/backend/internal/middleware"
	"github.com/isaacmain254/fleetstack/backend/internal/projects"
)

func NewRouter(db *sql.DB) *mux.Router {

	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()

	githubHandler := github.NewHandler()

	projectsHandler := projects.NewHandler()

	api.HandleFunc("/health", HealthCheckHandler).Methods("GET")
	api.HandleFunc("/redirect", githubHandler.RedirectHandler)
	api.HandleFunc("/callback", githubHandler.CallbackHandler)
	api.HandleFunc("/clone", githubHandler.CloneRepoHandler).Methods("POST")
	api.HandleFunc("/build", projectsHandler.BuildProjectHandler).Methods("POST")

	r.Use(middleware.Logging)

	log.Println("Starting server on :8080")

	return r

}
