package main

type Article struct {
	Id     string `json:"art_id"`
	Amount int    `json:"amount_of,string"`
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
