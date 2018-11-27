package email

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/husainaloos/notfy/dto"
	"github.com/husainaloos/notfy/messaging"
	"github.com/husainaloos/notfy/status"
)

func Test_sendEmailHandler(t *testing.T) {
	t.Run("should return bad request if the body is invalid", func(t *testing.T) {
		statusAPI := status.NewAPI(status.NewInMemoryStorage())
		api := NewAPIHandler(&messaging.ErrPublisher{}, statusAPI)
		w := httptest.NewRecorder()
		body := strings.NewReader(`{"from" : "bademail.com", "to" : ["friend@gmail.com"]`)
		r := httptest.NewRequest(http.MethodPost, "http://localhost", body)
		api.sendEmailHandler(w, r)
		if w.Code != http.StatusBadRequest {
			t.Errorf("sendEmailHandler(): expected %d, but got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("should return bad request if the body is invalid json", func(t *testing.T) {
		statusAPI := status.NewAPI(status.NewInMemoryStorage())
		api := NewAPIHandler(&messaging.ErrPublisher{}, statusAPI)
		w := httptest.NewRecorder()
		body := strings.NewReader(`{"field" : "bad json"`)
		r := httptest.NewRequest(http.MethodPost, "http://localhost", body)
		api.sendEmailHandler(w, r)
		if w.Code != http.StatusBadRequest {
			t.Errorf("sendEmailHandler(): expected %d, but got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("should return internal server error if the publisher fails", func(t *testing.T) {
		statusAPI := status.NewAPI(status.NewInMemoryStorage())
		api := NewAPIHandler(&messaging.ErrPublisher{Err: errors.New("[testing] ErrPublisher is configured to return error")}, statusAPI)
		w := httptest.NewRecorder()
		body := strings.NewReader(`{"from" : "email@gmail.com", "to" : ["fiend@gmail.com"]}`)
		r := httptest.NewRequest(http.MethodPost, "http://localhost", body)
		api.sendEmailHandler(w, r)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("sendEmailHandler(): expected %d, but got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("should publish message if valid", func(t *testing.T) {
		statusAPI := status.NewAPI(status.NewInMemoryStorage())
		pf := messaging.NewInMemoryPublisher()
		api := NewAPIHandler(pf, statusAPI)
		w := httptest.NewRecorder()
		body := strings.NewReader(`{"from" : "email@gmail.com", "to" : ["fiend@gmail.com"], "cc": null, "bcc": [],  "body" : "body"}`)
		r := httptest.NewRequest(http.MethodPost, "http://localhost", body)
		api.sendEmailHandler(w, r)
		if w.Code != http.StatusOK {
			t.Errorf("sendEmailHandler(): expected %d, but got %d", http.StatusOK, w.Code)
		}
		msg := <-pf.C
		var model dto.PublishedEmail
		if err := proto.Unmarshal(msg, &model); err != nil {
			t.Errorf("sendEmailHandler(): cannot unmarshal binary: %v", err)
		}
		if model.From != "<email@gmail.com>" {
			t.Errorf("model.From: expected %s, but got %s", "email@gmail.com", model.GetFrom())
		}
		if len(model.GetTo()) != 1 {
			t.Errorf("model.GetTo(): expected %d elements, but got %d elements", 1, len(model.GetTo()))
		}
		if len(model.GetCc()) != 0 {
			t.Errorf("model.GetCc(): expected %d elements, but got %d elements", 0, len(model.GetCc()))
		}
		if len(model.GetBcc()) != 0 {
			t.Errorf("model.GetBcc(): expected %d elements, but got %d elements", 0, len(model.GetBcc()))
		}
		if model.GetSubject() != "" {
			t.Errorf("model.GetSubject(): expected %s, but got %s", "", model.GetSubject())
		}
		if model.GetBody() != "body" {
			t.Errorf("model.GetBody(): expected %s, but got %s", "body", model.GetBody())
		}

	})
}

func BenchmarkPostEmail(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		statusAPI := status.NewAPI(status.NewInMemoryStorage())
		api := NewAPIHandler(messaging.NilPublisher{}, statusAPI)
		w := httptest.NewRecorder()
		body := strings.NewReader(`
{
	"from" : "FROM <from@gmail.com>",
	"to" : ["TO <to@gmail.com>","TO <to@gmail.com>","TO <to@gmail.com>","TO <to@gmail.com>","TO <to@gmail.com>","TO <to@gmail.com>","TO <to@gmail.com>","TO <to@gmail.com>","TO <to@gmail.com>","TO <to@gmail.com>"],
	"cc" : ["CC <cc@gmail.com>","CC <cc@gmail.com>","CC <cc@gmail.com>","CC <cc@gmail.com>","CC <cc@gmail.com>","CC <cc@gmail.com>","CC <cc@gmail.com>","CC <cc@gmail.com>","CC <cc@gmail.com>","CC <cc@gmail.com>"],
	"bcc" : ["BBC <bcc@gmail.com>", "BBC <bcc@gmail.com>","BBC <bcc@gmail.com>","BBC <bcc@gmail.com>","BBC <bcc@gmail.com>","BBC <bcc@gmail.com>","BBC <bcc@gmail.com>","BBC <bcc@gmail.com>","BBC <bcc@gmail.com>","BBC <bcc@gmail.com>"],
	"subject" : "Neque porro quisquam est qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit",
	"body" : "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum. Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum. Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum. Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."
}`)
		r := httptest.NewRequest(http.MethodPost, "http://localhost", body)
		b.StartTimer()
		api.sendEmailHandler(w, r)
		b.StopTimer()
		if w.Code != http.StatusOK {
			b.Fatal("not OK")
		}
		b.StartTimer()
	}
}
