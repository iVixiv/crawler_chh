package server

import "github.com/op/go-logging"

var (
	SERVE_PORT = "8000"
	log        = logging.MustGetLogger("server")
)


type Result struct {
	Code    uint16
	Message string
	Data    interface{}
}
