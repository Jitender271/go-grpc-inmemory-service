package service

import (
	"context"
	//"errors"
	//"github.com/go-grpc-inmemory-service/internal/cache"
	"github.com/go-grpc-inmemory-service/internal/config"
	"github.com/go-grpc-inmemory-service/internal/dao"
	"github.com/go-grpc-inmemory-service/internal/models"
	//"github.com/go-grpc-inmemory-service/internal/service/service_helper"
	"github.com/go-grpc-inmemory-service/resources/moviepb"
)

const (
	duplicateBookingError = "movie already exists"
	bookingError          = "error validating movie details"
)

type BookingService interface {
	CreateBooking(ctx context.Context, request *moviepb.BookingRequest) (*models.BookMyMovie, error)
	GetBookings(ctx context.Context, request *moviepb.GetBookingRequest) ([]*models.BookMyMovie, error)
	//GetAllMovies(ctx context.Context) ([]*models.Movie, error)
	//UpdateMovies(ctx context.Context, request *moviepb.UpdateMovieRequest) (*models.Movie, error)
}

type BookingServiceImpl struct {
	bookingDao dao.BookMovieDao
	//movieCache cache.MovieCacheService
}

func NewBookingImpl(dbConfigs config.DbConfigs, cacheConfigs config.AppConfig) BookingService {
	return &BookingServiceImpl{
		bookingDao: dao.NewBookingDaoImpl(dbConfigs),
		//movieCache: cache.NewMovieCacheService(cacheConfigs, dao.NewMovieDaoImpl(dbConfigs)),
	}
}

func (m *BookingServiceImpl) CreateBooking(ctx context.Context, request *moviepb.BookingRequest) (*models.BookMyMovie, error) {
	//isDuplicate, getMovieDetailsErr := service_helper.IsDuplicateMovie(ctx, m.bookingDao, request.GetMovie())
	//
	//if getMovieDetailsErr != nil {
	//	return nil, errors.New(movieError)
	//}
	//
	//if isDuplicate {
	//	return nil, errors.New(duplicateMovieError)
	//}

	movie, err := m.bookingDao.InsertBooking(ctx, getBookingModel(request))
	if err != nil {
		return nil, err
	}
	//m.movieCache.DeleteMovieCacheKey(ctx)
	return movie, nil
}

func (m *BookingServiceImpl) GetBookings(ctx context.Context, request *moviepb.GetBookingRequest) ([]*models.BookMyMovie, error) {

}

func getBookingModel(request *moviepb.BookingRequest) *models.BookMyMovie {
	return &models.BookMyMovie{
		MovieName:   request.GetMovieName(),
		TheatreName: request.GetTheatreName(),
	}
}
