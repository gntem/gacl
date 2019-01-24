package main

import (
	"github.com/gin-gonic/gin"
)

func (envCtx *Env) getAllUsers(ctx *gin.Context)      { ctx.JSON(200, gin.H{"n": 1}) }
func (envCtx *Env) getUser(ctx *gin.Context)          { ctx.JSON(200, gin.H{"n": 1}) }
func (envCtx *Env) createUser(ctx *gin.Context)       { ctx.JSON(200, gin.H{"n": 1}) }
func (envCtx *Env) removeUser(ctx *gin.Context)       { ctx.JSON(200, gin.H{"n": 1}) }
func (envCtx *Env) updateUser(ctx *gin.Context)       { ctx.JSON(200, gin.H{"n": 1}) }
func (envCtx *Env) upsertUser(ctx *gin.Context)       { ctx.JSON(200, gin.H{"n": 1}) }
func (envCtx *Env) grantPermission(ctx *gin.Context)  { ctx.JSON(200, gin.H{"n": 1}) }
func (envCtx *Env) revokePermission(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) }
