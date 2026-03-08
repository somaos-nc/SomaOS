package hardware

import (
	"testing"
)

func TestFPGADriver(t *testing.T) {
	driver := NewFPGADriver()

	if driver.Phase != 0.0 {
		t.Errorf("Expected initial phase 0.0, got %f", driver.Phase)
	}

	// Poll should advance the phase
	driver.Poll()
	if driver.Phase <= 0.0 {
		t.Errorf("Expected phase to advance, got %f", driver.Phase)
	}

	state := driver.GetHardwareData()
	if state.ThermalLoad < 35.0 || state.ThermalLoad > 40.0 {
		t.Errorf("Expected thermal load to simulate realistic hardware temps (35-40C), got %f", state.ThermalLoad)
	}
}
