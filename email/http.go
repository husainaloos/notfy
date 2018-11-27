package email

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/golang/protobuf/proto"
	"github.com/husainaloos/notfy/dto"
	"github.com/husainaloos/notfy/messaging"
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
)

type emailDto struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	CC      []string `json:"cc"`
	BCC     []string `json:"bcc"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

type statusInfo struct {
	ID        int       `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// HTTPHandler is the handler for Email
type HTTPHandler struct {
	p         messaging.Publisher
	statusAPI *status.API
}

// NewAPIHandler creates a new handler for email reqeusts
func NewAPIHandler(p messaging.Publisher, statusAPI *status.API) *HTTPHandler {
	return &HTTPHandler{p: p, statusAPI: statusAPI}
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
	email, err := New(model.From, model.To, model.CC, model.BCC, model.Subject, model.Body)
	if err != nil {
		api.writeErr(w, errBadRequest(err), http.StatusBadRequest)
		return
	}
	info, err := api.statusAPI.Create(status.Queued)
	if err != nil {
		api.writeErr(w, errStatusCreateFailed, http.StatusInternalServerError)
		logrus.WithFields(lf).Errorf("cannot insert status: %v", err)
		return
	}
	if err := api.publish(email, info.ID()); err != nil {
		api.writeErr(w, errPublishFailed, http.StatusInternalServerError)
		logrus.WithFields(lf).Errorf("cannot publish message: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(statusInfo{
		ID:        info.ID(),
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

func (api *HTTPHandler) publish(m Email, statusID int) error {
	pe := &dto.PublishedEmail{
		From:    m.from.String(),
		To:      make([]string, 0),
		Cc:      make([]string, 0),
		Bcc:     make([]string, 0),
		Subject: m.Subject(),
		Body:    m.Body(),
		Id:      int64(statusID),
	}
	for _, v := range m.To() {
		pe.To = append(pe.To, v.String())
	}
	for _, v := range m.CC() {
		pe.Cc = append(pe.Cc, v.String())
	}
	for _, v := range m.BCC() {
		pe.Bcc = append(pe.Bcc, v.String())
	}
	b, err := proto.Marshal(pe)
	if err != nil {
		return err
	}
	if err := api.p.Publish(b); err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{"size": len(b)}).Info("published message")
	return nil
}
