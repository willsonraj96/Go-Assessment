package main

import (
	"time"
	"log"
	"net/http"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"encoding/json"
)

var jwtSecret = []byte("your_jwt_secret_key")

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

var users = map[string]string{}  

func main() {
	http.HandleFunc("/signup", signUpHandler)
	http.HandleFunc("/signin", signInHandler)
	http.HandleFunc("/refresh", refreshTokenHandler)
	http.HandleFunc("/protected", protectedHandler)

	log.Println("Server Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	log.Println("Handling signup request")
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if _, exists := users[user.Email]; exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	users[user.Email] = user.Password
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("User created successfully")
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if password, ok := users[user.Email]; !ok || password != user.Password {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := createToken(user.Email)
	if err != nil {
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(Token{Token: token})
}


func createToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp": time.Now().Add(time.Minute * 5).Unix(), 
	})
	return token.SignedString(jwtSecret)
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	claims, err := parseToken(tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"email": claims["email"].(string)})
}


func parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	return token.Claims.(jwt.MapClaims), nil
}



func refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	// Remove "Bearer " prefix if it exists
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	claims, err := parseToken(tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Generate a new token
	email := claims["email"].(string)
	newToken, err := createToken(email)
	if err != nil {
		http.Error(w, "Failed to refresh token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(Token{Token: newToken})
}
