package controller_room

import (
	"database/sql"

	"github.com/google/uuid"
	controller_interfaces "github.com/rochaeduardo997/irede_golang_dev/internal/controller/interfaces"
	model_movie "github.com/rochaeduardo997/irede_golang_dev/internal/model/movie"
	model_room "github.com/rochaeduardo997/irede_golang_dev/internal/model/room"
)

type ControllerRoom struct {
	Db              *sql.DB
	MovieController controller_interfaces.IGenericController[model_movie.Movie]
}

func NewControllerRoom(vm *ControllerRoom) (result controller_interfaces.IGenericController[model_room.Room], err error) {
	result = vm
	return
}

func (cm *ControllerRoom) Create(r *model_room.Room) (result string, err error) {
	tx, err := cm.Db.Begin()
	if err != nil {
		return "", err
	}
	query := `
		INSERT INTO rooms(id, number, description)
		VALUES(?,?,?)
	`
	r.Id = uuid.NewString()
	_, err = tx.Query(query, &r.Id, &r.Number, &r.Description)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	err = cm.InsertRoomMovies(r.Id, r.Movies, tx)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	tx.Commit()

	return r.Id, nil
}

func (cm *ControllerRoom) InsertRoomMovies(roomId string, ms []*model_movie.Movie, tx *sql.Tx) (err error) {
	query := `
		INSERT INTO room_movies(fk_room_id, fk_movie_id)
		VALUES(?,?)
	`
	for _, movie := range ms {
		_, err = tx.Query(query, &roomId, &movie.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cm *ControllerRoom) FindBy(id string) (result *model_room.Room, err error) {
	query := `
		SELECT id, number, description
		FROM rooms
		WHERE id = ?
		LIMIT 1
	`
	rows, err := cm.Db.Query(query, &id)
	if err != nil {
		return nil, err
	}
	result = &model_room.Room{}
	for rows.Next() {
		rows.Scan(&result.Id, &result.Number, &result.Description)
	}
	result.Movies = cm.GetAssociatedMoviesBy(result.Id)
	err = result.IsValid()
	if err != nil {
		return nil, err
	}

	return
}

func (cm *ControllerRoom) GetAssociatedMoviesBy(roomId string) (result []*model_movie.Movie) {
	query := `
		SELECT fk_movie_id
		FROM room_movies
		WHERE fk_room_id = ?
	`
	rows, err := cm.Db.Query(query, &roomId)
	if err != nil {
		return nil
	}
	result = []*model_movie.Movie{}
	for rows.Next() {
		var movieId *string
		rows.Scan(&movieId)
		movie, err := cm.MovieController.FindBy(*movieId)
		if err != nil {
			continue
		}
		result = append(result, movie)
	}

	return
}

func (cm *ControllerRoom) FindAll(page uint16) (result *controller_interfaces.FindAllResponse[model_room.Room], err error) {
	query := `
		SELECT id, number, description
		FROM rooms
		LIMIT ?
		OFFSET ?
	`
	limit := uint16(10)
	offset := limit * (page - 1)
	rows, err := cm.Db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	result = &controller_interfaces.FindAllResponse[model_room.Room]{}
	for rows.Next() {
		var target model_room.Room
		rows.Scan(&target.Id, &target.Number, &target.Description)
		target.Movies = cm.GetAssociatedMoviesBy(target.Id)
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

func (cm *ControllerRoom) GetTotal() (result uint32, err error) {
	query := `SELECT COUNT(1) FROM rooms`
	rows, err := cm.Db.Query(query)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		rows.Scan(&result)
	}
	return
}

func (cm *ControllerRoom) UpdateBy(id string, m *model_room.Room) (result bool, err error) {
	updateQuery := `
		UPDATE rooms
		SET
			number = ?,
			description = ?
		WHERE id = ?;
	`
	tx, err := cm.Db.Begin()
	if err != nil {
		return false, err
	}
	_, err = tx.Query(updateQuery, &m.Number, &m.Description, id)
	if err != nil {
		return false, err
	}
	deleteAllRoomMoviesQuery := `DELETE FROM room_movies WHERE fk_room_id = ?`
	_, err = tx.Query(deleteAllRoomMoviesQuery, &id)
	if err != nil {
		return false, err
	}
	err = cm.InsertRoomMovies(id, m.Movies, tx)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (cm *ControllerRoom) DeleteBy(id string) (result bool, err error) {
	query := `DELETE FROM rooms WHERE id = ?`
	_, err = cm.Db.Query(query, &id)
	if err != nil {
		return false, err
	}

	return true, nil
}
