package models

import (
	models "github.com/Dparty/dao/restaurant"
)

type PutRestaurantRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Restaurant struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Items       []Item  `json:"items"`
	Tables      []Table `json:"tables"`
}

type PutItemRequest struct {
	Tags       []string           `json:"tags"`
	Printers   []string           `json:"printers"`
	Name       string             `json:"name"`
	Pricing    int64              `json:"pricing"`
	Attributes []models.Attribute `json:"attributes"`
	Images     []string           `json:"images"`
}

type Item struct {
	Tags       []string           `json:"tags"`
	Printers   []string           `json:"printers"`
	Id         string             `json:"id"`
	Name       string             `json:"name"`
	Pricing    int64              `json:"pricing"`
	Attributes []models.Attribute `json:"attributes"`
	Images     []string           `json:"images"`
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
}

type Printer struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	Sn          string `json:"sn"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
