package handler

import (
	"encoding/json"
	"net/http"

	"before/user"
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserHandler struct{}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	data := CreateUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u := user.User{
		ID:    uuid.New(),
		Email: data.Email,
	}

	if err := u.SetPassword(data.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := u.Save(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	emailSubject := "Bienvenue !"
	emailBody := "Merci d'avoir souscrit à notre service."

	if err := u.SendEmail(r.Context(), emailSubject, emailBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
