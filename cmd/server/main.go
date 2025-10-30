package main

import (
	"filmfolk/internal/config"
	"fmt"
	"log"
	"os"
)

func init() {
	cfg, source, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Configuration Loading Error: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Configuration loaded from: %s\n", source)
	fmt.Printf("APP Name is %s \n", cfg.App.Name)
}

func main() {

}
