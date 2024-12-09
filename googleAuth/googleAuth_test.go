package googleauth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleLogin(t *testing.T) {
	req, err := http.NewRequest("GET", "/google/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleLogin)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusTemporaryRedirect)
	}

	if location := rr.Header().Get("Location"); location == "" {
		t.Error("handler did not redirect to the login URL")
	}
}

func TestHandleCallbackInvalidState(t *testing.T) {
	req, err := http.NewRequest("GET", "/google/callback?state=invalid", nil)
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
