package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Dparty/common/fault"
	"github.com/Dparty/common/utils"
	apiModels "github.com/Dparty/restaurant-api/models"
	restaurantModels "github.com/Dparty/restaurant-services/models"
	"github.com/chenyunda218/golambda"
	"github.com/gin-gonic/gin"
)

var DEFAULT_OFFSET = 10

type RestaurantApi struct{}

func (RestaurantApi) CreateRestaurant(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	var request apiModels.PutRestaurantRequest
	ctx.ShouldBindJSON(&request)
	restaurant := restaurantService.CreateRestaurant(account, request.Name, request.Description)
	ctx.JSON(http.StatusCreated, RestaurantConvert(restaurant))
}

func (RestaurantApi) UpdateRestaurant(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	id := ctx.Param("id")
	var request apiModels.PutRestaurantRequest
	ctx.ShouldBindJSON(&request)
	restaurant, err := restaurantService.UpdateRestaurant(utils.StringToUint(id), request.Name, request.Description, request.Categories)
	if err != nil {
		fault.GinHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, RestaurantConvert(restaurant))
}

func (RestaurantApi) CreateItem(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	restaurant, _ := restaurantService.GetRestaurant(utils.StringToUint(ctx.Param("id")))
	if !account.Own(&restaurant) {
		fault.GinHandler(ctx, fault.ErrUnauthorized)
		return
	}
	var putItemRequest apiModels.PutItemRequest
	ctx.ShouldBindJSON(&putItemRequest)
	item, err := restaurant.CreateItem(putItemRequest.Name,
		putItemRequest.Pricing,
		putItemRequest.Attributes,
		putItemRequest.Images,
		putItemRequest.Tags,
		golambda.Map(putItemRequest.Printers, func(_ int, printerId string) uint {
			return utils.StringToUint(printerId)
		}),
		putItemRequest.Status,
	)
	if err != nil {
		fault.GinHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, ItemConvert(item))
}

func (RestaurantApi) UpdateItem(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	id := ctx.Param("id")
	item, err := itemService.GetById(utils.StringToUint(id))
	if err != nil {
		fault.GinHandler(ctx, err)
	}
	if item.Owner().Owner().ID() != account.ID() {
		fault.GinHandler(ctx, fault.ErrPermissionDenied)
		return
	}
	var request apiModels.PutItemRequest
	ctx.ShouldBindJSON(&request)
	// TODO: refactor
	entity := item.Entity()
	entity.Name = request.Name
	entity.Tags = request.Tags
	entity.Attributes = request.Attributes
	entity.Images = request.Images
	entity.Pricing = request.Pricing
	entity.Printers = golambda.Map(request.Printers, func(_ int, printerId string) uint {
		return utils.StringToUint(printerId)
	})
	entity.Status = request.Status
	item.SetEntity(entity)
	item.Save()
	ctx.JSON(http.StatusCreated, ItemConvert(item))
}

func (RestaurantApi) ListRestaurant(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	accountId := account.ID()
	restaurants := restaurantService.List(&accountId)
	ctx.JSON(http.StatusOK,
		apiModels.RestaurantList{Data: golambda.Map(restaurants,
			func(_ int, restaurant restaurantModels.Restaurant) apiModels.Restaurant {
				return RestaurantConvert(restaurant)
			})})
}

func (RestaurantApi) CreateTable(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	restaurant, err := restaurantService.GetRestaurant(utils.StringToUint(ctx.Param("id")))
	if err != nil {
		fault.GinHandler(ctx, fault.ErrNotFound)
		return
	}
	if !account.Own(&restaurant) {
		fault.GinHandler(ctx, fault.ErrUnauthorized)
		return
	}
	var request apiModels.PutTableRequest
	ctx.ShouldBindJSON(&request)
	table, err := restaurant.CreateTable(request.Label, request.X, request.Y)
	if err != nil {
		fault.GinHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, TableBackward(table))
}

func (RestaurantApi) UpdateTable(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	id := ctx.Param("id")
	table, err := tableService.GetById(utils.StringToUint(id))
	if err != nil {
		fault.GinHandler(ctx, err)
		return
	}
	if table.Owner().Owner().ID() != account.ID() {
		fault.GinHandler(ctx, fault.ErrPermissionDenied)
		return
	}
	var request apiModels.PutTableRequest
	ctx.ShouldBindJSON(&request)
	if !table.Update(request.Label, request.X, request.Y) {
		fault.GinHandler(ctx, fault.ErrCreateTableConflict)
		return
	}
	ctx.JSON(http.StatusOK, apiModels.Table{
		Id:    utils.UintToString(table.ID()),
		X:     request.X,
		Y:     request.Y,
		Label: request.Label,
	})
}

func (RestaurantApi) GetRestaurant(ctx *gin.Context) {
	restaurant, err := restaurantService.GetRestaurant(utils.StringToUint(ctx.Param("id")))
	if err != nil {
		fault.GinHandler(ctx, fault.ErrNotFound)
		return
	}
	ctx.JSON(http.StatusOK, RestaurantConvert(restaurant))
}

func (RestaurantApi) CreatePrinter(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	restaurant, err := restaurantService.GetRestaurant(utils.StringToUint(ctx.Param("id")))
	if err != nil {
		fault.GinHandler(ctx, fault.ErrNotFound)
		return
	}
	if !account.Own(&restaurant) {
		fault.GinHandler(ctx, fault.ErrUnauthorized)
		return
	}
	var putPrinterRequest apiModels.PutPrinterRequest
	ctx.ShouldBindJSON(&putPrinterRequest)
	table, err := restaurant.CreatePrinter(putPrinterRequest.Type, putPrinterRequest.Sn, putPrinterRequest.Name, putPrinterRequest.Description, putPrinterRequest.Model)
	if err != nil {
		fault.GinHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, PrinterBackward(table))
}

func (RestaurantApi) UpdatePrinter(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	id := ctx.Param("id")
	printer, err := printerService.GetById(utils.StringToUint(id))
	if err != nil {
		fault.GinHandler(ctx, err)
		return
	}
	if account.ID() != printer.Owner().Owner().ID() {
		fault.GinHandler(ctx, fault.ErrPermissionDenied)
		return
	}
	var updateRequest apiModels.PutPrinterRequest
	ctx.ShouldBindJSON(&updateRequest)
	printer.Update(updateRequest.Name, updateRequest.Description, updateRequest.Sn, updateRequest.Type, updateRequest.Model)
	ctx.JSON(http.StatusOK, PrinterBackward(*printer))
}

func (RestaurantApi) ListPrinter(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	restaurant, err := restaurantService.GetRestaurant(utils.StringToUint(ctx.Param("id")))
	if err != nil {
		fault.GinHandler(ctx, fault.ErrNotFound)
		return
	}
	if restaurant.Owner().ID() != account.ID() {
		fault.GinHandler(ctx, fault.ErrPermissionDenied)
		return
	}
	ctx.JSON(http.StatusOK, golambda.Map(restaurant.Printers(),
		func(_ int, printer restaurantModels.Printer) apiModels.Printer {
			return PrinterBackward(printer)
		}))
}

func (RestaurantApi) DeletePrinter(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	printer, err := printerService.GetById(utils.StringToUint(ctx.Param("id")))
	if err != nil {
		fault.GinHandler(ctx, err)
		return
	}
	if r := printer.Owner(); r.Owner().ID() != account.ID() {
		fault.GinHandler(ctx, fault.ErrPermissionDenied)
		return
	}
	if !printer.Delete() {
		ctx.JSON(http.StatusConflict, gin.H{
			"code":    "60010",
			"message": "printer is in using",
		})
		return
	}
	ctx.JSON(http.StatusNoContent, "")
}

func (RestaurantApi) DeleteItem(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	item, err := itemService.GetById(utils.StringToUint(ctx.Param("id")))
	if err != nil {
		fault.GinHandler(ctx, err)
		return
	}
	res, _ := restaurantService.GetRestaurant(item.Owner().ID())
	if !account.Own(&res) {
		fault.GinHandler(ctx, fault.ErrPermissionDenied)
		return
	}
	item.Delete()
	ctx.JSON(http.StatusNoContent, "")
}

func (RestaurantApi) DeleteTable(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	table, err := tableService.GetById(utils.StringToUint(ctx.Param("id")))
	if err != nil {
		fault.GinHandler(ctx, err)
	}
	res, _ := restaurantService.GetRestaurant(table.Owner().ID())
	if !account.Own(&res) {
		fault.GinHandler(ctx, fault.ErrPermissionDenied)
		return
	}
	table.Delete()
	ctx.JSON(http.StatusNoContent, "")
}

func (RestaurantApi) CreateOrder(ctx *gin.Context) {
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

func (RestaurantApi) PrintBills(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	var printBillRequest apiModels.PrintBillsRequest
	ctx.ShouldBindJSON(&printBillRequest)
	billService.PrintBills(account.ID(),
		golambda.Map(printBillRequest.BillIdList,
			func(_ int, id string) uint {
				return utils.StringToUint(id)
			}), printBillRequest.Offset)
}

func (RestaurantApi) SetBills(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	var setBillRequest apiModels.SetBillsRequest
	ctx.ShouldBindJSON(&setBillRequest)
	billService.SetBill(account.ID(),
		golambda.Map(setBillRequest.BillIdList,
			func(_ int, id string) uint {
				return utils.StringToUint(id)
			}), setBillRequest.Offset, setBillRequest.Status)
}

func (RestaurantApi) UploadItemCover(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		fault.GinHandler(ctx, fault.ErrUnauthorized)
		return
	}
	itemId := ctx.Param("id")
	item, err := itemService.GetById(utils.StringToUint(itemId))
	if err != nil {
		fault.GinHandler(ctx, err)
		return
	}

	if item.Owner().Owner().ID() != account.ID() {
		fault.GinHandler(ctx, fault.ErrPermissionDenied)
		return
	}

	file, _ := ctx.FormFile("file")
	url := item.UploadImage(file)
	ctx.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

func (RestaurantApi) ListBills(ctx *gin.Context) {
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

func (RestaurantApi) CancelItems(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	billId := ctx.Param("id")
	var createBillRequest apiModels.CreateBillRequest
	fmt.Println(billId, createBillRequest)
}
