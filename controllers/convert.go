package controllers

import (
	"github.com/Dparty/common/utils"
	model "github.com/Dparty/dao/restaurant"
	apiModels "github.com/Dparty/restaurant-api/models"
	serviceModels "github.com/Dparty/restaurant-services/models"
	"github.com/chenyunda218/golambda"
)

func arrayConvert[T any](array []T) []T {
	if len(array) == 0 {
		return make([]T, 0)
	}
	return array
}

func RestaurantConvert(restaurant serviceModels.Restaurant) apiModels.Restaurant {
	return apiModels.Restaurant{
		Id:          utils.UintToString(restaurant.ID()),
		Name:        restaurant.Name(),
		Description: restaurant.Description(),
		Items: golambda.Map(restaurant.Items(), func(_ int, i serviceModels.Item) apiModels.Item {
			return ItemConvert(i)
		}),
		Tables: golambda.Map(restaurant.Tables(), func(_ int, i serviceModels.Table) apiModels.Table {
			return TableBackward(i)
		}),
	}
}

func TableBackward(table serviceModels.Table) apiModels.Table {
	return apiModels.Table{
		Id:    utils.UintToString(table.ID()),
		Label: table.Label(),
		X:     table.X(),
		Y:     table.Y(),
	}
}

func ItemConvert(item serviceModels.Item) apiModels.Item {
	entity := item.Entity()
	return apiModels.Item{
		Id:         utils.UintToString(item.ID()),
		Name:       entity.Name,
		Pricing:    item.Entity().Pricing,
		Attributes: arrayConvert(entity.Attributes),
		Images:     arrayConvert(entity.Images),
		Tags:       arrayConvert(entity.Tags),
		Printers: golambda.Map(entity.Printers, func(_ int, id uint) string {
			return utils.UintToString(id)
		}),
	}
}

func AttributeForward(attribute model.Attribute) apiModels.Attribute {
	return apiModels.Attribute{
		Label:   attribute.Label,
		Options: attribute.Options,
	}
}

func ItemForward(item apiModels.PutItemRequest) model.Item {
	return model.Item{
		Name:       item.Name,
		Pricing:    item.Pricing,
		Attributes: arrayConvert(item.Attributes),
		Tags:       arrayConvert(item.Tags),
		Printers: golambda.Map(item.Printers, func(_ int, printer string) uint {
			return utils.StringToUint(printer)
		}),
	}
}

func PrinterBackward(printer serviceModels.Printer) apiModels.Printer {
	return apiModels.Printer{
		Id:          utils.UintToString(printer.ID()),
		Type:        printer.Type(),
		Sn:          printer.Sn(),
		Name:        printer.Name(),
		Description: printer.Description(),
	}
}

func BillBackward(bill serviceModels.Bill) apiModels.Bill {
	return apiModels.Bill{
		Status:     bill.Entity().Status,
		ID:         utils.UintToString(bill.ID()),
		PickUpCode: bill.PickUpCode(),
		Orders: golambda.Map(bill.Orders(), func(_ int, order model.Order) apiModels.Order {
			return apiModels.Order{
				Item:          ItemConvert(serviceModels.NewItem(order.Item)),
				Specification: order.Specification,
			}
		}),
	}
}
