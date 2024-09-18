package ports

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/MarinDmitrii/notes-service/internal/common"
	"github.com/MarinDmitrii/notes-service/internal/user/builder"
	"github.com/MarinDmitrii/notes-service/internal/user/domain"
	"github.com/MarinDmitrii/notes-service/internal/user/usecase"
)

type HttpUserHandler struct {
	app *builder.Application
}

func NewHttpUserHandler(app *builder.Application) HttpUserHandler {
	return HttpUserHandler{app: app}
}

func (h HttpUserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	l := log.Default()

	if r.Method != http.MethodPost {
		h.mapToResponse(w, http.StatusMethodNotAllowed, nil, http.StatusText(http.StatusMethodNotAllowed))
		l.Println(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	request := &PostUser{}

	if r.Header.Get("Content-Type") == "application/json" {
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			h.mapToResponse(w, http.StatusBadRequest, nil, err.Error())
			return
		}
	} else {
		h.mapToResponse(w, http.StatusBadRequest, nil, http.StatusText(http.StatusBadRequest))
		return
	}

	createdUser, _ := h.app.GetUserByEmail.Execute(r.Context(), request.Email)

	if createdUser.ID != 0 {
		h.mapToResponse(w, http.StatusBadRequest, nil, "User already exists")
		return
	}

	createdUser, err := h.app.SaveUser.Execute(
		r.Context(),
		usecase.SaveUser{
			Email:    request.Email,
			Password: request.Password,
		},
	)
	if err != nil {
		h.mapToResponse(w, http.StatusServiceUnavailable, nil, err.Error())
		return
	}

	h.mapToResponse(w, http.StatusOK, NewUser(createdUser), "")
}

type PostUser struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Created time.Time `json:"created"`
	Email   string    `json:"email"`
	ID      int       `json:"id"`
}

func NewUser(user domain.User) User {
	return User{
		Created: user.CreateDt,
		Email:   user.Email,
		ID:      user.ID,
	}
}

func (h HttpUserHandler) mapToResponse(w http.ResponseWriter, statusCode int, data interface{}, errMessage string) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	response := make(map[string]interface{})

	if statusCode >= 200 && statusCode < 300 {
		response["result"] = data
	} else {
		response["error"] = errMessage
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CustomRegisterHandlers(router *http.ServeMux, h HttpUserHandler) {
	router.HandleFunc("/auth/sign-up", common.LogMiddleware(h.CreateUser))
}
