package main

import (
	"context"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"log"
	"presentation_go_pg/tools"
	"time"
)

type Product struct {
	tableName struct{} `pg:"products,alias:prd"`

	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Type  string    `json:"type"`
	Price int       `json:"price"`

	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at" pg:"default:now()"`
}

func main() {

	ctx := context.Background()

	pool, dockerResource, db := tools.InitPostgres(ctx, "file://./migrations/2.InsertAndGet")

	prdInsert := Product{
		Name:  "chair",
		Type:  "furniture",
		Price: 5500,
	}
	_, err := db.Model(&prdInsert).Insert()
	if err != nil {
		fmt.Println(err)
		return
	}

	chairString := "chair"
	var prdSelect Product
	err = db.Model(&prdSelect).Where("name = ?", chairString).Select()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("prdSelect: \n",
		"id: ", prdSelect.Id,
		"\n name: ", prdSelect.Name,
		"\n type: ", prdSelect.Type,
		"\n price: ", prdSelect.Price,
		"\n updated_at: ", prdSelect.UpdatedAt,
		"\n created_at:", prdSelect.CreatedAt)

	if err := pool.Purge(dockerResource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

}
