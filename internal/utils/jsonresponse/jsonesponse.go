package jsonresponse

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	success = "success"
	failed  = "failed"
)

const (
	BadRequestMessage          = "bad request"
	UnauthorizedMessage        = "unauthorized"
	InternalServerErrorMessage = "internal server error"
)

type response struct {
	Status  string      `json:"status"`
	Message *string     `json:"message"`
	Data    interface{} `json:"data"`
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
