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
	Manufacturer   []Manufacturer `json:"manufacturer" pg:"rel:has-many"`

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

	Owner []Owner `json:"owner" pg:"rel:has-many,fk:id,join_fk:manufacturer_id"`
}

func main() {

	ctx := context.Background()

	pool, dockerResource, db := tools.InitPostgres(ctx, "file://./migrations/10.RelationHasManyAnyColumn")

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

	var ownsInsert []Owner

	ownsInsert = append(ownsInsert, Owner{
		ManufacturerId: mft.Id,
		Name:           "Peter",
		Phone:          "+38099000000",
		Email:          "fakemail@mail.fake",
	})
	ownsInsert = append(ownsInsert, Owner{
		ManufacturerId: mft.Id,
		Name:           "German",
		Phone:          "+38095000000",
		Email:          "fakeailGerman@mail.fake",
	})

	_, err = db.Model(&ownsInsert).Insert()
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
		"\n  relation: \n")

	for i, own := range mftSelect.Owner {
		fmt.Println("owner", i , ": ",
			"\n id: ", own.Id,
			"\n name: ", own.Name,
			"\n type: ", own.Phone,
			"\n price: ", own.Email,
			"\n updated_at: ", own.UpdatedAt,
			"\n created_at:", own.CreatedAt)
	}
	if err := pool.Purge(dockerResource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

}
