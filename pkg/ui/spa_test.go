package ui

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestSpaHandlerNonExistent(t *testing.T) {
	// Test we can handle rest spec
	req := httptest.NewRequest(http.MethodGet, "/doesnotexist.json/", nil)
	w := httptest.NewRecorder()
	handler := spaHandler{staticPath: "/ui/build", indexPath: "/ui/build/index.html", fileSystem: pkgerFS{}}
	handler.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	assert.Equal(t, 500, result.StatusCode, "expected status code to be 500 but got %d", result.StatusCode)
}

func TestSpaHandlerInvalidPath(t *testing.T) {
	// Test we can handle rest spec
	req := httptest.NewRequest(http.MethodGet, "/data/doesnotexist.json", nil)
	w := httptest.NewRecorder()
	handler := spaHandler{staticPath: "/ui/build", indexPath: "/ui/build/index.html", fileSystem: pkgerFS{}}
	handler.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	assert.Equal(t, 500, result.StatusCode, "expected status code to be 500 but got %d", result.StatusCode)
}

type mockFileSystem struct {
}

func (m *mockFileSystem) Stat(name string) (os.FileInfo, error) {
	return nil, nil
}

type mockHTTPFileSystem struct{}

func (f mockHTTPFileSystem) Open(name string) (http.File, error) {
	appFS := afero.NewMemMapFs()
	if err := appFS.MkdirAll("/ui/build", 0755); err != nil {
		return nil, err
	}
	if err := appFS.MkdirAll("/ui/build/data/", 0755); err != nil {
		return nil, err
	}
	if err := afero.WriteFile(appFS, "/ui/build/data/services.json", []byte("file a"), 0644); err != nil {
		return nil, err
	}
	if err := afero.WriteFile(appFS, "/ui/build/index.html", []byte("file b"), 0644); err != nil {
		return nil, err
	}
	return appFS.Open(name)
}

func (m *mockFileSystem) Dir(name string) http.FileSystem {
	appFS := afero.NewMemMapFs()
	if err := appFS.MkdirAll("/ui/build", 0755); err != nil {
		return nil
	}
	if err := afero.WriteFile(appFS, "/ui/build/index.html", []byte("file b"), 0644); err != nil {
		return nil
	}
	return mockHTTPFileSystem{}
}
func TestSpaHandlerValidPath(t *testing.T) {
	// Test we can handle spa
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mockFileSystem := new(mockFileSystem)
	handler := spaHandler{staticPath: "/ui/build", indexPath: "/ui/build/index.html", fileSystem: mockFileSystem}
	handler.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	assert.Equal(t, 200, result.StatusCode, "expected status code to be 200 but got %d", result.StatusCode)
}

func TestSpaHandlerFolderPath(t *testing.T) {
	// Test we can handle folder path
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mockFileSystem := new(mockFileSystem)
	handler := spaHandler{staticPath: "/ui/build", indexPath: "/ui/build/index.html", fileSystem: mockFileSystem}
	handler.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	assert.Equal(t, 200, result.StatusCode, "expected status code to be 200 but got %d", result.StatusCode)
}

// This test currently fails when it runs on windows. Uncomment when pkger has been fixed or replaced
// func TestSpaHandlerWorksOnWindows(t *testing.T) {
// 	// Test we can handle rest spec
// 	req := httptest.NewRequest(http.MethodGet, "/", nil)
// 	w := httptest.NewRecorder()
// 	handler := spaHandler{staticPath: "/ui/build", indexPath: "/ui/build/index.html", fileSystem: pkgerFS{}}
// 	handler.ServeHTTP(w, req)
// 	defer w.Result().Body.Close()
// 	assert.Equal(t, 200, w.Result().StatusCode, "expected status code to be 200 but got %d", w.Result().StatusCode)
// }
