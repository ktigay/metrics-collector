package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/ktigay/metrics-collector/internal/metric"
	"github.com/ktigay/metrics-collector/internal/server/handler/mocks"
	"github.com/ktigay/metrics-collector/internal/server/repository"
	"github.com/ktigay/metrics-collector/internal/server/service"
)

func TestServer_CollectHandler(t *testing.T) {
	type args struct {
		request     string
		contentType string
	}
	tests := []struct {
		collector       func(controller *gomock.Controller) CollectorInterface
		args            args
		name            string
		wantContentType string
		wantStatus      int
	}{
		{
			name: "Positive_test_gauge",
			args: args{
				request:     "/update/gauge/Alloc/122.1",
				contentType: "text/plain",
			},
			collector: func(mockCtrl *gomock.Controller) CollectorInterface {
				st := mocks.NewMockCollectorInterface(mockCtrl)
				v := 122.1
				st.EXPECT().Save(gomock.Any(), gomock.Eq(metric.Metrics{
					ID:    "Alloc",
					Type:  "gauge",
					Value: &v,
				})).Return(nil).Times(1)
				return st
			},
			wantStatus:      http.StatusOK,
			wantContentType: "text/plain",
		},
		{
			name: "Positive_test_counter",
			args: args{
				request:     "/update/counter/PollCount/12345",
				contentType: "text/plain",
			},
			collector: func(mockCtrl *gomock.Controller) CollectorInterface {
				st := mocks.NewMockCollectorInterface(mockCtrl)
				v := int64(12345)
				st.EXPECT().Save(gomock.Any(), gomock.Eq(metric.Metrics{
					ID:    "PollCount",
					Type:  "counter",
					Delta: &v,
				})).Return(nil).Times(1)
				return st
			},
			wantStatus:      http.StatusOK,
			wantContentType: "text/plain",
		},
		{
			name: "Not_found_test",
			args: args{
				request:     "/update/gauge/",
				contentType: "text/plain",
			},
			collector: func(mockCtrl *gomock.Controller) CollectorInterface {
				st := mocks.NewMockCollectorInterface(mockCtrl)
				st.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).Times(0)
				return st
			},
			wantStatus:      http.StatusNotFound,
			wantContentType: "text/plain; charset=utf-8",
		},
		{
			name: "Not_found_test_with_value",
			args: args{
				request:     "/update/gauge/222.33",
				contentType: "text/plain",
			},
			collector: func(mockCtrl *gomock.Controller) CollectorInterface {
				st := mocks.NewMockCollectorInterface(mockCtrl)
				st.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).Times(0)
				return st
			},
			wantStatus:      http.StatusNotFound,
			wantContentType: "text/plain; charset=utf-8",
		},
		{
			name: "Not_found_test_with_wrong_url",
			args: args{
				request:     "/update/gauge/Alloc/222.33/111",
				contentType: "text/plain",
			},
			collector: func(mockCtrl *gomock.Controller) CollectorInterface {
				st := mocks.NewMockCollectorInterface(mockCtrl)
				st.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).Times(0)
				return st
			},
			wantStatus:      http.StatusNotFound,
			wantContentType: "text/plain; charset=utf-8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			h := NewMetricHandler(tt.collector(mockCtrl), zap.NewNop().Sugar())

			router := mux.NewRouter()
			router.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", tt.args.contentType)
					next.ServeHTTP(w, r)
				})
			})
			router.HandleFunc("/update/{type}/{name}/{value}", h.CollectHandler)

			srv := httptest.NewServer(router)
			defer srv.Close()

			resp, _ := http.Post(srv.URL+tt.args.request, "text/plain", nil)
			require.Equal(t, tt.wantStatus, resp.StatusCode)
			require.Equal(t, tt.wantContentType, resp.Header.Get("Content-Type"))

			_ = resp.Body.Close()
		})
	}
}

func TestServer_UpdateJSONHandler(t *testing.T) {
	newCollector := func() *service.MetricCollector {
		st, _ := repository.NewMemRepository(nil, zap.NewNop().Sugar())
		return service.NewMetricCollector(st, zap.NewNop().Sugar())
	}

	type fields struct {
		collector CollectorInterface
	}
	type args struct {
		contentType string
		request     []byte
	}
	type want struct {
		contentType string
		response    string
		statusCode  int
		wantErr     bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "Positive_test_gauge",
			fields: fields{
				collector: newCollector(),
			},
			args: args{
				request:     []byte(`{"id":"TestSet90","type":"gauge","delta":0,"value":10}`),
				contentType: "application/json",
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
				response:    "{\"value\":10,\"id\":\"TestSet90\",\"type\":\"gauge\"}\n",
				wantErr:     false,
			},
		},
		{
			name: "Positive_test_counter",
			fields: fields{
				collector: service.NewMetricCollector(
					&repository.MemMetricRepository{
						Metrics: map[string]repository.MetricEntity{
							"counter:TestSet91": {
								Key:   "counter:TestSet91",
								Name:  "TestSet91",
								Type:  "counter",
								Delta: int64(10),
							},
						},
					},
					zap.NewNop().Sugar(),
				),
			},
			args: args{
				request:     []byte(`{"id":"TestSet91","type":"counter","delta":15,"value":0}`),
				contentType: "application/json",
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
				response:    "{\"delta\":25,\"id\":\"TestSet91\",\"type\":\"counter\"}\n",
				wantErr:     false,
			},
		},
		{
			name: "Bad_Request_Wrong_ContentType",
			fields: fields{
				collector: newCollector(),
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
				collector: newCollector(),
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
				collector: newCollector(),
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
			h := NewMetricHandler(
				tt.fields.collector,
				zap.NewNop().Sugar(),
			)

			srv := httptest.NewServer(http.HandlerFunc(h.UpdateJSONHandler))
			defer srv.Close()

			resp, err := http.Post(srv.URL+"/update/", tt.args.contentType, bytes.NewReader(tt.args.request))
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
	newCollector := func() *service.MetricCollector {
		logger := zap.NewNop().Sugar()
		st, _ := repository.NewMemRepository(nil, logger)
		return service.NewMetricCollector(st, logger)
	}

	type fields struct {
		collector CollectorInterface
	}
	type args struct {
		contentType string
		request     []byte
	}
	type want struct {
		contentType string
		response    string
		statusCode  int
		wantErr     bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "Positive_test_gauge",
			fields: fields{
				collector: service.NewMetricCollector(
					&repository.MemMetricRepository{
						Metrics: map[string]repository.MetricEntity{
							"gauge:TestSet90": {
								Key:   "counter:TestSet90",
								Name:  "TestSet90",
								Type:  "gauge",
								Value: 15.444,
							},
						},
					},
					zap.NewNop().Sugar(),
				),
			},
			args: args{
				request:     []byte(`{"id":"TestSet90","type":"gauge","delta":0,"value":10}`),
				contentType: "application/json",
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
				response:    "{\"value\":15.444,\"id\":\"TestSet90\",\"type\":\"gauge\"}\n",
				wantErr:     false,
			},
		},
		{
			name: "Positive_test_counter",
			fields: fields{
				collector: service.NewMetricCollector(
					&repository.MemMetricRepository{
						Metrics: map[string]repository.MetricEntity{
							"counter:TestSet91": {
								Key:   "counter:TestSet91",
								Name:  "TestSet91",
								Type:  "counter",
								Delta: int64(10),
							},
						},
					},
					zap.NewNop().Sugar(),
				),
			},
			args: args{
				request:     []byte(`{"id":"TestSet91","type":"counter"}`),
				contentType: "application/json",
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
				response:    "{\"delta\":10,\"id\":\"TestSet91\",\"type\":\"counter\"}\n",
				wantErr:     false,
			},
		},
		{
			name: "Bad_Request_Not_Found",
			fields: fields{
				collector: newCollector(),
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
				collector: newCollector(),
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
				collector: newCollector(),
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
				collector: newCollector(),
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
			h := NewMetricHandler(
				tt.fields.collector,
				zap.NewNop().Sugar(),
			)

			srv := httptest.NewServer(http.HandlerFunc(h.GetJSONValueHandler))
			defer srv.Close()

			resp, err := http.Post(srv.URL+"/value/", tt.args.contentType, bytes.NewReader(tt.args.request))
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

func TestMetricHandler_UpdatesJSONHandler(t *testing.T) {
	newCollector := func() *service.MetricCollector {
		logger := zap.NewNop().Sugar()
		st, _ := repository.NewMemRepository(nil, logger)
		return service.NewMetricCollector(st, logger)
	}
	type fields struct {
		collector CollectorInterface
	}
	type args struct {
		contentType string
		request     []byte
	}
	type want struct {
		contentType string
		statusCode  int
		wantErr     bool
	}
	tests := []struct {
		args   args
		name   string
		fields fields
		want   want
	}{
		{
			name: "Positive_test_gauge",
			fields: fields{
				collector: newCollector(),
			},
			args: args{
				request:     []byte(`[{"id":"TestSet90","type":"gauge","delta":0,"value":10}]`),
				contentType: "application/json",
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
				wantErr:     false,
			},
		},
		{
			name: "Positive_test_counter",
			fields: fields{
				collector: newCollector(),
			},
			args: args{
				request:     []byte(`[{"id":"TestSet91","type":"counter","delta":15,"value":0}]`),
				contentType: "application/json",
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
				wantErr:     false,
			},
		},
		{
			name: "Bad_Request_Wrong_ContentType",
			fields: fields{
				collector: newCollector(),
			},
			args: args{
				request:     []byte(`[{"id":"TestSet92","type":"gauge","delta":0,"value":10}]`),
				contentType: "text/plain; charset=utf-8",
			},
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "",
				wantErr:     false,
			},
		},
		{
			name: "Bad_Request_Broken_Body",
			fields: fields{
				collector: newCollector(),
			},
			args: args{
				request:     []byte(`[{"id":"TestSet93","type":"gauge","delta":0,"value":10]`),
				contentType: "application/json",
			},
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "application/json",
				wantErr:     false,
			},
		},
		{
			name: "Bad_Request_Wrong_Body_Type",
			fields: fields{
				collector: newCollector(),
			},
			args: args{
				request:     []byte(`{"id":"TestSet92","type":"gauge","delta":0,"value":10}`),
				contentType: "application/json",
			},
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "application/json",
				wantErr:     false,
			},
		},
		{
			name: "Bad_Request_Wrong_Type",
			fields: fields{
				collector: newCollector(),
			},
			args: args{
				request:     []byte(`[{"id":"TestSet94","type":"wrongType","delta":0,"value":10}]`),
				contentType: "application/json",
			},
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewMetricHandler(
				tt.fields.collector,
				zap.NewNop().Sugar(),
			)
			srv := httptest.NewServer(http.HandlerFunc(h.UpdatesJSONHandler))
			defer srv.Close()

			resp, err := http.Post(srv.URL+"/updates/", tt.args.contentType, bytes.NewReader(tt.args.request))
			if (err != nil) != tt.want.wantErr {
				t.Errorf("UpdatesJSONHandler() error = %v, wantErr %v", err, tt.want.wantErr)
				return
			}
			defer func() {
				if resp != nil && resp.Body != nil {
					_ = resp.Body.Close()
				}
			}()

			require.Equal(t, tt.want.statusCode, resp.StatusCode)
			require.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}
