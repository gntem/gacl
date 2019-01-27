package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// GACLAPIError standard error
type GACLAPIError struct {
	Message string
}

// User model
type User struct {
	gorm.Model
	ID   int64
	Name string `gorm:"type:varchar(255);unique;not null"`
}

// Group model
type Group struct {
	gorm.Model
	ID          int64
	Name        string        `gorm:"type:varchar(255);unique;not null"`
	Permissions []*Permission `gorm:"many2many:group_permissions;"`
	Users       []*User       `gorm:"many2many:group_users;"`
}

// Permission model
type Permission struct {
	gorm.Model
	Name string `gorm:"type:varchar(255);unique;not null"`
}

type Pagination struct {
	Page   uint64 `validate:"gte=0"`
	Offset uint64 `validate:"gte=0"`
	Limit  uint64 `validate:"gte=0"`
	SortBy string `validate:"oneof=created_at id"`
	Order  string `validate:"oneof=desc asc"`
}

type UserCreateRequest struct {
	Name string `form:"name" validate:"min=4,max=255" binding:"required"`
}

type GroupCreateRequest struct {
	Name string `form:"name" validate:"min=4,max=255" binding:"required"`
}

type UserUpdateRequest struct {
	Name string `form:"name" validate:"min=4,max=255"`
}

type GroupUpdateRequest struct {
	Name string `form:"name" validate:"min=4,max=255"`
}

func authorizationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}

func main() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=gacl password='' sslmode=disable")
	db.LogMode(true)

	if err != nil {
		log.Fatal("Unable to connect to db", err.Error())
		panic("Failed connecting to db")
	}

	defer db.Close()

	//db.DropTableIfExists(&User{}, &Group{}, &Permission{})
	//db.AutoMigrate(&User{}, &Group{}, &Permission{})

	/*
		usera := User{Name: "usera"}
		userb := User{Name: "userb"}

		db.Create(usera)
		db.Create(userb)
	*/

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Group
	// with=['users','permissions']
	router.GET("/group/:groupID", func(ctx *gin.Context) {
		groupIDArg := ctx.Query("groupID")

		var result []Group

		if len(groupIDArg) == 0 {
			ctx.JSON(http.StatusBadRequest,
				gin.H{"error": GACLAPIError{Message: "Missing :groupID query param value"},
					"result": nil})
		}

		groupID, parseError := strconv.ParseInt(groupIDArg, 8, 64)

		if parseError != nil {
			panic(parseError)
		}

		dbError := db.First(&result, groupID)

		if dbError != nil {
			panic(dbError)
		}
		ctx.JSON(http.StatusOK, gin.H{"error": nil, "result": result})
	})
	router.GET("/groups", func(ctx *gin.Context) {
		var result []Group
		var rp Pagination

		if bindingError := ctx.ShouldBindQuery(&rp); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": bindingError.Error()})
			return
		}

		sortBy := fmt.Sprintf("%s %s", rp.SortBy, rp.Order)

		dbError := db.Limit(rp.Limit).Offset(rp.Offset).Order(sortBy).Find(&result)

		if dbError != nil {
			panic(dbError)
		}

		ctx.JSON(http.StatusOK, gin.H{"error": nil, "result": result})
	})
	router.POST("/group", func(ctx *gin.Context) {
		var rgroup []GroupCreateRequest

		if bindingError := ctx.ShouldBindJSON(&rgroup); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": bindingError.Error()})
			return
		}

		dbError := db.Create(&rgroup)

		if dbError != nil {
			panic(dbError)
		}

		ctx.JSON(http.StatusCreated, gin.H{"error": nil, "result": rgroup})
	})
	router.DELETE("/group/:groupID", func(ctx *gin.Context) {
		groupIDArg := ctx.Query("groupID")

		if len(groupIDArg) == 0 {
			ctx.JSON(http.StatusBadRequest,
				gin.H{"error": GACLAPIError{Message: "Missing :groupID query param value"},
					"result": nil})
		}

		groupID, parseError := strconv.ParseInt(groupIDArg, 8, 64)

		if parseError != nil {
			panic(parseError)
		}

		dbError := db.Delete(&Group{ID: groupID})

		if dbError != nil {
			panic(dbError)
		}

		ctx.JSON(http.StatusOK, gin.H{"error": nil, "result": nil})
	})
	router.PUT("/group/:groupID", func(ctx *gin.Context) {
		groupIDArg := ctx.Query("groupID")

		if len(groupIDArg) == 0 {
			ctx.JSON(http.StatusBadRequest,
				gin.H{"error": GACLAPIError{Message: "Missing :groupID query param value"},
					"result": nil})
			return
		}

		var rgroupUpdate GroupUpdateRequest

		if bindingError := ctx.ShouldBindJSON(&rgroupUpdate); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": bindingError.Error()})
			return
		}

		groupID, parseError := strconv.ParseInt(groupIDArg, 8, 64)

		if parseError != nil {
			ctx.JSON(http.StatusBadRequest,
				gin.H{"error": GACLAPIError{Message: "Internal server error"},
					"result": nil})
			panic(parseError)
		}

		dbError := db.Model(&Group{ID: groupID}).Updates(&Group{Name: rgroupUpdate.Name})

		if dbError != nil {
			panic(dbError)
		}

		ctx.JSON(http.StatusOK, gin.H{"error": nil, "result": nil})

	})
	router.POST("/group/user/add", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) })
	router.DELETE("/group/user/:userId", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) })

	// User
	router.GET("/user/:userID", func(ctx *gin.Context) {
		userIDArg := ctx.Query("userID")

		var result []User

		if len(userIDArg) == 0 {
			ctx.JSON(http.StatusBadRequest,
				gin.H{"error": GACLAPIError{Message: "Missing :userID query param value"},
					"result": nil})
		}

		userID, parseError := strconv.ParseInt(userIDArg, 8, 64)

		if parseError != nil {
			panic(parseError)
		}

		dbError := db.First(&result, userID)

		if dbError != nil {
			panic(dbError)
		}
		ctx.JSON(http.StatusOK, gin.H{"error": nil, "result": result})
	})

	router.GET("/users", func(ctx *gin.Context) {

		var result []User
		var rp Pagination

		if bindingError := ctx.ShouldBindQuery(&rp); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": bindingError.Error()})
			return
		}

		sortBy := fmt.Sprintf("%s %s", rp.SortBy, rp.Order)

		dbError := db.Limit(rp.Limit).Offset(rp.Offset).Order(sortBy).Find(&result)

		if dbError != nil {
			panic(dbError)
		}

		ctx.JSON(http.StatusOK, gin.H{"error": nil, "result": result})
	})

	router.POST("/user", func(ctx *gin.Context) {
		var ruser []UserCreateRequest

		if bindingError := ctx.ShouldBindJSON(&ruser); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": bindingError.Error()})
			return
		}

		dbError := db.Create(&ruser)

		if dbError != nil {
			panic(dbError)
		}

		ctx.JSON(http.StatusCreated, gin.H{"error": nil, "result": ruser})
	})

	router.DELETE("/user/:userID", func(ctx *gin.Context) {
		userIDArg := ctx.Query("userID")

		if len(userIDArg) == 0 {
			ctx.JSON(http.StatusBadRequest,
				gin.H{"error": GACLAPIError{Message: "Missing :userID query param value"},
					"result": nil})
		}

		userID, parseError := strconv.ParseInt(userIDArg, 8, 64)

		if parseError != nil {
			panic(parseError)
		}

		dbError := db.Delete(&User{ID: userID})

		if dbError != nil {
			panic(dbError)
		}

		ctx.JSON(http.StatusOK, gin.H{"error": nil, "result": nil})

	})
	router.PUT("/user/:userID", func(ctx *gin.Context) {
		userIDArg := ctx.Query("userID")

		if len(userIDArg) == 0 {
			ctx.JSON(http.StatusBadRequest,
				gin.H{"error": GACLAPIError{Message: "Missing :userID query param value"},
					"result": nil})
			return
		}

		var rUserUpdate UserUpdateRequest

		if bindingError := ctx.ShouldBindJSON(&rUserUpdate); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": bindingError.Error()})
			return
		}

		userID, parseError := strconv.ParseInt(userIDArg, 8, 64)

		if parseError != nil {
			ctx.JSON(http.StatusBadRequest,
				gin.H{"error": GACLAPIError{Message: "Internal server error"},
					"result": nil})
			panic(parseError)
		}

		dbError := db.Model(&User{ID: userID}).Updates(&User{Name: rUserUpdate.Name})

		if dbError != nil {
			panic(dbError)
		}

		ctx.JSON(http.StatusOK, gin.H{"error": nil, "result": nil})

	})
	router.PUT("/user/:userID/permissions/grant", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) })
	router.PUT("/user/:userID/permissions/revoke", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) })

	// Permission
	router.GET("/permission", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) })
	router.DELETE("/permission/:id", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) })
	router.PUT("/permission/:id", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) })

	router.Run()
}
