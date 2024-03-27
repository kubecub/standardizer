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
		configPath = os.Getenv("CONFIG_PATH")
	}

	if configPath == "" {
		configPath = "config.yaml"
		if _, err := os.Stat(".github/composite.yaml"); err == nil {
			configPath = ".github/composite.yaml"
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
		fmt.Println("Error during check:", err)
		os.Exit(1)
	}

	// if len(c.Errors) > 0 {
	// 	fmt.Println("Found errors:")
	// 	for _, errMsg := range c.Errors {
	// 		fmt.Println("-", errMsg)
	// 	}
	// 	os.Exit(1)
	// }

	summaryJSON, err := json.MarshalIndent(c.Summary, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling summary:", err)
		return
	}

	fmt.Println(string(summaryJSON))
}
