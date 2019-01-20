package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// User model
type User struct {
	gorm.Model
	Name string `gorm:"type:varchar(255);unique;not null"`
}

// Group model
type Group struct {
	gorm.Model
	Name        string        `gorm:"type:varchar(255);unique;not null"`
	Permissions []*Permission `gorm:"many2many:group_permissions;"`
	Users       []*User       `gorm:"many2many:group_users;"`
}

// Permission model
type Permission struct {
	gorm.Model
	Name string `gorm:"type:varchar(255);unique;not null"`
}

func main() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=gacl password='' sslmode=disable")
	db.LogMode(true)

	if err != nil {
		log.Fatal("Unable to connect to db", err.Error())
		panic("Failed connecting to db")
	}

	db.DropTableIfExists(&User{}, &Group{}, &Permission{})
	db.AutoMigrate(&User{}, &Group{}, &Permission{})

	defer db.Close()

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/group", getGroups)
	router.POST("/group", postGroup)
	router.DELETE("/group/:id", deleteGroup)
	router.PATCH("/group/:id", updateGroup)

	router.GET("/user", getUsers)
	router.POST("/user", postUser)
	router.DELETE("/user/:id", deleteUser)
	router.PATCH("/user/:id", updateUser)

	router.GET("/permission", getPermissions)
	router.POST("/permission", postPermission)
	router.DELETE("/permission/:id", deletePermission)
	router.PATCH("/permission/:id", updatePermission)

	router.Run()
}
