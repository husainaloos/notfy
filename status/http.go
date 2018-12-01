package status

import (
	"context"
	"encoding/json"
	"fmt"
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
	errCannotReadBody     = errModel{"cannot read json body", 201}
	errMalformedJSON      = errModel{"malformed json", 202}
	errInsertFailed       = errModel{"an error has occured", 203}
	errBuildMessageFailed = errModel{"an error has occured", 204}
	errWriteMessageFailed = errModel{"an error has occured", 205}
	errInvalid            = func(err error) errModel { return errModel{fmt.Sprintf("invalid request: %v", err), 206} }
	errRetreiveStatus     = errModel{"an error has occured", 206}
)

type getModel struct {
	ID           int       `json:"id"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	LastUpdateAt time.Time `json:"last_update_at"`
}

type APIInterface interface {
	Get(ctx context.Context, id int) (Info, error)
}

// HTTPHandler is the http handler for REST calls to the status API
type HTTPHandler struct {
	api APIInterface
}

// NewHTTPHandler creates new http handler for status
func NewHTTPHandler(api APIInterface) *HTTPHandler {
	return &HTTPHandler{api}
}

// Route configures the routes
func (h *HTTPHandler) Route(r chi.Router) {
	r.Get("/{id}", h.getStatusHandler)
}

func (h *HTTPHandler) getStatusHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogEntry(r)
	idStr := chi.URLParam(r, "id")
	log.WithField("id", idStr).Debug("id value from URL")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.WithField("id", idStr).Debugf("id passed is not a valid integer: %v", err)
		return
	}
	info, err := h.api.Get(r.Context(), id)
	if err != nil {
		switch err {
		case ErrStatusNotFound:
			w.WriteHeader(http.StatusNotFound)
			log.WithField("id", id).Debugf("status not found")
			return
		default:
			h.writeErr(w, r, errRetreiveStatus, http.StatusInternalServerError)
			log.Errorf("failed to retreive the status: %v", err)
			return
		}
	}
	log.WithFields(logrus.Fields{"id": id, "status": info.Status().String()}).Debugf("status found")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(getModel{info.ID(), info.Status().String(), info.CreatedAt(), info.LastUpdateAt()})
}

func (h *HTTPHandler) writeErr(w http.ResponseWriter, r *http.Request, e errModel, status int) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log := logger.GetLogEntry(r)
		log.Errorf("failed writing response for an error: %v", err)
		return
	}
	w.WriteHeader(status)
}
