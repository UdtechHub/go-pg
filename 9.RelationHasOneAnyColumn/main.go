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

type Owner struct {
	tableName struct{} `pg:"owners,alias:own"`

	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Phone string    `json:"phone"`
	Email string    `json:"email"`

	ManufacturerId uuid.UUID     `json:"manufacturer_id"`
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

	Owner *Owner `json:"owner" pg:"rel:has-one,fk:id,join_fk:manufacturer_id"`
}

func main() {

	ctx := context.Background()

	pool, dockerResource, db := tools.InitPostgres(ctx, "file://./migrations/9.RelationHasOneAnyColumn")

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
		"\n created_at:", mft.CreatedAt, "\n\n")

	ownInsert := Owner{
		ManufacturerId: mft.Id,
		Name:           "Peter",
		Phone:          "+38099000000",
		Email:          "fakemail@mail.fake",
	}
	_, err = db.Model(&ownInsert).Insert()
	if err != nil {
		fmt.Println(err)
		return
	}

	var mftSelect Manufacturer
	err = db.Model(&mftSelect).Where("mft.name = ?", mftInsert.Name).Relation("Owner").Select()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("mftSelect: \n",
		"id: ", mftSelect.Id,
		"\n name: ", mftSelect.Name,
		"\n phone: ", mftSelect.Phone,
		"\n email: ", mftSelect.Email,
		"\n updated_at: ", mftSelect.UpdatedAt,
		"\n created_at:", mftSelect.CreatedAt,
		"\n  relation: \n",
		" Owner: \n",
		"id: ", mftSelect.Owner.Id,
		"\n name: ", mftSelect.Owner.Name,
		"\n type: ", mftSelect.Owner.Phone,
		"\n price: ", mftSelect.Owner.Email,
		"\n updated_at: ", mftSelect.Owner.UpdatedAt,
		"\n created_at:", mftSelect.Owner.CreatedAt)

	if err := pool.Purge(dockerResource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

}
