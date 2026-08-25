package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("goFinalProject-jwt-secret")

const tokenTTL = 8 * time.Hour

var authPassword string

func initAuth() {
	authPassword = os.Getenv("TODO_PASSWORD")
}

type passwordClaims struct {
	Hash string `json:"hash"`
	jwt.RegisteredClaims
}

func passwordHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный JSON: "+err.Error())
		return
	}

	if authPassword == "" || req.Password != authPassword {
		writeError(w, http.StatusUnauthorized, "Неверный пароль")
		return
	}

	claims := passwordClaims{
		Hash: passwordHash(authPassword),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		},
	}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"token": signed})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if authPassword == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			writeError(w, http.StatusUnauthorized, "Требуется аутентификация")
			return
		}

		claims := &passwordClaims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (any, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid || claims.Hash != passwordHash(authPassword) {
			writeError(w, http.StatusUnauthorized, "Требуется аутентификация")
			return
		}

		next(w, r)
	}
}
