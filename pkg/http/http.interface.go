package http_adapter

import (
	"net/http"

	"github.com/gorilla/mux"
)

type IHTTP interface {
	AddRoute(method, url string, callback func(w http.ResponseWriter, r *http.Request))
	Listen()
	Route() *mux.Router
}
