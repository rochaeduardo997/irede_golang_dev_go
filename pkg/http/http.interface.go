package http_adapter

import "net/http"

type IHTTP interface {
	AddRoute(method, url string, callback func(w http.ResponseWriter, r *http.Request))
	Listen()
}
