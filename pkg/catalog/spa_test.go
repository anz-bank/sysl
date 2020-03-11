package catalog

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpaHandler(t *testing.T) {
	// Test we can handle rest spec
	req := httptest.NewRequest(http.MethodGet, "/data/services.json", nil)
	w := httptest.NewRecorder()
	handler := spaHandler{staticPath: "/ui/build", indexPath: "/ui/build/index.html"}
	handler.ServeHTTP(w, req)
	//nolint:bodyclose
	if w.Result().StatusCode != 200 {
		t.Errorf("Not returning 200")
	}
}

func TestSpaHandlerNonExistent(t *testing.T) {
	// Test we can handle rest spec
	req := httptest.NewRequest(http.MethodGet, "/doesnotexist.json", nil)
	w := httptest.NewRecorder()
	handler := spaHandler{staticPath: "/ui/build", indexPath: "/ui/build/index.html"}
	handler.ServeHTTP(w, req)
	//nolint:bodyclose
	assert.Equal(t, 500, w.Result().StatusCode, "expected status code to be equal")
}
