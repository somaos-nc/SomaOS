package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"soma_server/hardware"
	"sync"
	"time"
)

type SynthesisJob struct {
	ID         string `json:"id"`
	Status     string `json:"status"` // "running", "success", "error"
	Output     string `json:"output"` // Transpiler output
	VivadoLogs string `json:"vivado"`
	Error      bool   `json:"error"`
}

type Server struct {
	driver      *hardware.FPGADriver
	activeJobs  map[string]*SynthesisJob
	jobsMu      sync.RWMutex
}

func NewServer(driver *hardware.FPGADriver) *Server {
	return &Server{
		driver:     driver,
		activeJobs: make(map[string]*SynthesisJob),
	}
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
	mux.HandleFunc("/api/synthesize", s.handleSynthesize)
	mux.HandleFunc("/api/synthesize/status", s.handleSynthesisStatus)

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

func (s *Server) handleSynthesize(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var payload struct {
		Code string `json:"code"`
		Mode string `json:"mode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	jobID := fmt.Sprintf("job_%d", time.Now().Unix())
	job := &SynthesisJob{ID: jobID, Status: "running"}
	
	s.jobsMu.Lock()
	s.activeJobs[jobID] = job
	s.jobsMu.Unlock()

	fmt.Printf("[SYNTHESIS] Job %s initiated for mode: %s\n", jobID, payload.Mode)

	// Run synthesis in a background goroutine
	go func() {
		wd, _ := os.Getwd()
		projectRoot := filepath.Dir(wd)
		tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("soma_%s.cljv", jobID))
		os.WriteFile(tmpFile, []byte(payload.Code), 0644)

		outputVerilog := filepath.Join(projectRoot, "build", "rtl", "sphy_core.v")
		outputVerilogDir := filepath.Dir(outputVerilog)

		// 1. Transpilation
		fmt.Printf("[SYNTHESIS] Job %s: Running ClojureV Transpiler...\n", jobID)
		toolchainDir := filepath.Join(projectRoot, "ClojureV", "toolchain", "go")
		cmdV := exec.Command("go", "run", "./cmd/clojurev", "-target=verilog", "-in="+tmpFile, "-out="+outputVerilog)
		cmdV.Dir = toolchainDir
		outV, errV := cmdV.CombinedOutput()
		
		fmt.Printf("[SYNTHESIS] Job %s: Transpiler Finished. Result: %v\n", jobID, errV == nil)
		if len(outV) > 0 {
			fmt.Printf("[SYNTHESIS] Job %s: Output Snippet: %s\n", jobID, string(outV))
		}

		s.jobsMu.Lock()
		job.Output = string(outV)
		if errV != nil {
			job.Status = "error"
			job.Error = true
			fmt.Printf("[ERROR] Job %s: Transpilation Fracture detected.\n", jobID)
			s.jobsMu.Unlock()
			return
		}
		s.jobsMu.Unlock()

		// 2. Real Vivado Synthesis
		fmt.Printf("[VIVADO] Job %s: Starting Hardware Manifestation (Vivado Batch Mode)...\n", jobID)
		vivadoPath := "/home/noam/vivado/2025.2/Vivado/bin/vivado"
		tclScript := filepath.Join(projectRoot, "build", "synthesize_soma_os.tcl")
		
		exec.Command("cp", filepath.Join(projectRoot, "ClojureV", "src", "soma", "hardware", "dac_i2c_injector.v"), outputVerilogDir).Run()
		exec.Command("cp", filepath.Join(projectRoot, "ClojureV", "src", "soma", "hardware", "geometric_qubit.v"), outputVerilogDir).Run()
		exec.Command("cp", filepath.Join(projectRoot, "ClojureV", "src", "soma", "hardware", "top_quantum_virtualizer.v"), outputVerilogDir).Run()

		cmdVivado := exec.Command(vivadoPath, "-mode", "batch", "-source", tclScript)
		cmdVivado.Dir = projectRoot
		vivadoOut, errVivado := cmdVivado.CombinedOutput()

		fmt.Printf("[VIVADO] Job %s: Hardware Synthesis Finished. Result: %v\n", jobID, errVivado == nil)

		s.jobsMu.Lock()
		job.VivadoLogs = string(vivadoOut)
		if errVivado != nil {
			job.Status = "error"
			job.Error = true
			fmt.Printf("[ERROR] Job %s: Vivado Implementation Failure.\n", jobID)
		} else {
			// 3. PHYSICAL DEPLOYMENT VIA JTAG
			fmt.Printf("[DEPLOY] Job %s: Flashing bitstream to ALINX Board Silicon...\n", jobID)
			deployCmd := exec.Command("bash", "./deploy_mabel.sh")
			deployCmd.Dir = projectRoot
			deployOut, errDeploy := deployCmd.CombinedOutput()
			
			job.VivadoLogs += "\n\n=== JTAG DEPLOYMENT LOGS ===\n" + string(deployOut)
			
			if errDeploy != nil {
				job.Status = "error"
				job.Error = true
				fmt.Printf("[ERROR] Job %s: JTAG Deployment failed.\n", jobID)
			} else {
				job.Status = "success"
				fmt.Printf("[SUCCESS] Job %s: SPHY Manifold manifest on Silicon. Reconfiguring Routing.\n", jobID)
				s.driver.SetRoutingMode(payload.Mode)
			}
		}
		s.jobsMu.Unlock()
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"job_id": jobID, "status": "started"})
}

func (s *Server) handleSynthesisStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	jobID := r.URL.Query().Get("id")
	if jobID == "" {
		http.Error(w, "Missing job id", http.StatusBadRequest)
		return
	}

	s.jobsMu.RLock()
	job, ok := s.activeJobs[jobID]
	s.jobsMu.RUnlock()

	if !ok {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
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
