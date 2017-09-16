package server

import (
	"net/http"
	"crawler_chh/cralwer/mod"
	"encoding/json"
	"io/ioutil"
)

func handleAPIChhMod(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("content-type", "text/html; charset=utf-8")

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Response(w, HandlerSuccessResult(nil, err))
		return
	}

	var params mod.RequestParams
	if err = json.Unmarshal([]byte(string(body)), &params); err != nil {
		Response(w, HandlerSuccessResult(nil, err))
		return
	}
	Response(w, HandlerSuccessResult(new(mod.Mod).Crawler(params.ModUrl)))
}
