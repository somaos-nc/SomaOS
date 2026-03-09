#!/bin/bash

# Exit on any error
set -e

echo "=== SomaOS: Macroscopic GHZ State Simulation ==="

# We compile the raw gate logic exactly as specified in the new whitepaper
echo ">> Synthesizing 4-Qubit GHZ Hardware Logic (iverilog)..."
iverilog -o ghz_sim.vvp geometric_qubit.v tb_ghz_4qubit.v

# Run the Simulation
echo ">> Executing Topological Entanglement (vvp)..."
vvp ghz_sim.vvp

echo "=== Pipeline Complete ==="
echo "You can view the macroscopic probability waves using: gtkwave ghz_state.vcd"
