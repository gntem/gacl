package main

import "github.com/gin-gonic/gin"

func getPermissions(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) }

func postPermission(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) }

func deletePermission(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) }

func updatePermission(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) }
