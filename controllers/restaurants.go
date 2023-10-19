package controllers

import (
	"fmt"

	authServices "github.com/Dparty/auth-services"
	"github.com/Dparty/common/fault"
	"github.com/Dparty/restaurant-api/models"
	"github.com/gin-gonic/gin"
)

func CreateRestaurant(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	var request models.PutRestaurantRequest
	ctx.ShouldBindJSON(&request)
	restaurant := restaurantService.CreateRestaurant(account, request.Name, request.Description)
	fmt.Println(restaurant)
}

func getAccount(ctx *gin.Context) *authServices.Account {
	accountInterface, ok := ctx.Get("account")
	fmt.Println(accountInterface)
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
