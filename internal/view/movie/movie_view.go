package view_movie

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	controller_interfaces "github.com/rochaeduardo997/irede_golang_dev/internal/controller/interfaces"
	model_movie "github.com/rochaeduardo997/irede_golang_dev/internal/model/movie"
	http_adapter "github.com/rochaeduardo997/irede_golang_dev/pkg/http"
)

type ViewMovie struct {
	Db              *sql.DB
	HTTPAdapter     http_adapter.IHTTP
	ControllerMovie controller_interfaces.IGenericController[model_movie.Movie]
}

func NewViewMovie(vm *ViewMovie) (result *ViewMovie) {
	result = vm

	result.HTTPAdapter.AddRoute("post", "/api/v1/movies", vm.CreateHandler)
	result.HTTPAdapter.AddRoute("get", "/api/v1/movies/{id}", vm.FindByIdHandler)
	result.HTTPAdapter.AddRoute("get", "/api/v1/movies/all/{page}", vm.FindAllHandler)
	result.HTTPAdapter.AddRoute("put", "/api/v1/movies/{id}", vm.UpdateByIdHandler)
	result.HTTPAdapter.AddRoute("delete", "/api/v1/movies/{id}", vm.DeleteByIdHandler)

	return
}

func (vm *ViewMovie) CreateHandler(w http.ResponseWriter, r *http.Request) {
	movie := &model_movie.Movie{}
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = movie.IsValid()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := vm.ControllerMovie.Create(movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

func (vm *ViewMovie) FindByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, "id must be provided", http.StatusBadRequest)
		return
	}
	result, err := vm.ControllerMovie.FindBy(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res := map[string]any{}
	res["id"] = result.Id
	res["name"] = result.Name
	res["director"] = result.Director
	res["durationInSeconds"] = result.DurationInSeconds
	res["durationInHours"] = result.DurationInHours()
	resJSON, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %s\n", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func (vm *ViewMovie) FindAllHandler(w http.ResponseWriter, r *http.Request) {
	page := mux.Vars(r)["page"]
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		http.Error(w, "page must be provided", http.StatusBadRequest)
		return
	}
	result, err := vm.ControllerMovie.FindAll(uint16(pageInt))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	registers := []map[string]any{}
	for _, movie := range result.Registers {
		target := map[string]any{}
		target["id"] = movie.Id
		target["name"] = movie.Name
		target["director"] = movie.Director
		target["durationInSeconds"] = movie.DurationInSeconds
		target["durationInHours"] = movie.DurationInHours()
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

func (vm *ViewMovie) UpdateByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	movie := &model_movie.Movie{}
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	movie.Id = id
	err = movie.IsValid()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := vm.ControllerMovie.UpdateBy(id, movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res := strconv.FormatBool(result)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func (vm *ViewMovie) DeleteByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, "id must be provided", http.StatusBadRequest)
		return
	}
	result, err := vm.ControllerMovie.DeleteBy(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res := strconv.FormatBool(result)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}
