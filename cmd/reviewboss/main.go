package main

import (
	"log"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	gitlab "github.com/xanzy/go-gitlab"
)

func main() {
	cfg := loadConfig()

	git := gitlab.NewClient(nil, cfg.GitlabToken)
	err := git.SetBaseURL(cfg.GitlabBaseURL)
	if err != nil {
		log.Fatalf("Can not set gitlab base url: %v", err)
	}

	db, err := sql.Open("sqlite3", cfg.DbPath)
	if err != nil {
		log.Fatalf("Can not open sqlite3 database with path %q: %v", cfg.DbPath, err)
	}
	defer db.Close()

	mengine := newMigrationEngine(db)
	err = mengine.run()
	if err != nil {
		log.Fatalf("Can not apply migrations: %v", err)
	}
}
