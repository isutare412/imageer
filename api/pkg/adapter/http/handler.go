package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func getGreeting(w http.ResponseWriter, r *http.Request) {
	const NAME = "name"

	name, ok := mux.Vars(r)[NAME]
	if !ok {
		http.Error(w, fmt.Sprintf("%q is mandatory field", NAME), http.StatusBadRequest)
		return
	}

	resBody := getGreetingResp{
		Message: fmt.Sprintf("Hello, %s!", name),
	}

	res, err := json.Marshal(&resBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
