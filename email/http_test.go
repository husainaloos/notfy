package email

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/husainaloos/notfy/status"
)

type mockAPI struct {
	queue func(Email) (Email, status.Info, error)
	get   func(int) (Email, status.Info, error)
}

func newMockAPI(
	queue func(Email) (Email, status.Info, error),
	get func(int) (Email, status.Info, error)) *mockAPI {
	return &mockAPI{queue, get}
}
func (api *mockAPI) Queue(e Email) (Email, status.Info, error) { return api.queue(e) }
func (api *mockAPI) Get(id int) (Email, status.Info, error)    { return api.get(id) }

func Test_sendEmailHandler(t *testing.T) {
	passQueue := func(Email) (Email, status.Info, error) { return Email{}, status.Info{}, nil }
	failQueue := func(Email) (Email, status.Info, error) { return Email{}, status.Info{}, errors.New("queue failed") }
	passGet := func(int) (Email, status.Info, error) { return Email{}, status.Info{}, nil }
	tt := []struct {
		name           string
		queuef         func(Email) (Email, status.Info, error)
		getf           func(int) (Email, status.Info, error)
		body           string
		expectedStatus int
	}{
		{
			name:           "should return bad request if the body is invalid",
			queuef:         passQueue,
			getf:           passGet,
			body:           `{"from" : "bademail.com", "to" : ["friend@gmail.com"]`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "should return bad request if the body is invalid json",
			queuef:         passQueue,
			getf:           passGet,
			body:           `{"field" : "bad json"`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "should return 500 if the API fails",
			queuef:         failQueue,
			getf:           passGet,
			body:           `{"from" : "email@gmail.com", "to" : ["fiend@gmail.com"]}`,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "should return 200 if message is valid",
			queuef:         passQueue,
			getf:           passGet,
			body:           `{"from" : "email@gmail.com", "to" : ["fiend@gmail.com"], "cc": null, "bcc": [],  "body" : "body"}`,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tst := range tt {
		t.Run(tst.name, func(t *testing.T) {
			api := NewHTTPHandler(newMockAPI(tst.queuef, tst.getf))
			w := httptest.NewRecorder()
			body := strings.NewReader(tst.body)
			r := httptest.NewRequest(http.MethodPost, "http://localhost", body)
			api.sendEmailHandler(w, r)
			if w.Code != tst.expectedStatus {
				t.Errorf("sendEmailHandler(): got %d when was expecting %d", w.Code, tst.expectedStatus)
			}
		})
	}
}
