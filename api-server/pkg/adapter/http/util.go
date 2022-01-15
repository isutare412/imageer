package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func responseError(w http.ResponseWriter, code int, format string, values ...interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	msg := fmt.Sprintf(format, values...)
	resBytes, _ := json.Marshal(&errorRes{
		Code:    code,
		Message: msg,
	})
	w.Write(resBytes)
}

// TODO: response by Text, JSON
// TODO: change OAS response error type
