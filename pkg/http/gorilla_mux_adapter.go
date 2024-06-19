package http_adapter

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	PORT := os.Getenv("API_PORT")
	log.Printf("server running on http://localhost:%s", PORT)
	http.ListenAndServe(fmt.Sprintf(":%s", PORT), gm.Router)
}

func (gm *GorillaMux) Route() *mux.Router {
	return gm.Router
}
