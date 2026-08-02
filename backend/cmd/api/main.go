package main

import (
	"log"
	"net/http"
	// "os"

	"github.com/isaacmain254/fleetstack/backend/internal/api"
	"github.com/isaacmain254/fleetstack/backend/internal/config"
	"github.com/isaacmain254/fleetstack/backend/internal/database"
)


func main() {
	conf := config.Load()

	cfg := database.GetDefaultConfig(conf)
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		// os.Exit(1)
	}
	defer db.Close()

	// run migrations after establishing the database connection
	// err = database.RunMigrations(os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	log.Fatalf("failed to run migrations: %v", err)
	// }

	router := api.NewRouter(db)
	log.Fatal(http.ListenAndServe(":8080", router))

}

