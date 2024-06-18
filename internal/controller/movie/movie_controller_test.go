package controller_movie_test

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/joho/godotenv"
	controller_movie "github.com/rochaeduardo997/irede_golang_dev/internal/controller/movie"
	"github.com/rochaeduardo997/irede_golang_dev/internal/infra/database"
	model_movie "github.com/rochaeduardo997/irede_golang_dev/internal/model/movie"
	"github.com/stretchr/testify/assert"
)

func instanceMovie() (result *model_movie.Movie) {
	expected := &model_movie.Movie{
		Id:                "id",
		Name:              "name",
		Director:          "director",
		DurationInSeconds: 3600,
	}
	result, _ = model_movie.NewMovie(expected)
	return
}

func instanceDB() (result *sql.DB) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatal("Error loading .env file, err: ", err)
	}
	result, _ = database.NewDatabaseConnection()
	result.Query("DELETE FROM room_movies")
	result.Query("DELETE FROM rooms")
	result.Query("DELETE FROM movies")
	return
}

func TestInsert(t *testing.T) {
	db := instanceDB()
	controllerMovie, _ := controller_movie.NewControllerMovie(&controller_movie.ControllerMovie{Db: db})
	movie := instanceMovie()
	result, err := controllerMovie.Create(movie)
	assert.Nil(t, err)
	assert.Greater(t, len(result), 10)
}

func TestFindById(t *testing.T) {
	db := instanceDB()
	controllerMovie, _ := controller_movie.NewControllerMovie(&controller_movie.ControllerMovie{Db: db})
	movie := instanceMovie()
	id, _ := controllerMovie.Create(movie)
	result, err := controllerMovie.FindBy(id)
	assert.Nil(t, err)
	assert.Equal(t, id, result.Id)
	assert.Equal(t, movie.Name, result.Name)
	assert.Equal(t, movie.Director, result.Director)
	assert.Equal(t, movie.DurationInSeconds, result.DurationInSeconds)
	assert.Equal(t, movie.DurationInHours(), result.DurationInHours())
}

func TestFindAll(t *testing.T) {
	db := instanceDB()
	controllerMovie, _ := controller_movie.NewControllerMovie(&controller_movie.ControllerMovie{Db: db})
	movie := instanceMovie()
	id, _ := controllerMovie.Create(movie)
	result, err := controllerMovie.FindAll()
	fmt.Println(result)
	assert.Nil(t, err)
	assert.Equal(t, id, result[0].Id)
	assert.Equal(t, movie.Name, result[0].Name)
	assert.Equal(t, movie.Director, result[0].Director)
	assert.Equal(t, movie.DurationInSeconds, result[0].DurationInSeconds)
	assert.Equal(t, movie.DurationInHours(), result[0].DurationInHours())
}

func TestUpdate(t *testing.T) {
	db := instanceDB()
	controllerMovie, _ := controller_movie.NewControllerMovie(&controller_movie.ControllerMovie{Db: db})
	movie := instanceMovie()
	id, _ := controllerMovie.Create(movie)
	movie.Name = "new_name"
	movie.Director = "new_director"
	movie.DurationInSeconds = 50
	result, err := controllerMovie.UpdateBy(id, movie)
	assert.Nil(t, err)
	assert.Equal(t, true, result)
}

func TestDelete(t *testing.T) {
	db := instanceDB()
	controllerMovie, _ := controller_movie.NewControllerMovie(&controller_movie.ControllerMovie{Db: db})
	movie := instanceMovie()
	id, _ := controllerMovie.Create(movie)
	result, err := controllerMovie.DeleteBy(id)
	assert.Nil(t, err)
	assert.Equal(t, true, result)
}
