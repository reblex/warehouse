package main

type ArticleId string

type Article struct {
	Id     ArticleId `json:"art_id"`
	Amount int       `json:"amount_of,string"`
}

type ProductName string

type Product struct {
	Name     ProductName `json:"name"`
	Articles []Article   `json:"contain_articles"`
	Price    float32     `json:"price"`
}

type ProductsDto struct {
	Products []Product `json:"products"`
}

type Order struct {
	ProductName ProductName `json:"product_name"`
	Amount      int         `json:"amount"`
}

type OrderDto struct {
	Orders []Order `json:"orders"`
}

type ArticleReservation struct {
	Id    ArticleId `json:"id"`
	Count int       `json:"count,string"`
}

type ArticleReservationsDto struct {
	Reservations []ArticleReservation `json:"reservations"`
}
