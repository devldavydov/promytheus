// Package middleware is a package for middleware server functionaluty.
package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"

	_http "github.com/devldavydov/promytheus/internal/common/http"
)

type gzipReader struct {
	io.ReadCloser

	Reader io.Reader
}

func (gr gzipReader) Read(p []byte) (n int, err error) {
	return gr.Reader.Read(p)
}

type gzipWriter struct {
	http.ResponseWriter

	gzWriter   io.Writer
	statusCode int
}

func newGzipWriter(rw http.ResponseWriter, w io.Writer) *gzipWriter {
	return &gzipWriter{
		ResponseWriter: rw,
		gzWriter:       w,
		statusCode:     http.StatusOK,
	}
}

func (gw *gzipWriter) Write(b []byte) (int, error) {
	if isContentTypeGzipSupported(gw.Header().Get("Content-Type")) {
		return gw.gzWriter.Write(b)
	}
	return gw.ResponseWriter.Write(b)
}

func (gw *gzipWriter) WriteHeader(statusCode int) {
	gw.ResponseWriter.WriteHeader(statusCode)
}

var supportedContentTypes = []string{
	_http.BaseContentTypeApplicationJS,
	_http.BaseContentTypeApplicationJSON,
	_http.BaseContentTypeCSS,
	_http.BaseContentTypeHTML,
	_http.BaseContentTextPlain,
	_http.BaseContentTypeXML,
}

var gzPool = sync.Pool{
	New: func() interface{} {
		w := gzip.NewWriter(io.Discard)
		gzip.NewWriterLevel(w, gzip.BestSpeed)
		return w
	},
}

// Gzip is a compression middleware.
func Gzip(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if shouldGzipDecodeRequest(r.Header) {
			gzr, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer gzr.Close()
			r.Body = gzipReader{ReadCloser: r.Body, Reader: gzr}
		}

		if shouldGzipEncodeResponse(r.Header) {
			gz := gzPool.Get().(*gzip.Writer)
			defer gzPool.Put(gz)

			gz.Reset(w)
			defer gz.Close()

			w.Header().Set("Content-Encoding", "gzip")
			next.ServeHTTP(newGzipWriter(w, gz), r)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func shouldGzipDecodeRequest(header http.Header) bool {
	return strings.Contains(header.Get("Content-Encoding"), "gzip")
}

func shouldGzipEncodeResponse(header http.Header) bool {
	return strings.Contains(header.Get("Accept-Encoding"), "gzip")
}

func isContentTypeGzipSupported(contentType string) bool {
	for _, v := range supportedContentTypes {
		if strings.Contains(contentType, v) {
			return true
		}
	}
	return false
}
