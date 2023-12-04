package controllers

import (
	authServices "github.com/Dparty/auth-services"
	"github.com/Dparty/common/server"
	restaurantServices "github.com/Dparty/restaurant-services"
	"github.com/gin-gonic/gin"
)

var authService = authServices.GetAuthService()
var restaurantService = restaurantServices.GetRestaurantService()
var printerService = restaurantServices.GetPrinterService()
var itemService = restaurantServices.GetItemService()
var tableService = restaurantServices.GetTableService()
var billService = restaurantServices.GetBillService()

func Init(addr ...string) {
	router := gin.Default()
	var restaurantApi RestaurantApi
	var billApi BillApi
	server.MetricsMiddleware(router)
	router.Use(authService.Auth())
	router.Use(server.CorsMiddleware())
	router.POST("/sessions", CreateSession)
	router.POST("/restaurants", restaurantApi.CreateRestaurant)
	router.GET("/restaurants", restaurantApi.ListRestaurant)
	router.GET("/restaurants/:id", restaurantApi.GetRestaurant)
	router.PUT("/restaurants/:id", restaurantApi.UpdateRestaurant)
	router.DELETE("/restaurants/:id", restaurantApi.DeleteRestaurant)
	router.POST("/restaurants/:id/items", restaurantApi.CreateItem)
	router.POST("/restaurants/:id/tables", restaurantApi.CreateTable)
	router.POST("/restaurants/:id/printers", restaurantApi.CreatePrinter)
	router.GET("/restaurants/:id/printers", restaurantApi.ListPrinter)
	router.DELETE("/printers/:id", restaurantApi.DeletePrinter)
	router.PUT("/printers/:id", restaurantApi.UpdatePrinter)
	router.DELETE("/items/:id", restaurantApi.DeleteItem)
	router.PUT("/items/:id", restaurantApi.UpdateItem)
	router.PUT("/items/:id/images", restaurantApi.UploadItemCover)
	router.DELETE("/tables/:id", restaurantApi.DeleteTable)
	router.PUT("/tables/:id", restaurantApi.UpdateTable)
	router.POST("/tables/:id/orders", billApi.CreateOrder)
	router.GET("/bills", billApi.ListBills)
	router.GET("/bills/subscription", billApi.BillSubscription)
	router.GET("/bills/:id", billApi.GetBill)
	router.PATCH("/bills/:id/items/cancel", billApi.CancelItems)
	router.POST("/bills/print", restaurantApi.PrintBills)
	router.POST("/bills/set", restaurantApi.SetBills)
	router.Run(addr...)
}
