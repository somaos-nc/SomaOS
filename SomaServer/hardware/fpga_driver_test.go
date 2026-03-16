package hardware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFPGADriverHPQC(t *testing.T) {
	driver := NewFPGADriver()
	driver.LiveMode = false // Test simulation first

	// 1. Initial State Check (Single Cube Mode)
	if driver.IsStation != false {
		t.Error("Expected initial state to be Single Cube Mode (IsStation=false)")
	}
	if driver.ActiveCells != 8 {
		t.Errorf("Expected 8 active cells in single cube mode, got %d", driver.ActiveCells)
	}

	// 2. Transition to Station Mode (64-Qubit HPQC)
	driver.SetRoutingMode("station")
	if driver.IsStation != true {
		t.Error("Expected IsStation=true after setting 'station' mode")
	}
	if driver.ActiveCells != 64 {
		t.Errorf("Expected 64 active cells in station mode, got %d", driver.ActiveCells)
	}

	// 3. Verification of Entangled GHZ State (Station Mode)
	foundTrue := false
	for i := 0; i < 100; i++ {
		driver.Poll()
		if driver.Register == 0xFFFFFFFFFFFFFFFF {
			foundTrue = true
			break
		}
	}
	if !foundTrue {
		t.Error("Register never reached 0xFFFFFFFFFFFFFFFF in station mode")
	}

	// 4. DPR Scaling in Station Mode
	driver.TriggerDPR("collapse")
	if driver.ActiveCells != 56 {
		t.Errorf("Expected active cells to decrease by 8, got %d", driver.ActiveCells)
	}
	driver.TriggerDPR("spawn")
	if driver.ActiveCells != 64 {
		t.Errorf("Expected active cells to increase by 8, got %d", driver.ActiveCells)
	}

	// 5. GetHardwareData Metrics Coverage
	state := driver.GetHardwareData()
	if state.ThermalLoad < 30.0 || state.ThermalLoad > 45.0 {
		t.Errorf("Unexpected ThermalLoad: %f", state.ThermalLoad)
	}
	if state.Fidelity < 0.0 || state.Fidelity > 1.0 {
		t.Errorf("Invalid Fidelity: %f", state.Fidelity)
	}
}

func TestFPGADriverLiveMode(t *testing.T) {
	// Mock ALINX Board Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/telemetry" {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"reg":  0xABCDEF,
				"temp": 42.5,
			})
		} else if r.URL.Path == "/reset" {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	driver := NewFPGADriver()
	driver.LiveMode = true
	// Extract IP from mock server URL (e.g., http://127.0.0.1:12345)
	driver.BoardIP = strings.TrimPrefix(server.URL, "http://")

	// Test Live Polling
	driver.Poll()
	if driver.Register != 0xABCDEF {
		t.Errorf("Expected Register 0xABCDEF from Live Board, got 0x%X", driver.Register)
	}
	if driver.BaseTemp != 42.5 {
		t.Errorf("Expected Temp 42.5 from Live Board, got %f", driver.BaseTemp)
	}

	// Test Live Reset
	driver.WindingNumber = 100
	driver.ResetHardware()
	if driver.WindingNumber != 0 {
		t.Error("WindingNumber should reset to 0")
	}

	// Test Live Failure Fallback
	driver.BoardIP = "invalid_ip:0"
	driver.Poll() // Should fallback to simulation
	if driver.Register == 0xABCDEF {
		t.Error("Expected register to change during fallback simulation")
	}
}
