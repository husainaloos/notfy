package email

import (
	"context"
	"errors"
	"flag"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi"
)

var (
	update     = flag.Bool("update", false, "update golden files")
	testFolder = "test-fixtures"
)

type mockAPI struct {
	queue func(Email) (Email, error)
	get   func(int) (Email, error)
}

func (api *mockAPI) Queue(ctx context.Context, e Email) (Email, error) { return api.queue(e) }
func (api *mockAPI) Get(ctx context.Context, id int) (Email, error)    { return api.get(id) }

func TestPostEmailHandler(t *testing.T) {
	email, _ := New(0, "from@example.com", []string{"to@example.com"}, nil, nil, "subject", "body")
	var (
		passQueue = func(Email) (Email, error) { return email, nil }
		failQueue = func(Email) (Email, error) { return email, errors.New("queue failed") }
	)
	tt := []struct {
		name   string
		queuef func(Email) (Email, error)
		getf   func(int) (Email, error)
		body   string
		status int
	}{
		{
			name:   "should return bad request if the body is invalid",
			queuef: passQueue,
			body:   `{"from" : "bademail.com", "to" : ["friend@gmail.com"]`,
			status: http.StatusBadRequest,
		},
		{
			name:   "should return bad request if the body is invalid json",
			queuef: passQueue,
			body:   `{"field" : "bad json"`,
			status: http.StatusBadRequest,
		},
		{
			name:   "should return 500 if queue fails",
			queuef: failQueue,
			body:   `{"from" : "email@gmail.com", "to" : ["fiend@gmail.com"]}`,
			status: http.StatusInternalServerError,
		},

		{
			name:   "should return 200 if message is valid",
			queuef: passQueue,
			body:   `{"from" : "email@gmail.com", "to" : ["fiend@gmail.com"], "cc": null, "bcc": [],  "body" : "body"}`,
			status: http.StatusOK,
		},
	}

	for _, tst := range tt {
		t.Run(tst.name, func(t *testing.T) {
			api := NewHTTPHandler(&mockAPI{tst.queuef, tst.getf})
			w := httptest.NewRecorder()
			body := strings.NewReader(tst.body)
			r := httptest.NewRequest(http.MethodPost, "http://localhost", body)
			api.sendEmailHandler(w, r)
			if w.Code != tst.status {
				t.Fatalf("sendEmailHandler(): got %d when was expecting %d", w.Code, tst.status)
			}
		})
	}
}

func TestGetEmailHandler(t *testing.T) {
	flag.Parse()
	at, _ := time.Parse(time.RFC3339, "2018-12-03T19:32:55.738296751Z")
	email, _ := New(10, "from@example.com", []string{"to@example.com"}, []string{"cc@example.com"}, []string{"bcc@example.com"}, "subject", "body")

	email.AddStatusEvent(MakeStatusEvent(Queued, at))
	var (
		passQueue   = func(Email) (Email, error) { return Email{}, nil }
		failGet     = func(int) (Email, error) { return Email{}, errors.New("get failed") }
		passGet     = func(int) (Email, error) { return email, nil }
		notFoundGet = func(int) (Email, error) { return Email{}, ErrItemNotFound }
	)
	tt := []struct {
		name   string
		id     string
		getf   func(int) (Email, error)
		want   string
		status int
	}{
		{
			name:   "should return 500 if the api returns a generic error",
			id:     "10",
			getf:   failGet,
			status: http.StatusInternalServerError,
		},
		{
			name:   "should return 404 if the id is invalid integer",
			id:     "string",
			getf:   passGet,
			status: http.StatusNotFound,
		},
		{
			name:   "should return 404 if the id does not exist",
			id:     "10",
			getf:   notFoundGet,
			status: http.StatusNotFound,
		},
		{
			name:   "should return 200 if the id exist",
			id:     "10",
			getf:   passGet,
			want:   "http_get_200.golden",
			status: http.StatusOK,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			h := NewHTTPHandler(&mockAPI{passQueue, test.getf})
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "http://localhost/emails/"+test.id, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", test.id)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			h.getEmailHandler(w, r)
			if w.Code != test.status {
				t.Fatalf("%s: got status %d, but expected %d", test.name, w.Code, test.status)
			}
			if w.Code != http.StatusOK {
				return
			}

			// test the body
			filepath := path.Join(testFolder, test.want)
			got, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Fatalf("%s: failed to read body: %v", test.name, err)
			}

			if *update {
				if err := ioutil.WriteFile(filepath, got, 0644); err != nil {
					t.Fatalf("%s: failed to write file: %v", test.name, err)
				}
			}

			expected, err := ioutil.ReadFile(filepath)
			if err != nil {
				t.Fatalf("%s: failed to read file %s: %v", test.name, filepath, err)
			}
			if !reflect.DeepEqual(got, expected) {
				t.Fatalf("%s: got %s, expected %s", test.name, got, expected)
			}

		})
	}
}
