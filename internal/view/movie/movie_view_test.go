package view_movie_test

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
	"github.com/rochaeduardo997/irede_golang_dev/internal/infra/database"
	model_movie "github.com/rochaeduardo997/irede_golang_dev/internal/model/movie"
	view_movie "github.com/rochaeduardo997/irede_golang_dev/internal/view/movie"
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

func TestInsert(t *testing.T) {
	db := instanceDB()
	cm := instanceControllerMovie(db)
	httpAdapter, handler := http_adapter.NewGorillaMux()
	view_movie.NewViewMovie(&view_movie.ViewMovie{Db: db, HTTPAdapter: httpAdapter, ControllerMovie: cm})

	server := httptest.NewServer(handler)
	defer server.Close()

	movieBody := map[string]any{}
	movieBody["name"] = "name"
	movieBody["director"] = "director"
	movieBody["durationInSeconds"] = 3600
	bodyJSON, _ := json.Marshal(movieBody)
	payload := bytes.NewBuffer(bodyJSON)
	url := fmt.Sprintf("%s/api/v1/movies", server.URL)
	resp, err := http.Post(url, "application/json", payload)
	if err != nil {
		t.Fatal(err)
	}
	actual, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	movie, err := cm.FindBy(string(actual))
	assert.Nil(t, err)
	assert.Equal(t, movieBody["name"], movie.Name)
	assert.Equal(t, movieBody["director"], movie.Director)
	assert.Equal(t, movieBody["durationInSeconds"], int(movie.DurationInSeconds))
}

func TestFindById(t *testing.T) {
	db := instanceDB()
	cm := instanceControllerMovie(db)
	httpAdapter, handler := http_adapter.NewGorillaMux()
	view_movie.NewViewMovie(&view_movie.ViewMovie{Db: db, HTTPAdapter: httpAdapter, ControllerMovie: cm})

	movie := instanceMovie()
	id, _ := cm.Create(movie)

	server := httptest.NewServer(handler)
	defer server.Close()

	url := fmt.Sprintf("%s/api/v1/movies/%s", server.URL, id)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	actual, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	bodyRes := &model_movie.Movie{}
	json.Unmarshal(actual, &bodyRes)
	assert.Equal(t, movie, bodyRes)
}

func TestFindAll(t *testing.T) {
	db := instanceDB()
	cm := instanceControllerMovie(db)
	httpAdapter, handler := http_adapter.NewGorillaMux()
	view_movie.NewViewMovie(&view_movie.ViewMovie{Db: db, HTTPAdapter: httpAdapter, ControllerMovie: cm})

	movie := instanceMovie()
	cm.Create(movie)

	server := httptest.NewServer(handler)
	defer server.Close()

	url := fmt.Sprintf("%s/api/v1/movies/all/%d", server.URL, 1)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	actual, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	bodyRes := &controller_interfaces.FindAllResponse[model_movie.Movie]{}
	json.Unmarshal(actual, &bodyRes)
	assert.Equal(t, movie, bodyRes.Registers[0])
	assert.Equal(t, 1, int(bodyRes.Total))
	assert.Equal(t, 1, int(bodyRes.Page))
}

func TestUpdate(t *testing.T) {
	db := instanceDB()
	cm := instanceControllerMovie(db)
	httpAdapter, handler := http_adapter.NewGorillaMux()
	view_movie.NewViewMovie(&view_movie.ViewMovie{Db: db, HTTPAdapter: httpAdapter, ControllerMovie: cm})

	movie := instanceMovie()
	id, _ := cm.Create(movie)

	server := httptest.NewServer(handler)
	defer server.Close()

	movieBody := map[string]any{}
	movieBody["name"] = "new_name"
	movieBody["director"] = "new_director"
	movieBody["durationInSeconds"] = 50
	bodyJSON, _ := json.Marshal(movieBody)
	payload := bytes.NewBuffer(bodyJSON)
	url := fmt.Sprintf("%s/api/v1/movies/%s", server.URL, id)
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

	movie, err = cm.FindBy(string(id))
	assert.Nil(t, err)
	assert.Equal(t, "new_name", movie.Name)
	assert.Equal(t, "new_director", movie.Director)
	assert.Equal(t, 50, int(movie.DurationInSeconds))
}

func TestDelete(t *testing.T) {
	db := instanceDB()
	cm := instanceControllerMovie(db)
	httpAdapter, handler := http_adapter.NewGorillaMux()
	view_movie.NewViewMovie(&view_movie.ViewMovie{Db: db, HTTPAdapter: httpAdapter, ControllerMovie: cm})

	movie := instanceMovie()
	id, _ := cm.Create(movie)

	server := httptest.NewServer(handler)
	defer server.Close()

	url := fmt.Sprintf("%s/api/v1/movies/%s", server.URL, id)
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

	movies, err := cm.FindAll(1)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(movies.Registers))
}

func TestFailInsert(t *testing.T) {
	db := instanceDB()
	cm := instanceControllerMovie(db)
	httpAdapter, handler := http_adapter.NewGorillaMux()
	view_movie.NewViewMovie(&view_movie.ViewMovie{Db: db, HTTPAdapter: httpAdapter, ControllerMovie: cm})

	server := httptest.NewServer(handler)
	defer server.Close()

	movieBody := map[string]any{}
	bodyJSON, _ := json.Marshal(movieBody)
	payload := bytes.NewBuffer(bodyJSON)
	url := fmt.Sprintf("%s/api/v1/movies", server.URL)
	resp, err := http.Post(url, "application/json", payload)
	if err != nil {
		t.Fatal(err)
	}
	actual, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "movie name must be provided\n", string(actual))
}

func TestFailFindByIdWithInvalidId(t *testing.T) {
	db := instanceDB()
	cm := instanceControllerMovie(db)
	httpAdapter, handler := http_adapter.NewGorillaMux()
	view_movie.NewViewMovie(&view_movie.ViewMovie{Db: db, HTTPAdapter: httpAdapter, ControllerMovie: cm})

	movie := instanceMovie()
	cm.Create(movie)

	server := httptest.NewServer(handler)
	defer server.Close()

	url := fmt.Sprintf("%s/api/v1/movies/%s", server.URL, "1")
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	actual, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "movie not found\n", string(actual))
}

func TestFailUpdateWithInvalidId(t *testing.T) {
	db := instanceDB()
	cm := instanceControllerMovie(db)
	httpAdapter, handler := http_adapter.NewGorillaMux()
	view_movie.NewViewMovie(&view_movie.ViewMovie{Db: db, HTTPAdapter: httpAdapter, ControllerMovie: cm})

	movie := instanceMovie()
	cm.Create(movie)

	server := httptest.NewServer(handler)
	defer server.Close()

	movieBody := map[string]any{}
	movieBody["name"] = "new_name"
	movieBody["director"] = "new_director"
	movieBody["durationInSeconds"] = 50
	bodyJSON, _ := json.Marshal(movieBody)
	payload := bytes.NewBuffer(bodyJSON)
	url := fmt.Sprintf("%s/api/v1/movies/%s", server.URL, "1")
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
	assert.Equal(t, "movie not found\n", string(actual))
}

func TestFailDeleteWithInvalidId(t *testing.T) {
	db := instanceDB()
	cm := instanceControllerMovie(db)
	httpAdapter, handler := http_adapter.NewGorillaMux()
	view_movie.NewViewMovie(&view_movie.ViewMovie{Db: db, HTTPAdapter: httpAdapter, ControllerMovie: cm})

	movie := instanceMovie()
	cm.Create(movie)

	server := httptest.NewServer(handler)
	defer server.Close()

	url := fmt.Sprintf("%s/api/v1/movies/%s", server.URL, "1")
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
	assert.Equal(t, "movie not found\n", string(actual))
}
