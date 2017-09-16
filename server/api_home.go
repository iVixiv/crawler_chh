package server

import (
	"net/http"
	"crawler_chh/cralwer/home"
)

func handleAPIChhHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("content-type", "text/html; charset=utf-8")
	defer r.Body.Close()

	Response(w, HandlerSuccessResult(new(home.Home).Crawler()))
}
