package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func CreateUser(t *testing.T) {
	var jsonStr = []byte(`{"name":"test","email":"test@gmail.com","password":"password"}`)
	req, err := http.NewRequest("POST", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	
	// handler.ServeHTTP(rr, req)
	// if status := rr.Code; status != http.StatusOK {
	// 	t.Errorf("handler returned wrong status code: got %v want %v",
	// 		status, http.StatusOK)
	// }

	// // Check the response body is what we expect.
	// expected := `[{"name":"test","email":"test@gmail.com","password":"password"}]`
	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }
}
