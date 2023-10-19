package controllers

import (
	"net/http"

	authServices "github.com/Dparty/auth-services"
	"github.com/Dparty/common/fault"
	"github.com/Dparty/restaurant-api/models"
	"github.com/gin-gonic/gin"
)

func CreateSession(ctx *gin.Context) {
	var createSessionRequest models.CreateSessionRequest
	ctx.ShouldBindJSON(&createSessionRequest)
	token, err := authService.CreateSession(createSessionRequest.Email, createSessionRequest.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, "")
		return
	}
	ctx.JSON(http.StatusCreated, &models.Session{
		Token: token,
	})
}

func getAccount(ctx *gin.Context) *authServices.Account {
	accountInterface, ok := ctx.Get("account")
	if !ok {
		fault.GinHandler(ctx, fault.ErrUnauthorized)
		return nil
	}
	account, ok := accountInterface.(authServices.Account)
	if !ok {
		fault.GinHandler(ctx, fault.ErrUnauthorized)
		return nil
	}
	return &account
}
