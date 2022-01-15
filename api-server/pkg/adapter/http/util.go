package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
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

func responseJson(w http.ResponseWriter, res interface{}) {
	resBytes, err := json.Marshal(&res)
	if err != nil {
		log.Errorf("failed to marshal response: %v", err)
		responseError(w, http.StatusInternalServerError, "failed to marshal response")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(resBytes)
}

func responseText(w http.ResponseWriter, txt string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(txt))
}
