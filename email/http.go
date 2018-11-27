package email

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

type errModel struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

var (
	errCannotReadBody     = errModel{"cannot ready body", 101}
	errMalformedJSON      = errModel{"bad json", 102}
	errBadRequest         = func(e error) errModel { return errModel{fmt.Sprintf("invalid request: %v", e), 103} }
	errPublishFailed      = errModel{"failed to publish message", 104}
	errStatusCreateFailed = errModel{"failed to create status", 105}
	errCannotQueue        = errModel{"failed to schedule the email", 106}
)

type emailDto struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	CC      []string `json:"cc"`
	BCC     []string `json:"bcc"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

type emailStatus struct {
	ID        int       `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// HTTPHandler is the handler for Email
type HTTPHandler struct {
	emailAPI APIInterface
}

// NewHTTPHandler creates a new handler for email reqeusts
func NewHTTPHandler(emailAPI APIInterface) *HTTPHandler {
	return &HTTPHandler{emailAPI}
}

// Route builds the routing for the email handlers
func (api *HTTPHandler) Route(r chi.Router) {
	r.Post("/", api.sendEmailHandler)
}

func (api *HTTPHandler) sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	lf := logrus.Fields{"reqID": reqID}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.writeErr(w, errCannotReadBody, http.StatusInternalServerError)
		logrus.WithFields(lf).Errorf("failed to read request body: %v", err)
		return
	}
	defer r.Body.Close()
	var model emailDto
	if err := json.Unmarshal(body, &model); err != nil {
		api.writeErr(w, errMalformedJSON, http.StatusBadRequest)
		return
	}
	e, err := New(model.From, model.To, model.CC, model.BCC, model.Subject, model.Body)
	if err != nil {
		api.writeErr(w, errBadRequest(err), http.StatusBadRequest)
		return
	}
	email, info, err := api.emailAPI.Queue(e)
	if err != nil {
		api.writeErr(w, errCannotQueue, http.StatusInternalServerError)
		logrus.WithFields(lf).Errorf("could not queue email: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(emailStatus{
		ID:        email.ID(),
		Status:    info.Status().String(),
		CreatedAt: info.CreatedAt(),
	})
}

func (api *HTTPHandler) writeErr(w http.ResponseWriter, e errModel, status int) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Errorf("failed to encode response: %v", err)
		return
	}
}
