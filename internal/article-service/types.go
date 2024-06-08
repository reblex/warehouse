package main

type ArticleId string
type Article struct {
	Id    ArticleId `json:"art_id"`
	Name  string    `json:"name"`
	Stock int       `json:"stock,string"`
}

type Reservation struct {
	Id    ArticleId `json:"id"`
	Count int       `json:"count"`
}

type ArticlesDto struct {
	Articles []Article `json:"articles"`
}
