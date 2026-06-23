package v2handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (u *UserHandler) GetUsersV2(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "list all user (v2)"})
}

func (u *UserHandler) GetUsersByIdV2(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "get user by id (v2)"})
}

func (u *UserHandler) PostUsersV2(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, gin.H{"message": "create user (v2)"})
}

func (u *UserHandler) PutUsersV2(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "update user by id (v2)"})
}

func (u *UserHandler) DeleteUsersV2(ctx *gin.Context) {
	ctx.JSON(http.StatusNoContent, gin.H{"message": "delete user by id (v2)"})
}
