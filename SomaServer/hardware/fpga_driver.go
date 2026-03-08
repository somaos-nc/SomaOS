package hardware

import (
	"math/rand"
)

// FPGADriver simulates the hardware interface to the Universal NAND Topology.
// In a real environment, this would read from /dev/mem or a PCIe mapping
// to interact directly with the XADC and the Zynq PL logic gates.
type FPGADriver struct {
	StateQ0 bool
	StateQ1 bool
	Phase   float64
}

// NewFPGADriver initializes the hardware link
func NewFPGADriver() *FPGADriver {
	return &FPGADriver{
		Phase: 0.0,
	}
}

// Poll reads the current hardware state of the 2x2 logic block
func (f *FPGADriver) Poll() {
	// Simulate the continuous "sloshing" of the superposition knot
	// and the injection of the SPHY Phase tuning field.
	f.Phase += 0.1
	if f.Phase > 6.28 { // 2*PI reset
		f.Phase = 0.0
	}

	// Simulated thermal noise interference
	noise := rand.Float64() * 0.2

	// Update simulated NAND knot states based on phase and noise
	f.StateQ0 = (f.Phase + noise) > 3.14
	f.StateQ1 = (f.Phase - noise) < 3.14
}

// GetHardwareData returns the structured state for the API
func (f *FPGADriver) GetHardwareData() HardwareState {
	return HardwareState{
		Q0:          f.StateQ0,
		Q1:          f.StateQ1,
		ThermalLoad: rand.Float64() * 5.0 + 35.0, // 35C to 40C
		PhaseField:  f.Phase,
	}
}

// HardwareState represents the payload sent to the visualization frontend
type HardwareState struct {
	Q0          bool    `json:"q0"`
	Q1          bool    `json:"q1"`
	ThermalLoad float64 `json:"thermal_load"` // Read from XADC
	PhaseField  float64 `json:"phase_field"`  // The SPHY tuning injection
}
