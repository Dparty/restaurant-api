package controllers

import (
	"fmt"

	authServices "github.com/Dparty/auth-services"
	"github.com/Dparty/common/server"
	"github.com/Dparty/dao"
	restaurantServices "github.com/Dparty/restaurant-services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var authService authServices.AuthService
var restaurantService restaurantServices.RestaurantService
var printerService restaurantServices.PrinterService
var itemService restaurantServices.ItemService
var tableService restaurantServices.TableService
var router *gin.Engine

func Init(addr ...string) {
	var err error
	viper.SetConfigName(".env.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("databases fatal error config file: %w", err))
	}
	user := viper.GetString("database.user")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	database := viper.GetString("database.database")
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
	router = gin.Default()
	router.Use(authService.Auth())
	router.Use(server.CorsMiddleware())
	router.POST("/sessions", CreateSession)
	router.POST("/restaurants", CreateRestaurant)
	router.GET("/restaurants", ListRestaurant)
	router.GET("/restaurants/:id", GetRestaurant)
	router.POST("/restaurants/:id/items", CreateItem)
	router.POST("/restaurants/:id/tables", CreateTable)
	router.POST("/restaurants/:id/printers", CreatePrinter)
	router.GET("/restaurants/:id/printers", ListPrinter)
	router.DELETE("/printers/:id", DeletePrinter)
	router.DELETE("/items/:id", DeleteItem)
	router.PUT("/items/:id/images", UploadItemCover)
	router.DELETE("/tables/:id", DeleteTable)
	router.POST("/tables/:id/orders", CreateOrder)
	router.Run(addr...)
}
