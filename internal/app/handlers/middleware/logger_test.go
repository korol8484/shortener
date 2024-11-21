package middleware

import (
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
)

func TestWriter(t *testing.T) {
	w := &writer{
		ResponseWriter: httptest.NewRecorder(),
		code:           0,
		bytes:          0,
		wroteHeader:    false,
	}

	_, err := w.Write([]byte("1"))
	require.NoError(t, err)

	w.WriteHeader(1)
}
