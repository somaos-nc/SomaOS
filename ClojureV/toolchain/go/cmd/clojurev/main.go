package main

import (
	"clojurev"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	inputPath := flag.String("in", "", "Path to the .cljv file")
	outputPath := flag.String("out", "", "Path to the output file")
	target := flag.String("target", "go", "Target language (go, verilog, javascript, python, wasm, dart)")
	pkgName := flag.String("pkg", "main", "Go package name")
	flag.Parse()

	if *inputPath == "" {
		log.Fatal("Input path is required (-in)")
	}

	cljvData, err := ioutil.ReadFile(*inputPath)
	if err != nil {
		log.Fatalf("Failed to read input file: %v", err)
	}

	output, err := clojurev.Transpile(string(cljvData), clojurev.Target(*target), *pkgName)
	if err != nil {
		log.Fatalf("Transpilation fracture: %v", err)
	}

	if *outputPath != "" {
		err = ioutil.WriteFile(*outputPath, []byte(output), 0644)
		if err != nil {
			log.Fatalf("Failed to write output file: %v", err)
		}
		fmt.Printf("Successfully manifested %s logic in %s (package: %s)\n", *target, *outputPath, *pkgName)
	} else {
		fmt.Println(output)
	}
}
