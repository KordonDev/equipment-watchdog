package models

type Group string

const (
	Friday Group = "friday"
	Monday Group = "monday"
	Mini   Group = "mini"
)

type EquipmentType string

const (
	Helmet   EquipmentType = "helmet"
	Jacket   EquipmentType = "jacket"
	Gloves   EquipmentType = "gloves"
	Trousers EquipmentType = "trousers"
	Boots    EquipmentType = "boots"
	TShirt   EquipmentType = "tshirt"
)

var GroupWithEquipment = map[Group][]EquipmentType{
	Friday: {Helmet, Jacket, Gloves, Trousers, Boots, TShirt},
	Monday: {Helmet, Jacket, Gloves, Trousers, Boots, TShirt},
	Mini:   {Helmet, Gloves, TShirt},
}
