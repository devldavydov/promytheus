package metric

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rsa"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/devldavydov/promytheus/internal/common/cipher"
	_http "github.com/devldavydov/promytheus/internal/common/http"
	"github.com/devldavydov/promytheus/internal/server/middleware"
	"github.com/devldavydov/promytheus/internal/server/mocks"
	"github.com/devldavydov/promytheus/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

type testResponse struct {
	headers     map[string][]string
	body        string
	contentType string
	code        int
}

type testRequest struct {
	body        io.Reader
	contentType *string
	hmacKey     *string
	encryption  bool
	headers     map[string][]string
	method      string
	url         string
}

type testItem struct {
	stgInitFunc  func(storage.Storage)
	stgCheckFunc func() []storage.StorageItem
	stgMockFunc  func(*mocks.MockStorage)
	req          testRequest
	name         string
	resp         testResponse
	xfail        bool
	dbStg        bool
}

var (
	privKey *rsa.PrivateKey
	pubKey  *rsa.PublicKey
)

func TestMain(m *testing.M) {
	privKey, pubKey, _ = cipher.GenerateKeyPair(2048)
	os.Exit(m.Run())
}

func runTests(t *testing.T, tests []testItem) {
	for _, tt := range tests {
		tt := tt
		if tt.xfail {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			logger := logrus.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var stg storage.Storage

			if tt.dbStg {
				pgStg := mocks.NewMockStorage(ctrl)
				if tt.stgMockFunc != nil {
					tt.stgMockFunc(pgStg)
				}
				stg = pgStg
			} else {
				memStg, _ := storage.NewMemStorage(context.TODO(), logger, storage.NewPersistSettings(0, "", false))
				if tt.stgInitFunc != nil {
					tt.stgInitFunc(memStg)
				}
				stg = memStg
			}

			router := chi.NewRouter()

			var cryptoPrivKey *rsa.PrivateKey
			if tt.req.encryption {
				cryptoPrivKey = privKey
			}
			mdlwrDecr := middleware.NewDecrpyt(cryptoPrivKey)

			router.Use(middleware.Gzip, mdlwrDecr.Handle)

			NewHandler(router, stg, tt.req.hmacKey, logger)
			ts := httptest.NewServer(router)
			defer ts.Close()

			statusCode, contentType, body, headers := doTestRequest(t, ts, tt.req)

			assert.Equal(t, tt.resp.code, statusCode)
			assert.Equal(t, tt.resp.contentType, contentType)

			if contentType == _http.ContentTypeApplicationJSON {
				assert.JSONEq(t, tt.resp.body, body)
			} else {
				assert.Equal(t, tt.resp.body, body)
			}

			for expName, expHeaders := range tt.resp.headers {
				for _, expHeader := range expHeaders {
					assert.True(t, slices.Contains(headers[expName], expHeader))
				}
			}

			if tt.stgCheckFunc != nil {
				expMetrics, _ := stg.GetAllMetrics()
				assert.Equal(t, expMetrics, tt.stgCheckFunc())
			}
		})
	}
}

func doTestRequest(t *testing.T, ts *httptest.Server, testReq testRequest) (int, string, string, map[string][]string) {
	req, err := http.NewRequest(testReq.method, ts.URL+testReq.url, testReq.body)
	require.NoError(t, err)

	if testReq.contentType != nil {
		req.Header.Set("Content-Type", *testReq.contentType)
	}

	for name, headers := range testReq.headers {
		for _, header := range headers {
			req.Header.Set(name, header)
		}
	}

	client := &http.Client{
		Transport: &http.Transport{DisableCompression: true},
	}
	resp, err := client.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	var respBody []byte
	var bodyReader io.Reader

	if !strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		bodyReader = resp.Body
	} else {
		var gzReader *gzip.Reader
		gzReader, err = gzip.NewReader(resp.Body)
		require.NoError(t, err)
		bodyReader = gzReader
	}
	respBody, err = io.ReadAll(bodyReader)
	require.NoError(t, err)

	return resp.StatusCode, resp.Header.Get("Content-Type"), string(respBody), resp.Header
}

func bodyStringReader(val string) io.Reader {
	return bytes.NewBuffer([]byte(val))
}

func bodyGzipReader(val string) io.Reader {
	var buf bytes.Buffer
	zw, _ := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	zw.Write([]byte(val))
	zw.Close()
	return &buf
}

func encryptString(val string) string {
	encBuf := cipher.NewEncBuffer(pubKey)
	encBuf.Write([]byte(val))
	encData, _ := io.ReadAll(encBuf)
	return string(encData)
}

func strPointer(s string) *string { return &s }
