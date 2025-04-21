package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"github.com/dgrijalva/jwt-go"
	"log"
)

// JWT Secret
var jwtKey = []byte("secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// JWT Token Claims
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

var users []User

func main() {
	data, _ := ioutil.ReadFile("users.json")
	json.Unmarshal(data, &users)

	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/manager", ManagerHandler)
	http.HandleFunc("/comments/get", GetComments)
	http.HandleFunc("/comments/add", AddComment)


	fmt.Println("Server çalışıyor :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "JSON formatında veri gönderin!", http.StatusBadRequest)
		return
	}

	// JSON decode hatasını kontrol edelim
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Geçersiz istek formatı!", http.StatusBadRequest)
		return
	}

	for _, u := range users {
		if u.Username == creds.Username && u.Password == creds.Password {
			expirationTime := time.Now().Add(1 * time.Hour)
			claims := &Claims{
				Username: u.Username,
				Role:     u.Role,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(),
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, _ := token.SignedString(jwtKey)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
			return
		}
	}
	http.Error(w, "Yetkisiz kullanıcı!", http.StatusUnauthorized)
}


func ManagerHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Token gerekli!", http.StatusUnauthorized)
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid || claims.Role != "admin" {
		http.Error(w, "Yetkiniz yok!", http.StatusForbidden)
		return
	}

	w.Write([]byte("🚩 FLAG{xxxxxxxxxx}"))
}

type Comment struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

// Yorumları Getir
func GetComments(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadFile("comments.json")
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func AddComment(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Token gerekli!", http.StatusUnauthorized)
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Geçersiz token!", http.StatusUnauthorized)
		return
	}

	var newComment Comment
	json.NewDecoder(r.Body).Decode(&newComment)
	newComment.Username = claims.Username

	data, _ := ioutil.ReadFile("comments.json")
	var comments []Comment
	json.Unmarshal(data, &comments)

	comments = append(comments, newComment)
	newData, _ := json.Marshal(comments)
	ioutil.WriteFile("comments.json", newData, 0644)

	w.WriteHeader(http.StatusCreated)
}

