package api

import (
        "encoding/json"
        "net/http"
        "soma_server/hardware"
        "time"
)

// Server handles the API endpoints
type Server struct {
        driver *hardware.FPGADriver
}

// NewServer initializes the API server with the hardware driver
func NewServer(driver *hardware.FPGADriver) *Server {
        return &Server{driver: driver}
}

// Start launches the polling loop and the HTTP server
func (s *Server) Start(addr string) error {
        // Background polling of the "hardware"
        go func() {
                for {
                        s.driver.Poll()
                        time.Sleep(100 * time.Millisecond) // Poll at 10Hz
                }
        }()

        mux := http.NewServeMux()
        mux.HandleFunc("/api/state", s.handleState)
        mux.HandleFunc("/api/dpr", s.handleDPR)

        return http.ListenAndServe(addr, mux)
}

func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
        }

        state := s.driver.GetHardwareData()
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(state)
}

func (s *Server) handleDPR(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
        }

        if r.Method != "POST" {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
        }

        var intent struct {
                Action string `json:"action"` // "spawn" or "collapse"
        }

        if err := json.NewDecoder(r.Body).Decode(&intent); err != nil {
                http.Error(w, "Invalid intent payload", http.StatusBadRequest)
                return
        }

        s.driver.TriggerDPR(intent.Action)

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"status": "DPR sequence initiated", "action": intent.Action})
}
