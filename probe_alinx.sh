#!/bin/bash
# SomaOS: Robust ALINX Hardware Probe
# Scans the local network to find the Zynq-7000 FPGA board.

echo "========================================================"
echo "    SomaOS: Probing for ALINX 7020 Hardware...         "
echo "========================================================"

# 1. Identify local IP and Subnet
INTERFACE=$(ip route | grep default | awk '{print $5}')
LOCAL_IP=$(ip -o -4 addr list "$INTERFACE" | awk '{print $4}' | cut -d/ -f1)
SUBNET=$(echo "$LOCAL_IP" | cut -d. -f1-3)

if [ -z "$SUBNET" ]; then
    echo ">> [ERROR] Could not identify local subnet. Check your network connection."
    exit 1
fi

echo ">> Local Interface: $INTERFACE"
echo ">> Local IP: $LOCAL_IP"
echo ">> Scanning Subnet: $SUBNET.0/24"

# 2. Fast Ping Scan
# We ping every address in the /24 subnet to populate the ARP cache.
# We use a very short timeout (-W 1) and run in background for speed.
echo ">> Pinging all addresses in $SUBNET.1-254... (This will take ~10 seconds)"
for i in {1..254}; do
    ping -c 1 -W 1 "$SUBNET.$i" >/dev/null 2>&1 &
done
wait

echo ">> Scan complete. Checking ARP table for ALINX/Xilinx signatures..."

# 3. Filter for Xilinx/ALINX OUI (00:0a:35)
ALINX_CANDIDATE=$(ip neighbor show | grep -i "00:0a:35")

if [ -n "$ALINX_CANDIDATE" ]; then
    echo ">> [MATCH FOUND] Potential ALINX board detected:"
    echo "$ALINX_CANDIDATE"
    
    BOARD_IP=$(echo "$ALINX_CANDIDATE" | awk '{print $1}')
    echo "--------------------------------------------------------"
    echo ">> SUCCESS: Target IP identified as: $BOARD_IP"
    echo ">> Update BoardIP in SomaServer/hardware/fpga_driver.go"
    echo "--------------------------------------------------------"
else
    echo ">> [WARNING] No exact OUI match found."
    echo ">> Listing all reachable devices for manual identification:"
    ip neighbor show | grep REACHABLE
fi
