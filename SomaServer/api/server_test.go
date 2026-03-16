package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"soma_server/hardware"
	"strings"
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

func TestHandleSynthesize(t *testing.T) {
	driver := hardware.NewFPGADriver()
	server := NewServer(driver)

	payload := `{"code": "(ns ClojureV.qurq) (defn-ai test [in] (qurq/assign out in))", "mode": "grover"}`
	req, err := http.NewRequest("POST", "/api/synthesize", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handleSynthesize)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&resp)

	if resp["status"] != "success" {
		t.Errorf("Expected success status, got %v", resp["status"])
	}
	if resp["python"] == "" || resp["javascript"] == "" || resp["vivado"] == "" {
		t.Errorf("Missing reference language or vivado logs in response")
	}
}

func TestHandleReconfigure(t *testing.T) {
	driver := hardware.NewFPGADriver()
	server := NewServer(driver)

	payload := `{"mode": "station"}`
	req, err := http.NewRequest("POST", "/api/reconfigure", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handleReconfigure)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if driver.RoutingMode != "station" {
		t.Errorf("Expected routing mode 'station', got %s", driver.RoutingMode)
	}
}

func TestHandleDPR(t *testing.T) {
	driver := hardware.NewFPGADriver()
	server := NewServer(driver)

	payload := `{"action": "spawn"}`
	req, err := http.NewRequest("POST", "/api/dpr", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handleDPR)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if driver.ActiveCells != 16 {
		t.Errorf("Expected 16 active cells after spawn, got %d", driver.ActiveCells)
	}
}
