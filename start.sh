#!/bin/bash

# =====================================================================
# SomaOS: Unified System Launcher
# =====================================================================
# This script launches both the Go FPGA Hardware Simulator 
# and the React 3D Visualizer concurrently.

echo "========================================================"
echo "    SomaOS v3.0: Initializing Visualization Environment "
echo "========================================================"

# Trap SIGINT (Ctrl+C) to clean up both background processes
trap "echo 'Shutting down SomaOS...'; pkill -P $$; exit" SIGINT SIGTERM

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 1. Start the Go Hardware Server
echo ">> [1/2] Booting Go FPGA Hardware Proxy on port 8081..."
cd "$PROJECT_ROOT/SomaServer"
go run main.go &
GO_PID=$!

# Wait a moment for the server to initialize
sleep 2

# 2. Start the React Frontend
echo ">> [2/2] Booting React WebGL Visualizer..."
cd "$PROJECT_ROOT/SomaUI/soma_web"
npm run dev &
VITE_PID=$!

echo "========================================================"
echo " SUCCESS: SomaOS Environment is Live!"
echo " 👉 Visualizer: http://localhost:5173"
echo " 👉 API:        http://localhost:8081/api/state"
echo " "
echo " Press Ctrl+C to terminate both servers."
echo "========================================================"

# Keep the script running to hold the trap active
wait $GO_PID $VITE_PID
