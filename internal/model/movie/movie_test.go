package model_movie_test

import (
	"testing"

	model_movie "github.com/rochaeduardo997/irede_golang_dev/internal/model/movie"
	"github.com/stretchr/testify/assert"
)

func TestMovieInstance(t *testing.T) {
	expected := &model_movie.Movie{
		Id:                "id",
		Name:              "name",
		Director:          "director",
		DurationInSeconds: 3600,
	}
	movie, err := model_movie.NewMovie(expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, movie)
}

func TestMovieDurationInHoursFunction(t *testing.T) {
	movie := &model_movie.Movie{
		Id:                "id",
		Name:              "name",
		Director:          "director",
		DurationInSeconds: 3600,
	}
	expected := "01:00:00"
	assert.Equal(t, expected, movie.DurationInHours())
}

func TestFailMovieInstanceWithoutName(t *testing.T) {
	movie, err := model_movie.NewMovie(&model_movie.Movie{
		Id:                "id",
		Director:          "director",
		DurationInSeconds: 3600,
	})
	assert.Nil(t, movie)
	assert.EqualError(t, err, "movie name must be provided")
}

func TestFailMovieInstanceWithoutDirector(t *testing.T) {
	movie, err := model_movie.NewMovie(&model_movie.Movie{
		Id:                "id",
		Name:              "name",
		DurationInSeconds: 3600,
	})
	assert.Nil(t, movie)
	assert.EqualError(t, err, "movie director must be provided")
}

func TestFailMovieInstanceWithoutDuration(t *testing.T) {
	movie, err := model_movie.NewMovie(&model_movie.Movie{
		Id:       "id",
		Name:     "name",
		Director: "director",
	})
	assert.Nil(t, movie)
	assert.EqualError(t, err, "movie duration must be provided")
}
