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

// GroupGetByID , fetch group by id
func GroupGetByID(ctx *gin.Context) {
	db := ctx.MustGet("database").(*gorm.DB)
	groupIDArg := ctx.Param("groupID")

	var result []models.Group

	if len(groupIDArg) == 0 {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": structs.GACLAPIError{Message: "Missing :groupID query param value"},
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
}

// GroupGetAll , with optional pagination
func GroupGetAll(ctx *gin.Context) {
	db := ctx.MustGet("database").(*gorm.DB)
	var result []models.Group
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

// GroupCreate , Create an group
func GroupCreate(ctx *gin.Context) {
	db := ctx.MustGet("database").(*gorm.DB)
	var rgroup []structs.GroupCreateRequest

	if err := ctx.ShouldBindJSON(&rgroup); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trx := db.Begin()

	dbError := trx.Create(&rgroup)

	if dbError.Error != nil {
		trx.Rollback()
		panic(dbError.Error)
	}

	trx.Commit()

	ctx.JSON(http.StatusCreated, gin.H{"error": nil, "result": rgroup})
}

// GroupDeleteByID delete group by id
func GroupDeleteByID(ctx *gin.Context) {
	db := ctx.MustGet("database").(*gorm.DB)
	groupIDArg := ctx.Param("groupID")

	if len(groupIDArg) == 0 {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": structs.GACLAPIError{Message: "Missing :groupID query param value"},
				"result": nil})
	}

	groupID, parseError := strconv.ParseInt(groupIDArg, 8, 64)

	if parseError != nil {
		panic(parseError)
	}

	trx := db.Begin()

	dbError := trx.Delete(&models.Group{ID: groupID})

	if dbError.Error != nil {
		trx.Rollback()
		panic(dbError.Error)
	}

	trx.Commit()

	ctx.JSON(http.StatusOK, gin.H{"error": nil, "result": nil})
}

// GroupUpdateByID , update group entry by id
func GroupUpdateByID(ctx *gin.Context) {
	db := ctx.MustGet("database").(*gorm.DB)
	groupIDArg := ctx.Param("groupID")

	if len(groupIDArg) == 0 {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": structs.GACLAPIError{Message: "Missing :groupID query param value"},
				"result": nil})
		return
	}

	var rgroupUpdate structs.GroupUpdateRequest

	if err := ctx.ShouldBindJSON(&rgroupUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	groupID, parseError := strconv.ParseInt(groupIDArg, 8, 64)

	if parseError != nil {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": structs.GACLAPIError{Message: "Internal server error"},
				"result": nil})
		panic(parseError)
	}

	trx := db.Begin()

	dbError := trx.Model(&models.Group{ID: groupID}).Updates(&models.Group{Name: rgroupUpdate.Name})

	if dbError.Error != nil {
		trx.Rollback()
		panic(dbError.Error)
	}

	trx.Commit()

	ctx.JSON(http.StatusOK, gin.H{"error": nil, "result": nil})
}
