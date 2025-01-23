package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/amirdaraby/go-todo-list-api/internal/auth"
	"github.com/amirdaraby/go-todo-list-api/internal/db"
	"github.com/amirdaraby/go-todo-list-api/internal/models"
	"github.com/amirdaraby/go-todo-list-api/internal/utils/jsonresponse"
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
		jsonresponse.New().SetMessage(jsonresponse.BadRequestMessage).Failed(w, http.StatusBadRequest)
		return
	}

	response := userShowResponse{
		ID:       user.ID,
		UserName: user.UserName,
	}

	jsonresponse.New().SetData(response).Success(w, http.StatusOK)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	var updatedUser userUpdateRequest

	err := json.NewDecoder(r.Body).Decode(&updatedUser)

	if err != nil {
		jsonresponse.New().SetMessage(jsonresponse.BadRequestMessage).Failed(w, http.StatusBadRequest)
		return
	}

	validate := validator.New()

	err = validate.Struct(updatedUser)

	if err != nil {
		errors := err.(validator.ValidationErrors)
		jsonresponse.New().SetMessage(fmt.Sprintf("Validation failed: %s", errors)).Failed(w, http.StatusUnprocessableEntity)
		return
	}

	var user models.User

	gorm := db.GetDb()

	tx := gorm.Model(&user).Where("id = ?", r.Context().Value(auth.AuthIdKey("user_id"))).First(&user)

	if tx.RowsAffected != 1 {
		jsonresponse.New().SetMessage(jsonresponse.BadRequestMessage).Failed(w, http.StatusBadRequest)
		return
	}

	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)

		if err != nil {
			jsonresponse.New().SetMessage(jsonresponse.BadRequestMessage).Failed(w, http.StatusBadRequest)
			return
		}
		user.Password = string(hashedPassword)
	}

	if updatedUser.UserName != "" && updatedUser.UserName != user.UserName {
		var userWithUserName models.User

		tx = gorm.Model(&userWithUserName).Where("user_name = ?", updatedUser.UserName)

		if tx.RowsAffected != 0 {
			jsonresponse.New().SetMessage("user_name is picked by other user").Failed(w, http.StatusBadRequest)
			return
		}
		user.UserName = updatedUser.UserName
	}

	gorm.Save(&user)

	jsonresponse.New().SetMessage("user updated").Success(w, http.StatusAccepted)
}
