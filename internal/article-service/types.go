package main

type ArticleId string
type Article struct {
	Id    ArticleId `json:"art_id"`
	Name  string    `json:"name"`
	Stock int       `json:"stock,string"`
}

type Reservation struct {
	Id    ArticleId `json:"id"`
	Count int       `json:"count,string"`
}

type ArticlesDto struct {
	Articles []Article `json:"articles"`
}

type ReservationsDto struct {
	Reservations []Reservation `json:"reservations"`
}

type AvailabilityDto struct {
	Availability int `json:"availability"`
}
