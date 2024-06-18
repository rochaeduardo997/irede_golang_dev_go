package view_interfaces

import model_movie "github.com/rochaeduardo997/irede_golang_dev/internal/model/movie"

type IControllerMovie interface {
	Create(m *model_movie.Movie) (result string, err error)
	FindBy(id string) (result *model_movie.Movie, err error)
	FindAll() (result []*model_movie.Movie, err error)
	UpdateBy(id string, m *model_movie.Movie) (result bool, err error)
	DeleteBy(id string) (result bool, err error)
}
