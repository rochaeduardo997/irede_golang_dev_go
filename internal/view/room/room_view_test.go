package view_room_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	controller_interfaces "github.com/rochaeduardo997/irede_golang_dev/internal/controller/interfaces"
	controller_movie "github.com/rochaeduardo997/irede_golang_dev/internal/controller/movie"
	controller_room "github.com/rochaeduardo997/irede_golang_dev/internal/controller/room"
	"github.com/rochaeduardo997/irede_golang_dev/internal/infra/database"
	model_movie "github.com/rochaeduardo997/irede_golang_dev/internal/model/movie"
	model_room "github.com/rochaeduardo997/irede_golang_dev/internal/model/room"
	view_room "github.com/rochaeduardo997/irede_golang_dev/internal/view/room"
	http_adapter "github.com/rochaeduardo997/irede_golang_dev/pkg/http"
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
	result, _ = model_room.NewRoom(&model_room.Room{
		Id:          "id",
		Number:      200,
		Description: "description",
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

func instanceControllerMovie(db *sql.DB) (result controller_interfaces.IGenericController[model_movie.Movie]) {
	result, _ = controller_movie.NewControllerMovie(&controller_movie.ControllerMovie{Db: db})
	return
}

func instanceControllerRoom(db *sql.DB, mc controller_interfaces.IGenericController[model_movie.Movie]) (result controller_interfaces.IGenericController[model_room.Room]) {
	result, _ = controller_room.NewControllerRoom(&controller_room.ControllerRoom{Db: db, MovieController: mc})
	return
}

func TestInsert(t *testing.T) {
	db := instanceDB()
	cm := instanceControllerMovie(db)
	cr := instanceControllerRoom(db, cm)
	httpAdapter, handler := http_adapter.NewGorillaMux()
	view_room.NewViewRoom(&view_room.ViewRoom{Db: db, HTTPAdapter: httpAdapter, ControllerRoom: cr, ControllerMovie: cm})

	server := httptest.NewServer(handler)
	defer server.Close()

	movie := instanceMovie()
	movieId, _ := cm.Create(movie)

	movieBody := map[string]any{}
	movieBody["number"] = 300
	movieBody["description"] = "description"
	movieBody["moviesId"] = []string{movieId}
	bodyJSON, _ := json.Marshal(movieBody)
	payload := bytes.NewBuffer(bodyJSON)
	url := fmt.Sprintf("%s/api/v1/rooms", server.URL)
	resp, err := http.Post(url, "application/json", payload)
	if err != nil {
		t.Fatal(err)
	}
	actual, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	room, err := cr.FindBy(string(actual))
	assert.Nil(t, err)
	assert.Equal(t, movieBody["number"], int(room.Number))
	assert.Equal(t, movieBody["description"], room.Description)
	assert.Equal(t, movie, room.Movies[0])
}

func TestFindById(t *testing.T) {
	db := instanceDB()
	cm := instanceControllerMovie(db)
	cr := instanceControllerRoom(db, cm)
	httpAdapter, handler := http_adapter.NewGorillaMux()
	view_room.NewViewRoom(&view_room.ViewRoom{Db: db, HTTPAdapter: httpAdapter, ControllerRoom: cr, ControllerMovie: cm})

	server := httptest.NewServer(handler)
	defer server.Close()

	movie := instanceMovie()
	movieId, _ := cm.Create(movie)
	movie.Id = movieId
	room := instanceRoom()
	room.Movies = append(room.Movies, movie)
	cr.Create(room)

	url := fmt.Sprintf("%s/api/v1/rooms/%s", server.URL, room.Id)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	actual, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	bodyRes := &model_room.Room{}
	json.Unmarshal(actual, &bodyRes)
	assert.Equal(t, room, bodyRes)
}

func TestFindAll(t *testing.T) {
	db := instanceDB()
	cm := instanceControllerMovie(db)
	cr := instanceControllerRoom(db, cm)
	httpAdapter, handler := http_adapter.NewGorillaMux()
	view_room.NewViewRoom(&view_room.ViewRoom{Db: db, HTTPAdapter: httpAdapter, ControllerRoom: cr, ControllerMovie: cm})

	server := httptest.NewServer(handler)
	defer server.Close()

	movie := instanceMovie()
	movieId, _ := cm.Create(movie)
	movie.Id = movieId
	room := instanceRoom()
	room.Movies = append(room.Movies, movie)
	cr.Create(room)

	url := fmt.Sprintf("%s/api/v1/rooms/all/%d", server.URL, 1)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	actual, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	bodyRes := &controller_interfaces.FindAllResponse[model_room.Room]{}
	json.Unmarshal(actual, &bodyRes)
	assert.Equal(t, room, bodyRes.Registers[0])
	assert.Equal(t, 1, int(bodyRes.Total))
	assert.Equal(t, 1, int(bodyRes.Page))
}

func TestUpdate(t *testing.T) {
	db := instanceDB()
	cm := instanceControllerMovie(db)
	cr := instanceControllerRoom(db, cm)
	httpAdapter, handler := http_adapter.NewGorillaMux()
	view_room.NewViewRoom(&view_room.ViewRoom{Db: db, HTTPAdapter: httpAdapter, ControllerRoom: cr, ControllerMovie: cm})

	server := httptest.NewServer(handler)
	defer server.Close()

	movie := instanceMovie()
	movieId, _ := cm.Create(movie)
	movie2 := instanceMovie()
	cm.Create(movie2)
	movie.Id = movieId
	room := instanceRoom()
	room.Movies = append(room.Movies, movie)
	cr.Create(room)

	movieBody := map[string]any{}
	movieBody["number"] = 300
	movieBody["description"] = "new_description"
	movieBody["movie"] = "new_description"
	movieBody["moviesId"] = []string{movie2.Id}
	bodyJSON, _ := json.Marshal(movieBody)
	payload := bytes.NewBuffer(bodyJSON)
	url := fmt.Sprintf("%s/api/v1/rooms/%s", server.URL, room.Id)
	req, err := http.NewRequest(http.MethodPut, url, payload)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	assert.Nil(t, err)

	room, err = cr.FindBy(room.Id)
	assert.Nil(t, err)
	assert.Equal(t, 300, int(room.Number))
	assert.Equal(t, "new_description", room.Description)
	assert.Equal(t, movie2, room.Movies[0])
}

func TestDelete(t *testing.T) {
	db := instanceDB()
	cm := instanceControllerMovie(db)
	cr := instanceControllerRoom(db, cm)
	httpAdapter, handler := http_adapter.NewGorillaMux()
	view_room.NewViewRoom(&view_room.ViewRoom{Db: db, HTTPAdapter: httpAdapter, ControllerRoom: cr, ControllerMovie: cm})

	server := httptest.NewServer(handler)
	defer server.Close()

	movie := instanceMovie()
	movieId, _ := cm.Create(movie)
	movie.Id = movieId
	room := instanceRoom()
	room.Movies = append(room.Movies, movie)
	cr.Create(room)

	url := fmt.Sprintf("%s/api/v1/rooms/%s", server.URL, room.Id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	assert.Nil(t, err)

	rooms, err := cr.FindAll(1)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(rooms.Registers))
}

func TestFailInsert(t *testing.T) {
	db := instanceDB()
	cm := instanceControllerMovie(db)
	cr := instanceControllerRoom(db, cm)
	httpAdapter, handler := http_adapter.NewGorillaMux()
	view_room.NewViewRoom(&view_room.ViewRoom{Db: db, HTTPAdapter: httpAdapter, ControllerRoom: cr, ControllerMovie: cm})

	server := httptest.NewServer(handler)
	defer server.Close()

	movieBody := map[string]any{}
	bodyJSON, _ := json.Marshal(movieBody)
	payload := bytes.NewBuffer(bodyJSON)
	url := fmt.Sprintf("%s/api/v1/rooms", server.URL)
	resp, err := http.Post(url, "application/json", payload)
	if err != nil {
		t.Fatal(err)
	}
	actual, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "room number must be provided\n", string(actual))
}

func TestFailFindByIdWithInvalidId(t *testing.T) {
	db := instanceDB()
	cm := instanceControllerMovie(db)
	cr := instanceControllerRoom(db, cm)
	httpAdapter, handler := http_adapter.NewGorillaMux()
	view_room.NewViewRoom(&view_room.ViewRoom{Db: db, HTTPAdapter: httpAdapter, ControllerRoom: cr, ControllerMovie: cm})

	server := httptest.NewServer(handler)
	defer server.Close()

	url := fmt.Sprintf("%s/api/v1/rooms/%s", server.URL, "1")
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	actual, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "room not found\n", string(actual))
}

func TestFailUpdateWithInvalidId(t *testing.T) {
	db := instanceDB()
	cm := instanceControllerMovie(db)
	cr := instanceControllerRoom(db, cm)
	httpAdapter, handler := http_adapter.NewGorillaMux()
	view_room.NewViewRoom(&view_room.ViewRoom{Db: db, HTTPAdapter: httpAdapter, ControllerRoom: cr, ControllerMovie: cm})

	server := httptest.NewServer(handler)
	defer server.Close()

	movieBody := map[string]any{}
	bodyJSON, _ := json.Marshal(movieBody)
	payload := bytes.NewBuffer(bodyJSON)
	url := fmt.Sprintf("%s/api/v1/rooms/%s", server.URL, "1")
	req, err := http.NewRequest(http.MethodPut, url, payload)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()
	actual, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "room not found\n", string(actual))
}

func TestFailDeleteWithInvalidId(t *testing.T) {
	db := instanceDB()
	cm := instanceControllerMovie(db)
	cr := instanceControllerRoom(db, cm)
	httpAdapter, handler := http_adapter.NewGorillaMux()
	view_room.NewViewRoom(&view_room.ViewRoom{Db: db, HTTPAdapter: httpAdapter, ControllerRoom: cr, ControllerMovie: cm})

	server := httptest.NewServer(handler)
	defer server.Close()

	url := fmt.Sprintf("%s/api/v1/rooms/%s", server.URL, "1")
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()
	actual, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "room not found\n", string(actual))
}
