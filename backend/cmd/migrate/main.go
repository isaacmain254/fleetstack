package main

import (
	"fmt"
	"os"

	"github.com/isaacmain254/fleetstack/backend/internal/config"
	"github.com/isaacmain254/fleetstack/backend/internal/database"
)

func main() {
	conf := config.Load()
	if len(os.Args) < 2 {
		fmt.Println("usage:")
		fmt.Println("create <migration_name>")
		fmt.Println("up")
		return
	}

	switch os.Args[1] {

	case "create":

		if len(os.Args) < 3 {
			fmt.Println("migration name required")
			return
		}

		err := database.CreateMigration(os.Args[2])
		if err != nil {
			panic(err)
		}

	case "up":

		err := database.RunMigrations(conf.DatabaseURL.URL)
		if err != nil {
			panic(err)
		}

	case "down":

		err := database.RunDownMigrations(conf.DatabaseURL.URL)
		if err != nil {
			panic(err)
		}

	default:
		fmt.Println("unknown command")
	}
}
