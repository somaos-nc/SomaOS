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
        ActiveCells int 
        RoutingMode string // "idle", "grover", "shor", "bell"
}

// NewFPGADriver initializes the hardware link
func NewFPGADriver() *FPGADriver {
        return &FPGADriver{
                Phase:       0.0,
                ActiveCells: 8,
                RoutingMode: "idle",
        }
}

// Poll reads the current hardware state of the 3D Macro-Cube
func (f *FPGADriver) Poll() {
        f.Phase += 0.1
        if f.Phase > 6.28 { f.Phase = 0.0 }

        noise := rand.Float64() * 0.2
        f.C0 = (f.Phase + noise) > 3.14

        // GHZ State: Constant Entanglement
        f.C1 = f.C0; f.C2 = f.C0; f.C3 = f.C0; f.C4 = f.C0; f.C5 = f.C0; f.C6 = f.C0; f.C7 = f.C0
}

func (f *FPGADriver) SetRoutingMode(mode string) {
        f.RoutingMode = mode
        fmt.Printf("[HARDWARE RECONFIG] Silicon routing set to: %s\n", mode)
}

func (f *FPGADriver) TriggerDPR(action string) {
        if action == "spawn" { f.ActiveCells++ } else if action == "collapse" && f.ActiveCells > 1 { f.ActiveCells-- }
}

func (f *FPGADriver) GetHardwareData() HardwareState {
        return HardwareState{
                C0: f.C0, C1: f.C1, C2: f.C2, C3: f.C3, C4: f.C4, C5: f.C5, C6: f.C6, C7: f.C7,
                ThermalLoad: rand.Float64()*5.0 + 35.0,
                PhaseField: f.Phase,
                ActiveCells: f.ActiveCells,
                RoutingMode: f.RoutingMode,
        }
}

type HardwareState struct {
        C0 bool `json:"c0"`; C1 bool `json:"c1"`; C2 bool `json:"c2"`; C3 bool `json:"c3"`
        C4 bool `json:"c4"`; C5 bool `json:"c5"`; C6 bool `json:"c6"`; C7 bool `json:"c7"`
        ThermalLoad float64 `json:"thermal_load"`
        PhaseField  float64 `json:"phase_field"`
        ActiveCells int     `json:"active_cells"`
        RoutingMode string  `json:"routing_mode"`
}
