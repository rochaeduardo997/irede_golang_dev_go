package controller_movie

import (
	"database/sql"

	"github.com/google/uuid"
	controller_interfaces "github.com/rochaeduardo997/irede_golang_dev/internal/controller/interfaces"
	model_movie "github.com/rochaeduardo997/irede_golang_dev/internal/model/movie"
)

type ControllerMovie struct{ Db *sql.DB }

func NewControllerMovie(cm *ControllerMovie) (result controller_interfaces.IGenericController[model_movie.Movie], err error) {
	result = cm
	return
}

func (cm *ControllerMovie) Create(m *model_movie.Movie) (result string, err error) {
	query := `
		INSERT INTO movies(id, name, director, duration_in_seconds)
		VALUES(?,?,?,?)
	`
	m.Id = uuid.NewString()
	_, err = cm.Db.Query(query, &m.Id, &m.Name, &m.Director, &m.DurationInSeconds)
	if err != nil {
		return "", err
	}

	return m.Id, nil
}

func (cm *ControllerMovie) FindBy(id string) (result *model_movie.Movie, err error) {
	query := `
		SELECT id, name, director, duration_in_seconds
		FROM movies
		WHERE id = ?
		LIMIT 1
	`
	rows, err := cm.Db.Query(query, &id)
	if err != nil {
		return nil, err
	}
	result = &model_movie.Movie{}
	for rows.Next() {
		rows.Scan(&result.Id, &result.Name, &result.Director, &result.DurationInSeconds)
	}
	err = result.IsValid()
	if err != nil {
		return nil, err
	}

	return
}

func (cm *ControllerMovie) FindAll(page uint16) (result *controller_interfaces.FindAllResponse[model_movie.Movie], err error) {
	query := `
		SELECT id, name, director, duration_in_seconds
		FROM movies
		LIMIT ?
		OFFSET ?
	`
	limit := uint16(10)
	offset := limit * (page - 1)
	rows, err := cm.Db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	result = &controller_interfaces.FindAllResponse[model_movie.Movie]{}
	for rows.Next() {
		var target model_movie.Movie
		rows.Scan(&target.Id, &target.Name, &target.Director, &target.DurationInSeconds)
		err = target.IsValid()
		if err != nil {
			continue
		}
		result.Registers = append(result.Registers, &target)
	}
	result.Total, err = cm.GetTotal()
	if err != nil {
		return nil, err
	}
	result.Page = page
	return
}

func (cm *ControllerMovie) GetTotal() (result uint32, err error) {
	query := `SELECT COUNT(1) FROM movies`
	rows, err := cm.Db.Query(query)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		rows.Scan(&result)
	}
	return
}

func (cm *ControllerMovie) UpdateBy(id string, m *model_movie.Movie) (result bool, err error) {
	query := `
		UPDATE movies
		SET 
			name = ?,
			director = ?,
			duration_in_seconds = ?
		WHERE id = ?; 
	`
	_, err = cm.Db.Query(query, &m.Name, &m.Director, &m.DurationInSeconds, id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (cm *ControllerMovie) DeleteBy(id string) (result bool, err error) {
	query := `DELETE FROM movies WHERE id = ?`
	_, err = cm.Db.Query(query, &id)
	if err != nil {
		return false, err
	}

	return true, nil
}
