package server

import (
	"net/http"
)

func Start(addr string) {
	http.HandleFunc("/home", handleAPIChhHome)
	http.HandleFunc("/mod", handleAPIChhMod)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
