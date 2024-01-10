package models

type Movie struct {
	Id     string
	Name   string
	Genre  string
	Desc   string
	Rating string
}

type BookMyMovie struct {
	Id          string
	MovieName   string
	TheatreName string
	StartTime   int64
	EndTime     int64
}

type GetBooking struct {
	MovieName   string
	TheatreName string
}
