#!/bin/bash

# =====================================================================
# SomaOS: Live Synthesis Integration Test
# =====================================================================
# This script validates the full pipeline: 
# ClojureV -> Go Transpiler -> Verilog -> Vivado -> Bitstream

set -e

PROJECT_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
BUILD_DIR="$PROJECT_ROOT/build"
RTL_DIR="$BUILD_DIR/rtl"
API_URL="http://localhost:8081/api/synthesize"

echo "========================================================"
echo "    SomaOS: Testing Live Hardware Manifestation        "
echo "========================================================"

# 1. Internal Toolchain Sanity Check
echo ">> [1/4] Validating Go Transpiler..."
cd "$PROJECT_ROOT/ClojureV/toolchain/go"
go run ./cmd/clojurev -in="../../examples/bell_state_entanglement.cljv" -target=verilog
echo ">> Transpiler: OK"

# 2. Server Connectivity
echo ">> [2/4] Checking SomaServer availability..."
if ! curl -s "http://localhost:8081/api/state" > /dev/null; then
    echo ">> [FAIL] SomaServer not found on :8081. Please run ./start.sh first."
    exit 1
fi
echo ">> SomaServer: ONLINE"

# 3. Trigger Asynchronous Synthesis Job
echo ">> [3/4] Triggering Real Vivado Synthesis (Asynchronous)..."
TEST_CODE="(ns ClojureV.qurq) (defn-ai sphy_core [clk rst_n in_flux] (qurq/assign out in_flux))"
RESPONSE=$(curl -s -X POST -H "Content-Type: application/json" -d "{\"code\": \"$TEST_CODE\", \"mode\": \"idle\"}" "$API_URL")
JOB_ID=$(echo $RESPONSE | grep -oP '(?<="job_id":")[^"]*')

if [ -z "$JOB_ID" ]; then
    echo ">> [FAIL] Failed to initiate synthesis job. Response: $RESPONSE"
    exit 1
fi
echo ">> Job Started: $JOB_ID"

# 4. Polling for Completion
echo ">> [4/4] Polling Vivado status (this may take several minutes)..."
MAX_RETRIES=60 # 10 minutes max
COUNT=0
while [ $COUNT -lt $MAX_RETRIES ]; do
    STATUS_RESP=$(curl -s "http://localhost:8081/api/synthesize/status?id=$JOB_ID")
    STATUS=$(echo $STATUS_RESP | grep -oP '(?<="status":")[^"]*')
    
    if [ "$STATUS" == "success" ]; then
        echo ""
        echo ">> [INFO] Verifying Deployment Artifacts..."
        if [ ! -f "$BUILD_DIR/mabel_x8c.bit" ]; then
            echo ">> [FAIL] Bitstream file missing after successful job report."
            exit 1
        fi
        
        # Check if JTAG logs were produced
        HAS_JTAG=$(echo $STATUS_RESP | grep "JTAG DEPLOYMENT LOGS" || true)
        if [ -z "$HAS_JTAG" ]; then
            echo ">> [FAIL] Synthesis succeeded but JTAG deployment logs are missing."
            exit 1
        fi

        echo "========================================================"
        echo " SUCCESS: Full Intent-to-Silicon Pipeline Verified."
        echo " 1. Transpilation: OK"
        echo " 2. Vivado Synthesis: OK"
        echo " 3. JTAG Manifestation: OK"
        echo "========================================================"
        exit 0
    elif [ "$STATUS" == "error" ]; then
        echo ""
        echo ">> [FAIL] Vivado Synthesis failed. Logs:"
        echo $STATUS_RESP | python3 -c "import sys, json; print(json.load(sys.stdin).get('vivado', 'No logs'))"
        exit 1
    fi
    
    echo -n "."
    sleep 10
    COUNT=$((COUNT+1))
done

echo ">> [FAIL] Synthesis timed out after 10 minutes."
exit 1
