package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"master/config"
	"master/db/postgres_db"
	"master/types"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

// func (s *Server) WithJWTAuth(h http.HandlerFunc) http.HandlerFunc{
//
// }
//

func (s *Server) RegisterUser(w http.ResponseWriter, r *http.Request) {
	//get mail and password
	var user_credentials types.User_Credentials
	err := json.NewDecoder(r.Body).Decode(&user_credentials)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	fmt.Println("imp ",user_credentials)

	//validate
	if len(user_credentials.Password) < 8 {
		http.Error(w, "password must be longe than 8 chars", http.StatusBadRequest)
		return
	}
	if user_credentials.UserMail == "" {
		http.Error(w, "user name can not be empty", http.StatusBadRequest)
		return
	}

	if exists, err := s.queries.UserNameExists(context.Background(), pgtype.Text{
		String: user_credentials.UserMail,
		Valid:  true,
	}); err != nil || exists == true {
		fmt.Println("faild with", err)
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	}

	hash, err := HashPassword(user_credentials.Password)
	if err != nil {
		fmt.Println("failed to hash password", err)
		http.Error(w, "can not hash password", http.StatusInternalServerError)
		return
	}

	parms := postgres_db.CreateUserAndReturnIdParams{
		Username: pgtype.Text{
			String: user_credentials.UserMail,
			Valid:  true,
		},
		Email: pgtype.Text{
			String: user_credentials.UserMail,
			Valid:  true,
		},
		EncryptedPassword: pgtype.Text{
			String: hash,
			Valid:  true,
		},
		Role: pgtype.Text{
			String: "user",
			Valid:  true,
		},
		CreatedAt: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
		LastOnline: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
	}
	_, err = s.queries.CreateUserAndReturnId(context.Background(), parms)
	if err != nil {
		fmt.Println("can not create user", err)
		http.Error(w, "can not create user", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("user created"))
}

func (s *Server) UserLogin(w http.ResponseWriter, r *http.Request) {
	//get mail and password
	var user_credentials types.User_Credentials
	err := json.NewDecoder(r.Body).Decode(&user_credentials)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	user, err := s.queries.GetUser(context.Background(), pgtype.Text{
		String: user_credentials.UserMail,
		Valid:  true,
	})
	if err != nil {
		fmt.Println("filde to login ", err)
		http.Error(w, "login failed", http.StatusInternalServerError)
		return
	}

	if CheckPassword(user_credentials.Password, user.EncryptedPassword.String) != nil {
		http.Error(w, "username or password incorrect", http.StatusBadRequest)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(int(user.ID)),
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		fmt.Println("token err ",err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
