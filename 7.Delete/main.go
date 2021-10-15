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

	CreatedAt time.Time `json:"created_at" pg:"default:now()"`
	UpdatedAt time.Time `json:"updated_at"`
}

func main() {

	ctx := context.Background()

	pool, dockerResource, db := tools.InitPostgres(ctx, "file://./migrations/7.Delete")

	var prdsInsert []Product

	prdsInsert = append(prdsInsert, Product{
		Name:  "chair",
		Type:  "furniture",
		Price: 5500,
	})
	prdsInsert = append(prdsInsert, Product{
		Name:  "table",
		Type:  "furniture",
		Price: 6200,
	})
	prdsInsert = append(prdsInsert, Product{
		Name:  "cupboard",
		Type:  "furniture",
		Price: 6500,
	})

	_, err := db.Model(&prdsInsert).Insert()
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
	fmt.Println("\n\n")

	var chair Product
	chairString := "chair"
	err = db.Model(&chair).Where("name = ?", chairString).Select()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("chairSelect: \n",
		"id: ", chair.Id,
		"\n name: ", chair.Name,
		"\n type: ", chair.Type,
		"\n price: ", chair.Price,
		"\n updated_at: ", chair.UpdatedAt,
		"\n created_at:", chair.CreatedAt,"\n\n")


	_, err = db.Model(&chair).WherePK().Delete()
	if err != nil {
		fmt.Println(err)
		return
	}

	var prdsSelectAfterDeleteChair []Product
	err = db.Model(&prdsSelectAfterDeleteChair).Where("type = ?", furnitureString).Select()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i, prd := range prdsSelectAfterDeleteChair {
		fmt.Println("prdSelectAfterDeleteChair", i , ": ",
			"\n id: ", prd.Id,
			"\n name: ", prd.Name,
			"\n type: ", prd.Type,
			"\n price: ", prd.Price,
			"\n updated_at: ", prd.UpdatedAt,
			"\n created_at:", prd.CreatedAt)
	}
	fmt.Println("\n\n")

	_, err = db.Model(&Product{}).Where("name = 'table'").Delete()
	if err != nil {
		fmt.Println(err)
		return
	}

	var prdsSelectAfterDeleteTable []Product
	err = db.Model(&prdsSelectAfterDeleteTable).Where("type = ?", furnitureString).Select()
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, prd := range prdsSelectAfterDeleteTable {
		fmt.Println("prdSelectAfterDelete", i , ": ",
			"\n id: ", prd.Id,
			"\n name: ", prd.Name,
			"\n type: ", prd.Type,
			"\n price: ", prd.Price,
			"\n updated_at: ", prd.UpdatedAt,
			"\n created_at:", prd.CreatedAt)
	}
	fmt.Println("\n\n")

	if err := pool.Purge(dockerResource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

}
