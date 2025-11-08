package main

import (
	"flag"
	"fmt"
	"github.com/stlesnik/goph_keeper/internal/config"
	"os"

	"github.com/stlesnik/goph_keeper/internal/client"
)

var (
	version = "1.0.0"
	build   = "dev"
	date    = "27-10-2025"
)

func main() {
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showVersion, "v", false, "Show version information")
	flag.Parse()

	if showVersion {
		fmt.Printf("GophKeeper Client v%s\n", version)
		fmt.Printf("Build: %s\n", build)
		fmt.Printf("Date: %s\n", date)
		os.Exit(0)
	}

	cfg, err := config.LoadClientConfig()
	if err != nil {
		fmt.Printf("Error loading client config: %s\n", err)
		os.Exit(1)
	}

	app := client.NewApp(cfg)
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
