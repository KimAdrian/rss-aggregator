package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerReadiness(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/ready", nil)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerReadiness)
	handler.ServeHTTP(responseRecorder, req)

	//Check expected status code
	if status := responseRecorder.Code; status != 200 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, 200)
	}
	//Check expected body
	expectedBody := "{}"
	if responseRecorder.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			responseRecorder.Body.String(), expectedBody)
	}
}
