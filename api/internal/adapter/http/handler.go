package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// @Summary Say greeting
// @Description Greeting by given name
// @Tags Greeting
// @Router /api/v1/greeting/{name} [get]
// @Param name path string true "name for greeting"
// @Accept json
// @Produce json
// @Success 200 {object} getGreetingResp "ok"
// @Failure 400 {string} string "error"
// @Failure 500 {string} string "error"
func getGreeting(w http.ResponseWriter, r *http.Request) {
	name, ok := mux.Vars(r)["name"]
	if !ok {
		http.Error(w, "'name' is mandatory field", http.StatusBadRequest)
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
