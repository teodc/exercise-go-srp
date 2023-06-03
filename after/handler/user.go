package handler

import (
	"encoding/json"
	"net/http"

	"after/email"
	"after/password"
	"after/user"
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserHandler struct {
	userRepository *user.Repository
	passwordHasher *password.Hasher
	mailer         *email.SesMailer
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	data := CreateUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	encodedPassword, err := h.passwordHasher.HashPassword(data.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u := user.User{
		ID:           uuid.New(),
		Email:        data.Email,
		PasswordHash: encodedPassword,
	}

	if err := h.userRepository.Save(r.Context(), u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	emailSubject := "Bienvenue !"
	emailBody := "Merci d'avoir souscrit à notre service."

	if err := h.mailer.SendEmail(r.Context(), u.Email, emailSubject, emailBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}
