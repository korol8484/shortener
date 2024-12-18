package middleware

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestCompressor(t *testing.T) {
	r := chi.NewRouter()

	compressor := NewCompressor()
	r.Use(compressor.Handler)

	r.Get("/get-json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("htmlstring"))
	})

	r.Get("/get-html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("htmlstring"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	type want struct {
		path              string
		expectedEncoding  string
		acceptedEncodings []string
		responseBody      string
	}

	tests := []struct {
		name string
		want
	}{
		{
			name: "no_expected_encodings_html",
			want: want{
				path:              "/get-html",
				acceptedEncodings: nil,
				expectedEncoding:  "",
				responseBody:      "htmlstring",
			},
		},
		{
			name: "gzip_only_encoding_html",
			want: want{
				path:              "/get-html",
				acceptedEncodings: []string{"gzip"},
				expectedEncoding:  "gzip",
				responseBody:      "htmlstring",
			},
		},
		{
			name: "no_expected_encodings_json",
			want: want{
				path:              "/get-json",
				acceptedEncodings: nil,
				expectedEncoding:  "",
				responseBody:      "htmlstring",
			},
		},
		{
			name: "gzip_only_encoding_json",
			want: want{
				path:              "/get-json",
				acceptedEncodings: []string{"gzip", "deflate"},
				expectedEncoding:  "gzip",
				responseBody:      "htmlstring",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			resp, respString := testRequestWithAcceptedEncodings(t, ts, "GET", tc.path, tc.acceptedEncodings...)
			defer resp.Body.Close()

			assert.Equal(t, tc.responseBody, respString)
			assert.Equal(t, tc.expectedEncoding, resp.Header.Get("Content-Encoding"))
		})
	}
}

func testRequestWithAcceptedEncodings(t *testing.T, ts *httptest.Server, method, path string, encodings ...string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	if len(encodings) > 0 {
		encodingsString := strings.Join(encodings, ",")
		req.Header.Set("Accept-Encoding", encodingsString)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody := decodeResponseBody(t, resp)

	return resp, respBody
}

func decodeResponseBody(t *testing.T, resp *http.Response) string {
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		var err error
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
	default:
		reader = resp.Body
	}

	respBody, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
		return ""
	}
	reader.Close()

	return string(respBody)
}

func TestCompressReader(t *testing.T) {
	in := bytes.NewReader([]byte{
		0x1f, 0x8b, 0x08, 0x08, 0xf7, 0x5e, 0x14, 0x4a,
		0x00, 0x03, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e,
		0x74, 0x78, 0x74, 0x00, 0x03, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	})

	rd, err := newCompressReader(io.NopCloser(in))
	if err != nil {
		t.Fatal(err)
	}

	_, err = rd.Read(make([]byte, 20))
	if err != nil && !errors.Is(err, io.EOF) {
		t.Fatal(err)
	}

	err = rd.Close()
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatal(err)
	}
}
