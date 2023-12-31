package models

import (
	"github.com/Dparty/common/data"
	models "github.com/Dparty/dao/restaurant"
)

type PutRestaurantRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Categories  []string `json:"categories"`
}

type Restaurant struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Items       []Item   `json:"items"`
	Tables      []Table  `json:"tables"`
	Categories  []string `json:"categories"`
}

type PutItemRequest struct {
	Tags       []string           `json:"tags"`
	Printers   []string           `json:"printers"`
	Name       string             `json:"name"`
	Pricing    int64              `json:"pricing"`
	Attributes []models.Attribute `json:"attributes"`
	Images     []string           `json:"images"`
	Status     string             `json:"status"`
}

type Item struct {
	Id         string             `json:"id"`
	Tags       []string           `json:"tags"`
	Printers   []string           `json:"printers"`
	Name       string             `json:"name"`
	Pricing    int64              `json:"pricing"`
	Attributes []models.Attribute `json:"attributes"`
	Images     []string           `json:"images"`
	Status     string             `json:"status"`
}

type Attribute struct {
	Label   string          `json:"label"`
	Options []models.Option `json:"options"`
}

type Option struct {
	Label string `json:"label"`
	Extra int64  `json:"extra"`
}

type PutTableRequest struct {
	X     int64  `json:"x"`
	Y     int64  `json:"y"`
	Label string `json:"label"`
}

type Table struct {
	X     int64  `json:"x"`
	Y     int64  `json:"y"`
	Id    string `json:"id"`
	Label string `json:"label"`
}

type PutPrinterRequest struct {
	Type        string `json:"type"`
	Sn          string `json:"sn"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Model       string `json:"model" default:"58mm"`
}

type Printer struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	Sn          string `json:"sn"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Model       string `json:"model"`
}

type CreateBillRequest struct {
	Specifications []Specification `json:"specifications"`
}

type Specification struct {
	ItemId  string                      `json:"itemId"`
	Options []data.Pair[string, string] `json:"options"`
}

type Order struct {
	Item          Item                        `json:"item"`
	Specification []data.Pair[string, string] `json:"specification"`
}

type Bill struct {
	ID         string  `json:"id"`
	Orders     []Order `json:"orders"`
	PickUpCode int64   `json:"pickUpCode"`
	Status     string  `json:"status"`
	TableLabel string  `json:"tableLabel"`
	Offset     int64   `json:"offset"`
	CreatedAt  int64   `json:"createdAt"`
	Total      int64   `json:"total"`
}

type PrintBillRequest struct {
	Offset int64 `json:"offset"`
}

type RestaurantList struct {
	Data []Restaurant `json:"data"`
}

type PrintBillsRequest struct {
	Offset     int64    `json:"offset"`
	BillIdList []string `json:"billIdList"`
}

type SetBillsRequest struct {
	Offset     int64    `json:"offset"`
	BillIdList []string `json:"billIdList"`
	Status     string   `json:"status"`
}

type CreateDiscountRequest struct {
	Label  string `json:"label"`
	Offset int64  `json:"offset"`
}

type Discount struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Offset int64  `json:"offset"`
}
