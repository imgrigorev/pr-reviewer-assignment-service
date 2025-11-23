package database

import (
	"context"
	"log"
	"os"

	"pr-reviewer-assigment-service/internal/client/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Addr     string `yaml:"addr"`
	DbName   string `yaml:"db_name"`
	User     string `yaml:"user_name"`
	Password string `yaml:"password"`

	PoolMaxConns string `yaml:"pool_max_conns"`
}

func Connect(c Database) db.Client {
	var err error

	poolConfig, err := pgxpool.ParseConfig("postgres://" + c.User + ":" + c.Password + "@" + c.Addr + "/" + c.DbName + "?pool_max_conns=" + c.PoolMaxConns)
	if err != nil {
		log.Fatal("Unable to parse DATABASE_URL", "error", err)
		os.Exit(1)
	}

	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatal("Unable to create connection pool", "error", err)
		os.Exit(1)
	}

	// Check if is alive
	_, err = db.Exec(context.Background(), "SELECT 1")
	if err != nil {
		log.Fatal("Database Error", err)
	}

	return db.Client{}
}
