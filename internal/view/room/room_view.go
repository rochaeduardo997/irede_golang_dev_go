package view_room

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	controller_interfaces "github.com/rochaeduardo997/irede_golang_dev/internal/controller/interfaces"
	model_movie "github.com/rochaeduardo997/irede_golang_dev/internal/model/movie"
	model_room "github.com/rochaeduardo997/irede_golang_dev/internal/model/room"
	http_adapter "github.com/rochaeduardo997/irede_golang_dev/pkg/http"
)

type InputRoomReq struct {
	Number      uint16
	Description string
	MoviesId    []string
}

type FindAll struct {
	Total     uint32
	Page      uint16
	Registers []*model_room.Room
}

type ViewRoom struct {
	Db              *sql.DB
	HTTPAdapter     http_adapter.IHTTP
	ControllerRoom  controller_interfaces.IGenericController[model_room.Room]
	ControllerMovie controller_interfaces.IGenericController[model_movie.Movie]
}

func NewViewRoom(rm *ViewRoom) (result *ViewRoom) {
	result = rm

	result.HTTPAdapter.AddRoute("post", "/api/v1/rooms", rm.CreateHandler)
	result.HTTPAdapter.AddRoute("get", "/api/v1/rooms/{id}", rm.FindByIdHandler)
	result.HTTPAdapter.AddRoute("get", "/api/v1/rooms/all/{page}", rm.FindAllHandler)
	result.HTTPAdapter.AddRoute("put", "/api/v1/rooms/{id}", rm.UpdateByIdHandler)
	result.HTTPAdapter.AddRoute("delete", "/api/v1/rooms/{id}", rm.DeleteByIdHandler)

	return
}

// @Summary      Create a movie
// @Tags         Rooms
// @Param        data body InputRoomReq true "body"
// @Success      201  {object}  model_room.Room
// @Router       /rooms [post]
func (rm *ViewRoom) CreateHandler(w http.ResponseWriter, r *http.Request) {
	input := &InputRoomReq{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	room := &model_room.Room{Number: input.Number, Description: input.Description}
	for _, movieId := range input.MoviesId {
		movie, err := rm.ControllerMovie.FindBy(movieId)
		if err != nil {
			continue
		}
		room.Movies = append(room.Movies, movie)
	}
	err = room.IsValid()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := rm.ControllerRoom.Create(room)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(result))
}

// @Summary      Get room by id
// @Tags         Rooms
// @Param        id   path      string true  "Room ID"
// @Success      200  {object} model_room.Room
// @Router       /rooms/{id} [get]
func (rm *ViewRoom) FindByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, "id must be provided", http.StatusBadRequest)
		return
	}
	result, err := rm.ControllerRoom.FindBy(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res := map[string]any{}
	res["id"] = result.Id
	res["number"] = result.Number
	res["description"] = result.Description
	roomMovies := []any{}
	for _, movie := range result.Movies {
		target := map[string]any{}
		target["id"] = movie.Id
		target["name"] = movie.Name
		target["director"] = movie.Director
		target["durationInSeconds"] = movie.DurationInSeconds
		target["durationInHours"] = movie.DurationInHours()
		roomMovies = append(roomMovies, target)
	}
	res["movies"] = roomMovies
	resJSON, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %s\n", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

// @Summary      Get all rooms
// @Tags         Rooms
// @Param        page   path      string true  "Page"
// @Success      200  {object}    FindAll
// @Router       /rooms/all/{page} [get]
func (rm *ViewRoom) FindAllHandler(w http.ResponseWriter, r *http.Request) {
	page := mux.Vars(r)["page"]
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		http.Error(w, "page must be provided", http.StatusBadRequest)
		return
	}
	result, err := rm.ControllerRoom.FindAll(uint16(pageInt))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	registers := []map[string]any{}
	for _, room := range result.Registers {
		target := map[string]any{}
		target["id"] = room.Id
		target["number"] = room.Number
		target["description"] = room.Description
		roomMovies := []any{}
		for _, movie := range room.Movies {
			targetMovie := map[string]any{}
			targetMovie["id"] = movie.Id
			targetMovie["name"] = movie.Name
			targetMovie["director"] = movie.Director
			targetMovie["durationInSeconds"] = movie.DurationInSeconds
			targetMovie["durationInHours"] = movie.DurationInHours()
			roomMovies = append(roomMovies, targetMovie)
		}
		target["movies"] = roomMovies
		registers = append(registers, target)
	}
	res := map[string]any{}
	res["total"] = result.Total
	res["page"] = result.Page
	res["registers"] = registers
	resJSON, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %s\n", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

// @Summary      Update room by id
// @Tags         Rooms
// @Param        id   path      string true  "Room ID"
// @Param        data body InputRoomReq true "body"
// @Success      200  {boolean} boolean true
// @Router       /rooms/{id} [put]
func (rm *ViewRoom) UpdateByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	input := &InputRoomReq{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	room := &model_room.Room{Number: input.Number, Description: input.Description}
	for _, movieId := range input.MoviesId {
		movie, err := rm.ControllerMovie.FindBy(movieId)
		if err != nil {
			continue
		}
		room.Movies = append(room.Movies, movie)
	}
	result, err := rm.ControllerRoom.UpdateBy(id, room)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res := strconv.FormatBool(result)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

// @Summary      Delete a room by id
// @Tags         Rooms
// @Param        id   path      string true  "Room ID"
// @Success      200  {boolean} boolean true
// @Router       /rooms/{id} [delete]
func (rm *ViewRoom) DeleteByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, "id must be provided", http.StatusBadRequest)
		return
	}
	result, err := rm.ControllerRoom.DeleteBy(id)
	fmt.Println(result, err)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res := strconv.FormatBool(result)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}
