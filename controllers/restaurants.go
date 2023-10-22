package controllers

import (
	"net/http"

	"github.com/Dparty/common/fault"
	"github.com/Dparty/common/utils"
	apiModels "github.com/Dparty/restaurant-api/models"
	restaurantModels "github.com/Dparty/restaurant-services/models"
	"github.com/chenyunda218/golambda"
	"github.com/gin-gonic/gin"
)

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
	)
	if err != nil {
		fault.GinHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, ItemConvert(item))
}

func (RestaurantApi) ListRestaurant(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		return
	}
	accountId := account.ID()
	restaurants := restaurantService.List(&accountId)
	ctx.JSON(http.StatusOK, golambda.Map(restaurants,
		func(_ int, restaurant restaurantModels.Restaurant) apiModels.Restaurant {
			return RestaurantConvert(restaurant)
		}))
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
	table, err := restaurant.CreatePrinter(putPrinterRequest.Type, putPrinterRequest.Sn, putPrinterRequest.Name, "")
	if err != nil {
		fault.GinHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, PrinterBackward(table))
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
	bill, err := billService.CreateBill(*table, specifications)
	if err != nil {
		fault.GinHandler(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, apiModels.Bill{
		ID:         utils.UintToString(bill.ID()),
		PickUpCode: bill.PickUpCode(),
	})
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
