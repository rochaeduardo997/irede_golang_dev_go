package model_room

import (
	"errors"

	model_movie "github.com/rochaeduardo997/irede_golang_dev/internal/model/movie"
)

type Room struct {
	Id          string
	Number      uint16
	Description string
	Movies      []*model_movie.Movie
}

func NewRoom(r *Room) (result *Room, err error) {
	result = r
	err = result.IsValid()
	if err != nil {
		return nil, err
	}
	for _, movie := range result.Movies {
		err = movie.IsValid()
		if err != nil {
			return nil, err
		}
	}
	return
}

func (r *Room) IsValid() (err error) {
	if r.Number == 0 {
		return errors.New("room number must be provided")
	}
	if r.Description == "" {
		return errors.New("room description must be provided")
	}
	return
}
