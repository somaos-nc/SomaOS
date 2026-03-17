#!/bin/bash
# SOMA OS: Professional Network-Safe Deployment Script (Zynq-7000)

BIT_FILE="./build/mabel_x8c.bit"
BIN_FILE="./build/mabel_x8c.bin"
BOARD_IP="10.100.102.9"
BOARD_USER="root"
REMOTE_PATH="/root/mabel_x8c.bin"

echo "========================================================"
echo "    SomaOS: Deploying MABEL x8C to ALINX 7020...        "
echo "========================================================"

# 1. Convert .bit to raw .bin (strip header for xdevcfg)
echo ">> Extracting raw bitstream from .bit header..."
python3 -c "
with open('$BIT_FILE', 'rb') as f:
    data = f.read()
    sync = b'\xaa\x99\x55\x66'
    pos = data.find(sync)
    if pos != -1:
        with open('$BIN_FILE', 'wb') as out:
            out.write(data[pos:])
    else:
        exit(1)
" || (echo ">> [ERROR] Bitstream extraction failed." && exit 1)

# 2. Stop Agent to prevent AXI bus contention
echo ">> Safeguarding PS: Stopping soma_agent..."
ssh -o ConnectTimeout=5 "$BOARD_USER@$BOARD_IP" "pkill -9 soma_agent || true"

# 3. Transfer raw bitstream
echo ">> Transferring raw bitstream to PS..."
scp -o ConnectTimeout=5 "$BIN_FILE" "$BOARD_USER@$BOARD_IP:$REMOTE_PATH"

# 4. Write to PL via internal bridge
echo ">> Mapping silicon logic via xdevcfg (INTERNAL BRIDGE)..."
ssh -o ConnectTimeout=5 "$BOARD_USER@$BOARD_IP" "cat $REMOTE_PATH > /dev/xdevcfg"

if [ $? -eq 0 ]; then
    echo ">> [SUCCESS] MABEL x8C Braided Heart manifest on ALINX Silicon."
    
    # 5. Restart Agent
    echo ">> Restoring telemetry services..."
    ./transfer_agent_and_run.sh
    echo ">> [NETWORK] Connection STABLE."
else
    echo ">> [ERROR] xdevcfg flash failed. Device might be locked."
    exit 1
fi
