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

	ManufacturerId uuid.UUID    `json:"manufacturer_id"`
	Manufacturer   *Manufacturer `json:"manufacturer" pg:"rel:has-one"`

	CreatedAt time.Time `json:"created_at" pg:"default:now()"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Manufacturer struct {
	tableName struct{} `pg:"manufacturers,alias:mft"`

	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Phone string    `json:"phone"`
	Email string    `json:"email"`

	CreatedAt time.Time `json:"created_at" pg:"default:now()"`
	UpdatedAt time.Time `json:"updated_at"`
}

func main() {

	ctx := context.Background()

	pool, dockerResource, db := tools.InitPostgres(ctx, "file://./migrations/8.RelationHasOne")

	mftInsert := Manufacturer{
		Name:  "FakeManufacturer",
		Phone: "+380000000000",
		Email: "fakemail@fake.com",
	}
	_, err := db.Model(&mftInsert).Insert()
	if err != nil {
		fmt.Println(err)
		return
	}

	var mft Manufacturer
	err = db.Model(&mft).Where("name = ?", mftInsert.Name).Select()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Manufacturer: \n",
		"id: ", mft.Id,
		"\n name: ", mft.Name,
		"\n type: ", mft.Phone,
		"\n price: ", mft.Email,
		"\n updated_at: ", mft.UpdatedAt,
		"\n created_at:", mft.CreatedAt,"\n\n")

	prdInsert := Product{
		ManufacturerId: mft.Id,
		Name:  "chair",
		Type:  "furniture",
		Price: 5500,
	}
	_, err = db.Model(&prdInsert).Insert()
	if err != nil {
		fmt.Println(err)
		return
	}

	var prdSelect Product
	err = db.Model(&prdSelect).Where("prd.name = ?", prdInsert.Name).Relation("Manufacturer").Select()
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
		"\n created_at:", prdSelect.CreatedAt,
		"\n  relation: \n",
		" Manufacturer: \n",
		"id: ", prdSelect.Manufacturer.Id,
		"\n name: ", prdSelect.Manufacturer.Name,
		"\n type: ", prdSelect.Manufacturer.Phone,
		"\n price: ", prdSelect.Manufacturer.Email,
		"\n updated_at: ", prdSelect.Manufacturer.UpdatedAt,
		"\n created_at:", prdSelect.Manufacturer.CreatedAt)

	if err := pool.Purge(dockerResource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

}
