package hardware

import (
        "math/rand"
        "fmt"
)

// FPGADriver simulates the 8-qubit 3D Scalable hardware interface.
type FPGADriver struct {
        C0          bool
        C1          bool
        C2          bool
        C3          bool
        C4          bool
        C5          bool
        C6          bool
        C7          bool
        Phase       float64
        ActiveCells int // The current d-state (8 cells = d=256)
}

// NewFPGADriver initializes the hardware link
func NewFPGADriver() *FPGADriver {
        return &FPGADriver{
                Phase:       0.0,
                ActiveCells: 8, // Default to 8 cells (d=256 GHZ state)
        }
}

// Poll reads the current hardware state of the 3D Macro-Cube
func (f *FPGADriver) Poll() {
        f.Phase += 0.1
        if f.Phase > 6.28 { // 2*PI reset
                f.Phase = 0.0
        }

        // Simulated thermal noise interference
        noise := rand.Float64() * 0.2

        // C0 (Anchor) simulates the primary geometric knot
        f.C0 = (f.Phase + noise) > 3.14

        // GHZ State: The 3D Entanglement Bus fan-out ensures all target cells
        // follow the anchor instantaneously.
        f.C1 = f.C0
        f.C2 = f.C0
        f.C3 = f.C0
        f.C4 = f.C0
        f.C5 = f.C0
        f.C6 = f.C0
        f.C7 = f.C0
}

// TriggerDPR simulates streaming a new partial bitstream to the FPGA fabric
func (f *FPGADriver) TriggerDPR(action string) {
        if action == "spawn" {
                f.ActiveCells++
                fmt.Printf("[DPR EVENT] Injecting new Macro-Cell bitstream. Active cells: %d\n", f.ActiveCells)
        } else if action == "collapse" && f.ActiveCells > 1 {
                f.ActiveCells--
                fmt.Printf("[DPR EVENT] Removing Macro-Cell bitstream. Active cells: %d\n", f.ActiveCells)
        }
}

// GetHardwareData returns the structured state for the API
func (f *FPGADriver) GetHardwareData() HardwareState {
        return HardwareState{
                C0:          f.C0,
                C1:          f.C1,
                C2:          f.C2,
                C3:          f.C3,
                C4:          f.C4,
                C5:          f.C5,
                C6:          f.C6,
                C7:          f.C7,
                ThermalLoad: rand.Float64() * 5.0 + 35.0, // 35C to 40C
                PhaseField:  f.Phase,
                ActiveCells: f.ActiveCells,
        }
}

// HardwareState represents the 8-qubit register payload
type HardwareState struct {
        C0          bool    `json:"c0"`
        C1          bool    `json:"c1"`
        C2          bool    `json:"c2"`
        C3          bool    `json:"c3"`
        C4          bool    `json:"c4"`
        C5          bool    `json:"c5"`
        C6          bool    `json:"c6"`
        C7          bool    `json:"c7"`
        ThermalLoad float64 `json:"thermal_load"` // Read from XADC
        PhaseField  float64 `json:"phase_field"`  // The SPHY tuning injection
        ActiveCells int     `json:"active_cells"` // The dynamic topological dimension
}
