package grpcserver

import (
	"context"
	"errors"
	"github.com/go-grpc-inmemory-service/commons/constants"
	"github.com/go-grpc-inmemory-service/internal/config"
	"github.com/go-grpc-inmemory-service/internal/grpc"
	"github.com/go-grpc-inmemory-service/internal/models"
	"github.com/go-grpc-inmemory-service/internal/service"
	"github.com/go-grpc-inmemory-service/resources/moviepb"
)

type MovieGrpcServer struct {
	moviepb.UnimplementedMoviePlatformServer
	service.MovieService
	service.BookingService
}

func NewGrpcServer(server *grpc.Server) *MovieGrpcServer {
	dbConfig := config.Configurations.DbConfigs
	cacheConfig := config.Configurations
	s := &MovieGrpcServer{
		MovieService:   service.NewMovieImpl(dbConfig, cacheConfig),
		BookingService: service.NewBookingImpl(dbConfig, cacheConfig),
	}
	moviepb.RegisterMoviePlatformServer(server.GrpcServer, s)
	return s
}

func (m *MovieGrpcServer) CreateMovie(ctx context.Context, request *moviepb.MovieRequest) (*moviepb.MovieResponse, error) {
	movieResponse, err := m.MovieService.CreateMovie(ctx, request)
	if err != nil {
		return nil, err
	}
	if movieResponse == nil {
		return &moviepb.MovieResponse{
			Status: moviepb.MovieStatus_FAILED,
		}, nil
	}
	if movieResponse.Id == constants.EmptyString {
		return nil, errors.New("err not captured but nil response received")
	}

	return &moviepb.MovieResponse{
		MovieDetails: buildMovie(movieResponse),
		Status:       moviepb.MovieStatus_CREATED,
	}, nil
}

func (m *MovieGrpcServer) CreateBooking(ctx context.Context, request *moviepb.BookingRequest) (*moviepb.BookingResponse, error) {
	bookingResponse, err := m.BookingService.CreateBooking(ctx, request)
	if err != nil {
		return nil, err
	}
	if bookingResponse == nil {
		return &moviepb.BookingResponse{
			Status: moviepb.BookingStatus_FAILED_BOOKING,
		}, nil
	}
	if bookingResponse.Id == constants.EmptyString {
		return nil, errors.New("err not captured but nil response received")
	}

	return &moviepb.BookingResponse{
		BookingDetails: buildBooking(bookingResponse),
		Status:         moviepb.BookingStatus_CREATED_BOOKING,
	}, nil
}

func (m *MovieGrpcServer) GetMovie(ctx context.Context, request *moviepb.GetMovieRequest) (*moviepb.GetMovieResponse, error) {
	getMovieResponse, err := m.MovieService.GetMovie(ctx, request)
	if err != nil {
		return nil, err
	}
	if getMovieResponse == nil {
		return nil, errors.New("err not captured but nil response received")
	}

	return &moviepb.GetMovieResponse{
		MovieDetails: buildMovie(getMovieResponse),
	}, nil
}

func (m *MovieGrpcServer) GetBookings(ctx context.Context, request *moviepb.GetMovieRequest) (*moviepb.GetMovieResponse, error) {
	getMovieResponse, err := m.BookingService.(ctx, request)
	if err != nil {
		return nil, err
	}
	if getMovieResponse == nil {
		return nil, errors.New("err not captured but nil response received")
	}

	return &moviepb.GetMovieResponse{
		MovieDetails: buildMovie(getMovieResponse),
	}, nil
}

func (m *MovieGrpcServer) GetAllMovies(ctx context.Context, request *moviepb.GetAllMoviesRequest) (*moviepb.GetAllMoviesResponse, error) {
	getAllMovies, err := m.MovieService.GetAllMovies(ctx)
	if err != nil {
		return nil, err
	}
	if getAllMovies == nil {
		return nil, errors.New("err not captured but nil response received")
	}

	var movieList []*moviepb.MovieDetails
	for _, movie := range getAllMovies {
		movieList = append(movieList, buildMovie(movie))
	}
	return &moviepb.GetAllMoviesResponse{
		MovieDetails: movieList,
	}, nil
}

func (m *MovieGrpcServer) UpdateMovie(ctx context.Context, request *moviepb.UpdateMovieRequest) (*moviepb.UpdateMovieResponse, error) {
	updatedMovieResponse, err := m.MovieService.UpdateMovies(ctx, request)
	if err != nil {
		return nil, err
	}

	if updatedMovieResponse == nil {
		return &moviepb.UpdateMovieResponse{
			Status: moviepb.MovieStatus_FAILED,
		}, nil
	}
	return &moviepb.UpdateMovieResponse{
		MovieDetails: buildMovie(updatedMovieResponse),
		Status:       moviepb.MovieStatus_UPDATED,
	}, nil
}

func buildMovie(movie *models.Movie) *moviepb.MovieDetails {
	return &moviepb.MovieDetails{
		Id:          movie.Id,
		MovieName:   movie.Name,
		Genre:       movie.Genre,
		Description: movie.Desc,
		Ratings:     movie.Rating,
	}
}

func buildBooking(booking *models.BookMyMovie) *moviepb.BookingDetails {
	return &moviepb.BookingDetails{
		Id:          booking.Id,
		MovieName:   booking.MovieName,
		TheatreName: booking.TheatreName,
	}
}
