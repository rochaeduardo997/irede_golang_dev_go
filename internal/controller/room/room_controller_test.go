package controller_room_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/joho/godotenv"
	controller_movie "github.com/rochaeduardo997/irede_golang_dev/internal/controller/movie"
	controller_room "github.com/rochaeduardo997/irede_golang_dev/internal/controller/room"
	"github.com/rochaeduardo997/irede_golang_dev/internal/infra/database"
	model_movie "github.com/rochaeduardo997/irede_golang_dev/internal/model/movie"
	model_room "github.com/rochaeduardo997/irede_golang_dev/internal/model/room"
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

func instanceRoom() (result *model_room.Room) {
	movie := instanceMovie()
	result, _ = model_room.NewRoom(&model_room.Room{
		Id:          "id",
		Number:      200,
		Description: "description",
		Movies:      []*model_movie.Movie{movie},
	})
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
	controllerRoom, _ := controller_room.NewControllerRoom(&controller_room.ControllerRoom{Db: db, MovieController: controllerMovie})
	room := instanceRoom()
	controllerMovie.Create(room.Movies[0])
	result, err := controllerRoom.Create(room)
	assert.Nil(t, err)
	assert.Greater(t, len(result), 10)
}

func TestFindById(t *testing.T) {
	db := instanceDB()
	controllerMovie, _ := controller_movie.NewControllerMovie(&controller_movie.ControllerMovie{Db: db})
	controllerRoom, _ := controller_room.NewControllerRoom(&controller_room.ControllerRoom{Db: db, MovieController: controllerMovie})
	room := instanceRoom()
	controllerMovie.Create(room.Movies[0])
	id, _ := controllerRoom.Create(room)
	result, err := controllerRoom.FindBy(id)
	assert.Nil(t, err)
	assert.Equal(t, id, result.Id)
	assert.Equal(t, room.Number, result.Number)
	assert.Equal(t, room.Description, result.Description)
	assert.Equal(t, room.Movies[0].Id, result.Movies[0].Id)
}

func TestFindAll(t *testing.T) {
	db := instanceDB()
	controllerMovie, _ := controller_movie.NewControllerMovie(&controller_movie.ControllerMovie{Db: db})
	controllerRoom, _ := controller_room.NewControllerRoom(&controller_room.ControllerRoom{Db: db, MovieController: controllerMovie})
	room := instanceRoom()
	controllerMovie.Create(room.Movies[0])
	id, _ := controllerRoom.Create(room)
	result, err := controllerRoom.FindAll(1)
	assert.Nil(t, err)
	assert.Equal(t, uint32(1), result.Total)
	assert.Equal(t, uint16(1), result.Page)
	assert.Equal(t, id, result.Registers[0].Id)
	assert.Equal(t, room.Number, result.Registers[0].Number)
	assert.Equal(t, room.Description, result.Registers[0].Description)
	assert.Equal(t, room.Movies[0].Id, result.Registers[0].Movies[0].Id)
}

func TestUpdate(t *testing.T) {
	db := instanceDB()
	controllerMovie, _ := controller_movie.NewControllerMovie(&controller_movie.ControllerMovie{Db: db})
	controllerRoom, _ := controller_room.NewControllerRoom(&controller_room.ControllerRoom{Db: db, MovieController: controllerMovie})
	room := instanceRoom()
	controllerMovie.Create(room.Movies[0])
	id, _ := controllerRoom.Create(room)
	room.Number = 300
	room.Description = "new_description"
	result, err := controllerRoom.UpdateBy(id, room)
	assert.Nil(t, err)
	assert.Equal(t, true, result)
}

func TestDelete(t *testing.T) {
	db := instanceDB()
	controllerMovie, _ := controller_movie.NewControllerMovie(&controller_movie.ControllerMovie{Db: db})
	controllerRoom, _ := controller_room.NewControllerRoom(&controller_room.ControllerRoom{Db: db, MovieController: controllerMovie})
	room := instanceRoom()
	controllerMovie.Create(room.Movies[0])
	controllerRoom.Create(room)
	result, err := controllerRoom.DeleteBy(room.Id)
	assert.Nil(t, err)
	assert.Equal(t, true, result)
}
