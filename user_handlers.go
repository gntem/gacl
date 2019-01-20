package main

import "github.com/gin-gonic/gin"

func getUsers(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) }

func postUser(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) }

func deleteUser(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) }

func updateUser(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) }
