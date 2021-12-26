package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/isutare412/imageer/api/pkg/core/job"
	"github.com/isutare412/imageer/api/pkg/core/user"
)

// @Summary Say greeting
// @Description Greeting by given name
// @Tags Greeting
// @Router /api/v1/greeting/{name} [get]
// @Param name path string true "name for greeting"
// @Accept json
// @Produce json
// @Success 200 {object} getGreetingRes "ok"
// @Failure 400 {string} string "error"
// @Failure 500 {string} string "error"
func getGreeting(jSvc *job.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name, ok := mux.Vars(r)["name"]
		if !ok {
			http.Error(w, "'name' is mandatory field", http.StatusBadRequest)
			return
		}

		msg := fmt.Sprintf("Hello, %s!", name)

		if err := jSvc.Produce(r.Context(), msg); err != nil {
			errMsg := fmt.Sprintf("Failed to produce job: %v", err)
			log.Error(errMsg)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}

		res := getGreetingRes{
			Message: msg,
		}

		resBytes, err := json.Marshal(&res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(resBytes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// @Summary Create a user
// @Description Create a user with basic information
// @Tags User
// @Router /api/v1/users [post]
// @Param request body createUserReq true "request to create a new user"
// @Accept json
// @Produce json
// @Success 200 {object} createUserRes "ok"
// @Failure 400 {string} string "error"
// @Failure 500 {string} string "error"
func createUser(uSvc user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		reqBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Need body param", http.StatusBadRequest)
			return
		}

		var req createUserReq
		if err := json.Unmarshal(reqBytes, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user := req.into()

		userID, err := uSvc.Create(ctx, user)
		if err != nil {
			log.Errorf("Failed to create user: %v", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		user.ID = userID

		var res createUserRes
		res.from(user)
		resBytes, err := json.Marshal(&res)
		if err != nil {
			log.Errorf("Failed marshal response: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(resBytes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// @Summary Get a user
// @Description Get a user with given id
// @Tags User
// @Router /api/v1/users/{id} [get]
// @Param id path string true "user id"
// @Accept json
// @Produce json
// @Success 200 {object} getUserRes "ok"
// @Failure 400 {string} string "error"
// @Failure 500 {string} string "error"
func getUser(uSvc user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		idStr, ok := mux.Vars(r)["id"]
		if !ok {
			http.Error(w, "'id' is mandatory field", http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "'id' is not valid", http.StatusBadRequest)
			return
		}

		user, err := uSvc.GetByID(ctx, id)
		if err != nil {
			log.Errorf("Failed to get user: %v", err)
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
			return
		}

		var res getUserRes
		res.from(user)
		resBytes, err := json.Marshal(&res)
		if err != nil {
			log.Errorf("Failed marshal response: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(resBytes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
