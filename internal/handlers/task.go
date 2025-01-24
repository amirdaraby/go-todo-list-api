package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/amirdaraby/go-todo-list-api/internal/auth"
	"github.com/amirdaraby/go-todo-list-api/internal/db"
	"github.com/amirdaraby/go-todo-list-api/internal/models"
	"github.com/amirdaraby/go-todo-list-api/internal/paginator"
	"github.com/amirdaraby/go-todo-list-api/internal/utils/jsonresponse"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type updateTaskRequest struct {
	Title       *string `json:"title" validate:"omitempty,min=1,max=40"`
	Description *string `json:"description" validate:"omitempty,max=255"`
}

func IndexTask(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth.AuthIdKey("user_id"))

	gorm := db.GetDb()

	var tasks []models.Task

	tx := gorm.Scopes(paginator.Paginate(r)).Model(&models.Task{}).Where("user_id = ?", userId).Order("id DESC").Find(&tasks)

	if tx.RowsAffected < 1 {
		jsonresponse.New().SetMessage("no tasks found").Failed(w, http.StatusNotFound)
		return
	}

	jsonresponse.New().SetData(tasks).SetMessage("all tasks").SuccessWithPagination(w, r, http.StatusOK)
}

func IndexUnDoneTasks(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth.AuthIdKey("user_id"))

	gorm := db.GetDb()

	var tasks []models.Task

	tx := gorm.Scopes(paginator.Paginate(r)).Model(&models.Task{}).Where("user_id = ?", userId).Where("done_at IS NULL").Order("id DESC").Find(&tasks)

	if tx.RowsAffected < 1 {
		jsonresponse.New().SetMessage("no undone tasks found").Failed(w, http.StatusNotFound)
		return
	}

	jsonresponse.New().SetData(tasks).SetMessage("undone tasks").SuccessWithPagination(w, r, http.StatusOK)
}

func IndexDoneTasks(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth.AuthIdKey("user_id"))

	gorm := db.GetDb()

	var tasks []models.Task

	tx := gorm.Scopes(paginator.Paginate(r)).Model(&models.Task{}).Where("user_id = ?", userId).Where("done_at IS NOT NULL").Order("id DESC").Find(&tasks)

	if tx.RowsAffected < 1 {
		jsonresponse.New().SetMessage("no done tasks found").Failed(w, http.StatusNotFound)
		return
	}

	jsonresponse.New().SetData(tasks).SetMessage("done tasks").SuccessWithPagination(w, r, http.StatusOK)
}

func StoreTask(w http.ResponseWriter, r *http.Request) {
	userIdContext := r.Context().Value(auth.AuthIdKey("user_id"))

	userId, ok := userIdContext.(uint)

	if !ok {
		log.Printf("user_id type must be uint, uint assertion failed")
		jsonresponse.New().SetMessage(jsonresponse.BadRequestMessage).Failed(w, http.StatusBadRequest)
		return
	}

	var task models.Task

	err := json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		jsonresponse.New().SetMessage(jsonresponse.BadRequestMessage).Failed(w, http.StatusBadRequest)
		return
	}

	task.UserID = userId

	validate := validator.New()

	err = validate.Struct(task)

	if err != nil {
		errors := err.(validator.ValidationErrors)
		jsonresponse.New().SetMessage(fmt.Sprintf("validation failed %s", errors)).Failed(w, http.StatusUnprocessableEntity)
		return
	}

	gorm := db.GetDb()

	tx := gorm.Model(&task).Create(&task)

	if tx.RowsAffected != 1 {
		jsonresponse.New().SetMessage("something goes wrong in task creation").Failed(w, http.StatusInternalServerError)
		return
	}

	jsonresponse.New().Success(w, http.StatusCreated)
}

func ShowTask(w http.ResponseWriter, r *http.Request) {
	taskIdInPath := mux.Vars(r)["id"]

	if taskIdInPath == "" {
		jsonresponse.New().SetMessage(jsonresponse.NotFoundMessage).Failed(w, http.StatusNotFound)
		return
	}

	taskId, err := strconv.ParseUint(taskIdInPath, 10, 32)

	if err != nil {
		log.Println(err)
		jsonresponse.New().SetMessage(jsonresponse.BadRequestMessage).Failed(w, http.StatusBadRequest)
		return
	}

	userId := r.Context().Value(auth.AuthIdKey("user_id"))

	gorm := db.GetDb()

	var task models.Task

	tx := gorm.Model(&task).Where("id = ?", taskId).Where("user_id = ?", userId).First(&task)

	if tx.RowsAffected != 1 {
		jsonresponse.New().SetMessage(jsonresponse.NotFoundMessage).Failed(w, http.StatusNotFound)
		return
	}

	jsonresponse.New().SetData(task).Success(w, http.StatusOK)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskId := mux.Vars(r)["id"]

	if taskId == "" {
		jsonresponse.New().SetMessage(jsonresponse.NotFoundMessage).Failed(w, http.StatusNotFound)
		return
	}

	var updatedTask updateTaskRequest

	err := json.NewDecoder(r.Body).Decode(&updatedTask)

	if err != nil {
		log.Println(err)
		jsonresponse.New().SetMessage(jsonresponse.BadRequestMessage).Success(w, http.StatusBadRequest)
		return
	}

	userId := r.Context().Value(auth.AuthIdKey("user_id"))

	var task models.Task

	gorm := db.GetDb()

	tx := gorm.Model(&task).Where("id = ?", taskId).Where("user_id = ?", userId).First(&task)

	if tx.RowsAffected != 1 {
		jsonresponse.New().SetMessage(jsonresponse.NotFoundMessage).Failed(w, http.StatusNotFound)
		return
	}

	if updatedTask.Title != nil {
		task.Title = *updatedTask.Title
	}

	if updatedTask.Description != nil {
		task.Description = updatedTask.Description
	}

	gorm.Save(&task)

	jsonresponse.New().SetMessage("task updated").Success(w, http.StatusAccepted)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskId := mux.Vars(r)["id"]

	if taskId == "" {
		jsonresponse.New().SetMessage(jsonresponse.NotFoundMessage).Failed(w, http.StatusNotFound)
		return
	}

	userId := r.Context().Value(auth.AuthIdKey("user_id"))

	gorm := db.GetDb()

	tx := gorm.Delete(&models.Task{}, "id = ?", taskId, "user_id = ?", userId)

	if tx.RowsAffected != 1 {
		jsonresponse.New().SetMessage(jsonresponse.NotFoundMessage).Failed(w, http.StatusNotFound)
		return
	}

	jsonresponse.New().SetMessage("task deleted").Success(w, http.StatusAccepted)
}

func ToggleDoneTask(w http.ResponseWriter, r *http.Request) {
	taskId := mux.Vars(r)["id"]

	if taskId == "" {
		jsonresponse.New().SetMessage(jsonresponse.BadRequestMessage).Failed(w, http.StatusBadRequest)
		return
	}

	userId := r.Context().Value(auth.AuthIdKey("user_id"))

	gorm := db.GetDb()

	var task models.Task

	tx := gorm.Model(&task).Where("id = ?", taskId).Where("user_id = ?", userId).First(&task)

	if tx.RowsAffected != 1 {
		jsonresponse.New().SetMessage(jsonresponse.NotFoundMessage).Failed(w, http.StatusNotFound)
		return
	}

	var message string

	if task.DoneAt == nil {
		now := time.Now()
		task.DoneAt = &now
		message = "task marked as done"
	} else {
		task.DoneAt = nil
		message = "task marked as undone"
	}

	gorm.Save(&task)

	jsonresponse.New().SetMessage(message).Success(w, http.StatusAccepted)
}
