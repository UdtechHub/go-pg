package main

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
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

	CreatedAt time.Time `json:"created_at" pg:"default:now()"`
	UpdatedAt time.Time `json:"updated_at"`
}

func main() {

	ctx := context.Background()

	pool, dockerResource, db := tools.InitPostgres(ctx, "file://./migrations/3.MultiInsertAndGet")

	var prdInsert []Product

	prdInsert = append(prdInsert, Product{
		Name:  "chair",
		Type:  "furniture",
		Price: 5500,
	})
	prdInsert = append(prdInsert, Product{
		Name:  "table",
		Type:  "furniture",
		Price: 6200,
	})
	prdInsert = append(prdInsert, Product{
		Name:  "cupboard",
		Type:  "furniture",
		Price: 6500,
	})

	_, err := db.Model(&prdInsert).Insert()
	if err != nil {
		fmt.Println(err)
		return
	}

	furnitureString := "furniture"
	var prdsSelect []Product
	err = db.Model(&prdsSelect).Where("type = ?", furnitureString).Select()
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, prd := range prdsSelect {
		fmt.Println("prdSelect", i , ": ",
			"\n id: ", prd.Id,
			"\n name: ", prd.Name,
			"\n type: ", prd.Type,
			"\n price: ", prd.Price,
			"\n updated_at: ", prd.UpdatedAt,
			"\n created_at:", prd.CreatedAt)
	}


	err = db.Model(&prdsSelect).Where("name IN (?)", pg.In([]string{"chair","table","cupboard"})).Select()
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, prd := range prdsSelect {
		fmt.Println("prdSelectByArray", i , ": ",
			"\n id: ", prd.Id,
			"\n name: ", prd.Name,
			"\n type: ", prd.Type,
			"\n price: ", prd.Price,
			"\n updated_at: ", prd.UpdatedAt,
			"\n created_at:", prd.CreatedAt)
	}

	if err := pool.Purge(dockerResource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

}
