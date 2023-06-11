package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerErr(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "/error", nil)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerErr)
	handler.ServeHTTP(responseRecorder, request)

	//Check for expected status code
	if statusCode := responseRecorder.Code; statusCode != 400 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			statusCode, 400)
	}

	//Check for expected body
	expected := []byte(`{"error":"Something went wrong"}`)
	if responseRecorder.Body.String() != string(expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			responseRecorder.Body.String(), expected)
	}

}
