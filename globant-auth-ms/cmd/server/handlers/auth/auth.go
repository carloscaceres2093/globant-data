package handlers

import (
	"encoding/json"
	"globant-auth-ms/local-lib/web"
	"net/http"
	"strings"

	"globant-auth-ms/internal/auth"

	"github.com/go-chi/chi/v5"
	"github.com/labstack/gommon/log"
)

const (
	ErrInvalidBody         = "invalid body"
	ErrInternalServerError = "internal server error"
	ErrInvalidName         = "invalid user name"
	ErrNameDuplicated      = "user name duplicated"
	ErrInvalidToken        = "invalid token"
	ErrTokenEmpty          = "token empty"
	ErrInvalidUser         = "invalid user"
)

type AuthHandler struct {
	AuthService
}

type AuthService interface {
	ValidateToken(user string, token string) (auth.AuthResponse, error)
	CreateUser(userReq auth.AuthRequest) (auth.AuthResponse, error)
}

func NewAuthHandler(service AuthService) AuthHandler {
	return AuthHandler{
		service,
	}
}

func (rh *AuthHandler) CreateUser(w http.ResponseWriter, req *http.Request) {
	var request auth.AuthRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		log.Errorf("[err:%+v]", err)
		web.RespondJSON(w, web.Error{Message: ErrInvalidBody}, http.StatusBadRequest)
		return
	}

	if request.UserName == "" {
		log.Errorf("[err:%+v]", ErrInvalidName)
		web.RespondJSON(w, web.Error{Message: ErrInvalidName}, http.StatusBadRequest)
		return
	}

	authResponse, err := rh.AuthService.CreateUser(request)
	if err != nil {
		log.Errorf("[err:%+v]", err)
		if err == auth.ErrNameDuplicated {
			web.RespondJSON(w, web.Error{Message: ErrNameDuplicated}, http.StatusBadRequest)
			return

		}
		if strings.Contains(err.Error(), "duplicate key value") {
			web.RespondJSON(w, web.Error{Message: ErrNameDuplicated}, http.StatusBadRequest)
			return
		}
		web.RespondJSON(w, web.Error{Message: ErrInternalServerError}, http.StatusInternalServerError)
		return
	}

	web.RespondJSON(w, authResponse, http.StatusCreated)
}

func (rh *AuthHandler) ValidateToken(w http.ResponseWriter, req *http.Request) {
	token := req.URL.Query().Get("token")
	userCode := chi.URLParam(req, "user_code")
	if token == "" {
		web.RespondJSON(w, web.Error{Message: ErrTokenEmpty}, http.StatusBadRequest)
		return
	}
	if userCode == "" {
		web.RespondJSON(w, web.Error{Message: ErrInvalidUser}, http.StatusBadRequest)
		return
	}
	userResponse, err := rh.AuthService.ValidateToken(userCode, token)

	if err != nil {
		log.Errorf("[err:%+v]", err)
		if err == auth.ErrInvalidToken {
			web.RespondJSON(w, web.Error{Message: ErrInvalidToken}, http.StatusBadRequest)
			return
		}
		if err == auth.ErrInvalidUser {
			web.RespondJSON(w, web.Error{Message: ErrInvalidUser}, http.StatusBadRequest)
			return
		}
		if strings.Contains(err.Error(), "duplicate key value") {
			web.RespondJSON(w, web.Error{Message: ErrNameDuplicated}, http.StatusBadRequest)
			return
		}
		web.RespondJSON(w, web.Error{Message: ErrInternalServerError}, http.StatusInternalServerError)
		return
	}

	web.RespondJSON(w, userResponse, http.StatusOK)
}
