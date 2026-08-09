package main

import (
	"log"

	"github.com/isaacmain254/fleetstack/backend/internal/config"
	"github.com/isaacmain254/fleetstack/backend/internal/database"
	"github.com/isaacmain254/fleetstack/backend/internal/projects"
)

func main() {
	conf := config.Load()

	cfg := database.GetDefaultConfig(conf)
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	handler := projects.NewHandler(db)
	log.Println("worker started; waiting for queued deployments")
	handler.ProcessPendingDeployments()
}
