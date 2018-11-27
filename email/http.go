package email

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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
	errGetEmailFailed     = errModel{"an error has occured", 107}
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
	StatusID  int       `json:"status_id"`
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
func (h *HTTPHandler) Route(r chi.Router) {
	r.Post("/", h.sendEmailHandler)
	r.Get("/{id}", h.getEmailHandler)
}

func (h *HTTPHandler) sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	lf := logrus.Fields{"reqID": reqID}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.writeErr(w, errCannotReadBody, http.StatusInternalServerError)
		logrus.WithFields(lf).Errorf("failed to read request body: %v", err)
		return
	}
	defer r.Body.Close()
	var model emailDto
	if err := json.Unmarshal(body, &model); err != nil {
		h.writeErr(w, errMalformedJSON, http.StatusBadRequest)
		return
	}
	e, err := New(model.From, model.To, model.CC, model.BCC, model.Subject, model.Body)
	if err != nil {
		h.writeErr(w, errBadRequest(err), http.StatusBadRequest)
		return
	}
	email, info, err := h.emailAPI.Queue(e)
	if err != nil {
		h.writeErr(w, errCannotQueue, http.StatusInternalServerError)
		logrus.WithFields(lf).Errorf("could not queue email: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(emailStatus{
		ID:        email.ID(),
		StatusID:  info.ID(),
		Status:    info.Status().String(),
		CreatedAt: info.CreatedAt(),
	})
}

func (h *HTTPHandler) getEmailHandler(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	lf := logrus.Fields{"reqID": reqID}
	idStr := chi.URLParam(r, "id")
	logrus.WithFields(lf).WithFields(logrus.Fields{"idStr": idStr}).Debug("id value from URL")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logrus.Debugf("id passed is not a valid integer: %v", err)
		return
	}
	_, _, err = h.emailAPI.Get(id)
	if err != nil {
		h.writeErr(w, errGetEmailFailed, http.StatusInternalServerError)
		logrus.WithFields(lf).Errorf("failed to retreive email: %v", err)
		return
	}
	return
}

func (h *HTTPHandler) writeErr(w http.ResponseWriter, e errModel, status int) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Errorf("failed to encode response: %v", err)
		return
	}
}
