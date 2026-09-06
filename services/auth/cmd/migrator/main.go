package main

import (
	"flag"
	"log"

	"github.com/ritchieridanko/apotekly/services/auth/configs"
	"github.com/ritchieridanko/apotekly/services/shared/infra/database"
)

func main() {
	fu := flag.Bool("up", false, "Apply all up migrations")
	fd := flag.Int("down", -1, "Apply N down migrations")
	flag.Parse()

	// Flags Validation
	if *fu && *fd >= 0 {
		log.Fatalln("[FATAL]: failed to apply migrations: -up and -down cannot be used together")
	}
	if !*fu && *fd < 0 {
		log.Fatalln("[FATAL]: failed to apply migrations: specify either -up or -down")
	}

	// Config Initialization
	cfg, err := configs.Init("./configs")
	if err != nil {
		log.Fatalln("[FATAL]:", err)
	}

	// Migrator Initialization
	m, err := database.NewMigrator(&cfg.Database, "./migrations")
	if err != nil {
		log.Fatalln("[FATAL]:", err)
	}
	defer func(m *database.Migrator) {
		if err := m.Close(); err != nil {
			log.Println("[WARN]:", err)
		}
	}(m)

	// DB Migrations Execution
	switch {
	case *fu:
		if err := m.Up(); err != nil {
			log.Fatalln("[FATAL]:", err)
		}
	case *fd >= 0:
		if err := m.Down(*fd); err != nil {
			log.Fatalln("[FATAL]:", err)
		}
	}
}
