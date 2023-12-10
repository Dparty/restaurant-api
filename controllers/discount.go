package controllers

import (
	"net/http"

	"github.com/Dparty/common/utils"
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
