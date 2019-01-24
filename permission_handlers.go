package main

import "github.com/gin-gonic/gin"

func (envCtx *Env) getPermissions(ctx *gin.Context)   { ctx.JSON(200, gin.H{"n": 1}) }
func (envCtx *Env) postPermission(ctx *gin.Context)   { ctx.JSON(200, gin.H{"n": 1}) }
func (envCtx *Env) deletePermission(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) }
func (envCtx *Env) updatePermission(ctx *gin.Context) { ctx.JSON(200, gin.H{"n": 1}) }
