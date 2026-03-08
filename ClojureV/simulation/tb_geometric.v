`timescale 1ns / 1ps

module tb_geometric;

    // Inputs
    reg clk;
    reg rst_n;
    reg [23:0] in_state;

    // Outputs
    wire [23:0] out;

    // Instantiate the Unit Under Test (UUT)
    // The name matches the defn-ai from the ClojureV file
    topological_knot uut (
        .clk(clk), 
        .rst_n(rst_n), 
        .in_state(in_state), 
        .out(out)
    );

    // Clock generation (100MHz -> 10ns period)
    always #5 clk = ~clk;

    initial begin
        // Setup waveform dumping for GTKWave
        $dumpfile("topological_knot.vcd");
        $dumpvars(0, tb_geometric);

        // Initialize Inputs
        clk = 0;
        rst_n = 0;
        in_state = 24'h000000;

        // Reset the system
        #20;
        rst_n = 1;

        // Test Case 1: Inject a base quantum state
        #10;
        in_state = 24'h555555;
        $display("Time=%0t | State Injected: %h", $time, in_state);

        // Wait for clock cycles to process through the manifold
        #20;
        $display("Time=%0t | Output State: %h", $time, out);

        // Test Case 2: Inject a complex harmonic seed
        #10;
        in_state = 24'hF0F0F0;
        $display("Time=%0t | Harmonic Seed Injected: %h", $time, in_state);

        #20;
        $display("Time=%0t | Final Entangled Output: %h", $time, out);

        // Complete the simulation
        #50;
        $display("Simulation Complete.");
        $finish;
    end
      
endmodule
