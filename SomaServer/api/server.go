package api

import (
        "encoding/json"
        "net/http"
        "soma_server/hardware"
        "time"
)

type Server struct {
        driver *hardware.FPGADriver
}

func NewServer(driver *hardware.FPGADriver) *Server {
        return &Server{driver: driver}
}

func (s *Server) Start(addr string) error {
        go func() {
                for {
                        s.driver.Poll()
                        time.Sleep(100 * time.Millisecond)
                }
        }()

        mux := http.NewServeMux()
        mux.HandleFunc("/api/state", s.handleState)
        mux.HandleFunc("/api/dpr", s.handleDPR)
        mux.HandleFunc("/api/reconfigure", s.handleReconfigure)

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

func (s *Server) handleReconfigure(w http.ResponseWriter, r *http.Request) {
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
                Mode string `json:"mode"`
        }

        if err := json.NewDecoder(r.Body).Decode(&intent); err != nil {
                http.Error(w, "Invalid intent payload", http.StatusBadRequest)
                return
        }

        s.driver.SetRoutingMode(intent.Mode)

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"status": "Reconfiguration complete", "mode": intent.Mode})
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
                Action string `json:"action"`
        }

        if err := json.NewDecoder(r.Body).Decode(&intent); err != nil {
                http.Error(w, "Invalid intent payload", http.StatusBadRequest)
                return
        }

        s.driver.TriggerDPR(intent.Action)

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"status": "DPR sequence initiated", "action": intent.Action})
}
