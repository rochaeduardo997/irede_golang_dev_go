package main

import (
	"database/sql"
	"log"

	"github.com/joho/godotenv"
	controller_interfaces "github.com/rochaeduardo997/irede_golang_dev/internal/controller/interfaces"
	controller_movie "github.com/rochaeduardo997/irede_golang_dev/internal/controller/movie"
	controller_room "github.com/rochaeduardo997/irede_golang_dev/internal/controller/room"
	_ "github.com/rochaeduardo997/irede_golang_dev/internal/docs"
	"github.com/rochaeduardo997/irede_golang_dev/internal/infra/database"
	model_movie "github.com/rochaeduardo997/irede_golang_dev/internal/model/movie"
	model_room "github.com/rochaeduardo997/irede_golang_dev/internal/model/room"
	view_docs "github.com/rochaeduardo997/irede_golang_dev/internal/view/docs"
	view_movie "github.com/rochaeduardo997/irede_golang_dev/internal/view/movie"
	view_room "github.com/rochaeduardo997/irede_golang_dev/internal/view/room"
	http_adapter "github.com/rochaeduardo997/irede_golang_dev/pkg/http"
)

// @title           Movies
// @version         1.0
// @description
// @termsOfService

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file, err: ", err)
	}

	db := instanceDB()

	cm := instanceControllerMovie(db)
	cr := instanceControllerRoom(db, cm)

	httpAdapter, _ := http_adapter.NewGorillaMux()
	view_docs.NewDocsView(&view_docs.DocsView{HTTPAdapter: httpAdapter})
	view_movie.NewViewMovie(&view_movie.ViewMovie{Db: db, HTTPAdapter: httpAdapter, ControllerMovie: cm})
	view_room.NewViewRoom(&view_room.ViewRoom{Db: db, HTTPAdapter: httpAdapter, ControllerRoom: cr, ControllerMovie: cm})

	httpAdapter.Listen()
}

func instanceDB() (result *sql.DB) {
	result, _ = database.NewDatabaseConnection()
	return
}

func instanceControllerMovie(db *sql.DB) (result controller_interfaces.IGenericController[model_movie.Movie]) {
	result, _ = controller_movie.NewControllerMovie(&controller_movie.ControllerMovie{Db: db})
	return
}

func instanceControllerRoom(db *sql.DB, cm controller_interfaces.IGenericController[model_movie.Movie]) (result controller_interfaces.IGenericController[model_room.Room]) {
	result, _ = controller_room.NewControllerRoom(&controller_room.ControllerRoom{Db: db, MovieController: cm})
	return
}
