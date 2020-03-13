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
	defer w.Result().Body.Close()
	assert.Equal(t, 500, w.Result().StatusCode, "expected status code to be 500 but got %d", w.Result().StatusCode)
}

func TestSpaHandlerInvalidPath(t *testing.T) {
	// Test we can handle rest spec
	req := httptest.NewRequest(http.MethodGet, "/data/doesnotexist.json", nil)
	w := httptest.NewRecorder()
	handler := spaHandler{staticPath: "/ui/build", indexPath: "/ui/build/index.html", fileSystem: pkgerFS{}}
	handler.ServeHTTP(w, req)
	defer w.Result().Body.Close()
	assert.Equal(t, 500, w.Result().StatusCode, "expected status code to be 500 but got %d", w.Result().StatusCode)
}

type mockFileSystem struct {
}

func (m *mockFileSystem) Stat(name string) (os.FileInfo, error) {
	return nil, nil
}

type mockHttpFileSystem struct{}

func (f mockHttpFileSystem) Open(name string) (http.File, error) {
	appFS := afero.NewMemMapFs()
	appFS.MkdirAll("/ui/build", 0755)
	appFS.MkdirAll("/ui/build/data/", 0755)
	afero.WriteFile(appFS, "/ui/build/data/services.json", []byte("file a"), 0644)
	afero.WriteFile(appFS, "/ui/build/index.html", []byte("file b"), 0644)
	return appFS.Open(name)
}

func (m *mockFileSystem) Dir(name string) http.FileSystem {
	appFS := afero.NewMemMapFs()
	appFS.MkdirAll("/ui/build", 0755)
	afero.WriteFile(appFS, "/ui/build/index.html", []byte("file b"), 0644)
	return mockHttpFileSystem{}
}
func TestSpaHandlerValidPath(t *testing.T) {
	// Test we can handle spa
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mockFileSystem := new(mockFileSystem)
	handler := spaHandler{staticPath: "/ui/build", indexPath: "/ui/build/index.html", fileSystem: mockFileSystem}
	handler.ServeHTTP(w, req)
	defer w.Result().Body.Close()
	assert.Equal(t, 200, w.Result().StatusCode, "expected status code to be 200 but got %d", w.Result().StatusCode)
}

func TestSpaHandlerFolderPath(t *testing.T) {
	// Test we can handle folder path
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mockFileSystem := new(mockFileSystem)
	handler := spaHandler{staticPath: "/ui/build", indexPath: "/ui/build/index.html", fileSystem: mockFileSystem}
	handler.ServeHTTP(w, req)
	defer w.Result().Body.Close()
	assert.Equal(t, 200, w.Result().StatusCode, "expected status code to be 200 but got %d", w.Result().StatusCode)
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
