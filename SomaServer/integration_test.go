package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"soma_server/api"
	"soma_server/hardware"
	"strings"
	"testing"
	"time"
)

// TestEndToEndIntegration tests the full pipeline from API request to ClojureV compilation
// and hardware state updates in the Go FPGADriver.
func TestEndToEndIntegration(t *testing.T) {
	// 1. Setup simulated hardware driver
	driver := hardware.NewFPGADriver()
	driver.LiveMode = false // Ensure we are in pure simulation mode
	
	server := api.NewServer(driver)
	
	// Start server on a test port
	go func() {
		_ = server.Start(":8099")
	}()
	
	// Give server a moment to start
	time.Sleep(200 * time.Millisecond)
	
	// 2. Prepare Synthesis Request
	cljvCode := `
		(ns ClojureV.qurq)
		(defn-ai grover_oracle [clk rst_n in]
			(let [target 0xABCDEF]
				(if (= in target)
					(qurq/phi-scale out in -1.0)
					(qurq/assign out in))))
	`
	payload := map[string]string{
		"code": cljvCode,
		"mode": "grover",
	}
	payloadBytes, _ := json.Marshal(payload)
	
	// 3. Fire Synthesis Request
	resp, err := http.Post("http://localhost:8099/api/synthesize", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatalf("Failed to execute synthesis POST request: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", resp.StatusCode)
	}
	
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)
	
	// Ensure the Go compiler successfully parsed and generated the Verilog
	if !strings.Contains(bodyStr, "success") {
		t.Errorf("Expected compilation success, got: %s", bodyStr)
	}
	
	// Read the output file to verify its content
	wd, _ := os.Getwd()
	projectRoot := filepath.Dir(wd)
	outputVerilog := filepath.Join(projectRoot, "build", "rtl", "sphy_core.v")
	verilogData, err := os.ReadFile(outputVerilog)
	if err != nil {
		t.Fatalf("Failed to read generated Verilog file: %v", err)
	}
	verilogStr := string(verilogData)
	if !strings.Contains(verilogStr, "module grover_oracle") {
		t.Errorf("Expected output file to contain 'module grover_oracle', got: %s", verilogStr)
	}
	
	// 4. Verify Hardware State update (Routing Mode should now be 'grover')
	stateResp, err := http.Get("http://localhost:8099/api/state")
	if err != nil {
		t.Fatalf("Failed to execute state GET request: %v", err)
	}
	defer stateResp.Body.Close()
	
	var state hardware.HardwareState
	if err := json.NewDecoder(stateResp.Body).Decode(&state); err != nil {
		t.Fatalf("Failed to decode state JSON: %v", err)
	}
	
	if state.RoutingMode != "grover" {
		t.Errorf("Expected RoutingMode to be updated to 'grover', got '%s'", state.RoutingMode)
	}
}
