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

// PermissionGetByID get permissiong row
func PermissionGetByID(ctx *gin.Context) {
	db := ctx.MustGet("database").(*gorm.DB)
	permissionIDArg := ctx.Param("permissionID")

	var result []models.Permission

	if len(permissionIDArg) == 0 {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": structs.GACLAPIError{Message: "Missing :permissionID query param value"},
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
}

// PermissionGetAll get all permissions row, with default pagination.
func PermissionGetAll(ctx *gin.Context) {
	db := ctx.MustGet("database").(*gorm.DB)
	var result []models.Permission
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

// PermissionDeleteByID delete permission by id.
func PermissionDeleteByID(ctx *gin.Context) {
	db := ctx.MustGet("database").(*gorm.DB)
	permissionIDArg := ctx.Param("permissionID")

	if len(permissionIDArg) == 0 {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": structs.GACLAPIError{Message: "Missing :permissionID query param value"},
				"result": nil})
	}

	permissionID, parseError := strconv.ParseInt(permissionIDArg, 8, 64)

	if parseError != nil {
		panic(parseError)
	}

	trx := db.Begin()

	dbError := trx.Delete(&models.Permission{ID: permissionID})

	if dbError.Error != nil {
		trx.Rollback()
		panic(dbError.Error)
	}

	trx.Commit()

	ctx.JSON(http.StatusOK, gin.H{"error": nil, "result": nil})

}

// PermissionUpdateByID update existing permission
func PermissionUpdateByID(ctx *gin.Context) {
	db := ctx.MustGet("database").(*gorm.DB)
	permissionIDArg := ctx.Param("permissionID")

	if len(permissionIDArg) == 0 {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": structs.GACLAPIError{Message: "Missing :permissionID query param value"},
				"result": nil})
		return
	}

	var rPermissionUpdate structs.PermissionUpdateRequest

	if err := ctx.ShouldBindJSON(&rPermissionUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permissionID, parseError := strconv.ParseInt(permissionIDArg, 8, 64)

	if parseError != nil {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": structs.GACLAPIError{Message: "Internal server error"},
				"result": nil})
		panic(parseError)
	}

	trx := db.Begin()

	dbError := trx.Model(&models.Permission{ID: permissionID}).Updates(&models.Permission{Name: rPermissionUpdate.Name})

	if dbError.Error != nil {
		trx.Rollback()
		panic(dbError.Error)
	}

	trx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"error": nil, "result": nil})

}
