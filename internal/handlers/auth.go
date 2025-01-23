package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/amirdaraby/go-todo-list-api/internal/auth"
	"github.com/amirdaraby/go-todo-list-api/internal/db"
	"github.com/amirdaraby/go-todo-list-api/internal/models"
	"github.com/amirdaraby/go-todo-list-api/internal/utils/jsonresponse"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		jsonresponse.New().SetMessage(jsonresponse.BadRequestMessage).Failed(w, http.StatusBadRequest)
		return
	}

	validate := validator.New()

	err = validate.Struct(user)

	if err != nil {
		errors := err.(validator.ValidationErrors)
		jsonresponse.New().SetMessage(fmt.Sprintf("Validation failed: %s", errors)).Failed(w, http.StatusUnprocessableEntity)
		return
	}

	gorm := db.GetDb()

	tx := gorm.Model(user).Where("user_name = ?", user.UserName).First(&user)

	if tx.RowsAffected != 0 {
		jsonresponse.New().SetMessage("user_name is picked by other user").Failed(w, http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		jsonresponse.New().SetMessage(jsonresponse.InternalServerErrorMessage).Failed(w, http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)

	tx = gorm.Model(user).Create(&user)

	if tx.RowsAffected != 1 {
		jsonresponse.New().SetMessage(jsonresponse.InternalServerErrorMessage).Failed(w, http.StatusInternalServerError)
		return
	}

	jsonresponse.New().SetMessage("user created").Success(w, http.StatusCreated)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {

	bodyString, err := io.ReadAll(r.Body)

	if err != nil {
		jsonresponse.New().SetMessage(jsonresponse.BadRequestMessage).Failed(w, http.StatusBadRequest)
		return
	}

	body := struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}{}

	err = json.Unmarshal(bodyString, &body)

	if err != nil {
		jsonresponse.New().SetMessage(jsonresponse.BadRequestMessage).Failed(w, http.StatusBadRequest)
		return
	}

	gorm := db.GetDb()

	user := models.User{}

	tx := gorm.Model(user).Where("user_name = ?", body.UserName).First(&user)

	if tx.RowsAffected != 1 {
		jsonresponse.New().SetMessage(jsonresponse.BadRequestMessage).Failed(w, http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		jsonresponse.New().SetMessage(jsonresponse.BadRequestMessage).Failed(w, http.StatusBadRequest)
		return
	}

	token, err := auth.NewToken(user.ID)

	if err != nil {
		jsonresponse.New().SetMessage(jsonresponse.InternalServerErrorMessage).Failed(w, http.StatusInternalServerError)
		return
	}

	response := struct {
		Token string `json:"token"`
	}{}

	response.Token = token

	jsonresponse.New().SetMessage("login successful").SetData(response).Success(w, http.StatusOK)
}
