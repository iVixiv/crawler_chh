package main

import (
	"github.com/op/go-logging"
	"fmt"
	"crawler_chh/server"
)

var log = logging.MustGetLogger("main")

func main() {
	addr := fmt.Sprintf(":%s", server.SERVE_PORT)
	log.Info("Server Listening %s", addr)

	server.Start(addr)
}
