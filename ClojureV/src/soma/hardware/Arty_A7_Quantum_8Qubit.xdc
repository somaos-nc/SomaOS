## SOMA OS: Physical Constraint Mapping (8-Qubit 3D Edition)
## Target: Arty A7-35T FPGA

# 1. System Clock
set_property -dict { PACKAGE_PIN E3    IOSTANDARD LVCMOS33 } [get_ports { CLK100MHZ }];
create_clock -add -name sys_clk_pin -period 10.00 -waveform {0 5} [get_ports { CLK100MHZ }];

# 2. Phase Tuning Injector (PMOD JA - I2C to DAC)
set_property -dict { PACKAGE_PIN G13   IOSTANDARD LVCMOS33 } [get_ports { ja_sda }]; 
set_property -dict { PACKAGE_PIN B11   IOSTANDARD LVCMOS33 } [get_ports { ja_scl }]; 

# 3. 8-Cell Macro-Cube Outputs (PMOD JB & JC)
# Using PMOD JB (Top/Bottom) and JC (Top/Bottom) for full 8-bit observation
set_property -dict { PACKAGE_PIN D15   IOSTANDARD LVCMOS33 } [get_ports { jb_c0_out }]; 
set_property -dict { PACKAGE_PIN C15   IOSTANDARD LVCMOS33 } [get_ports { jb_c1_out }]; 
set_property -dict { PACKAGE_PIN J17   IOSTANDARD LVCMOS33 } [get_ports { jb_c2_out }]; 
set_property -dict { PACKAGE_PIN J18   IOSTANDARD LVCMOS33 } [get_ports { jb_c3_out }]; 

set_property -dict { PACKAGE_PIN K15   IOSTANDARD LVCMOS33 } [get_ports { jb_c4_out }]; 
set_property -dict { PACKAGE_PIN J15   IOSTANDARD LVCMOS33 } [get_ports { jb_c5_out }]; 
set_property -dict { PACKAGE_PIN K16   IOSTANDARD LVCMOS33 } [get_ports { jb_c6_out }]; 
set_property -dict { PACKAGE_PIN J16   IOSTANDARD LVCMOS33 } [get_ports { jb_c7_out }]; 

# 4. Thermodynamic Observer (Internal XADC Auxiliary)
set_property -dict { PACKAGE_PIN C5    IOSTANDARD LVCMOS33 } [get_ports { vauxp0 }]; 
set_property -dict { PACKAGE_PIN A5    IOSTANDARD LVCMOS33 } [get_ports { vauxn0 }]; 

# 5. Critical Synthesis Directives
# Allow combinatorial feedback loops for the geometric knots
set_property ALLOW_COMBINATORIAL_LOOPS TRUE [get_nets -of_objects [get_cells *Cell/*]]
