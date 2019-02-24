package routes

import (
	"fmt"
	"gacl/models"
	"gacl/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// UserGetByID get user by id
func UserGetByID(ctx *gin.Context) {
	db := ctx.MustGet("database").(*gorm.DB)
	userIDArg := ctx.Param("userID")

	var result []models.User

	if len(userIDArg) == 0 {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": structs.GACLAPIError{Message: "Missing :userID query param value"},
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

}

// UsersGetAll get all user rows with default pagination
func UsersGetAll(ctx *gin.Context) {
	db := ctx.MustGet("database").(*gorm.DB)
	var result []models.User
	var rp structs.PaginationQuery

	if err := ctx.ShouldBindQuery(&rp); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sortBy := fmt.Sprintf("%s %s", rp.SortBy, rp.Order)

	dbError := db.Limit(rp.Limit).Offset(rp.Page).Order(sortBy).Find(&result)

	if dbError.Error != nil {
		panic(dbError.Error)
	}

	ctx.JSON(http.StatusOK, gin.H{"error": nil, "result": result})
}

// UserCreate create a user
func UserCreate(ctx *gin.Context) {
	db := ctx.MustGet("database").(*gorm.DB)
	var ruser []structs.UserCreateRequest

	if err := ctx.ShouldBindJSON(&ruser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	trx := db.Begin()
	dbError := trx.Create(&ruser)

	if dbError.Error != nil {
		trx.Rollback()
		panic(dbError.Error)
	}

	trx.Commit()

	ctx.JSON(http.StatusCreated, gin.H{"error": nil, "result": ruser})
}

// UserDeleteByID delete user row by id
func UserDeleteByID(ctx *gin.Context) {
	db := ctx.MustGet("database").(*gorm.DB)
	userIDArg := ctx.Param("userID")

	if len(userIDArg) == 0 {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": structs.GACLAPIError{Message: "Missing :userID query param value"},
				"result": nil})
	}

	userID, parseError := strconv.ParseInt(userIDArg, 8, 64)

	if parseError != nil {
		panic(parseError)
	}

	trx := db.Begin()

	dbError := trx.Delete(&models.User{ID: userID})

	if dbError.Error != nil {
		trx.Rollback()
		panic(dbError.Error)
	}

	trx.Commit()

	ctx.JSON(http.StatusOK, gin.H{"error": nil, "result": nil})
}

// UserUpdateByID update user by id
func UserUpdateByID(ctx *gin.Context) {
	db := ctx.MustGet("database").(*gorm.DB)
	userIDArg := ctx.Param("userID")

	if len(userIDArg) == 0 {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": structs.GACLAPIError{Message: "Missing :userID query param value"},
				"result": nil})
		return
	}

	var rUserUpdate structs.UserUpdateRequest

	if err := ctx.ShouldBindJSON(&rUserUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, parseError := strconv.ParseInt(userIDArg, 8, 64)

	if parseError != nil {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": structs.GACLAPIError{Message: "Internal server error"},
				"result": nil})
		panic(parseError)
	}
	trx := db.Begin()

	dbError := trx.Model(&models.User{ID: userID}).Updates(&models.User{Name: rUserUpdate.Name})

	if dbError.Error != nil {
		trx.Rollback()
		panic(dbError.Error)
	}

	trx.Commit()

	ctx.JSON(http.StatusOK, gin.H{"error": nil, "result": nil})

}
