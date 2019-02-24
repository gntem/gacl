package main

import (
	"gacl/models"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Automigrate using database connection.
func Automigrate(db *gorm.DB) {
	db.DropTableIfExists(&models.User{}, &models.Group{}, &models.Permission{})
	db.AutoMigrate(&models.User{}, &models.Group{}, &models.Permission{})
}

func main() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgresql dbname=gacl password='' sslmode=disable")
	db.LogMode(true)
	if err != nil {
		log.Fatal("Unable to connect to db", err.Error())
		panic("Failed connecting to db")
	}

	defer db.Close()

	Automigrate(db)
}
