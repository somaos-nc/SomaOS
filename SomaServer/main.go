package main

import (
	"log"
	"soma_server/api"
	"soma_server/hardware"
)

func main() {
	log.Println("Starting SomaOS FPGA Interface Server...")
	
	driver := hardware.NewFPGADriver()
	server := api.NewServer(driver)
	
	addr := ":8081"
	log.Printf("Listening on http://localhost%s\n", addr)
	if err := server.Start(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
