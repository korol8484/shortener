package middleware

import (
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/korol8484/shortener/internal/app/util"
)

var compressibleContentTypes = []string{
	"text/html",
	"application/json",
}

// Compressor is a middleware that compresses response
// body of a given content types to a data format based
// on Accept-Encoding request header.
type Compressor struct {
	allowedTypes map[string]struct{}
}

// NewCompressor - Factory for GZIP middleware
func NewCompressor() *Compressor {
	allowedTypes := make(map[string]struct{})

	for _, t := range compressibleContentTypes {
		allowedTypes[t] = struct{}{}
	}

	return &Compressor{allowedTypes: allowedTypes}
}

type compressWriter struct {
	http.ResponseWriter
	zw           io.Writer
	allowedTypes map[string]struct{}
	encoding     string
	compressible bool
	writeHeader  bool
}

type compressReader struct {
	r  io.ReadCloser
	zr io.ReadCloser
}

func (c *compressWriter) isCompressible() bool {
	if c.encoding == "" {
		return false
	}

	cT := util.FilterContentType(c.Header().Get("Content-Type"))
	if _, ok := c.allowedTypes[cT]; ok {
		return true
	}

	return false
}

// WriteHeader gzip header interceptor
func (c *compressWriter) WriteHeader(code int) {
	c.writeHeader = true
	defer c.ResponseWriter.WriteHeader(code)

	if c.Header().Get("Content-Encoding") != "" {
		return
	}

	if !c.isCompressible() {
		c.compressible = false
		return
	}

	c.compressible = true
	c.Header().Set("Content-Encoding", c.encoding)
	c.Header().Add("Vary", "Accept-Encoding")

	c.Header().Del("Content-Length")
}

// Write gzip writes the data to the connection as part of an HTTP reply.
func (c *compressWriter) Write(p []byte) (int, error) {
	if !c.writeHeader {
		c.WriteHeader(http.StatusOK)
	}

	return c.writer().Write(p)
}

func (c *compressWriter) writer() io.Writer {
	if c.compressible {
		return c.zw
	}

	return c.ResponseWriter
}

// Close is the interface that wraps the basic Close method.
func (c *compressWriter) Close() error {
	if c, ok := c.writer().(io.WriteCloser); ok {
		return c.Close()
	}

	return errors.New("writeCloser is unavailable")
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read gzip header interceptor
func (c *compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close is the interface that wraps the basic Close method.
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}

	return c.zr.Close()
}

// Handler returns a new middleware that will compress the response based on the
// current Compressor.
func (c *Compressor) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		baseEnc := "gzip"
		wr := &compressWriter{
			ResponseWriter: w,
			zw:             gzip.NewWriter(w),
			allowedTypes:   c.allowedTypes,
			encoding:       baseEnc,
			compressible:   false,
		}

		enc := r.Header.Get("Accept-Encoding")
		if !strings.Contains(enc, baseEnc) {
			wr.encoding = ""
		}

		defer func(wr *compressWriter) {
			_ = wr.Close()
		}(wr)

		contentEncoding := r.Header.Get("Content-Encoding")
		if strings.Contains(contentEncoding, baseEnc) {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr
			defer func(cr *compressReader) {
				_ = cr.Close()
			}(cr)
		}

		next.ServeHTTP(wr, r)
	})
}
