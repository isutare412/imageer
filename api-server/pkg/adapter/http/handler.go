package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
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
			log.Warnf("failed to read body")
			responseError(w, http.StatusInternalServerError, "failed to read body")
			return
		}

		var req signInReq
		if err := json.Unmarshal(reqBytes, &req); err != nil {
			log.Warnf("invalid body param")
			responseError(w, http.StatusBadRequest, "invalid body param")
			return
		}

		userEntity, err := uSvc.GetByEmailPassword(ctx, req.Email, req.Password)
		if errors.Is(err, user.ErrUserNotFound) {
			log.Warnf("user not found by email[%s]", req.Email)
			responseError(w, http.StatusBadRequest, "invalid email or password")
			return
		} else if errors.Is(err, user.ErrPasswordNotCorrect) {
			log.Warnf("password[%s] mismatch", req.Password)
			responseError(w, http.StatusBadRequest, "invalid email or password")
			return
		} else if err != nil {
			log.Errorf("failed to get user: %v", err)
			responseError(w, http.StatusInternalServerError, "failed to get user")
			return
		}

		sess := auth.Session{
			Id:        strconv.Itoa(int(userEntity.ID)),
			Privilege: string(userEntity.Privilege),
		}
		token, err := authSvc.SignToken(&sess)
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
		responseText(w, "success")
	}
}

// @Summary Say greeting
// @Description Greeting by given name
// @Tags Greeting
// @Router /api/v1/greetings/{name} [get]
// @Param Authorization header string false "bearer authorization" extensions(x-example=Bearer your_jwt_token)
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
			log.Warnf("name field not given")
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
// @Param Authorization header string false "bearer authorization" extensions(x-example=Bearer your_jwt_token)
// @Accept json
// @Produce json
// @Success 200 {object} getUserRes "ok"
// @Failure 400 {object} errorRes "error"
// @Failure 500 {object} errorRes "error"
func getUser(uSvc user.Service, authSvc auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sess, err := authSvc.SessionFromContext(ctx)
		if err != nil {
			log.Errorf("failed to get session from context: %v", err)
			responseError(w, http.StatusInternalServerError, "invalid session in request")
			return
		}
		id, err := strconv.ParseInt(string(sess.Id), 10, 64)
		if err != nil {
			log.Errorf("failed to parse session id[%s]", sess.Id)
			responseError(w, http.StatusInternalServerError, "failed to parse session id")
			return
		}

		usr, err := uSvc.GetByID(ctx, id)
		if errors.Is(err, user.ErrUserNotFound) {
			log.Warnf("user not found by id[%d]", id)
			responseError(w, http.StatusBadRequest, "id[%d] is invalid", id)
			return
		} else if err != nil {
			log.Errorf("failed to get user: %v", err)
			responseError(w, http.StatusInternalServerError, "failed to get user")
			return
		}

		res := getUserRes{
			ID:         usr.ID,
			Privilege:  string(usr.Privilege),
			GivenName:  usr.GivenName,
			FamilyName: usr.FamilyName,
			Email:      usr.Email,
			Credit:     usr.Credit,
		}
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
			log.Warnf("failed to read body")
			responseError(w, http.StatusInternalServerError, "failed to read body")
			return
		}

		var req createUserReq
		if err := json.Unmarshal(reqBytes, &req); err != nil {
			log.Warnf("invalid body param")
			responseError(w, http.StatusBadRequest, "invalid body param")
			return
		}

		usr := user.User{
			GivenName:  req.GivenName,
			FamilyName: req.FamilyName,
			Email:      req.Email,
		}
		userID, err := uSvc.Create(ctx, &usr, req.Password)
		if errors.Is(err, user.ErrDuplicate) {
			log.Warnf("failed to create user: %v", err)
			responseError(w, http.StatusBadRequest, err.Error())
			return
		} else if err != nil {
			log.Errorf("failed to create user: %v", err)
			responseError(w, http.StatusInternalServerError, "failed to create user")
			return
		}
		usr.ID = userID

		res := createUserRes{
			ID:         usr.ID,
			GivenName:  usr.GivenName,
			FamilyName: usr.FamilyName,
			Email:      usr.Email,
			Credit:     usr.Credit,
		}
		responseJson(w, &res)
	}
}

// @Summary Get a user
// @Description Get a user with given id
// @Tags User
// @Router /api/v1/users/{id} [get]
// @Param Authorization header string false "bearer authorization" extensions(x-example=Bearer your_jwt_token)
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
			log.Warnf("id field not given")
			responseError(w, http.StatusBadRequest, "'id' is mandatory field")
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Warnf("failed to parse int[%s]", idStr)
			responseError(w, http.StatusBadRequest, "id[%s] is invalid", idStr)
			return
		}

		usr, err := uSvc.GetByID(ctx, id)
		if errors.Is(err, user.ErrUserNotFound) {
			log.Warnf("user not found by id[%d]", id)
			responseError(w, http.StatusBadRequest, "user not exists")
			return
		} else if err != nil {
			log.Errorf("failed to get user: %v", err)
			responseError(w, http.StatusInternalServerError, "failed to get user")
			return
		}

		res := getUserRes{
			ID:         usr.ID,
			Privilege:  string(usr.Privilege),
			GivenName:  usr.GivenName,
			FamilyName: usr.FamilyName,
			Email:      usr.Email,
			Credit:     usr.Credit,
		}
		responseJson(w, &res)
	}
}

// @Summary Archive an image
// @Description Archive an image to object storage
// @Tags Image
// @Router /api/v1/images/archives [post]
// @Param Authorization header string false "bearer authorization" extensions(x-example=Bearer your_jwt_token)
// @Param image formData file true "image files to be archived"
// @Accept multipart/form-data
// @Success 200 {string} string "ok"
// @Failure 400 {object} errorRes "error"
// @Failure 500 {object} errorRes "error"
func archiveImage(jSvc job.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := r.ParseMultipartForm(32 << 20); err != nil {
			log.Errorf("failed parse multipart form: %v", err)
			responseError(w, http.StatusBadRequest, "failed parse multipart form")
			return
		}

		images := r.MultipartForm.File["image"]
		if len(images) == 0 {
			responseError(w, http.StatusBadRequest, "need at least one image")
			return
		}

		var files []io.ReadSeekCloser
		for _, img := range images {
			f, err := img.Open()
			if err != nil {
				log.Errorf("cannot open file[%s]: %v", img.Filename, err)
				responseError(w, http.StatusBadRequest, "cannot open file[%s]", img.Filename)
				return
			}
			defer f.Close()
			files = append(files, f)
		}

		for _, f := range files {
			if err := jSvc.Archive(ctx, "png", f); err != nil {
				log.Errorf("failed to produce by file: %v", err)
				responseError(w, http.StatusInternalServerError, "failed to produce by file")
				return
			}
		}
	}
}
