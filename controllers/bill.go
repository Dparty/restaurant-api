package controllers

import (
	"net/http"

	"github.com/Dparty/common/fault"
	"github.com/Dparty/common/utils"
	apiModels "github.com/Dparty/restaurant-api/models"
	restaurantModels "github.com/Dparty/restaurant-services"
	"github.com/chenyunda218/golambda"
	"github.com/gin-gonic/gin"
)

type BillApi struct {
}

func (BillApi) GetBill(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	billId := ctx.Param("id")
	bill, err := billService.GetById(utils.StringToUint(billId))
	if err != nil {
		fault.GinHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, BillBackward(bill))
}

func (BillApi) CancelItems(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	billId := ctx.Param("id")
	var createBillRequest apiModels.CreateBillRequest
	ctx.ShouldBindJSON(&createBillRequest)
	bill, _ := billService.GetById(utils.StringToUint(billId))
	bill.CancelItems(golambda.Map(createBillRequest.Specifications,
		func(_ int, specification apiModels.Specification) restaurantModels.Specification {
			return restaurantModels.Specification{ItemId: specification.ItemId, Options: specification.Options}
		}))
	ctx.JSON(http.StatusOK, BillBackward(bill))
}

func (BillApi) CreateOrder(ctx *gin.Context) {
	tableId := ctx.Param("id")
	var createBillRequest apiModels.CreateBillRequest
	ctx.ShouldBindJSON(&createBillRequest)
	if len(createBillRequest.Specifications) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "No Specifications",
			"code":    "70001",
		})
		return
	}
	table, err := tableService.GetById(utils.StringToUint(tableId))
	if err != nil {
		fault.GinHandler(ctx, err)
		return
	}
	var specifications []restaurantModels.Specification
	for _, specification := range createBillRequest.Specifications {
		specifications = append(specifications, restaurantModels.Specification{
			ItemId:  specification.ItemId,
			Options: specification.Options,
		})
	}
	bill, err := billService.CreateBill(*table, specifications, int64(DEFAULT_OFFSET))
	if err != nil {
		fault.GinHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, BillBackward(*bill))
}
