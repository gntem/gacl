package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Env locals
type Env struct {
	db *gorm.DB
}

func main() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=gacl password='' sslmode=disable")
	db.LogMode(true)

	if err != nil {
		log.Fatal("Unable to connect to db", err.Error())
		panic("Failed connecting to db")
	}

	defer db.Close()

	db.DropTableIfExists(&User{}, &Group{}, &Permission{})
	db.AutoMigrate(&User{}, &Group{}, &Permission{})

	envCtx := &Env{db: db}

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Group
	// with=['users','permissions']
	router.GET("/group/:groupId", envCtx.getGroup)
	router.GET("/group/all", envCtx.getAllGroups)
	router.POST("/group", envCtx.createGroup)
	router.DELETE("/group/:groupId", envCtx.removeGroup)
	router.PATCH("/group/:groupId", envCtx.updateGroup)
	router.PUT("/group/:groupId?", envCtx.upsertGroup)
	router.POST("/group/:groupId/user/add", envCtx.addUserToGroup)
	router.DELETE("/group/:groupId/user/:userId", envCtx.removeUserFromGroup)

	// User
	router.GET("/user/all", envCtx.getAllUsers)
	// ?with=['permissions','groups']
	router.GET("/user/:userId", envCtx.getUser)
	router.POST("/user", envCtx.createUser)
	router.DELETE("/user/:userId", envCtx.removeUser)
	router.PATCH("/user/:userId", envCtx.updateUser)
	router.PUT("/user/:userId", envCtx.upsertUser)
	router.PUT("/user/:userId/permissions/grant", envCtx.grantPermission)
	router.PUT("/user/:userId/permissions/revoke", envCtx.revokePermission)

	// Permission
	router.GET("/permission", envCtx.getPermissions)
	router.POST("/permission", envCtx.postPermission)
	router.DELETE("/permission/:id", envCtx.deletePermission)
	router.PATCH("/permission/:id", envCtx.updatePermission)

	router.Run()
}
