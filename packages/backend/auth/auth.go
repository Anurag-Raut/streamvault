package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"streamvault/postgres"
	"streamvault/utils"
	"strings"
	"time"

	"context"

	"github.com/gofor-little/env"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	oauth2pkg "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

// func SignUpWithEmail(email string) {
// Sender data.'
var secret = []byte("eat shit")

func SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var signInReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&signInReq); err != nil {
		fmt.Println("error hello", err.Error())
		utils.SendError(w, fmt.Sprintf("error decoding json %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Println(signInReq.Username, signInReq.Password)
	isPasscordCorrect, id, err := postgres.CheckUsernamePassword(signInReq.Username, signInReq.Password)
	if err != nil {
		fmt.Println("error hello", err.Error())
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isPasscordCorrect {
		fmt.Println("error hello", err.Error())
		utils.SendError(w, fmt.Sprintf("password not matching %v", err), http.StatusUnauthorized)
		return
	}

	fmt.Println(signInReq.Username)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": signInReq.Username,
		"userId":   id,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secret)

	if err != nil {
		fmt.Println("errurr", err.Error())
		utils.SendError(w, fmt.Sprintf("Error Signing token: %v", err), http.StatusInternalServerError)
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
	var response = "ok"
	resp, _ := json.MarshalIndent(response, "", "  ")
	w.Write(resp)

}

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var signUpReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var id string
	if err := json.NewDecoder(r.Body).Decode(&signUpReq); err != nil {
		utils.SendError(w, fmt.Sprintf("error decoding json %v", err), http.StatusInternalServerError)
		return
	}
	if strings.Contains(signUpReq.Username, " ") {
		utils.SendError(w, "Username cannot contain spaces", http.StatusBadRequest)
		return
	}

	fmt.Println(signUpReq.Username)

	id, err := postgres.CreateUserWithPassword(signUpReq.Username, signUpReq.Password)

	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": signUpReq.Username,
		"userId":   id,
	})

	if err != nil {
		utils.SendError(w, fmt.Sprintf("error creating user %v", err), http.StatusInternalServerError)
		return
	}

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secret)

	if err != nil {
		// utils.SendError(w, fmt.Sprintf("Error Signing token: %v", err), http.StatusInternalServerError)
		utils.SendError(w, fmt.Sprintf("Error Signing token: %v", err), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24 * 10), // Set expiration time same as token
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}
	
	http.SetCookie(w, &cookie)
	var response = "ok"
	resp, _ := json.MarshalIndent(response, "", "  ")
	w.Write(resp)
}

func GetUserDetails(w http.ResponseWriter, r *http.Request) {
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
	var userDetails postgres.UserDetails
	userDetails, err = postgres.GetUserDetailsFromDatabase(userId)
	if err != nil {
		fmt.Println("error getting user details")
		response.IsLoggedIn = false
		resp, _ := json.MarshalIndent(response, "", "  ")
		w.Write(resp)
		return
	}

	response.UserDetails = userDetails
	response.IsLoggedIn = true
	resp, _ := json.MarshalIndent(response, "", "  ")
	w.Write(resp)

}

func SignOut(w http.ResponseWriter, r *http.Request) {
	fmt.Println("signing out")
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}


	http.SetCookie(w, &cookie)
	// w.Write([]byte("ok"))
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte("ok"))
	var response string
	response = "ok"
	resp, _ := json.MarshalIndent(response, "", "  ")
	w.Write(resp)

}

func GetGoogleUrl(w http.ResponseWriter, r *http.Request) {
	GOOGLE_CLIENT_ID, err := env.MustGet("GOOGLE_CLIENT_ID")
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	GOOGLE_CLIENT_SECRET, err := env.MustGet("GOOGLE_CLIENT_SECRET")
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var conf = &oauth2.Config{
		ClientID:     GOOGLE_CLIENT_ID,
		ClientSecret: GOOGLE_CLIENT_SECRET,
		RedirectURL:  fmt.Sprintf("%s/auth/signIn", env.Get("FRONTEND_URL", "https://streamvault.vercel.app")),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	// // Redirect user to Google's consent page to ask for permission
	// // for the scopes specified above.
	url := conf.AuthCodeURL("state")
	// fmt.Printf("Visit the URL for the auth dialog: %v", url)
	response, err := json.MarshalIndent(url, "", "  ")
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(response)

}

func LoginWithGoogle(w http.ResponseWriter, r *http.Request) {
	GOOGLE_CLIENT_ID, err := env.MustGet("GOOGLE_CLIENT_ID")
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	GOOGLE_CLIENT_SECRET, err := env.MustGet("GOOGLE_CLIENT_SECRET")
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var conf = &oauth2.Config{
		ClientID:     GOOGLE_CLIENT_ID,
		ClientSecret: GOOGLE_CLIENT_SECRET,
		RedirectURL:  fmt.Sprintf("%s/auth/signIn", env.Get("FRONTEND_URL", "https://streamvault.vercel.app")),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	var code string
	err = json.NewDecoder(r.Body).Decode(&code)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		// log.Fatal("fuck me")
		utils.SendError(w, "Unable to exchange ", http.StatusInternalServerError)
	}
	ctx = context.Background()

	oauth2Service, err := oauth2pkg.NewService(ctx, option.WithScopes("https://www.googleapis.com/auth/userinfo.profile"), option.WithTokenSource(conf.TokenSource(ctx, tok)))
	if err != nil {
		// log.Fatalf("Fuck you")
		// log.Fatalf("Unable to create Oauth2 service: %v", err)
		utils.SendError(w, "Unable to create Oauth2 service", http.StatusInternalServerError)
	}
	fmt.Println(tok.AccessToken, "tokennn")
	userinfoService := oauth2pkg.NewUserinfoService(oauth2Service)

	userInfo, err := userinfoService.Get().Do(googleapi.QueryParameter("access_token", tok.AccessToken))

	id, err := postgres.CreateUser(userInfo.Name, &userInfo.Picture)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": userInfo.Name,
		"userId":   id,
	})

	if err != nil {
		utils.SendError(w, fmt.Sprintf("error creating user %v", err), http.StatusInternalServerError)
		return
	}

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secret)

	if err != nil {
		// utils.SendError(w, fmt.Sprintf("Error Signing token: %v", err), http.StatusInternalServerError)
		utils.SendError(w, fmt.Sprintf("Error Signing token: %v", err), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24 * 10), // Set expiration time same as token
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	http.SetCookie(w, &cookie)
	var response string
	response = fmt.Sprintf("User %s created", userInfo.Name)
	resp, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(resp)

}
