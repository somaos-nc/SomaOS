package hardware

import (
	"testing"
)

func TestFPGADriverHPQC(t *testing.T) {
	driver := NewFPGADriver()

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
	// Poll until we hit a 'true' master state
	foundTrue := false
	for i := 0; i < 100; i++ {
		driver.Poll()
		if driver.Register == 0xFFFFFFFFFFFFFFFF {
			foundTrue = true
			break
		}
	}
	if !foundTrue {
		t.Error("Register never reached 0xFFFFFFFFFFFFFFFF in station mode after 100 polls")
	}

	// Ensure all bits are identical (GHZ correlation)
	if driver.Register != 0 && driver.Register != 0xFFFFFFFFFFFFFFFF {
		t.Errorf("Broken GHZ Entanglement: Register state 0x%X contains mixed bits", driver.Register)
	}

	// 4. DPR Scaling in Station Mode
	driver.TriggerDPR("collapse")
	if driver.ActiveCells != 56 {
		t.Errorf("Expected active cells to decrease by 8 (cube collapse), got %d", driver.ActiveCells)
	}

	// 5. Revert to Single Cube Mode
	driver.SetRoutingMode("grover")
	if driver.IsStation != false {
		t.Error("Expected IsStation=false after switching to 'grover' mode")
	}
	if driver.ActiveCells != 8 {
		t.Errorf("Expected 8 active cells in single cube mode, got %d", driver.ActiveCells)
	}
}
