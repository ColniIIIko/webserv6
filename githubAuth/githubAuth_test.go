package githubauth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleLogin(t *testing.T) {
	// req, err := http.NewRequest("GET", "/github/login", nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// rr := httptest.NewRecorder()
	// handler := http.HandlerFunc(HandleLogin)

	// handler.ServeHTTP(rr, req)

	// if status := rr.Code; status != http.StatusTemporaryRedirect {
	// 	t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusTemporaryRedirect)
	// }

	// if location := rr.Header().Get("Location"); location == "" {
	// 	t.Error("handler did not redirect to the login URL")
	// }
	t.Error("test")
}

func TestHandleCallbackInvalidState(t *testing.T) {
	req, err := http.NewRequest("GET", "/github/callback?state=invalid", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleCallback)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
