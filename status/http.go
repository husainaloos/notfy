package status

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

type errModel struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

var (
	errCannotReadBody     = errModel{"cannot read json body", 201}
	errMalformedJSON      = errModel{"malformed json", 202}
	errInsertFailed       = errModel{"an error has occured", 203}
	errBuildMessageFailed = errModel{"an error has occured", 204}
	errWriteMessageFailed = errModel{"an error has occured", 205}
	errInvalid            = func(err error) errModel { return errModel{fmt.Sprintf("invalid request: %v", err), 206} }
	errRetreiveStatus     = errModel{"an error has occured", 206}
)

type insertModel struct {
	Status string `json:"status"`
}

func (m *insertModel) validate() error {
	switch m.Status {
	case "Sent":
	case "Queued":
	case "Failed":
		return nil
	default:
		return fmt.Errorf("not valid status: %s. valid status = ", m.Status)
	}
	return nil
}

type getModel struct {
	ID           int       `json:"id"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	LastUpdateAt time.Time `json:"last_update_at"`
}

// HTTPHandler is the http handler for REST calls to the status API
type HTTPHandler struct {
	api *API
}

// NewHTTPHandler creates new http handler for status
func NewHTTPHandler(api *API) *HTTPHandler {
	return &HTTPHandler{api}
}

// Route configures the routes
func (h *HTTPHandler) Route(r chi.Router) {
	r.Post("/", h.postStatusHandler)
	r.Get("/{id}", h.getStatusHandler)
}

func (h *HTTPHandler) postStatusHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.writeErr(w, errCannotReadBody, http.StatusInternalServerError)
		logrus.Errorf("could not read body: %v", err)
		return
	}
	var im insertModel
	if err := json.Unmarshal(body, &im); err != nil {
		h.writeErr(w, errMalformedJSON, http.StatusBadRequest)
		logrus.Debugf("could not unmarshal message: %v", err)
		return
	}
	if err := im.validate(); err != nil {
		h.writeErr(w, errInvalid(err), http.StatusBadRequest)
		logrus.Debugf("failed validation: %v", err)
		return
	}
	si, err := h.api.Create(Sent)
	if err != nil {
		h.writeErr(w, errInsertFailed, http.StatusInternalServerError)
		logrus.Errorf("failed creating info: %v", err)
		return
	}
	gm := getModel{
		ID:        si.ID(),
		Status:    si.Status().String(),
		CreatedAt: si.CreatedAt(),
	}
	j, err := json.Marshal(gm)
	if err != nil {
		h.writeErr(w, errBuildMessageFailed, http.StatusInternalServerError)
		logrus.Errorf("failed json.Marshal: %v", err)
		return
	}
	if _, err := w.Write(j); err != nil {
		h.writeErr(w, errWriteMessageFailed, http.StatusInternalServerError)
		logrus.Errorf("failed writing response: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *HTTPHandler) writeErr(w http.ResponseWriter, e errModel, status int) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Errorf("failed writing response for an error: %v", err)
		return
	}
	w.WriteHeader(status)
}

func (h *HTTPHandler) getStatusHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	logrus.WithFields(logrus.Fields{"idStr": idStr}).Debug("id value from URL")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logrus.Debugf("id passed is not a valid integer: %v", err)
		return
	}
	info, err := h.api.Get(id)
	if err != nil {
		switch err {
		case errNotFound:
			w.WriteHeader(http.StatusNotFound)
			return
		default:
			h.writeErr(w, errRetreiveStatus, http.StatusInternalServerError)
			logrus.Errorf("failed to retreive the status: %v", err)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(getModel{info.ID(), info.Status().String(), info.CreatedAt(), info.LastUpdateAt()})
}
