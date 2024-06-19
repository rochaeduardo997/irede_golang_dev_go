package http_adapter

import (
	"net/http"

	"github.com/gorilla/mux"
)

type GorillaMux struct {
	Router *mux.Router
}

func NewGorillaMux() (result IHTTP, handler *mux.Router) {
	gm := &GorillaMux{}
	gm.Router = mux.NewRouter()
	result = gm
	return result, gm.Router
}

func (gm *GorillaMux) AddRoute(method, url string, callback func(w http.ResponseWriter, r *http.Request)) {
	gm.Router.HandleFunc(url, callback).Methods(method)
}

func (gm *GorillaMux) Listen() {
	http.ListenAndServe(":3000", gm.Router)
}
