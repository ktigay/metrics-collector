package middleware

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/ktigay/metrics-collector/internal/crypto"
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

func TestDecryptRequestHandler(t *testing.T) {
	privateKey := `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`

	publicKey := `-----BEGIN RSA PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAqpp6tG2TErmUxfX87VYU
HuCPxIy1dJbqL4UTr2UIs8vw0k+V7PTsy+/niUMWl5bY7kWL0nwxndZEhmkJLQmp
lgy9ZlmSWCk42hsYJVx3XO0IVr7748zVV1PChpZyBPYhG2sdQhQ+8oPB4Pj5Ea1N
XvIbc7IAumiXqCz/mvo50i7DokL5iE7my9MHVd2jSPSio8lOfL2hdDSMVFaH7opx
CNwecuS77a5D4Xr3oLUEmgpBqgpndzWEGGQDlhsUeonHBMTrcnBX6X75guECeu+R
X9ix+5v930Df0XI39WUmV8qcnPdBtMGxdiNN1Pf1ru3ovbccPu/Bn2Nt3Ox3IWHT
kdCd3iNxNHz/IZ96kjk4v9YYxECwiG3TG8+By8oz7Zg1yv6YV8TMGuU4aEPahwZ0
LesbCQkleQO1SfxDn+rfkYsKSH4DEdb8PN/zPJrsFFvciFljdtFJbQC48IHfGR+c
PYm/lfA/ICnWHuh3E5JVLFAtlC4IvfP0+sQxRNuY7QKzHTWwcfonMYKdikEdjZtK
AsIW1g3cJu8qmaJZQR7E8JSJtRSWIdIw+O0o/8reBG46ExqsCtvAYUuBnURy5iqb
CJ0g/skKBnZoY+Gyzcdj2tDB+lt9M3VcFfuSf9w2EmOZbsnY5XxUkNQGK7LPa+kn
oTdBbK732TjFe2lO4y+wJ1cCAwEAAQ==
-----END RSA PUBLIC KEY-----
`
	type args struct {
		k    *crypto.PrivateKey
		body []byte
	}

	encrypt := func(str string) []byte {
		k, err := crypto.NewPublicKey(strings.NewReader(publicKey))
		if err != nil {
			panic(err)
		}
		b, err := rsa.EncryptPKCS1v15(rand.Reader, k.Key, []byte(str))
		if err != nil {
			panic(err)
		}
		return b
	}

	tests := []struct {
		name string
		want string
		args args
	}{
		{
			name: "Positive_test_decrypt",
			args: args{
				k: func() *crypto.PrivateKey {
					k, err := crypto.NewPrivateKey(strings.NewReader(privateKey))
					if err != nil {
						panic(err)
					}
					return k
				}(),
				body: encrypt("test test"),
			},
			want: "test test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			router := mux.NewRouter()
			router.Use(DecryptRequestHandler(zap.NewNop().Sugar(), tt.args.k.Key))
			router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
				b, err := io.ReadAll(request.Body)
				if err != nil {
					panic(err)
				}
				assert.Equal(t, tt.want, string(b))
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

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				_ = resp.Body.Close()
			}()
		})
	}
}
