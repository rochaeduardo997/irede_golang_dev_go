package model_movie

import (
	"errors"
	"time"
)

type Movie struct {
	Id                string
	Name              string
	Director          string
	DurationInSeconds uint16
}

func NewMovie(m *Movie) (result *Movie, err error) {
	result = m
	err = result.IsValid()
	if err != nil {
		return nil, err
	}
	return
}

func (m *Movie) IsValid() (err error) {
	if m.Name == "" {
		return errors.New("movie name must be provided")
	}
	if m.Director == "" {
		return errors.New("movie director must be provided")
	}
	if m.DurationInSeconds == 0 {
		return errors.New("movie duration must be provided")
	}
	return
}

func (m *Movie) DurationInHours() (result string) {
	var t time.Time
	t = t.Add(time.Duration(m.DurationInSeconds) * time.Second)
	return t.Format("15:00:00")
}
