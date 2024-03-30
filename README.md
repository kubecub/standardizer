# README for Standardizer

## Overview

**Standardizer** is a versatile GitHub Action tool designed to ensure your project's file and directory names adhere to specified naming conventions. It automates the process of checking the conformity of naming standards within your repository, saving time and maintaining consistency across your codebase. Standardizer can be easily integrated into your GitHub workflows or used locally to inspect your project's structure for any deviations from the defined naming rules.

## Features

- **Flexible Configuration**: Customize naming rules for files and directories through a straightforward YAML configuration.
- **Regex Support for Ignoring Formats**: Utilize regular expressions to define exceptions for file formats that do not require standardization.
- **Environment Variables & Flags Support**: Specify the configuration file's location using environment variables or command-line flags for versatile setup options.
- **Local & CI Integration**: Designed for both local development environments and continuous integration workflows, ensuring consistency across different stages of development.

## Getting Started

### Installation

For local use, you can either download the latest release from [GitHub Releases](https://github.com/kubecub/standardizer/releases/) or install directly using Go:

```shell
go install github.com/kubecub/standardizer@latest
```

### Configuration

Create a `standardizer.yml` configuration file in your project's root directory or under `.github/`. Alternatively, specify the configuration file's path through environment variables or flags.

**Example `standardizer.yml` Configuration:**

```yaml
baseConfig:
  searchDirectory: "./"
  ignoreCase: false

directoryNaming:
  allowHyphens: true
  allowUnderscores: false
  mustBeLowercase: true

fileNaming:
  allowHyphens: true
  allowUnderscores: true
  mustBeLowercase: true

ignoreFormats:
  - "\\.log$"
  - "\\.env$"
  - "README\\.md$"
  - "_test\\.go$"
  - "\\.md$"
  - LICENSE

ignoreDirectories:
  - "vendor"
  - ".git"
  - "node_modules"
  - "logs"
  - "CHANGELOG"
  - "components"
  - "_output"
  - "tools/openim-web"
  - "CHANGELOG"
  - "examples/Test_directory"

fileTypeSpecificNaming:
  ".yaml":
    allowHyphens: true
    allowUnderscores: false
    mustBeLowercase: true
  ".go":
    allowHyphens: false
    allowUnderscores: true
    mustBeLowercase: true
```

### Configuration Details

- **baseConfig**: General settings for the scan, such as the directory to search (`searchDirectory`) and whether to ignore case sensitivity (`ignoreCase`).
- **directoryNaming** & **fileNaming**: Define the allowed naming conventions for directories and files, respectively, including the use of hyphens, underscores, and case sensitivity.
- **ignoreFormats**: A list of regex patterns for file formats that should be ignored during the naming checks. This is useful for excluding specific file types or naming patterns.
- **ignoreDirectories**: Directories to exclude from the scan, ensuring that vendor or third-party directories won't affect your compliance results.
- **fileTypeSpecificNaming**: Allows for specifying naming conventions on a per-file-extension basis, offering greater flexibility for projects that use a variety of file types with different standards.

### Using Standardizer in GitHub Actions

To integrate Standardizer into your GitHub Actions workflow, add the following step to your `.github/workflows/main.yml`:

```yaml
- name: Conformity Checker for Project
  uses: kubecub/standardizer@main # or use a specific tag like @v1.0.0
```

### Local Usage

After configuring `standardizer.yml`, run Standardizer locally to check your project's naming conventions:

```shell
standardizer -config .github/standardizer.yml
```

### Success Output

A successful run will output the number of directories and files checked, along with any identified issues:

```json
{
  "CheckedDirectories": 4,
  "CheckedFiles": 15,
  "Issues": null
}
```

- **CheckedDirectories** & **CheckedFiles**: The count of directories and files that were checked.
- **Issues**: Lists any naming convention violations found. `null` indicates no issues were detected.

## Contributing

Contributions to Standardizer are welcome! Please refer to the project's GitHub page for contribution guidelines: [https://github.com/kubecub/standardizer](https://github.com/kubecub/standardizer).

## License

Standardizer is released under the MIT License. See the LICENSE file for more details.