package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	h "github.com/ktigay/metrics-collector/internal/http"
)

func TestCheckSumRequestHandler(t *testing.T) {
	type args struct {
		hashKey  string
		checksum string
		body     []byte
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "Positive_test_checksum",
			args: args{
				hashKey:  "sha256",
				checksum: "20a8636b989d82f2f6b0bc108f5ccd1b22d44f9f8f3281e0b1ad070aedf24aba",
				body:     []byte("hello world"),
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Positive_test_no_checksum",
			args: args{
				hashKey:  "sha2563dd322",
				checksum: "",
				body:     []byte("hello world"),
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Negative_test_wrong_checksum",
			args: args{
				hashKey:  "sha2563322",
				checksum: "20a8636b989d82f2f6b0bc108f5ccd1b22d44f9f8f3281e0b1ad070aedf24aba",
				body:     []byte("hello world"),
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Negative_test_invalid_byte_error",
			args: args{
				hashKey:  "sha256",
				checksum: "r0a8636b222d8doi76b0bc108f5ccd1b22d44f9glhu771e0b1ad070aedf24dda",
				body:     []byte("hello world"),
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			router := mux.NewRouter()
			router.Use(CheckSumRequestHandler(zap.NewNop().Sugar(), tt.args.hashKey))
			router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
			})

			srv := httptest.NewServer(router)
			defer srv.Close()

			req, err := http.NewRequest(http.MethodPost, srv.URL+"/", bytes.NewReader(tt.args.body))
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				_ = req.Body.Close()
			}()

			if tt.args.checksum != "" {
				req.Header[h.HashSHA256Header] = []string{tt.args.checksum}
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				_ = resp.Body.Close()
			}()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}
