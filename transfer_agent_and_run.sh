#!/bin/bash

# =====================================================================
# SomaOS: C-Agent SoC Deployment & Execution Script (ARM Cortex-A9)
# =====================================================================
# This script cross-compiles the C-Agent, transfers it to the ALINX board,
# and automatically starts the telemetry service.

set -e

# Setup Paths
PROJECT_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
AGENT_SRC="$PROJECT_ROOT/SilenceProtocol/src/soma_agent.c"
AGENT_BIN="$PROJECT_ROOT/SilenceProtocol/soma_agent"

# ALINX Board Credentials
BOARD_USER="root"
BOARD_IP="10.100.102.9" 
TARGET_PATH="/root/soma_agent"

echo "========================================================"
echo "    SomaOS: Deploying & Running C-Agent on SoC...      "
echo "========================================================"

# 1. Cross-Compile for ARM
echo ">> [1/3] Cross-compiling for ARM Cortex-A9..."
if ! command -v arm-linux-gnueabihf-gcc &> /dev/null; then
    echo ">> [ERROR] Cross-compiler 'arm-linux-gnueabihf-gcc' not found."
    exit 1
fi

cd "$PROJECT_ROOT/SilenceProtocol"
arm-linux-gnueabihf-gcc -O2 -Wall src/soma_agent.c -o soma_agent -static
echo ">> Build Successful: $AGENT_BIN"

# 2. Prepare Board
echo ">> [2/3] Preparing SoC (Stopping existing agent)..."
ssh -o ConnectTimeout=5 "$BOARD_USER@$BOARD_IP" "pkill -9 soma_agent || true" || echo ">> [WARNING] Pre-kill skipped (SSH unreachable)."

# 3. Transfer to Board
echo ">> [3/3] Transferring binary to ALINX board at $BOARD_IP..."
# Delete the destination first to avoid "Text file busy" errors
ssh -o ConnectTimeout=5 "$BOARD_USER@$BOARD_IP" "rm -f $TARGET_PATH" || true
scp -o ConnectTimeout=5 "$AGENT_BIN" "$BOARD_USER@$BOARD_IP:$TARGET_PATH"

# 4. Remote Execution
echo ">> [4/4] Starting SomaOS Agent on remote hardware..."
ssh "$BOARD_USER@$BOARD_IP" "chmod +x $TARGET_PATH; nohup $TARGET_PATH > /root/soma_agent.log 2>&1 &"

if [ $? -eq 0 ]; then
    echo "========================================================"
    echo " SUCCESS: SomaOS C-Agent is LIVE on the SoC."
    echo " Port: 8080"
    echo " Logs: /root/soma_agent.log"
    echo "========================================================"
else
    echo ">> [ERROR] Remote execution failed."
    exit 1
fi
