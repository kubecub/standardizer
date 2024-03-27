package checker

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kubecub/standardizer/config"
)

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
}

func (c *Checker) Check() error {
	return filepath.Walk(c.Config.BaseConfig.SearchDirectory, c.checkPath)
}

func (c *Checker) checkPath(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	relativePath, err := filepath.Rel(c.Config.BaseConfig.SearchDirectory, path)
	if err != nil {
		return err
	}

	if relativePath == "." {
		return nil
	}

	if info.IsDir() {
		c.Summary.CheckedDirectories++
		if c.isIgnoredDirectory(relativePath) {
			c.Summary.Issues = append(c.Summary.Issues, Issue{
				Type:    "ignoredDirectory",
				Path:    path,
				Message: "This directory has been ignored",
			})
			return filepath.SkipDir
		}
		if !c.checkDirectoryName(relativePath) {
			c.Summary.Issues = append(c.Summary.Issues, Issue{
				Type:    "directoryNaming",
				Path:    path,
				Message: "The directory name is invalid",
			})
		}
	} else {
		if c.isIgnoredFile(path) {
			return nil
		}
		c.Summary.CheckedFiles++
		if !c.checkFileName(relativePath) {
			c.Summary.Issues = append(c.Summary.Issues, Issue{
				Type:    "fileNaming",
				Path:    path,
				Message: "The file name does not comply with the specification",
			})
		}
	}

	return nil
}

func (c *Checker) isIgnoredDirectory(path string) bool {
	for _, ignoredDir := range c.Config.IgnoreDirectories {
		if strings.Contains(path, ignoredDir) {
			return true
		}
	}
	return false
}

func (c *Checker) isIgnoredFile(path string) bool {
	ext := filepath.Ext(path)
	for _, format := range c.Config.IgnoreFormats {
		if ext == format {
			return true
		}
	}
	return false
}

func (c *Checker) checkDirectoryName(path string) bool {
	dirName := filepath.Base(path)
	if c.Config.DirectoryNaming.MustBeLowercase && (dirName != strings.ToLower(dirName)) {
		return false
	}
	if !c.Config.DirectoryNaming.AllowHyphens && strings.Contains(dirName, "-") {
		return false
	}
	if !c.Config.DirectoryNaming.AllowUnderscores && strings.Contains(dirName, "_") {
		return false
	}
	return true
}

func (c *Checker) checkFileName(path string) bool {
	fileName := filepath.Base(path)
	ext := filepath.Ext(fileName)

	allowHyphens := c.Config.FileNaming.AllowHyphens
	allowUnderscores := c.Config.FileNaming.AllowUnderscores
	mustBeLowercase := c.Config.FileNaming.MustBeLowercase

	if specificNaming, ok := c.Config.FileTypeSpecificNaming[ext]; ok {
		allowHyphens = specificNaming.AllowHyphens
		allowUnderscores = specificNaming.AllowUnderscores
		mustBeLowercase = specificNaming.MustBeLowercase
	}

	if mustBeLowercase && (fileName != strings.ToLower(fileName)) {
		return false
	}
	if !allowHyphens && strings.Contains(fileName, "-") {
		return false
	}
	if !allowUnderscores && strings.Contains(fileName, "_") {
		return false
	}

	return true
}
