package email

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/husainaloos/notfy/status"
)

func Test_sendEmailHandler(t *testing.T) {
	passQueue := func(Email) (Email, status.Info, error) { return Email{}, status.Info{}, nil }
	failQueue := func(Email) (Email, status.Info, error) { return Email{}, status.Info{}, errors.New("queue failed") }
	t.Run("should return bad request if the body is invalid", func(t *testing.T) {
		api := NewHTTPHandler(NewMockAPI(passQueue))
		w := httptest.NewRecorder()
		body := strings.NewReader(`{"from" : "bademail.com", "to" : ["friend@gmail.com"]`)
		r := httptest.NewRequest(http.MethodPost, "http://localhost", body)
		api.sendEmailHandler(w, r)
		if w.Code != http.StatusBadRequest {
			t.Errorf("sendEmailHandler(): expected %d, but got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("should return bad request if the body is invalid json", func(t *testing.T) {
		api := NewHTTPHandler(NewMockAPI(passQueue))
		w := httptest.NewRecorder()
		body := strings.NewReader(`{"field" : "bad json"`)
		r := httptest.NewRequest(http.MethodPost, "http://localhost", body)
		api.sendEmailHandler(w, r)
		if w.Code != http.StatusBadRequest {
			t.Errorf("sendEmailHandler(): expected %d, but got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("should return 500 if the API fails", func(t *testing.T) {
		api := NewHTTPHandler(NewMockAPI(failQueue))
		w := httptest.NewRecorder()
		body := strings.NewReader(`{"from" : "email@gmail.com", "to" : ["fiend@gmail.com"]}`)
		r := httptest.NewRequest(http.MethodPost, "http://localhost", body)
		api.sendEmailHandler(w, r)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("sendEmailHandler(): expected %d, but got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("should return 200 if message is valid", func(t *testing.T) {
		api := NewHTTPHandler(NewMockAPI(passQueue))
		w := httptest.NewRecorder()
		body := strings.NewReader(`{"from" : "email@gmail.com", "to" : ["fiend@gmail.com"], "cc": null, "bcc": [],  "body" : "body"}`)
		r := httptest.NewRequest(http.MethodPost, "http://localhost", body)
		api.sendEmailHandler(w, r)
		if w.Code != http.StatusOK {
			t.Errorf("sendEmailHandler(): expected %d, but got %d", http.StatusOK, w.Code)
		}
	})
}
