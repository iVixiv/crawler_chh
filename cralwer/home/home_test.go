package home

import (
	"testing"
	"encoding/json"
)

func Test_Home(t *testing.T) {
	info, err := new(Home).Crawler()
	if err != nil {
		log.Error("", err)
	} else {
		jsonStr, _ := json.Marshal(info)
		log.Notice("Crawler output: %s", string(jsonStr))
	}
}
