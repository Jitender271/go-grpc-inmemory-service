package dao

import (
	"context"
	"errors"
	"github.com/go-grpc-inmemory-service/commons/constants"
	"github.com/go-grpc-inmemory-service/internal/log"
	"go.uber.org/zap"
	"strings"

	daomodels "github.com/go-grpc-inmemory-service/internal/dao/dao_models"
	"github.com/go-grpc-inmemory-service/internal/db"
	"github.com/scylladb/gocqlx/v2/table"
)

var (
	movieColumns         = []string{"movie_id", "name", "genre", "description", "rating"}
	moviePartitionKeys   = []string{"name"}
	updateMovieColumns   = []string{"genre", "description", "rating"}
	movieSortKeys        []string
	bookingColumns       = []string{"id", "movie_name", "theatre_name", "start_time", "end_time"}
	singlePartitionKeys  = []string{"movie_name"}
	allPartitionKeys     = []string{"movie_name", "theatre_name"}
	updateBookingColumns = []string{"theatre_name", "start_time", "end_time"}
	bookingSortKeys      []string
)

const (
	movieTableName   = "movies"
	bookingTableName = "book_my_movie"
	name             = "name"
	id               = "id"
	movieName        = "movie_name"
)

func getTable(tableName string, cols, partitionKeys, sortKeys []string) db.Table {
	tableMeta := table.Metadata{
		Name:    tableName,
		Columns: cols,
		PartKey: partitionKeys,
		SortKey: sortKeys,
	}
	return table.New(tableMeta)
}

func insertMoviesInDb(ctx context.Context, session db.SessionWrapperService, movies *daomodels.Movies) error {
	configTable := getTable(movieTableName, movieColumns, moviePartitionKeys, movieSortKeys)
	query := session.Query(configTable.Insert()).BindStruct(movies)
	if err := query.Exec(ctx); err != nil {
		log.Logger.Error("error inserting movie name in db ", zap.String("key", movies.Name))
		return err
	}
	log.Logger.Info("Movie insertion successful", zap.String(id, movies.MovieID))
	return nil
}

func getMovieFromDb(ctx context.Context, session db.SessionWrapperService, movieName string) (*daomodels.Movies, error) {
	var movies daomodels.Movies
	configTable := getTable(movieTableName, movieColumns, moviePartitionKeys, movieSortKeys)
	query := session.Query(configTable.Get()).BindMap(map[string]interface{}{name: movieName})
	err := query.GetRelease(ctx, &movies)
	if err != nil && strings.EqualFold(err.Error(), "not found") {
		log.Logger.Error("movie name not found in db ", zap.String("movie_name", movieName))
		return nil, err
	} else if err != nil {
		log.Logger.Error("error fetching movie details from db ", zap.Error(err))
		return nil, err
	}
	return &movies, nil
}

func getAllMoviesFromDb(ctx context.Context, session db.SessionWrapperService) ([]daomodels.Movies, error) {
	var movies []daomodels.Movies
	configTable := getTable(movieTableName, movieColumns, moviePartitionKeys, movieSortKeys)
	query := session.Query(configTable.SelectAll())
	err := query.SelectRelease(ctx, &movies)
	if err != nil {
		log.Logger.Error("error fetching all movie details from db ", zap.Error(err))
		return nil, err
	}
	return movies, nil
}

func updateMoviesInDb(ctx context.Context, session db.SessionWrapperService, movies *daomodels.Movies) error {
	configTable := getTable(movieTableName, movieColumns, moviePartitionKeys, movieSortKeys)
	query := session.Query(configTable.Update(updateMovieColumns...)).BindStruct(movies)
	if err := query.Exec(ctx); err != nil {
		log.Logger.Error("error updating movie", zap.String("key", movies.Name))
		return err
	}
	log.Logger.Info("Movie updated successful", zap.String(id, movies.MovieID))
	return nil
}

func insertBookingInDb(ctx context.Context, session db.SessionWrapperService, booking *daomodels.BookMyMovie) error {

	bookingTable := getTable(bookingTableName, bookingColumns, singlePartitionKeys, bookingSortKeys)
	query := session.Query(bookingTable.Insert()).BindStruct(booking)
	if err := query.Exec(ctx); err != nil {
		log.Logger.Error("error inserting movie name in db ", zap.String("key", booking.MovieName))
		return err
	}
	log.Logger.Info("Movie insertion successful", zap.String(id, booking.Id))
	return nil
}

func getBookingFromDb(ctx context.Context, session db.SessionWrapperService, movieName, theaterName string) ([]daomodels.BookMyMovie, error) {
	var query db.QueryXService
	var bookingDetails []daomodels.BookMyMovie

	bookingQueryObject := createGetBookingQueryObject(movieName, theaterName)
	partitionKeys := singlePartitionKeys
	if theaterName == constants.EmptyString {
		partitionKeys = allPartitionKeys
	}

	bookingTable := getTable(bookingTableName, bookingColumns, partitionKeys, bookingSortKeys)
	query = session.Query(bookingTable.Select()).BindStruct(bookingQueryObject)
	if err := query.SelectRelease(ctx, &bookingDetails); err != nil {
		log.Logger.Error("error querying movie name in db ", zap.String("movie_name", movieName), zap.String("theatre_name", theaterName), zap.Error(err))
		return nil, errors.New("error querying booking, movie: " + movieName)
	}

	if len(bookingDetails) == 0 {
		log.Logger.Info("no bookings found", zap.String("movie_name", movieName))
		return bookingDetails, nil
	}

	return bookingDetails, nil
}

func createGetBookingQueryObject(movieName, theaterName string) *daomodels.GetBookingQueryObject {
	return &daomodels.GetBookingQueryObject{
		MovieName:   movieName,
		TheatreName: theaterName,
	}
}
