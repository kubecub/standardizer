package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/kubecub/standardizer/checker"
	"github.com/kubecub/standardizer/config"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to the configuration file")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("STANDARDIZER_CONFIG_PATH")
	}

	if configPath == "" {
		configPath = "standardizer.yaml"
		if _, err := os.Stat(".github/standardizer.yaml"); err == nil {
			configPath = ".github/standardizer.yaml"
		} else if _, err := os.Stat(".github/standardizer.yml"); err == nil {
			configPath = ".github/standardizer.yml"
		}
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	c := &checker.Checker{Config: cfg}
	err = c.Check()
	if err != nil {
		fmt.Println()
		fmt.Println("===================================================================================================")
		fmt.Println("Please check whether the above file conforms to the specification, or check whether the configuration file is qualified")
		fmt.Println("!!!Error during check:", err)
		os.Exit(1)
	}

	summaryJSON, err := json.MarshalIndent(c.Summary, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling summary:", err)
		return
	}

	fmt.Println(string(summaryJSON))
}
