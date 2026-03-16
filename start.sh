#!/bin/bash

# =====================================================================
# SomaOS: Unified System Launcher (HPQC Master Logging Edition)
# =====================================================================
# This script launches the Go Hardware Proxy, the Flutter Web UI, 
# and the Central WebSocket Logger concurrently.

# Trap SIGINT (Ctrl+C) to clean up all background processes
cleanup() {
    echo "Shutting down SomaOS..."
    pkill -P $$
    exit
}
trap cleanup SIGINT SIGTERM

# Kill any existing processes on project ports (8081, 8082, 8083, 5173)
cleanup_stale() {
    echo ">> [INIT] Cleaning up stale SomaOS processes..."
    ports=(8081 8082 8083 5173)
    for port in "${ports[@]}"; do
        pid=$(lsof -t -i :$port || true)
        if [ -n "$pid" ]; then
            echo "   - Killing process $pid on port $port"
            kill -9 $pid >/dev/null 2>&1 || true
        fi
    done
}
cleanup_stale

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOG_FILE="$PROJECT_ROOT/master_execution.log"

# Clear old logs
echo ">> [INIT] Resetting Master Execution Log: $LOG_FILE"
echo "SomaOS Session Start: $(date)" > "$LOG_FILE"

# 1. Start the WebSocket Master Logger
echo ">> [1/3] Booting WebSocket Logger on port 8082..."
cd "$PROJECT_ROOT/logger"
go run main.go >> "$LOG_FILE" 2>&1 &
sleep 2

# 2. Start the Go Hardware Server
echo ">> [2/3] Booting Go FPGA Hardware Proxy on port 8081..."
cd "$PROJECT_ROOT/SomaServer"
stdbuf -oL go run main.go 2>&1 | while read line; do echo "[HARDWARE] $line" >> "$LOG_FILE"; done &
sleep 2

# 3. Start the SomaAI Cortex Router (Vertex AI Multimodal)
echo ">> [3/4] Booting SomaAI Cortex Router on port 8083..."
cd "$PROJECT_ROOT"
source SomaAI/.venv/bin/activate
stdbuf -oL python3 SomaAI/src/router.py 2>&1 | while read line; do echo "[AI-CORTEX] $line" >> "$LOG_FILE"; done &
sleep 2

# 4. Start the React/Vite Web Visualizer (HPQC Dashboard + IDE)
echo ">> [4/4] Booting React/Vite Visualizer on port 5173..."
cd "$PROJECT_ROOT/SomaUI/soma_web"
npm run dev -- --port 5173 --host 2>&1 | while read line; do echo "[WEB-VISUALIZER] $line" >> "$LOG_FILE"; done &

echo "========================================================"
echo " SUCCESS: SomaOS HPQC Environment is Live!"
echo " 👉 React Dashboard: http://localhost:5173"
echo " 👉 Master Log: master_execution.log"
echo "========================================================"

# Wait for background processes to finish
wait
