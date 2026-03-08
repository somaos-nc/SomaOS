#!/bin/bash

# Exit on any error
set -e

echo "=== SomaOS Topological Simulation Pipeline ==="

# 1. Compile the ClojureV DSL into Verilog
echo "[1/3] Transpiling ClojureV to Verilog Manifold..."
cd ../toolchain/go
go run ./cmd/clojurev -target=verilog -in=../../simulation/test_circuit.cljv -out=../../simulation/topological_knot.v
cd ../../simulation

# 2. Compile the Verilog and Testbench using Icarus Verilog
echo "[2/3] Synthesizing Hardware Logic (iverilog)..."
iverilog -o knot_sim.vvp topological_knot.v tb_geometric.v

# 3. Run the Simulation
echo "[3/3] Executing Topological Simulation (vvp)..."
vvp knot_sim.vvp

echo "=== Pipeline Complete ==="
echo "You can view the resulting waves using: gtkwave topological_knot.vcd"
