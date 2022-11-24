package members

type Group string

const (
	Friday Group = "friday"
	Monday Group = "monday"
	Mini   Group = "mini"
)

type Equipment string

const (
	Helmet   Equipment = "helmet"
	Jacket   Equipment = "jacket"
	Gloves   Equipment = "gloves"
	Trousers Equipment = "trousers"
	Boots    Equipment = "boots"
	TShirt   Equipment = "tshirt"
)

var GroupWithEquipment = map[Group][]Equipment{
	Friday: {Helmet, Jacket, Gloves, Trousers, Boots, TShirt},
	Monday: {Helmet, Jacket, Gloves, Trousers, Boots, TShirt},
	Mini:   {Helmet, Gloves, TShirt},
}
