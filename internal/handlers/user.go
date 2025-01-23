package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/amirdaraby/go-todo-list-api/internal/auth"
	"github.com/amirdaraby/go-todo-list-api/internal/db"
	"github.com/amirdaraby/go-todo-list-api/internal/models"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type userShowResponse struct {
	ID       uint   `json:"id"`
	UserName string `json:"user_name"`
}

type userUpdateRequest struct {
	UserName string `json:"user_name" validate:"omitempty,min=2,max=255"`
	Password string `json:"password" validate:"omitempty,min=8,max=255"`
}

func ShowUser(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(auth.AuthIdKey("user_id"))

	gorm := db.GetDb()

	user := models.User{}

	tx := gorm.Model(user).Where("id = ?", userId).First(&user)

	if tx.RowsAffected != 1 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	response := userShowResponse{
		ID:       user.ID,
		UserName: user.UserName,
	}

	marshalledResponse, err := json.Marshal(response)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(marshalledResponse)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	var updatedUser userUpdateRequest

	err := json.NewDecoder(r.Body).Decode(&updatedUser)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	validate := validator.New()

	err = validate.Struct(updatedUser)

	if err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation failed: %s", errors), http.StatusUnprocessableEntity)
		return
	}

	var user models.User

	gorm := db.GetDb()

	tx := gorm.Model(&user).Where("id = ?", r.Context().Value(auth.AuthIdKey("user_id"))).First(&user)

	if tx.RowsAffected != 1 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)

		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		user.Password = string(hashedPassword)
	}

	if updatedUser.UserName != "" && updatedUser.UserName != user.UserName {
		var userWithUserName models.User

		tx = gorm.Model(&userWithUserName).Where("user_name = ?", updatedUser.UserName)

		if tx.RowsAffected != 0 {
			http.Error(w, "username is picked by other user", http.StatusBadRequest)
			return
		}
		user.UserName = updatedUser.UserName
	}

	gorm.Save(&user)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}
