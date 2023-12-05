package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Dparty/common/fault"
	"github.com/Dparty/common/utils"
	apiModels "github.com/Dparty/restaurant-api/models"
	restaurantModels "github.com/Dparty/restaurant-services"
	"github.com/chenyunda218/golambda"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type BillApi struct{}

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

func (BillApi) BillSubscription(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	restaurantId := ctx.Query("restaurantId")
	restaurant, err := restaurantService.GetRestaurant(utils.StringToUint(restaurantId))
	if err != nil {
		fault.GinHandler(ctx, fault.ErrNotFound)
		return
	}
	if account.ID() != restaurant.Owner().ID() {
		fault.GinHandler(ctx, fault.ErrPermissionDenied)
		return
	}
	tableId := ctx.Query("tableId")
	status := ctx.Query("status")
	startAt := ctx.Query("startAt")
	endAt := ctx.Query("endAt")
	var _tableId *uint
	if tableId != "" {
		_tableId = golambda.Reference(utils.StringToUint(tableId))
	}
	var _status *string
	if status != "" {
		_status = &status
	}
	var _startAt *time.Time
	if startAt != "" {
		_startAt = golambda.Reference(time.Unix(int64(utils.StringToUint(startAt)), 0))
	}
	var _endAt *time.Time
	if endAt != "" {
		_endAt = golambda.Reference(time.Unix(int64(utils.StringToUint(endAt)), 0))
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	mess := make(chan string)
	var getBills = func() {
		j, _ := json.Marshal(
			golambda.Map(billService.ListBills(utils.StringToUint(restaurantId), _tableId, _status, _startAt, _endAt),
				func(_ int, bill restaurantModels.Bill) apiModels.Bill {
					return BillBackward(bill)
				}))
		conn.WriteMessage(websocket.TextMessage, j)
	}
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				close(mess)
				break
			}
		}
	}()
	c := pb.Subscribe("restaurant-" + restaurantId)
	ch := c.Channel()
	for {
		select {
		case <-ch:
			getBills()
		case _, ok := <-mess:
			if !ok {
				c.Close()
				return
			}
		case <-time.After(time.Second * 3):
			pb.Publish("restaurant-"+restaurantId, "polling")
		}
	}
}

func (BillApi) ListBills(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	restaurantId := ctx.Query("restaurantId")
	restaurant, err := restaurantService.GetRestaurant(utils.StringToUint(restaurantId))
	if err != nil {
		fault.GinHandler(ctx, fault.ErrNotFound)
		return
	}
	if account.ID() != restaurant.Owner().ID() {
		fault.GinHandler(ctx, fault.ErrPermissionDenied)
		return
	}
	tableId := ctx.Query("tableId")
	status := ctx.Query("status")
	startAt := ctx.Query("startAt")
	endAt := ctx.Query("endAt")
	var _tableId *uint
	if tableId != "" {
		_tableId = golambda.Reference(utils.StringToUint(tableId))
	}
	var _status *string
	if status != "" {
		_status = &status
	}
	var _startAt *time.Time
	if startAt != "" {
		_startAt = golambda.Reference(time.Unix(int64(utils.StringToUint(startAt)), 0))
	}
	var _endAt *time.Time
	if endAt != "" {
		_endAt = golambda.Reference(time.Unix(int64(utils.StringToUint(endAt)), 0))
	}
	bills := billService.ListBills(utils.StringToUint(restaurantId), _tableId, _status, _startAt, _endAt)
	ctx.JSON(http.StatusOK, golambda.Map(bills,
		func(_ int, bill restaurantModels.Bill) apiModels.Bill {
			return BillBackward(bill)
		}))
}
