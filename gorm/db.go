package gorm

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID          string
	Price       uint
	Title       string
	Description string
}

func InitDB(opts string) {
	db, err := gorm.Open(postgres.Open(opts), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Product{})

	db.Create(&Product{Price: 100, Description: "something", Title: "anything"})

	var product Product

	db.First(&product, 1)

	db.First(&product, "price = ?", 100)

	db.Model(&product).Update("Price", 200)

	db.Delete(&product, 1)

}
