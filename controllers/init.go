package controllers

import (
	"github.com/Depado/ginprom"
	authServices "github.com/Dparty/auth-services"
	"github.com/Dparty/common/config"
	"github.com/Dparty/common/server"
	"github.com/Dparty/dao"
	restaurantServices "github.com/Dparty/restaurant-services"
	"github.com/gin-gonic/gin"
)

var authService authServices.AuthService
var restaurantService *restaurantServices.RestaurantService
var printerService restaurantServices.PrinterService
var itemService restaurantServices.ItemService
var tableService restaurantServices.TableService
var billService restaurantServices.BillService

func Init(addr ...string) {
	user := config.GetString("database.user")
	password := config.GetString("database.password")
	host := config.GetString("database.host")
	port := config.GetString("database.port")
	database := config.GetString("database.database")
	db, err := dao.NewConnection(user, password, host, port, database)
	restaurantServices.Init(db)
	if err != nil {
		panic(err)
	}
	authService = authServices.NewAuthService(db)
	restaurantService = restaurantServices.NewRestaurantService(db)
	printerService = restaurantServices.NewPrinterService(db)
	itemService = restaurantServices.NewItemService(db)
	tableService = restaurantServices.NewTableService(db)
	billService = restaurantServices.NewBillService(db)
	router := gin.Default()
	var restaurantApi RestaurantApi
	p := ginprom.New(
		ginprom.Engine(router),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
	)
	router.Use(p.Instrument())
	router.Use(authService.Auth())
	router.Use(server.CorsMiddleware())
	router.POST("/sessions", CreateSession)
	router.POST("/restaurants", restaurantApi.CreateRestaurant)
	router.GET("/restaurants", restaurantApi.ListRestaurant)
	router.GET("/restaurants/:id", restaurantApi.GetRestaurant)
	router.PUT("/restaurants/:id", restaurantApi.UpdateRestaurant)
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
	router.POST("/tables/:id/orders", restaurantApi.CreateOrder)
	router.GET("/bills", restaurantApi.ListBills)
	router.PATCH("/bills/:id/items/cancel", restaurantApi.CancelItems)
	router.POST("/bills/print", restaurantApi.PrintBills)
	router.POST("/bills/set", restaurantApi.SetBills)
	router.Run(addr...)
}
