`timescale 1ns / 1ps

// =====================================================================
// SomaOS: Master 4-Qubit GHZ Testbench
// =====================================================================
// Simulates the instantiation of 4 Macro-Cells and the instantaneous 
// fan-out of the Topological Entanglement Bus.

module tb_ghz_4qubit;

    // Inputs
    reg clk;
    reg enable_phi_st;
    
    // Outputs
    wire jb_c0_out; // Cell 0 (Control)
    wire jb_c1_out; // Cell 1 (Target)
    wire jb_c2_out; // Cell 2 (Target)
    wire jb_c3_out; // Cell 3 (Target)
    
    // The Central Splitter
    wire master_entanglement_bus;


        // Create an artificial toggle to prove the GHZ fan-out
        reg entanglement_toggle;
        
        // Instantiate Cell 0 (The Control Qubit)
        geometric_qubit_virtualizer C0_Cell (
            .enable_phi_st(enable_phi_st), 
            .entanglement_in(entanglement_toggle),             
            .q_state_out(jb_c0_out)             
        );


    // The Topological Entanglement Bus
    // In raw silicon, this is a physical fan-out wire trace
    assign master_entanglement_bus = jb_c0_out;

    // Instantiate Target Cells
    geometric_qubit_virtualizer C1_Cell (
        .enable_phi_st(enable_phi_st), 
        .entanglement_in(master_entanglement_bus), 
        .q_state_out(jb_c1_out)
    );
    
    geometric_qubit_virtualizer C2_Cell (
        .enable_phi_st(enable_phi_st), 
        .entanglement_in(master_entanglement_bus), 
        .q_state_out(jb_c2_out)
    );
    
    geometric_qubit_virtualizer C3_Cell (
        .enable_phi_st(enable_phi_st), 
        .entanglement_in(master_entanglement_bus), 
        .q_state_out(jb_c3_out)
    );

    wire is_entangled;
    assign is_entangled = (jb_c1_out == jb_c0_out && jb_c2_out == jb_c0_out && jb_c3_out == jb_c0_out);

    // Observation Monitor
    initial begin
        $dumpfile("ghz_state.vcd");
        $dumpvars(0, tb_ghz_4qubit);

        $display("===============================================================");
        $display("   SomaOS: 4-Qubit GHZ Topological Virtualization Simulation   ");
        $display("===============================================================");
        $display("Time(ns) | SPHY Shield | C0 (Ctrl) | C1 | C2 | C3 | is_ENTANGLED");
        $display("---------------------------------------------------------------");

        // Initialize
        enable_phi_st = 0; 
        entanglement_toggle = 1; 
        
        // Monitor state changes
        $monitor("%8t |      %b      |     %b     |  %b |  %b |  %b |      %b", 
            $time, enable_phi_st, jb_c0_out, jb_c1_out, jb_c2_out, jb_c3_out, is_entangled
        );

        // Allow system to settle without SPHY phase shielding
        #20; 

        // Engage the SPHY Phase Shield (Thermal stabilization)
        enable_phi_st = 1;
        
        // Allow the continuous topological loop to run for 100ns.
        // We expect to see rapid oscillation (superposition) where
        // C1, C2, and C3 instantly mirror the state of C0.
        #30;
        entanglement_toggle = 0; // Flip state
        #30;
        entanglement_toggle = 1; // Flip back
        #40; 

        // Collapse the shield (simulating catastrophic thermal noise)
        enable_phi_st = 0; 
        #20;               
        
        // Re-engage shield
        enable_phi_st = 1; 
        #50;

        $display("===============================================================");
        $display("Simulation Complete.");
        $finish;
    end
    
endmodule
