package jsonresponse

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/amirdaraby/go-todo-list-api/internal/paginator"
)

const (
	success = "success"
	failed  = "failed"
)

const (
	BadRequestMessage          = "bad request"
	UnauthorizedMessage        = "unauthorized"
	InternalServerErrorMessage = "internal server error"
	NotFoundMessage            = "not found"
)

type response struct {
	Status  string      `json:"status"`
	Message *string     `json:"message"`
	Data    interface{} `json:"data"`
}

type paginatedResponse struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	response
}

func New() *response {
	return &response{}
}

func (resp *response) SetMessage(msg string) *response {
	resp.Message = &msg

	return resp
}

func (resp *response) SetData(data interface{}) *response {
	resp.Data = data

	return resp
}

func (resp *response) Success(w http.ResponseWriter, status int) {
	resp.Status = success

	marshalledResponse, err := json.Marshal(resp)

	if err != nil {
		log.Printf("marshalling failed err: %s", err)
		return
	}

	w.WriteHeader(status)
	w.Write(marshalledResponse)
}

func (resp *response) Failed(w http.ResponseWriter, status int) {
	resp.Status = failed

	marshalledResponse, err := json.Marshal(resp)

	if err != nil {
		log.Printf("marshalling failed err: %s", err)
		return
	}

	w.WriteHeader(status)
	w.Write(marshalledResponse)
}

func (resp *response) SuccessWithPagination(w http.ResponseWriter, r *http.Request, status int) {
	pageQuery := r.URL.Query().Get("page")

	page, _ := strconv.Atoi(pageQuery)

	if page <= 0 {
		page = 1
	}

	resp.Status = success
	paginatedResp := paginatedResponse{
		Page:     page,
		PerPage:  paginator.PerPage,
		response: *resp,
	}

	marshalledResponse, err := json.Marshal(paginatedResp)

	if err != nil {
		log.Printf("marshalling failed err: %s", err)
		return
	}

	w.WriteHeader(status)
	w.Write(marshalledResponse)
}
