package database

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CreateMigration(name string) error {
	name = strings.ReplaceAll(name, " ", "_")

	timestamp := time.Now().Format("20060102150405")

	up := filepath.Join(
		"migrations",
		fmt.Sprintf("%s_%s.up.sql", timestamp, name),
	)

	down := filepath.Join(
		"migrations",
		fmt.Sprintf("%s_%s.down.sql", timestamp, name),
	)

	if err := os.WriteFile(up, []byte("-- UP\n"), 0644); err != nil {
		return err
	}

	if err := os.WriteFile(down, []byte("-- DOWN\n"), 0644); err != nil {
		return err
	}

	fmt.Println("Created:")
	fmt.Println(up)
	fmt.Println(down)

	return nil
}