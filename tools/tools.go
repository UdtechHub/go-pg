package tools

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/golang-migrate/migrate/v4"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"time"
)

func InitPostgres(ctx context.Context, migrationPath string) (pool *dockertest.Pool, dockerResource *dockertest.Resource, db *pg.DB) {

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	dockerResource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_PASSWORD=12345",
			"POSTGRES_USER=postgres",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := dockerResource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://postgres:12345@%s/dbname?sslmode=disable", hostAndPort)

	dockerResource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second

	if err = pool.Retry(func() error {

		opt, err := pg.ParseURL(databaseUrl)
		if err != nil {
			panic(err)
		}

		db = pg.Connect(opt)

		//pg.Connect(&pg.Options{
		//		Addr:     viper.GetString("db.host") + viper.GetString("db.port"),
		//		User:     viper.GetString("db.user"),
		//		Password: viper.GetString("db.pass"),
		//		Database: viper.GetString("db.name"),
		//	})

		return db.Ping(ctx)

	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// "file://../src/infrastructure/storage/postgres/migrations"
	migrationInstance, err := migrate.New(migrationPath, databaseUrl)
	if err != nil {
		panic(err)
	}

	err = migrationInstance.Up()
	if err != nil {
		panic(err)
	}

	return pool, dockerResource, db

}
