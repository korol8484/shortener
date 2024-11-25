package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilterContentType(t *testing.T) {
	cont := FilterContentType("text/html; charset=utf-8")
	assert.Equal(t, cont, "text/html")

	cont = FilterContentType("multipart/form-data; boundary=ExampleBoundaryString")
	assert.Equal(t, cont, "multipart/form-data")

	cont = FilterContentType("multipart/form-data")
	assert.Equal(t, cont, "multipart/form-data")
}
