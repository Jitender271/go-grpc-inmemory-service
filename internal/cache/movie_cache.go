package cache

import (
	"context"
	"errors"
	"github.com/go-grpc-inmemory-service/internal/cache/localcache"
	"github.com/go-grpc-inmemory-service/internal/config"
	"github.com/go-grpc-inmemory-service/internal/dao"
	daomodels "github.com/go-grpc-inmemory-service/internal/dao/dao_models"
	"sync"
)

const (
	all_movie_cache_key = "all_movies"
	workers             = 8
	maxKeyLoadQueueSize = 2000
)

var once sync.Once
var movieCache MovieCacheService

type MovieCacheService interface {
	GetAllConfigs(ctx context.Context, cacheKey string) ([]daomodels.Movies, error)
	DeleteMovieCacheKey(ctx context.Context)
}

type MovieCacheImpl struct {
	cache localcache.CacheService
}

func NewMovieCacheService(cacheConfigs config.AppConfig, movieDao dao.MovieDao) MovieCacheService {
	once.Do(func() {
		cache := localcache.NewCache(&localcache.Config{
			Name:    "movie_cache",
			MaxKeys: 5000,
			Loader: func(ctx context.Context, key interface{}) (interface{}, error) {
				return reloadMovieValues(ctx, movieDao)
			},
			Workers:             workers,
			MaxKeyLoadQueueSize: maxKeyLoadQueueSize,
			TTL:                 cacheConfigs.CacheTTLInSECONDS,
			ReloadAfter:         cacheConfigs.RefreshTimeSeconds,
			AsyncLoad:           false,
			LoadTimeout:         cacheConfigs.CacheLoadTimeMillis,
		})
		movieCache = &MovieCacheImpl{cache: cache}
	})
	return movieCache
}

func reloadMovieValues(ctx context.Context, movieDao dao.MovieDao) (interface{}, error) {
	movie, err := movieDao.GetAllMovies(ctx)
	if err != nil {
		return nil, errors.New("error getting all movies from db")
	}
	if len(movie) == 0 {
		return nil, errors.New("no movies found in db")
	}
	return movie, nil
}

func (m *MovieCacheImpl) GetAllConfigs(ctx context.Context, cacheKey string) ([]daomodels.Movies, error) {
	movies, err := m.cache.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}
	return movies.([]daomodels.Movies), nil
}

func (m *MovieCacheImpl) DeleteMovieCacheKey(ctx context.Context) {
	m.cache.Delete(ctx, all_movie_cache_key)
}
