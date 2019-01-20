package main

import "github.com/gin-gonic/gin"

func getGroups(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) }

func postGroup(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) }

func deleteGroup(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) }

func updateGroup(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) }
