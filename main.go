package main

import (
	"log"

	"github.com/lmagdanello/bluebanquise-installer/cmd"
	"github.com/lmagdanello/bluebanquise-installer/internal/utils"
)

func main() {
	// Initialize logger
	if err := utils.InitLogger(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Execute the root command
	cmd.Execute()
}
