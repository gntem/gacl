package main

import (
	"gacl/middlewares"
	"gacl/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Env var
type Env struct {
	db gorm.DB
}

func main() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgresql dbname=gacl password='' sslmode=disable")
	db.LogMode(true)

	if err != nil {
		log.Fatal("Unable to connect to db", err.Error())
		panic("Failed connecting to db")
	}

	defer db.Close()

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.DisableFavicon())

	// Group
	router.GET("/group/:groupID", routes.GroupGetByID)
	/*
		router.GET("/groups", routes.GroupGetAll)
		router.POST("/group", routes.GroupCreate)
		router.DELETE("/group/:groupID", routes.GroupDeleteByID)
		router.PUT("/group/:groupID", routes.GroupUpdateByID)
		// router.POST("/group/users", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) })
		// router.DELETE("/group/users", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) })

		// User
		router.GET("/user/:userID", routes.UserGetByID)
		router.GET("/users", routes.UserGetByID)
		router.POST("/user", routes.UserCreate)
		router.DELETE("/user/:userID", routes.UserDeleteByID)
		router.PUT("/user/:userID", routes.UserUpdateByID)
		router.PUT("/user/:userID/permissions/grant", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) })
		router.PUT("/user/:userID/permissions/revoke", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) })

		// Permission
		router.GET("/permission/:permissionID", routes.PermissionGetByID)
		router.GET("/permissions", routes.PermissionGetAll)
		router.DELETE("/permission/:permissionID", routes.PermissiongDeleteById)
		router.PUT("/permission/:permissionID", routes.PermissionUpdateById)

	*/
	router.Run()
}
