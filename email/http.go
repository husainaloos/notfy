package email

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/husainaloos/notfy/logger"
	"github.com/husainaloos/notfy/status"
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

type postEmailDto struct {
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

type getEmailDto struct {
	From     string   `json:"from"`
	To       []string `json:"to"`
	CC       []string `json:"cc"`
	BCC      []string `json:"bcc"`
	Subject  string   `json:"subject"`
	Body     string   `json:"body"`
	StatusID int      `json:"status_id"`
	Status   string   `json:"status"`
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
	log := logger.GetLogEntry(r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.writeErr(w, errCannotReadBody, http.StatusInternalServerError)
		log.Errorf("failed to read request body: %v", err)
		return
	}
	defer r.Body.Close()
	var model postEmailDto
	if err := json.Unmarshal(body, &model); err != nil {
		h.writeErr(w, errMalformedJSON, http.StatusBadRequest)
		log.Debugf("failed to unmarshal json: %v", err)
		return
	}
	e, err := New(model.From, model.To, model.CC, model.BCC, model.Subject, model.Body)
	if err != nil {
		h.writeErr(w, errBadRequest(err), http.StatusBadRequest)
		log.Debugf("failed to create email due to validation: %v", err)
		return
	}
	email, info, err := h.emailAPI.Queue(e)
	if err != nil {
		h.writeErr(w, errCannotQueue, http.StatusInternalServerError)
		log.Errorf("could not queue email: %v", err)
		return
	}
	log.WithFields(logrus.Fields{
		"email_status":  info.Status().String(),
		"email_subject": email.Subject(),
	}).Debugf("email queued")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(emailStatus{
		ID:        email.ID(),
		StatusID:  info.ID(),
		Status:    info.Status().String(),
		CreatedAt: info.CreatedAt(),
	})
}

func (h *HTTPHandler) getEmailHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogEntry(r)
	idStr := chi.URLParam(r, "id")
	log.WithField("id", idStr).Debug("id value from URL")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.WithField("id", idStr).Debugf("id passed is not a valid integer: %v", err)
		return
	}
	email, info, err := h.emailAPI.Get(id)
	if err != nil {
		switch err {
		case ErrEmailNotFound:
			w.WriteHeader(http.StatusNotFound)
			log.WithField("id", id).Debugf("email not found")
			return
		default:
			h.writeErr(w, errGetEmailFailed, http.StatusInternalServerError)
			log.Errorf("failed to retreive email: %v", err)
			return
		}
	}

	model := h.buildGetEmailDto(email, info)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model)
}

func (h *HTTPHandler) buildGetEmailDto(e Email, i status.Info) getEmailDto {
	model := getEmailDto{}
	model.Body = e.Body()
	model.Subject = e.Subject()

	from := e.From()
	model.From = from.String()

	tos := []string{}
	for _, addr := range e.To() {
		tos = append(tos, addr.String())
	}

	ccs := []string{}
	for _, addr := range e.CC() {
		ccs = append(ccs, addr.String())
	}

	bccs := []string{}
	for _, addr := range e.BCC() {
		bccs = append(bccs, addr.String())
	}
	model.To = tos
	model.CC = ccs
	model.BCC = bccs
	model.StatusID = i.ID()
	model.Status = i.Status().String()
	return model
}

func (h *HTTPHandler) writeErr(w http.ResponseWriter, e errModel, status int) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Errorf("failed to encode response: %v", err)
		return
	}
}
