package main

type Order struct {
	ProductName string `json:"product_name"`
	Amount      int    `json:"amount"`
}

type OrderDto struct {
	Orders []Order `json:"orders"`
}
