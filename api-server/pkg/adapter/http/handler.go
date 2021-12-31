package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/isutare412/imageer/api-server/pkg/core/auth"
	"github.com/isutare412/imageer/api-server/pkg/core/job"
	"github.com/isutare412/imageer/api-server/pkg/core/user"
)

// @Summary Sign in
// @Description Sign in using email and password
// @Tags Authentication
// @Router /signIn [post]
// @Param request body signInReq true "request to sign in"
// @Accept json
// @Produce json
// @Success 200 {object} signInRes "ok"
// @Failure 400 {string} string "error"
// @Failure 500 {string} string "error"
func signIn(uSvc user.Service, authSvc auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		reqBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Need body param", http.StatusBadRequest)
			return
		}

		var req signInReq
		if err := json.Unmarshal(reqBytes, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userEntity, err := uSvc.GetByEmailPassword(ctx, req.Email, req.Password)
		if errors.Is(err, user.ErrUserNotFound) {
			log.Infof("Email not found: %v", req.Email)
			http.Error(w, "Invalid email or password", http.StatusBadRequest)
			return
		} else if errors.Is(err, user.ErrPasswordNotCorrect) {
			log.Infof("Password incorrect: %v", req.Password)
			http.Error(w, "Invalid email or password", http.StatusBadRequest)
			return
		} else if err != nil {
			log.Errorf("Failed to get user: %v", err)
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
			return
		}

		id := auth.ID(strconv.Itoa(int(userEntity.ID)))
		token, err := authSvc.SignToken(id)
		if err != nil {
			log.Errorf("Failed to sign token: %v", err)
			http.Error(w, "Failed to sign token", http.StatusInternalServerError)
			return
		}

		var res signInRes
		res.Token = string(token)
		resBytes, err := json.Marshal(&res)
		if err != nil {
			log.Errorf("Failed marshal response: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: string(token),
		})

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(resBytes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// @Summary Sign check
// @Description Debug sign in taking query or cookie
// @Tags Authentication
// @Router /signCheck [get]
// @Param token query string false "jwt token"
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "error"
// @Failure 500 {string} string "error"
func signCheck(authSvc auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			cookie, err := r.Cookie("token")
			if err != nil {
				msg := "Need token as query param or cookie"
				log.Info(msg)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
			token = cookie.Value
		}

		id, err := authSvc.VerifyToken(auth.Token(token))
		if errors.Is(err, auth.ErrTokenExpired) {
			msg := "Token expired"
			log.Info(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		} else if err != nil {
			log.Errorf("Failed to verify token: %v", err)
			http.Error(w, "Failed to verify token", http.StatusInternalServerError)
			return
		}
		log.Infof("Verified token: id(%v)", id)

		w.Header().Set("Content-Type", "plain/text")
		w.Write([]byte("Token verified"))
	}
}

// @Summary Say greeting
// @Description Greeting by given name
// @Tags Greeting
// @Router /api/v1/greetings/{name} [get]
// @Param name path string true "name for greeting"
// @Accept json
// @Produce json
// @Success 200 {object} getGreetingRes "ok"
// @Failure 400 {string} string "error"
// @Failure 500 {string} string "error"
func getGreeting(jSvc job.Service) http.HandlerFunc {
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
		userEntity := req.into()

		userID, err := uSvc.Create(ctx, userEntity, req.Password)
		if errors.Is(err, user.ErrDuplicate) {
			msg := fmt.Sprintf("Duplicate email: %v", userEntity.Email)
			log.Info(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		} else if err != nil {
			log.Errorf("Failed to create user: %v", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		userEntity.ID = userID

		var res createUserRes
		res.from(userEntity)
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
func getUserByID(uSvc user.Service) http.HandlerFunc {
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

		userEntity, err := uSvc.GetByID(ctx, id)
		if errors.Is(err, user.ErrUserNotFound) {
			msg := fmt.Sprintf("Invalid user id: %v", id)
			log.Info(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		} else if err != nil {
			log.Errorf("Failed to get user: %v", err)
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
			return
		}

		var res getUserRes
		res.from(userEntity)
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
