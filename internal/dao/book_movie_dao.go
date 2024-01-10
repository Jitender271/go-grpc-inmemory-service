package dao

import (
	"context"
	"errors"
	"github.com/go-grpc-inmemory-service/commons/utils"
	"github.com/go-grpc-inmemory-service/internal/config"
	daomodels "github.com/go-grpc-inmemory-service/internal/dao/dao_models"
	"github.com/go-grpc-inmemory-service/internal/db"
	"github.com/go-grpc-inmemory-service/internal/models"
	"github.com/gocql/gocql"
	"sync"
)

var (
	once       sync.Once
	bookingDao BookMovieDao
)

type BookMovieDao interface {
	InsertBooking(ctx context.Context, req *models.BookMyMovie) (*models.BookMyMovie, error)
	GetBookings(ctx context.Context, movieName, theatreName string) ([]daomodels.BookMyMovie, error)
	//GetAllMovies(ctx context.Context) ([]daomodels.Movies, error)
	//UpdateMovies(ctx context.Context, req *models.Movie) (*models.Movie, error)
}

type BookingImpl struct {
	SessionWrapper db.SessionWrapperService
}

func NewBookingDaoImpl(dbConfigs config.DbConfigs) BookMovieDao {
	once.Do(func() {
		bookingDao = &BookingImpl{
			SessionWrapper: db.GetSession(dbConfigs),
		}
	})
	return bookingDao
}

func (b *BookingImpl) InsertBooking(ctx context.Context, req *models.BookMyMovie) (*models.BookMyMovie, error) {
	req.Id = gocql.TimeUUID().String()
	req.StartTime = utils.GetCurrentTimestampInMillis()
	req.EndTime = utils.GetCurrentTimestampInMillis()
	daoMovie := convertToDaoBooking(req)
	if err := insertBookingInDb(ctx, b.SessionWrapper, daoMovie); err != nil {
		return nil, errors.New("error inserting key in db" + req.MovieName)
	}
	return req, nil
}

func (b *BookingImpl) GetBookings(ctx context.Context, movieName, theatre string) ([]daomodels.BookMyMovie, error) {
	return getBookingFromDb(ctx, b.SessionWrapper, movieName, theatre)
}

func convertToDaoBooking(req *models.BookMyMovie) *daomodels.BookMyMovie {
	return &daomodels.BookMyMovie{
		Id:          req.Id,
		MovieName:   req.MovieName,
		TheatreName: req.TheatreName,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	}

}
