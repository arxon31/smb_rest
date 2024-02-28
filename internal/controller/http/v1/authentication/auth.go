package authentication

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"git.spbec-mining.ru/arxon31/sambaMW/internal/entity"
	"net/http"
	"time"
)

type AuthService struct {
	token string
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (m *AuthService) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("X-Auth-Token")
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if token.Value != m.token {
			fmt.Printf("token: %s, expected: %s\n", token, m.token)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *AuthService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var model entity.User

	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = model.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b := make([]byte, 32)

	_, err = rand.Read(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m.token = hex.EncodeToString(b)

	http.SetCookie(w, &http.Cookie{
		Name:    "X-Auth-Token",
		Value:   m.token,
		Expires: time.Now().Add(365 * 24 * time.Hour),
	})

	w.WriteHeader(http.StatusOK)
}
