package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"streamvault/postgres"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// func SignUpWithEmail(email string) {
// Sender data.'
var secret = []byte("eat shit")

func SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var signInReq struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&signInReq); err != nil {
		http.Error(w, fmt.Sprintf("error decoding json %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println(signInReq.Username)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": signInReq.Username,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secret)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error Signing token: %v", err), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24), // Set expiration time same as token
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	w.Write([]byte("ok"))

}

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var signUpReq struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&signUpReq); err != nil {
		http.Error(w, fmt.Sprintf("error decoding json %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println(signUpReq.Username)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": signUpReq.Username,
	})
	err := postgres.CreateUser(signUpReq.Username)
	if err != nil {
		http.Error(w, fmt.Sprintf("error creating user %v", err), http.StatusInternalServerError)
		return
	}

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secret)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error Signing token: %v", err), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24), // Set expiration time same as token
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}
