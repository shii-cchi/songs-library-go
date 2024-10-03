package repository

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"songs-library-go/internal/config"
	"time"
)

const (
	successfulConnectionToDb = "successfully connected to db"
	errConnectingToDb        = "error connecting to db"
	mesReconnectingToDB      = "reconnecting to db"
	errRunningMigrations     = "error running migration"
	successfulRunMigrations  = "successfully executed migrations"
)

const (
	pingInterval = 10 * time.Second
	countTries   = 10
)

func Init(cfg *config.Config) *sql.DB {
	conn, err := sql.Open("postgres", fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName))
	if err != nil {
		log.WithError(err).Fatal(errConnectingToDb)
	}

	if err := conn.Ping(); err != nil {
		log.WithError(err).Fatal(errConnectingToDb)
	}

	go pingDatabase(conn)

	log.Info(successfulConnectionToDb)

	runMigrations(cfg)

	return conn
}

func runMigrations(cfg *config.Config) {
	cmd := exec.Command("goose", "postgres", fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName), "up")
	cmd.Dir = "./internal/repository/migrations"

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.WithError(fmt.Errorf("%s: %s", err, output)).Fatalf(errRunningMigrations)
	}

	log.Info(successfulRunMigrations)
}

func pingDatabase(conn *sql.DB) {
	count := 0
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := conn.Ping()
			if err != nil {
				count++

				if count == countTries {
					log.WithError(err).Fatal(errConnectingToDb)
				}

				log.WithError(err).Info(mesReconnectingToDB)
			}
		}
	}
}
