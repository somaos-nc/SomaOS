#!/bin/bash

# =====================================================================
# SomaOS: Unified Master Test Suite (v4.3)
# =====================================================================
# This script executes all automated tests across the entire SomaOS 
# codebase, including the Go compiler, the backend API, and the 
# hardware drivers.

set -e

# Setup Paths
PROJECT_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

echo "========================================================"
echo "    SomaOS HPQC: Initiating Master Test Protocol        "
echo "========================================================"

# 1. Test ClojureV Transpiler & Parser
echo ">> [1/3] Testing ClojureV Toolchain (Parser & Verilog Emitter)..."
cd "$PROJECT_ROOT/ClojureV/toolchain/go"
go test -v -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total | awk '{print ">> Transpiler Coverage: " $3}'
echo "--------------------------------------------------------"

# 2. Test SomaServer API
echo ">> [2/3] Testing SomaServer Backend API Routing..."
cd "$PROJECT_ROOT/SomaServer"
go test -v -coverprofile=api_coverage.out ./api/...
go tool cover -func=api_coverage.out | grep total | awk '{print ">> API Router Coverage: " $3}'
echo "--------------------------------------------------------"

# 3. Test SomaServer Hardware Driver (HPQC Logic)
echo ">> [3/3] Testing SomaServer Hardware Telemetry & Scaling..."
cd "$PROJECT_ROOT/SomaServer"
go test -v -coverprofile=hw_coverage.out ./hardware/...
go tool cover -func=hw_coverage.out | grep total | awk '{print ">> Hardware Driver Coverage: " $3}'
echo "--------------------------------------------------------"

echo "========================================================"
echo " [SUCCESS] All SomaOS subsystems verified and stable."
echo "========================================================"
