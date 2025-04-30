package server

import (
	"github.com/gorilla/mux"
	"github.com/ktigay/metrics-collector/internal/server/collector"
	"github.com/ktigay/metrics-collector/internal/server/storage"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServer_CollectHandler(t *testing.T) {
	type args struct {
		requests    []string
		contentType string
	}
	tests := []struct {
		name            string
		args            args
		wantStatus      int
		wantContentType string
	}{
		{
			name: "Positive_test",
			args: args{
				requests: []string{
					"/update/gauge/Alloc/122.1",
					"/update/gauge/Lookups/122.00",
					"/update/counter/PollCount/5",
					"/update/gauge/Alloc/222.21",
					"/update/gauge/Lookups/152.00",
					"/update/counter/PollCount/10",
				},
				contentType: "text/plain",
			},
			wantStatus:      http.StatusOK,
			wantContentType: "text/plain",
		},
		{
			name: "Not_found_test",
			args: args{
				requests: []string{
					"/update/gauge/",
				},
				contentType: "text/plain",
			},
			wantStatus:      http.StatusNotFound,
			wantContentType: "text/plain; charset=utf-8",
		},
		{
			name: "Not_found_test_#2",
			args: args{
				requests: []string{
					"/update/gauge/222.33",
				},
				contentType: "text/plain",
			},
			wantStatus:      http.StatusNotFound,
			wantContentType: "text/plain; charset=utf-8",
		},
		{
			name: "Not_found_test_#3",
			args: args{
				requests: []string{
					"/update/gauge/Alloc/222.33/111",
				},
				contentType: "text/plain",
			},
			wantStatus:      http.StatusNotFound,
			wantContentType: "text/plain; charset=utf-8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewServer(
				collector.NewMetricCollector(
					storage.NewMemStorage(),
				),
				zap.NewNop(),
			)

			router := mux.NewRouter()
			router.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", tt.args.contentType)
					next.ServeHTTP(w, r)
				})
			})
			router.HandleFunc("/update/{type}/{name}/{value}", h.CollectHandler)

			svr := httptest.NewServer(router)
			defer svr.Close()

			for _, req := range tt.args.requests {
				resp, _ := http.Post(svr.URL+req, "text/plain", strings.NewReader(""))
				require.Equal(t, tt.wantStatus, resp.StatusCode)
				require.Equal(t, tt.wantContentType, resp.Header.Get("Content-Type"))

				_ = resp.Body.Close()
			}
		})
	}
}
