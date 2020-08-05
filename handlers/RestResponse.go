package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type RestResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

type RestHttpResponse struct {
	StatusCode  int
	Status      string
	ContentType string
	Error       error
	Message     string
	Result      interface{}
}

func NewRestHttpResponse() *RestHttpResponse {
	return &RestHttpResponse{
		StatusCode: http.StatusOK,
		Status:     "success",
	}
}

func (r *RestHttpResponse) Write(w http.ResponseWriter) {
	res := &RestResponse{
		Status:  r.Status,
		Message: r.Message,
		Result:  r.Result,
	}

	if r.Error != nil {
		if r.StatusCode < 400 {
			r.StatusCode = http.StatusInternalServerError
		}
		res.Status = "error"
		res.Message = r.Error.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.StatusCode)

	b, err := json.Marshal(res)
	if err != nil {
		log.Print(err)
		return
	}

	if _, err := w.Write(b); err != nil {
		log.Print(err)
	}
}
