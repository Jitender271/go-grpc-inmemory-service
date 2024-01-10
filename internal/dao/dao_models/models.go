package daomodels

type Movies struct {
	MovieID     string `db:"movie_id"`
	Name        string `db:"name"`
	Genre       string `db:"genre"`
	Description string `db:"description"`
	Rating      string `db:"rating"`
}

type BookMyMovie struct {
	Id          string `db:"id"`
	MovieName   string `db:"movie_name"`
	TheatreName string `db:"theatre_name"`
	StartTime   int64  `db:"start_time"`
	EndTime     int64  `db:"end_time"`
}

type GetBookingQueryObject struct {
	MovieName   string `db:"movie_name"`
	TheatreName string `db:"theatre_name"`
}
