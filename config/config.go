package config

import (
	"os"

	"github.com/openimsdk/open-im-server/tools/codescan/config"
	"gopkg.in/yaml.v2"
)

type Config struct {
	BaseConfig struct {
		SearchDirectory string `yaml:"searchDirectory"`
		IgnoreCase      bool   `yaml:"ignoreCase"`
	} `yaml:"baseConfig"`
	DirectoryNaming struct {
		AllowHyphens     bool `yaml:"allowHyphens"`
		AllowUnderscores bool `yaml:"allowUnderscores"`
		MustBeLowercase  bool `yaml:"mustBeLowercase"`
	} `yaml:"directoryNaming"`
	FileNaming struct {
		AllowHyphens     bool `yaml:"allowHyphens"`
		AllowUnderscores bool `yaml:"allowUnderscores"`
		MustBeLowercase  bool `yaml:"mustBeLowercase"`
	} `yaml:"fileNaming"`
	IgnoreFormats          []string                          `yaml:"ignoreFormats"`
	IgnoreDirectories      []string                          `yaml:"ignoreDirectories"`
	FileTypeSpecificNaming map[string]FileTypeSpecificNaming `yaml:"fileTypeSpecificNaming"`
}

type FileTypeSpecificNaming struct {
	AllowHyphens     bool `yaml:"allowHyphens"`
	AllowUnderscores bool `yaml:"allowUnderscores"`
	MustBeLowercase  bool `yaml:"mustBeLowercase"`
}

type Issue struct {
	Type    string
	Path    string
	Message string
}

type Checker struct {
	Config  *config.Config
	Summary struct {
		CheckedDirectories int
		CheckedFiles       int
		Issues             []Issue
	}
	Errors []string
}

func LoadConfig(configPath string) (*Config, error) {
	var config Config

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
