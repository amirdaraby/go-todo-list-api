package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/amirdaraby/go-todo-list-api/internal/auth"
	"github.com/amirdaraby/go-todo-list-api/internal/db"
	"github.com/amirdaraby/go-todo-list-api/internal/models"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	validate := validator.New()

	err = validate.Struct(user)

	if err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation failed: %s", errors), http.StatusUnprocessableEntity)
		return
	}

	gorm := db.GetDb()

	tx := gorm.Model(user).Where("user_name = ?", user.UserName).First(&user)

	if tx.RowsAffected != 0 {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)

	tx = gorm.Model(user).Create(&user)

	if tx.RowsAffected != 1 {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User created successfuly"))
}

func LoginUser(w http.ResponseWriter, r *http.Request) {

	bodyString, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	body := struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}{}

	err = json.Unmarshal(bodyString, &body)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	gorm := db.GetDb()

	user := models.User{}

	tx := gorm.Model(user).Where("user_name = ?", body.UserName).First(&user)

	if tx.RowsAffected != 1 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	token, err := auth.NewToken(user.ID)

	if err != nil {
		log.Fatalln(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	response := struct {
		Token string `json:"token"`
	}{}

	response.Token = token

	marshalledResponse, err := json.Marshal(response)

	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(marshalledResponse)
}
