package email

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_sendEmailHandler(t *testing.T) {
	t.Run("should return bad request if the body is invalid", func(t *testing.T) {
		api := NewHTTPHandler(NewMockAPI())
		w := httptest.NewRecorder()
		body := strings.NewReader(`{"from" : "bademail.com", "to" : ["friend@gmail.com"]`)
		r := httptest.NewRequest(http.MethodPost, "http://localhost", body)
		api.sendEmailHandler(w, r)
		if w.Code != http.StatusBadRequest {
			t.Errorf("sendEmailHandler(): expected %d, but got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("should return bad request if the body is invalid json", func(t *testing.T) {
		api := NewHTTPHandler(NewMockAPI())
		w := httptest.NewRecorder()
		body := strings.NewReader(`{"field" : "bad json"`)
		r := httptest.NewRequest(http.MethodPost, "http://localhost", body)
		api.sendEmailHandler(w, r)
		if w.Code != http.StatusBadRequest {
			t.Errorf("sendEmailHandler(): expected %d, but got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("should return 500 if the API fails", func(t *testing.T) {
		api := NewHTTPHandler(NewMockAPI())
		w := httptest.NewRecorder()
		body := strings.NewReader(`{"from" : "email@gmail.com", "to" : ["fiend@gmail.com"]}`)
		r := httptest.NewRequest(http.MethodPost, "http://localhost", body)
		api.sendEmailHandler(w, r)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("sendEmailHandler(): expected %d, but got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("should return 200 if message is valid", func(t *testing.T) {
		api := NewHTTPHandler(NewMockAPI())
		w := httptest.NewRecorder()
		body := strings.NewReader(`{"from" : "email@gmail.com", "to" : ["fiend@gmail.com"], "cc": null, "bcc": [],  "body" : "body"}`)
		r := httptest.NewRequest(http.MethodPost, "http://localhost", body)
		api.sendEmailHandler(w, r)
		if w.Code != http.StatusOK {
			t.Errorf("sendEmailHandler(): expected %d, but got %d", http.StatusOK, w.Code)
		}
	})
}
