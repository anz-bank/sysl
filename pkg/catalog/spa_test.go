package catalog

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
