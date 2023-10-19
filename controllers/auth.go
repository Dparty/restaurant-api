package controllers

import (
	"net/http"

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
