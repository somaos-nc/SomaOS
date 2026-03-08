package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"soma_server/hardware"
	"testing"
)

func TestHandleState(t *testing.T) {
	driver := hardware.NewFPGADriver()
	server := NewServer(driver)

	req, err := http.NewRequest("GET", "/api/state", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handleState)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var state hardware.HardwareState
	err = json.NewDecoder(rr.Body).Decode(&state)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	if state.ThermalLoad == 0.0 {
		t.Errorf("Expected valid ThermalLoad, got 0.0")
	}
}
