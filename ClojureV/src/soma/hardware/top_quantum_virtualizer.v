`timescale 1ns / 1ps

// =====================================================================
// SOMA OS: Master Top-Level Virtualizer
// =====================================================================
// Bridges the virtual math (SPHY Core) with the physical thermal
// environment (XADC & DAC).

module top_quantum_virtualizer (
    input wire CLK100MHZ,       // 100MHz System Clock
    
    // I2C Output to DAC (Phase Injector)
    output wire ja_scl,         
    inout  wire ja_sda,         
    
    // Virtualized Superposition Outputs
    output wire jb_q0,          
    output wire jb_q1,          
    
    // FPGA Internal Thermal Sensors (XADC)
    input wire vauxp0,          
    input wire vauxn0           
);

    wire [15:0] xadc_temp_data_full;  
    wire [11:0] xadc_temp_data;       
    wire [11:0] calculated_psi_sc;    
    wire trigger_i2c;                 

    // 1. Internal Observer: Digilent XADC Core
    // Reads the ambient temperature / thermal noise of the routing matrix
    xadc_wiz_0 internal_observer (
        .daddr_in(7'h00),             
        .den_in(1'b1),                
        .di_in(16'b0), 
        .dwe_in(1'b0), 
        .do_out(xadc_temp_data_full),      
        .drdy_out(trigger_i2c),       
        .dclk_in(CLK100MHZ), 
        .vp_in(1'b0), 
        .vn_in(1'b0),
        .vauxp0(vauxp0),              
        .vauxn0(vauxn0)
    );

    // The XADC returns a 16-bit value where the top 12 bits are the ADC reading.
    assign xadc_temp_data = xadc_temp_data_full[15:4];

    // 2. The Inverse Phase Transformer
    // E_CA = -[V_SI + \Psi_SC]. We map the inverse thermal data.
    assign calculated_psi_sc = ~xadc_temp_data; 
    
    // 3. The DAC Phase Tuner (\Phi_{ST} Injector)
    dac_i2c_injector phase_tuner (
        .clk(CLK100MHZ),
        .trigger_injection(trigger_i2c),
        .psi_sc(calculated_psi_sc),
        .i2c_scl(ja_scl),
        .i2c_sda(ja_sda)
    );

    // 4. Connect to SPHY Core (Compiled from ClojureV)
    wire [23:0] s_core_in;
    wire [23:0] s_core_out;
    
    // Route thermal noise into SPHY Core for dampening
    assign s_core_in = {12'b0, xadc_temp_data};

    // Instantiate our dynamically compiled AST module!
    sphy_core engine (
        .clk(CLK100MHZ),
        .rst_n(1'b1),
        .in_flux(s_core_in),
        .out(s_core_out)
    );

    // Extract binary projection from the continuous 24-bit damped wave
    assign jb_q0 = s_core_out[0];
    assign jb_q1 = s_core_out[1];

endmodule
