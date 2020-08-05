package handlers

import (
	"errors"
	"golang-chat/lib"
	"io/ioutil"
	"log"
	"net/http"
)

type RestApiHandler struct {
	store  *lib.Store
	stream *lib.Stream
}

var _ http.Handler = (*RestApiHandler)(nil)

var (
	ErrMethodNotAllowed = errors.New("Method not allowed")
	ErrRequestBodyEmpty = errors.New("Request body is empty")
)

func NewRestApiHandler(store *lib.Store, stream *lib.Stream) *RestApiHandler {
	return &RestApiHandler{store, stream}
}

func (h *RestApiHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	res := NewRestHttpResponse()
	err := h.HandleRequest(req, res)
	if err != nil {
		log.Print(err)
		res.Error = err
	} else if res.Error != nil {
		log.Print(res.Error)
	}
	res.Write(w)
}

func (h *RestApiHandler) HandleRequest(req *http.Request, res *RestHttpResponse) error {
	if req.Method == "POST" {
		code, err := h.PostMessage(req)
		res.StatusCode = code
		return err
	}
	res.Result = h.GetMessages()
	return nil
}

func (h *RestApiHandler) GetMessages() []string {
	return h.store.FindAll()
}

func (h *RestApiHandler) PostMessage(req *http.Request) (int, error) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return 400, err
	}
	message := string(b)
	if message == "" {
		return 400, ErrRequestBodyEmpty
	}
	h.store.Save(message)
	h.stream.Fire(message)
	return 200, nil
}
