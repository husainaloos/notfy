package email

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/husainaloos/notfy/logger"
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
	errFailedToQueueEmail = errModel{"failed to queue email", 104}
	errFailedToInitEmail  = errModel{"an error has occured", 105}
	errStatusCreateFailed = errModel{"failed to create status", 106}
	errGetEmailFailed     = errModel{"an error has occured", 107}
)

type postEmailModel struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	CC      []string `json:"cc"`
	BCC     []string `json:"bcc"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

type getEmailModel struct {
	ID      int            `json:"id"`
	From    string         `json:"from"`
	To      []string       `json:"to"`
	CC      []string       `json:"cc"`
	BCC     []string       `json:"bcc"`
	Subject string         `json:"subject"`
	Body    string         `json:"body"`
	History []emailHistory `json:"history"`
}

type emailHistory struct {
	Status string    `json:"status"`
	At     time.Time `json:"at"`
}

type APIInterface interface {
	Queue(context.Context, Email) (Email, error)
	Get(context.Context, int) (Email, error)
}

// HTTPHandler is the handler for Email
type HTTPHandler struct {
	api APIInterface
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
		h.writeErr(w, r, errCannotReadBody, http.StatusInternalServerError)
		log.Errorf("failed to read request body: %v", err)
		return
	}
	defer r.Body.Close()
	var model postEmailModel
	if err := json.Unmarshal(body, &model); err != nil {
		h.writeErr(w, r, errMalformedJSON, http.StatusBadRequest)
		log.Debugf("failed to unmarshal json: %v", err)
		return
	}
	e, err := New(0, model.From, model.To, model.CC, model.BCC, model.Subject, model.Body)
	if err != nil {
		h.writeErr(w, r, errBadRequest(err), http.StatusBadRequest)
		log.Debugf("failed to create email due to validation: %v", err)
		return
	}
	e, err = h.api.Queue(r.Context(), e)
	if err != nil {
		h.writeErr(w, r, errFailedToQueueEmail, http.StatusInternalServerError)
		log.Errorf("failed to queue email: %v", err)
		return
	}
	log.WithFields(logrus.Fields{
		"email_from":    e.From(),
		"email_subject": e.Subject(),
	}).Debugf("email queued")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.buildGetEmailDto(e))
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
	email, err := h.api.Get(r.Context(), id)
	if err != nil {
		switch err {
		case ErrItemNotFound:
			w.WriteHeader(http.StatusNotFound)
			log.WithField("id", id).Debugf("email not found")
			return
		default:
			h.writeErr(w, r, errGetEmailFailed, http.StatusInternalServerError)
			log.Errorf("failed to retreive email: %v", err)
			return
		}
	}

	model := h.buildGetEmailDto(email)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model)
}

func (h *HTTPHandler) buildGetEmailDto(e Email) getEmailModel {
	model := getEmailModel{}
	model.ID = e.ID()
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

	history := make([]emailHistory, 0)
	for _, v := range e.StatusHistory() {
		history = append(history, emailHistory{v.Status().String(), v.At()})
	}
	model.History = history
	return model
}

func (h *HTTPHandler) writeErr(w http.ResponseWriter, r *http.Request, e errModel, status int) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		log := logger.GetLogEntry(r)
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("failed to encode response: %v", err)
		return
	}
}
