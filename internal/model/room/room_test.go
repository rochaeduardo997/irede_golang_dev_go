package model_room_test

import (
	"testing"

	model_movie "github.com/rochaeduardo997/irede_golang_dev/internal/model/movie"
	model_room "github.com/rochaeduardo997/irede_golang_dev/internal/model/room"
	"github.com/stretchr/testify/assert"
)

func TestRoomInstance(t *testing.T) {
	movies := []*model_movie.Movie{{
		Id:                "id",
		Name:              "name",
		Director:          "director",
		DurationInSeconds: 3600,
	}}
	expected := &model_room.Room{
		Id:          "id",
		Number:      200,
		Description: "description",
		Movies:      movies,
	}
	room, err := model_room.NewRoom(expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, room)
}

func TestFailRoomInstanceWithoutNumber(t *testing.T) {
	room, err := model_room.NewRoom(&model_room.Room{
		Id:          "id",
		Description: "description",
	})
	assert.Nil(t, room)
	assert.EqualError(t, err, "room number must be provided")
}

func TestFailRoomInstanceWithoutDescription(t *testing.T) {
	room, err := model_room.NewRoom(&model_room.Room{
		Id:     "id",
		Number: 200,
	})
	assert.Nil(t, room)
	assert.EqualError(t, err, "room description must be provided")
}

func TestFailRoomInstanceWithInvalidMovie(t *testing.T) {
	movies := []*model_movie.Movie{{}}
	expected := &model_room.Room{
		Id:          "id",
		Number:      200,
		Description: "description",
		Movies:      movies,
	}
	room, err := model_room.NewRoom(expected)
	assert.Nil(t, room)
	assert.EqualError(t, err, "movie name must be provided")
}
