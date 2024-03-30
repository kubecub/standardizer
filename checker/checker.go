package checker

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
	err := filepath.Walk(c.Config.BaseConfig.SearchDirectory, c.checkPath)
	if err != nil {
		return err
	}

	if len(c.Summary.Issues) > 0 {
		c.printIssues()
		return fmt.Errorf("found %d issues in the codebase", len(c.Summary.Issues))
	}
	return nil
}

func (c *Checker) printIssues() {
	fmt.Println("## Issues found:")
	fmt.Println("===================================================================================================")
	for _, issue := range c.Summary.Issues {
		fmt.Printf("Type: %s, Path: %s, Message: %s\n", issue.Type, issue.Path, issue.Message)
	}
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
		if c.isIgnoredDirectory(relativePath) {
			return filepath.SkipDir
		}
		c.Summary.CheckedDirectories++
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
		if strings.HasSuffix(path, ignoredDir) || strings.Contains(path, ignoredDir+"/") {
			return true
		}
	}
	return false
}

func (c *Checker) isIgnoredFile(path string) bool {
	for _, format := range c.Config.IgnoreFormats {
		matched, err := regexp.MatchString(format, path)
		if err != nil {
			fmt.Printf("Invalid regex pattern: %s, error: %s\n", format, err)
			continue
		}
		if matched {
			return true
		}
	}
	return false
}

func (c *Checker) checkDirectoryName(path string) bool {
	dirName := filepath.Base(path)
	errors := []string{}

	if c.Config.DirectoryNaming.MustBeLowercase && (dirName != strings.ToLower(dirName)) {
		errors = append(errors, "directory name must be lowercase")
	}
	if !c.Config.DirectoryNaming.AllowHyphens && strings.Contains(dirName, "-") {
		errors = append(errors, "hyphens are not allowed in directory names")
	}
	if !c.Config.DirectoryNaming.AllowUnderscores && strings.Contains(dirName, "_") {
		errors = append(errors, "underscores are not allowed in directory names")
	}

	if len(errors) > 0 {
		c.Summary.Issues = append(c.Summary.Issues, Issue{
			Type:    "directoryNaming",
			Path:    path,
			Message: fmt.Sprintf("Directory naming issues: %s. Example of valid directory name: '%s'", strings.Join(errors, "; "), c.exampleDirectoryName()),
		})
		return false
	}

	return true
}

func (c *Checker) checkFileName(path string) bool {
	fileName := filepath.Base(path)
	ext := filepath.Ext(fileName)
	errors := []string{}

	allowHyphens := c.Config.FileNaming.AllowHyphens
	allowUnderscores := c.Config.FileNaming.AllowUnderscores
	mustBeLowercase := c.Config.FileNaming.MustBeLowercase

	if specificNaming, ok := c.Config.FileTypeSpecificNaming[ext]; ok {
		allowHyphens = specificNaming.AllowHyphens
		allowUnderscores = specificNaming.AllowUnderscores
		mustBeLowercase = specificNaming.MustBeLowercase
	}

	if mustBeLowercase && (fileName != strings.ToLower(fileName)) {
		errors = append(errors, "file name must be lowercase")
	}
	if !allowHyphens && strings.Contains(fileName, "-") {
		errors = append(errors, "hyphens are not allowed in file names")
	}
	if !allowUnderscores && strings.Contains(fileName, "_") {
		errors = append(errors, "underscores are not allowed in file names")
	}

	if len(errors) > 0 {
		c.Summary.Issues = append(c.Summary.Issues, Issue{
			Type:    "fileNaming",
			Path:    path,
			Message: fmt.Sprintf("File naming issues: %s. Example of valid file name: '%s'", strings.Join(errors, "; "), c.exampleFileName(ext)),
		})
		return false
	}

	return true
}

func (c *Checker) exampleDirectoryName() string {
	exampleName := "ExampleDirectory"
	if c.Config.DirectoryNaming.MustBeLowercase {
		exampleName = strings.ToLower(exampleName)
	}
	if !c.Config.DirectoryNaming.AllowHyphens && !c.Config.DirectoryNaming.AllowUnderscores {
		exampleName = strings.ReplaceAll(exampleName, "-", "")
		exampleName = strings.ReplaceAll(exampleName, "_", "")
	} else if c.Config.DirectoryNaming.AllowHyphens {
		exampleName = strings.ReplaceAll(exampleName, "_", "-")
	} else if c.Config.DirectoryNaming.AllowUnderscores {
		exampleName = strings.ReplaceAll(exampleName, "-", "_")
	}
	return exampleName
}

func (c *Checker) exampleFileName(ext string) string {
	baseName := "example_file"
	if c.Config.FileNaming.MustBeLowercase {
		baseName = strings.ToLower(baseName)
	}

	if !c.Config.FileNaming.AllowHyphens && !c.Config.FileNaming.AllowUnderscores {
		baseName = strings.ReplaceAll(baseName, "-", "")
		baseName = strings.ReplaceAll(baseName, "_", "")
	} else if c.Config.FileNaming.AllowHyphens {
	} else if c.Config.FileNaming.AllowUnderscores {
		baseName = strings.ReplaceAll(baseName, "-", "_")
	}

	if specificNaming, ok := c.Config.FileTypeSpecificNaming[ext]; ok {
		if specificNaming.MustBeLowercase {
			baseName = strings.ToLower(baseName)
		}
		if !specificNaming.AllowHyphens && !specificNaming.AllowUnderscores {
			baseName = strings.ReplaceAll(baseName, "-", "")
			baseName = strings.ReplaceAll(baseName, "_", "")
		} else if specificNaming.AllowHyphens {
		} else if specificNaming.AllowUnderscores {
			baseName = strings.ReplaceAll(baseName, "-", "_")
		}
	}

	return fmt.Sprintf("%s%s", baseName, ext)
}
