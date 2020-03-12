package ui

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpaHandlerNonExistent(t *testing.T) {
	// Test we can handle rest spec
	req := httptest.NewRequest(http.MethodGet, "/doesnotexist.json/", nil)
	w := httptest.NewRecorder()
	handler := spaHandler{staticPath: "/ui/build", indexPath: ""}
	handler.ServeHTTP(w, req)
	defer w.Result().Body.Close()
	assert.Equal(t, 500, w.Result().StatusCode, "expected status code to be 500 but got %d", w.Result().StatusCode)
}

func TestSpaHandlerInvalidPath(t *testing.T) {
	// Test we can handle rest spec
	req := httptest.NewRequest(http.MethodGet, "/data/doesnotexist.json", nil)
	w := httptest.NewRecorder()
	handler := spaHandler{staticPath: "/ui/build", indexPath: ""}
	handler.ServeHTTP(w, req)
	defer w.Result().Body.Close()
	assert.Equal(t, 500, w.Result().StatusCode, "expected status code to be 500 but got %d", w.Result().StatusCode)
}

func TestSpaHandlerValidPath(t *testing.T) {
	// Test we can handle rest spec
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler := spaHandler{staticPath: "/ui/build", indexPath: ""}
	handler.ServeHTTP(w, req)
	defer w.Result().Body.Close()
	assert.Equal(t, 200, w.Result().StatusCode, "expected status code to be 200 but got %d", w.Result().StatusCode)
}

func TestSpaHandlerFolderPath(t *testing.T) {
	// Test we can handle rest spec
	req := httptest.NewRequest(http.MethodGet, "/data/", nil)
	w := httptest.NewRecorder()
	handler := spaHandler{staticPath: "/ui/build", indexPath: ""}
	handler.ServeHTTP(w, req)
	defer w.Result().Body.Close()
	assert.Equal(t, 200, w.Result().StatusCode, "expected status code to be 200 but got %d", w.Result().StatusCode)
}
