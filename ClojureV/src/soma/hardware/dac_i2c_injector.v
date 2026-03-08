`timescale 1ns / 1ps

// =====================================================================
// SOMA OS: I2C Thermodynamic Feedback Injector
// =====================================================================
// Injects the calculated Phase Tuning Field (\Phi_{ST}) back into the
// hardware substrate using an external MCP4725 I2C DAC.

module dac_i2c_injector (
    input wire clk,               // 100MHz System Clock
    input wire trigger_injection, // Goes high when XADC has fresh data
    input wire [11:0] psi_sc,     // The 12-bit Stochastic Compensation operator
    output reg i2c_scl,           // I2C Clock Line
    inout wire i2c_sda            // I2C Data Line
);

    // I2C Clock Divider: 100MHz -> ~100kHz I2C Bus Speed
    reg [9:0] clk_div = 0;
    reg i2c_clk_en = 0;
    
    always @(posedge clk) begin
        if (clk_div == 10'd499) begin
            clk_div <= 0;
            i2c_clk_en <= 1; 
        end else begin
            clk_div <= clk_div + 1;
            i2c_clk_en <= 0;
        end
    end

    reg [6:0] state = 0;
    reg [7:0] shift_reg;
    reg [3:0] bit_count;
    reg sda_out_reg = 1; 
    
    // Open-drain architecture for I2C data line
    assign i2c_sda = (sda_out_reg == 0) ? 1'b0 : 1'bz;

    always @(posedge clk) begin
        if (i2c_clk_en) begin
            case (state)
                0: begin // Idle
                    i2c_scl <= 1;
                    sda_out_reg <= 1;
                    if (trigger_injection) state <= 1;
                end
                
                // START Condition: SDA goes low while SCL is high
                1: begin 
                    sda_out_reg <= 0;
                    state <= 2;
                end
                
                // DAC Device Address Byte (Write mode: 11000000 for MCP4725)
                2: begin 
                    i2c_scl <= 0;
                    shift_reg <= 8'b1100_0000;
                    bit_count <= 7;
                    state <= 3;
                end
                
                // Write Bit Phase 1: Set SDA
                3: begin 
                    sda_out_reg <= shift_reg[bit_count];
                    i2c_scl <= 1; 
                    state <= 4;
                end
                
                // Write Bit Phase 2: Lower SCL
                4: begin 
                    i2c_scl <= 0;
                    if (bit_count == 0) state <= 5; 
                    else begin
                        bit_count <= bit_count - 1;
                        state <= 3;
                    end
                end
                
                // Wait for ACK
                5: begin 
                    sda_out_reg <= 1; // Release SDA to listen
                    i2c_scl <= 1;
                    state <= 6;
                end
                
                // Fast Write Command + Upper 4 bits of Data
                6: begin 
                    i2c_scl <= 0;
                    shift_reg <= {4'b0000, psi_sc[11:8]};
                    bit_count <= 7;
                    state <= 7; 
                end
                
                // Write Bit Phase 1 (Upper byte)
                7: begin 
                    sda_out_reg <= shift_reg[bit_count];
                    i2c_scl <= 1; 
                    state <= 8;
                end
                
                // Write Bit Phase 2 (Upper byte)
                8: begin 
                    i2c_scl <= 0;
                    if (bit_count == 0) state <= 9; 
                    else begin
                        bit_count <= bit_count - 1;
                        state <= 7;
                    end
                end
                
                // Wait for ACK
                9: begin 
                    sda_out_reg <= 1; 
                    i2c_scl <= 1;
                    state <= 10;
                end

                // Lower 8 bits of Data
                10: begin 
                    i2c_scl <= 0;
                    shift_reg <= psi_sc[7:0];
                    bit_count <= 7;
                    state <= 11; 
                end

                // Write Bit Phase 1 (Lower byte)
                11: begin 
                    sda_out_reg <= shift_reg[bit_count];
                    i2c_scl <= 1; 
                    state <= 12;
                end
                
                // Write Bit Phase 2 (Lower byte)
                12: begin 
                    i2c_scl <= 0;
                    if (bit_count == 0) state <= 13; 
                    else begin
                        bit_count <= bit_count - 1;
                        state <= 11;
                    end
                end
                
                // Wait for ACK
                13: begin 
                    sda_out_reg <= 1; 
                    i2c_scl <= 1;
                    state <= 14;
                end
                
                // STOP Condition: SDA goes high while SCL is high
                14: begin 
                    i2c_scl <= 1;
                    state <= 15;
                end
                15: begin
                    sda_out_reg <= 1;
                    state <= 0; 
                end
                
                default: state <= 0;
            endcase
        end
    end
endmodule
