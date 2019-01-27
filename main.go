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
	ID   int64
	Name string `gorm:"type:varchar(255);unique;not null"`
}

// Pagination struct
type Pagination struct {
	Page   uint64 `validate:"gte=0"`
	Offset uint64 `validate:"gte=0"`
	Limit  uint64 `validate:"gte=0"`
	SortBy string `validate:"oneof=created_at id"`
	Order  string `validate:"oneof=desc asc"`
}

// UserCreateRequest struct
type UserCreateRequest struct {
	Name string `form:"name" validate:"min=4,max=255" binding:"required"`
}

// GroupCreateRequest struct
type GroupCreateRequest struct {
	Name string `form:"name" validate:"min=4,max=255" binding:"required"`
}

// UserUpdateRequest struct
type UserUpdateRequest struct {
	Name string `form:"name" validate:"min=4,max=255"`
}

// GroupUpdateRequest struct
type GroupUpdateRequest struct {
	Name string `form:"name" validate:"min=4,max=255"`
}

// PermissionUpdateRequest struct
type PermissionUpdateRequest struct {
	Name string `form:"name" validate:"min=4,max=255"`
}

// Authorization middleware, authorize using vault.
func authorizationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}

func disableFavicon() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/favicon.ico" {
			return
		}
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
	router.Use(disableFavicon())

	// Group
	// with=['users','permissions']
	router.GET("/group/:groupID", func(ctx *gin.Context) {
		groupIDArg := ctx.Param("groupID")

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

		if dbError.Error != nil {
			panic(dbError.Error)
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

		if dbError.Error != nil {
			panic(dbError.Error)
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

		if dbError.Error != nil {
			panic(dbError.Error)
		}

		ctx.JSON(http.StatusCreated, gin.H{"error": nil, "result": rgroup})
	})
	router.DELETE("/group/:groupID", func(ctx *gin.Context) {
		groupIDArg := ctx.Param("groupID")

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

		if dbError.Error != nil {
			panic(dbError.Error)
		}

		ctx.JSON(http.StatusOK, gin.H{"error": nil, "result": nil})
	})
	router.PUT("/group/:groupID", func(ctx *gin.Context) {
		groupIDArg := ctx.Param("groupID")

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

		if dbError.Error != nil {
			panic(dbError.Error)
		}

		ctx.JSON(http.StatusOK, gin.H{"error": nil, "result": nil})

	})
	// router.POST("/group/users", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) })
	// router.DELETE("/group/users", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) })

	// User
	router.GET("/user/:userID", func(ctx *gin.Context) {
		userIDArg := ctx.Param("userID")

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

		if dbError.Error != nil {
			panic(dbError.Error)
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

		if dbError.Error != nil {
			panic(dbError.Error)
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

		if dbError.Error != nil {
			panic(dbError.Error)
		}

		ctx.JSON(http.StatusCreated, gin.H{"error": nil, "result": ruser})
	})

	router.DELETE("/user/:userID", func(ctx *gin.Context) {
		userIDArg := ctx.Param("userID")

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

		if dbError.Error != nil {
			panic(dbError.Error)
		}

		ctx.JSON(http.StatusOK, gin.H{"error": nil, "result": nil})

	})
	router.PUT("/user/:userID", func(ctx *gin.Context) {
		userIDArg := ctx.Param("userID")

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

		if dbError.Error != nil {
			panic(dbError.Error)
		}

		ctx.JSON(http.StatusOK, gin.H{"error": nil, "result": nil})

	})
	router.PUT("/user/:userID/permissions/grant", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) })
	router.PUT("/user/:userID/permissions/revoke", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) })

	// Permission
	router.GET("/permission/:permissionID", func(ctx *gin.Context) {
		permissionIDArg := ctx.Param("permissionID")

		var result []Permission

		if len(permissionIDArg) == 0 {
			ctx.JSON(http.StatusBadRequest,
				gin.H{"error": GACLAPIError{Message: "Missing :permissionID query param value"},
					"result": nil})
		}

		permissionID, parseError := strconv.ParseInt(permissionIDArg, 8, 64)

		if parseError != nil {
			panic(parseError)
		}

		dbError := db.First(&result, permissionID)

		if dbError.Error != nil {
			panic(dbError.Error)
		}
		ctx.JSON(http.StatusOK, gin.H{"error": nil, "result": result})
	})
	router.GET("/permissions", func(ctx *gin.Context) {
		var result []Permission
		var rp Pagination

		if bindingError := ctx.ShouldBindQuery(&rp); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": bindingError.Error()})
			return
		}

		sortBy := fmt.Sprintf("%s %s", rp.SortBy, rp.Order)

		dbError := db.Limit(rp.Limit).Offset(rp.Offset).Order(sortBy).Find(&result)

		if dbError.Error != nil {
			panic(dbError.Error)
		}

		ctx.JSON(http.StatusOK, gin.H{"error": nil, "result": result})
	})
	router.DELETE("/permission/:permissionID", func(ctx *gin.Context) {
		permissionIDArg := ctx.Param("permissionID")

		if len(permissionIDArg) == 0 {
			ctx.JSON(http.StatusBadRequest,
				gin.H{"error": GACLAPIError{Message: "Missing :permissionID query param value"},
					"result": nil})
		}

		permissionID, parseError := strconv.ParseInt(permissionIDArg, 8, 64)

		if parseError != nil {
			panic(parseError)
		}

		dbError := db.Delete(&Permission{ID: permissionID})

		if dbError.Error != nil {
			panic(dbError.Error)
		}

		ctx.JSON(http.StatusOK, gin.H{"error": nil, "result": nil})

	})
	router.PUT("/permission/:permissionID", func(ctx *gin.Context) {
		permissionIDArg := ctx.Param("permissionID")

		if len(permissionIDArg) == 0 {
			ctx.JSON(http.StatusBadRequest,
				gin.H{"error": GACLAPIError{Message: "Missing :permissionID query param value"},
					"result": nil})
			return
		}

		var rPermissionUpdate PermissionUpdateRequest

		if bindingError := ctx.ShouldBindJSON(&rPermissionUpdate); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": bindingError.Error()})
			return
		}

		permissionID, parseError := strconv.ParseInt(permissionIDArg, 8, 64)

		if parseError != nil {
			ctx.JSON(http.StatusBadRequest,
				gin.H{"error": GACLAPIError{Message: "Internal server error"},
					"result": nil})
			panic(parseError)
		}

		dbError := db.Model(&Permission{ID: permissionID}).Updates(&Permission{Name: rPermissionUpdate.Name})

		if dbError.Error != nil {
			panic(dbError.Error)
		}

		ctx.JSON(http.StatusOK, gin.H{"error": nil, "result": nil})

	})

	router.Run()
}
