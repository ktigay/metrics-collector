package server

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/ktigay/metrics-collector/internal/server/collector"
	"github.com/ktigay/metrics-collector/internal/server/storage"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
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
					storage.NewMemStorage(nil),
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

func TestServer_UpdateJSONHandler(t *testing.T) {
	type fields struct {
		collector CollectorInterface
		logger    *zap.Logger
	}
	type args struct {
		request     []byte
		contentType string
	}
	type want struct {
		statusCode  int
		contentType string
		response    string
		wantErr     bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "Positive_test",
			fields: fields{
				collector: collector.NewMetricCollector(
					storage.NewMemStorage(nil),
				),
				logger: zap.NewNop(),
			},
			args: args{
				request:     []byte(`{"id":"TestSet90","type":"gauge","delta":0,"value":10}`),
				contentType: "application/json",
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
				response:    "{\"id\":\"TestSet90\",\"type\":\"gauge\",\"value\":10}\n",
				wantErr:     false,
			},
		},
		{
			name: "Positive_test_#2",
			fields: fields{
				collector: collector.NewMetricCollector(
					&storage.MemStorage{
						Metrics: map[string]storage.Entity{
							"counter:TestSet91": {
								Key:   "counter:TestSet91",
								Name:  "TestSet91",
								Type:  "counter",
								Delta: int64(10),
							},
						},
					},
				),
				logger: zap.NewNop(),
			},
			args: args{
				request:     []byte(`{"id":"TestSet91","type":"counter","delta":15,"value":0}`),
				contentType: "application/json",
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
				response:    "{\"id\":\"TestSet91\",\"type\":\"counter\",\"delta\":25}\n",
				wantErr:     false,
			},
		},
		{
			name: "Bad_Request_Wrong_ContentType",
			fields: fields{
				collector: collector.NewMetricCollector(
					storage.NewMemStorage(nil),
				),
				logger: zap.NewNop(),
			},
			args: args{
				request:     []byte(`{"id":"TestSet92","type":"gauge","delta":0,"value":10}`),
				contentType: "text/plain; charset=utf-8",
			},
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "",
				response:    "",
				wantErr:     false,
			},
		},
		{
			name: "Bad_Request_Broken_Body",
			fields: fields{
				collector: collector.NewMetricCollector(
					storage.NewMemStorage(nil),
				),
				logger: zap.NewNop(),
			},
			args: args{
				request:     []byte(`{"id":"TestSet93","type":"gauge","delta":0,"value":10`),
				contentType: "application/json",
			},
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "application/json",
				response:    "",
				wantErr:     false,
			},
		},
		{
			name: "Bad_Request_Wrong_Type",
			fields: fields{
				collector: collector.NewMetricCollector(
					storage.NewMemStorage(nil),
				),
				logger: zap.NewNop(),
			},
			args: args{
				request:     []byte(`{"id":"TestSet94","type":"wrongType","delta":0,"value":10}`),
				contentType: "application/json",
			},
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "application/json",
				response:    "",
				wantErr:     false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewServer(
				tt.fields.collector,
				tt.fields.logger,
			)

			router := mux.NewRouter()
			router.HandleFunc("/update/", h.UpdateJSONHandler)

			svr := httptest.NewServer(router)
			defer svr.Close()

			resp, err := http.Post(svr.URL+"/update/", tt.args.contentType, bytes.NewReader(tt.args.request))
			if (err != nil) != tt.want.wantErr {
				t.Errorf("UpdateJSONHandler() error = %v, wantErr %v", err, tt.want.wantErr)
				return
			}
			defer func() {
				if resp != nil && resp.Body != nil {
					_ = resp.Body.Close()
				}
			}()

			require.Equal(t, tt.want.statusCode, resp.StatusCode)
			require.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))

			b, _ := io.ReadAll(resp.Body)
			require.Equal(t, tt.want.response, string(b))
		})
	}
}

func TestServer_GetJSONValueHandler(t *testing.T) {
	type fields struct {
		collector CollectorInterface
		logger    *zap.Logger
	}
	type args struct {
		request     []byte
		contentType string
	}
	type want struct {
		statusCode  int
		contentType string
		response    string
		wantErr     bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "Positive_test",
			fields: fields{
				collector: collector.NewMetricCollector(
					&storage.MemStorage{
						Metrics: map[string]storage.Entity{
							"gauge:TestSet90": {
								Key:   "counter:TestSet90",
								Name:  "TestSet90",
								Type:  "gauge",
								Value: 15.444,
							},
						},
					},
				),
				logger: zap.NewNop(),
			},
			args: args{
				request:     []byte(`{"id":"TestSet90","type":"gauge","delta":0,"value":10}`),
				contentType: "application/json",
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
				response:    "{\"id\":\"TestSet90\",\"type\":\"gauge\",\"value\":15.444}\n",
				wantErr:     false,
			},
		},
		{
			name: "Positive_test_#2",
			fields: fields{
				collector: collector.NewMetricCollector(
					&storage.MemStorage{
						Metrics: map[string]storage.Entity{
							"counter:TestSet91": {
								Key:   "counter:TestSet91",
								Name:  "TestSet91",
								Type:  "counter",
								Delta: int64(10),
							},
						},
					},
				),
				logger: zap.NewNop(),
			},
			args: args{
				request:     []byte(`{"id":"TestSet91","type":"counter"}`),
				contentType: "application/json",
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
				response:    "{\"id\":\"TestSet91\",\"type\":\"counter\",\"delta\":10}\n",
				wantErr:     false,
			},
		},
		{
			name: "Bad_Request_Not_Found",
			fields: fields{
				collector: collector.NewMetricCollector(
					storage.NewMemStorage(nil),
				),
				logger: zap.NewNop(),
			},
			args: args{
				request:     []byte(`{"id":"TestSet92","type":"gauge"}`),
				contentType: "application/json",
			},
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "application/json",
				response:    "",
				wantErr:     false,
			},
		},
		{
			name: "Bad_Request_Broken_Body",
			fields: fields{
				collector: collector.NewMetricCollector(
					storage.NewMemStorage(nil),
				),
				logger: zap.NewNop(),
			},
			args: args{
				request:     []byte(`{"id":"TestSet93","type":"gauge"`),
				contentType: "application/json",
			},
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "application/json",
				response:    "",
				wantErr:     false,
			},
		},
		{
			name: "Bad_Request_Wrong_Type",
			fields: fields{
				collector: collector.NewMetricCollector(
					storage.NewMemStorage(nil),
				),
				logger: zap.NewNop(),
			},
			args: args{
				request:     []byte(`{"id":"TestSet94","type":"wrongType"}`),
				contentType: "application/json",
			},
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "application/json",
				response:    "",
				wantErr:     false,
			},
		},
		{
			name: "Bad_Request_Wrong_ContentType",
			fields: fields{
				collector: collector.NewMetricCollector(
					storage.NewMemStorage(nil),
				),
				logger: zap.NewNop(),
			},
			args: args{
				request:     []byte(`{"id":"TestSet95","type":"gauge"}`),
				contentType: "text/plain; charset=utf-8",
			},
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "",
				response:    "",
				wantErr:     false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewServer(
				tt.fields.collector,
				tt.fields.logger,
			)

			router := mux.NewRouter()
			router.HandleFunc("/value/", h.GetJSONValueHandler)

			svr := httptest.NewServer(router)
			defer svr.Close()

			resp, err := http.Post(svr.URL+"/value/", tt.args.contentType, bytes.NewReader(tt.args.request))
			if (err != nil) != tt.want.wantErr {
				t.Errorf("GetJSONValueHandler() error = %v, wantErr %v", err, tt.want.wantErr)
				return
			}
			defer func() {
				if resp != nil && resp.Body != nil {
					_ = resp.Body.Close()
				}
			}()

			require.Equal(t, tt.want.statusCode, resp.StatusCode)
			require.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))

			b, _ := io.ReadAll(resp.Body)
			require.Equal(t, tt.want.response, string(b))
		})
	}
}
