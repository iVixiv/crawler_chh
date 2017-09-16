package server

import (
	"encoding/json"
	"net/http"
	"fmt"
)

func HandlerSuccessResult(data interface{}, err error) (string) {
	result := Result{}
	if err == nil {
		result.Code = 200
		result.Message = "success"
		result.Data = data
	} else {
		result.Code = 500
		result.Message = err.Error()
	}
	responseBody, _ := json.Marshal(result)
	return string(responseBody)
}

func Response(w http.ResponseWriter, data string) {
	fmt.Fprint(w, data)
}
