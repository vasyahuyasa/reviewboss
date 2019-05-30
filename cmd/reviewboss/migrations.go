package main

import (
	"database/sql"
	"fmt"
)

type migrateEngine struct {
	db *sql.DB
}

type migration struct {
	name string
	up   string
	down string
}

var migrations = []migration{}

func newMigrationEngine(db *sql.DB) *migrateEngine {
	return &migrateEngine{
		db: db,
	}
}

func (e *migrateEngine) applyMigration(m migration) error {
	_, err := e.db.Exec(m.up)
	return err
}

func (e *migrateEngine) maxBatch() (int, error) {
	max := 0
	row := e.db.QueryRow(`SELECT max(batch) from migrations`)
	err := row.Scan(&max)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	return max, nil
}

// get names of applied migrations
func (e *migrateEngine) appliedMigrations() (map[string]struct{}, error) {
	rows, err := e.db.Query(`SELECT migration FROM migrations`)
	if err != nil {
		return nil, err
	}

	list := map[string]struct{}{}
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		list[name] = struct{}{}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (e *migrateEngine) registerMigration(name string, batchId int) error {
	_, err := e.db.Exec(`INSERT INTO migrations (migration, batch) VALUES (?,?)`, name, batchId)
	return err
}

func (e *migrateEngine) run() error {
	// create migration table if does not exists
	m := migration{
		up: `CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			migration varchar(255) NOT NULL,
			batch int(11) NOT NULL
		)`,
	}
	err := e.applyMigration(m)
	if err != nil {
		return fmt.Errorf("can not create migration table: %v", err)
	}

	// get list of applied migrations and find not applied our migrations
	toApply := []migration{}
	applied, err := e.appliedMigrations()
	if err != nil {
		return fmt.Errorf("can not get list of applied migrations: %e", err)
	}

	for _, m := range migrations {
		_, ok := applied[m.name]
		if !ok {
			toApply = append(toApply, m)
		}
	}

	if len(toApply) == 0 {
		return nil
	}

	// find last migration batch id and calculate next id
	currentBatch, err := e.maxBatch()
	if err != nil {
		return fmt.Errorf("can not get max batch: %v", err)
	}

	nextBatch := currentBatch + 1
	for _, m := range toApply {
		err = e.applyMigration(m)
		if err != nil {
			return fmt.Errorf("can not apply migration %q: %v", m.name, err)
		}

		err = e.registerMigration(m.name, nextBatch)
		if err != nil {
			return fmt.Errorf("can not register migration %q: %v", m.name, err)
		}
	}

	return nil
}
