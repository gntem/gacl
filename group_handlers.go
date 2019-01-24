package main

import "github.com/gin-gonic/gin"

func (envCtx *Env) getAllGroups(ctx *gin.Context)         { ctx.JSON(200, gin.H{"n": 1}) }
func (envCtx *Env) getGroup(ctx *gin.Context)             { ctx.JSON(200, gin.H{"n": 1}) }
func (envCtx *Env) createGroup(ctx *gin.Context)          { ctx.JSON(200, gin.H{"n": 1}) }
func (envCtx *Env) removeGroup(ctx *gin.Context)          { ctx.JSON(200, gin.H{"n": 1}) }
func (envCtx *Env) updateGroup(ctx *gin.Context)          { ctx.JSON(200, gin.H{"n": 1}) }
func (envCtx *Env) upsertGroup(ctx *gin.Context)          { ctx.JSON(200, gin.H{"n": 1}) }
func (envCtx *Env) listGroupUsers(ctx *gin.Context)       { ctx.JSON(200, gin.H{"n": 1}) }
func (envCtx *Env) listGroupPermissions(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) }
func (envCtx *Env) addUserToGroup(ctx *gin.Context)       { ctx.JSON(200, gin.H{"n": 1}) }
func (envCtx *Env) removeUserFromGroup(ctx *gin.Context)  { ctx.JSON(200, gin.H{"n": 1}) }
