package models

type Group string

const (
	Friday Group = "friday"
	Monday Group = "monday"
	Mini   Group = "mini"
)

var GroupWithEquipment = map[Group][]EquipmentType{
	Friday: {Helmet, Jacket, Gloves, Trousers, Boots, TShirt},
	Monday: {Helmet, Jacket, Gloves, Trousers, Boots, TShirt},
	Mini:   {Helmet, Gloves, TShirt},
}
