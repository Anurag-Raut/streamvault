package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"streamvault/postgres"
	"streamvault/utils"
	"strings"
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
		Expires:  time.Now().Add(time.Hour * 24 * 10), // Set expiration time same as token
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
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
	if strings.Contains(signUpReq.Username, " ") {
		http.Error(w, "Username cannot contain spaces", http.StatusBadRequest)
		return
	}

	fmt.Println(signUpReq.Username)
	id, err := postgres.CreateUser(signUpReq.Username)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": signUpReq.Username,
		"userId":   id,
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("error creating user %v", err), http.StatusInternalServerError)
		return
	}

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secret)

	if err != nil {
		// http.Error(w, fmt.Sprintf("Error Signing token: %v", err), http.StatusInternalServerError)
		utils.SendError(w, fmt.Sprintf("Error Signing token: %v", err), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24 * 10), // Set expiration time same as token
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func GetUserDetails(w http.ResponseWriter,r *http.Request) {
	var response struct {
		postgres.UserDetails
		IsLoggedIn bool `json:"isLoggedIn"`

	}
	cookie, err := r.Cookie("jwt")
	if err != nil {
		fmt.Println("error getting cookie")
		response.IsLoggedIn = false
		resp, _ := json.MarshalIndent(response, "", "  ")
		w.Write(resp)
		return 
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		// return "", "", fmt.Errorf("error parsing token")
		// utils.SendError(w, "error parsing token", http.StatusInternalServerError)\
		fmt.Println("error parsing token")
		response.IsLoggedIn = false
		resp, _ := json.MarshalIndent(response, "", "  ")
		w.Write(resp)

		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		// return "", "", fmt.Errorf("error getting claims")
		fmt.Println("error getting claims")
		response.IsLoggedIn = false
		resp, _ := json.MarshalIndent(response, "", "  ")
		w.Write(resp)

		return

	}

	

	userId, ok := claims["userId"].(string)
	if !ok {
		// return "", "", fmt.Errorf("error getting userId")
		fmt.Println("error getting userId")
		response.IsLoggedIn = false
		resp, _ := json.MarshalIndent(response, "", "  ")
		w.Write(resp)
		return

	}


	response.UserId = userId
	response.IsLoggedIn = true
	var  userDetails postgres.UserDetails
	userDetails,err=postgres.GetUserDetailsFromDatabase(userId)
	if err != nil {
		fmt.Println("error getting user details")
		response.IsLoggedIn = false
		resp, _ := json.MarshalIndent(response, "", "  ")
		w.Write(resp)
		return
	}


	response.UserDetails=userDetails
	response.IsLoggedIn = true
	resp, _ := json.MarshalIndent(response, "", "  ")
	w.Write(resp)


	
	

	
}

