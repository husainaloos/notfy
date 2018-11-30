package status

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi"
)

func Test_Get(t *testing.T) {
	storage := NewInMemoryStorage()
	_, _ = storage.insert(MakeInfo(1, Queued))
	tt := []struct {
		name      string
		id        int
		expStatus int
		expBody   getModel
	}{
		{
			name:      "should return 200 when status exists",
			id:        1,
			expStatus: http.StatusOK,
			expBody: getModel{
				ID:     1,
				Status: "Queued",
			},
		},
		{
			name:      "should return 404 when status does not exists",
			id:        17,
			expStatus: http.StatusNotFound,
			expBody:   getModel{},
		},
	}

	for _, tst := range tt {
		t.Run(tst.name, func(t *testing.T) {
			api := NewAPI(storage)
			h := NewHTTPHandler(api)
			w := httptest.NewRecorder()
			url := fmt.Sprintf("http://localhost/status/%d", tst.id)
			r := httptest.NewRequest(http.MethodGet, url, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", strconv.Itoa(tst.id))
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			h.getStatusHandler(w, r)
			if w.Code != tst.expStatus {
				t.Errorf("Get(): got status %d but expected %d", w.Code, tst.expStatus)
			}
			var res getModel
			body, _ := ioutil.ReadAll(w.Body)
			json.Unmarshal(body, &res)
			if res.Status != tst.expBody.Status || res.ID != tst.expBody.ID {
				t.Errorf("Get(): got body %v, but expected %v", res, tst.expBody)
			}
		})
	}

}
