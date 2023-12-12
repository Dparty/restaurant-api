package controllers

import (
	"net/http"

	"github.com/Dparty/common/fault"
	"github.com/Dparty/common/utils"
	apiModels "github.com/Dparty/restaurant-api/models"
	"github.com/gin-gonic/gin"
)

type DiscountApi struct{}

func (DiscountApi) DeleteDiscount(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	id := ctx.Param("id")
	discountService.DeleteDiscount(utils.StringToUint(id))
	ctx.String(http.StatusNoContent, "")
}

func (DiscountApi) UpdateDiscount(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	id := ctx.Param("id")
	var putDiscount apiModels.CreateDiscountRequest
	if err := ctx.ShouldBindJSON(&putDiscount); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	discount := discountService.GetDiscount(utils.StringToUint(id))
	if discount == nil {
		ctx.String(http.StatusNotFound, "Discount not found")
		return
	}
	restaurant, err := restaurantService.GetRestaurant(discount.Owner().ID())
	if err != nil {
		fault.GinHandler(ctx, err)
		return
	}
	if !account.Own(&restaurant) {
		ctx.String(http.StatusForbidden, "Forbidden")
		return
	}
	discount.SetLabel(putDiscount.Label).SetOffset(putDiscount.Offset).Submit()
	ctx.JSON(http.StatusNoContent, DiscountBackward(*discount))
}
