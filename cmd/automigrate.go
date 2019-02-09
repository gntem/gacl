package cmd

import (
	"gacl/models"

	"github.com/jinzhu/gorm"
)

// Automigrate
func Automigrate(db *gorm.DB) {
	db.DropTableIfExists(&models.User{}, &models.Group{}, &models.Permission{})
	db.AutoMigrate(&models.User{}, &models.Group{}, &models.Permission{})
}

func main() {
	Automigrate()
}
