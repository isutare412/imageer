package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

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
// @Failure 400 {object} errorRes "error"
// @Failure 500 {object} errorRes "error"
func signIn(uSvc user.Service, authSvc auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		reqBytes, err := io.ReadAll(r.Body)
		if err != nil {
			msg := "failed to read body"
			log.Errorf(msg)
			responseError(w, http.StatusInternalServerError, msg)
			return
		}

		var req signInReq
		if err := json.Unmarshal(reqBytes, &req); err != nil {
			responseError(w, http.StatusBadRequest, "invalid body param")
			return
		}

		userEntity, err := uSvc.GetByEmailPassword(ctx, req.Email, req.Password)
		if errors.Is(err, user.ErrUserNotFound) {
			responseError(w, http.StatusBadRequest, "invalid email or password")
			return
		} else if errors.Is(err, user.ErrPasswordNotCorrect) {
			responseError(w, http.StatusBadRequest, "invalid email or password")
			return
		} else if err != nil {
			log.Errorf("failed to get user: %v", err)
			responseError(w, http.StatusInternalServerError, "failed to get user")
			return
		}

		id := auth.ID(strconv.Itoa(int(userEntity.ID)))
		token, err := authSvc.SignToken(id)
		if err != nil {
			log.Errorf("failed to sign token: %v", err)
			responseError(w, http.StatusInternalServerError, "failed to sign token")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: string(token),
		})
		res := signInRes{
			Token: string(token),
		}
		responseJson(w, &res)
	}
}

// @Summary Sign out
// @Description Sign out by deleting cookie
// @Tags Authentication
// @Router /signOut [get]
// @Success 200 {string} string "ok"
// @Failure 400 {object} errorRes "error"
// @Failure 500 {object} errorRes "error"
func signOut(uSvc user.Service, authSvc auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   "",
			Expires: time.Unix(0, 0),
		})
	}
}

// @Summary Sign in test
// @Description Sign in test using authorization header or cookie
// @Tags Authentication
// @Router /signTest [get]
// @Param Authorization header string false "bearer authorization" extensions(x-example=Bearer your_jwt_token)
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 400 {object} errorRes "error"
// @Failure 500 {object} errorRes "error"
func signCheck(authSvc auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var token string
		if rawAuth := r.Header.Get("Authorization"); rawAuth != "" {
			authSplit := strings.SplitN(rawAuth, "Bearer ", 2)
			if len(authSplit) < 2 {
				responseError(w, http.StatusBadRequest, "invalid authorization header")
				return
			}
			token = authSplit[1]
		} else {
			cookie, err := r.Cookie("token")
			if err != nil {
				responseError(w, http.StatusBadRequest, "need token from cookie or authorization header")
				return
			}
			token = cookie.Value
		}

		id, err := authSvc.VerifyToken(auth.Token(token))
		if errors.Is(err, auth.ErrTokenExpired) {
			responseError(w, http.StatusInternalServerError, "token expired")
			return
		} else if err != nil {
			log.Errorf("failed to verify token: %v", err)
			responseError(w, http.StatusInternalServerError, "failed to verify token")
			return
		}
		log.Infof("verified token: id(%v)", id)

		responseText(w, "token verified")
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
// @Failure 400 {object} errorRes "error"
// @Failure 500 {object} errorRes "error"
func getGreeting(jSvc job.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name, ok := mux.Vars(r)["name"]
		if !ok {
			responseError(w, http.StatusBadRequest, "'name' is mandatory field")
			return
		}

		msg := fmt.Sprintf("Hello, %s!", name)

		if err := jSvc.Produce(r.Context(), msg); err != nil {
			log.Errorf("failed to produce job: %v", err)
			responseError(w, http.StatusInternalServerError, "failed to produce job")
			return
		}

		res := getGreetingRes{
			Message: msg,
		}
		responseJson(w, &res)
	}
}

// @Summary Get an authenticated user
// @Description Get an user by header or cookie
// @Tags User
// @Router /api/v1/users [get]
// @Accept json
// @Produce json
// @Success 200 {object} getUserRes "ok"
// @Failure 400 {object} errorRes "error"
// @Failure 500 {object} errorRes "error"
func getUser(uSvc user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		idStr, err := auth.IDFromContext(ctx)
		if err != nil {
			log.Errorf("failed to get id from context: %v", err)
			responseError(w, http.StatusInternalServerError, "invalid ID in request")
			return
		}
		id, err := strconv.ParseInt(string(idStr), 10, 64)
		if err != nil {
			log.Errorf("failed to parse id[%s]", idStr)
			responseError(w, http.StatusInternalServerError, "failed to parse id")
			return
		}

		userEntity, err := uSvc.GetByID(ctx, id)
		if errors.Is(err, user.ErrUserNotFound) {
			responseError(w, http.StatusBadRequest, "id[%d] is invalid", id)
			return
		} else if err != nil {
			log.Errorf("failed to get user: %v", err)
			responseError(w, http.StatusInternalServerError, "failed to get user")
			return
		}

		var res getUserRes
		res.from(userEntity)
		responseJson(w, &res)
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
// @Failure 400 {object} errorRes "error"
// @Failure 500 {object} errorRes "error"
func createUser(uSvc user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		reqBytes, err := io.ReadAll(r.Body)
		if err != nil {
			msg := "failed to read body"
			log.Errorf(msg)
			responseError(w, http.StatusInternalServerError, msg)
			return
		}

		var req createUserReq
		if err := json.Unmarshal(reqBytes, &req); err != nil {
			responseError(w, http.StatusBadRequest, "invalid body param")
			return
		}
		userEntity := req.into()

		userID, err := uSvc.Create(ctx, userEntity, req.Password)
		if errors.Is(err, user.ErrDuplicate) {
			responseError(w, http.StatusBadRequest, "email[%s] duplicated", userEntity.Email)
			return
		} else if err != nil {
			log.Errorf("failed to marshal response: %v", err)
			responseError(w, http.StatusInternalServerError, "failed to marshal response")
			return
		}
		userEntity.ID = userID

		var res createUserRes
		res.from(userEntity)
		responseJson(w, &res)
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
// @Failure 400 {object} errorRes "error"
// @Failure 500 {object} errorRes "error"
func getUserByID(uSvc user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		idStr, ok := mux.Vars(r)["id"]
		if !ok {
			responseError(w, http.StatusBadRequest, "'id' is mandatory field")
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			responseError(w, http.StatusBadRequest, "id[%s] is invalid", idStr)
			return
		}

		userEntity, err := uSvc.GetByID(ctx, id)
		if errors.Is(err, user.ErrUserNotFound) {
			responseError(w, http.StatusBadRequest, "user not exists")
			return
		} else if err != nil {
			log.Errorf("failed to get user: %v", err)
			responseError(w, http.StatusInternalServerError, "failed to get user")
			return
		}

		var res getUserRes
		res.from(userEntity)
		responseJson(w, &res)
	}
}
