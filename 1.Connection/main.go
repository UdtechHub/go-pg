package main

import (
	"context"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"presentation_go_pg/tools"
)

func main() {

	ctx := context.Background()

	pool, dockerResource, db := tools.InitPostgres(ctx,"file://./migrations/1.Connection")

	fmt.Println(db.Ping(ctx))

	if err := pool.Purge(dockerResource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

}