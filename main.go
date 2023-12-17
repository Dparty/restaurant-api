package main

import (
	"github.com/Dparty/restaurant-api/controllers"
	_ "github.com/Dparty/restaurant-api/schedule"
)

func main() {
	controllers.Init(":8080")
}
