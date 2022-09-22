package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	r := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	health(w, r)

	want := http.StatusOK
	got := w.Code
	assert.Equal(t, want, got)
}
