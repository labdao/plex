package ipwl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/labdao/plex/internal/ipfs"
)

type ToolInput struct {
	Type    string   `json:"type"`
	Glob    []string `json:"glob"`
	Default string   `json:"default"`
}

type ToolOutput struct {
	Type string   `json:"type"`
	Item string   `json:"item"`
	Glob []string `json:"glob"`
}

type Tool struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
  Author      string                `json:"author"`
	BaseCommand []string              `json:"baseCommand"`
	Arguments   []string              `json:"arguments"`
	DockerPull  string                `json:"dockerPull"`
	GpuBool     bool                  `json:"gpuBool"`
	MemoryGB    *int                  `json:"memoryGB"`
	NetworkBool bool                  `json:"networkBool"`
	Inputs      map[string]ToolInput  `json:"inputs"`
	Outputs     map[string]ToolOutput `json:"outputs"`
}

func ReadToolConfig(toolPath string) (Tool, ToolInfo, error) {
	var tool Tool
	var toolInfo ToolInfo
	var toolFilePath string
	var cid string
	var err error

	// Check if toolPath is a key in CORE_TOOLS
	if cid, ok := CORE_TOOLS[toolPath]; ok {
		toolPath = cid
		toolInfo.IPFS = toolPath
	}

	if ipfs.IsValidCID(toolPath) {
		toolInfo.IPFS = toolPath
		toolFilePath, err = ipfs.DownloadToTempDir(toolPath)
		if err != nil {
			return tool, toolInfo, err
		}

		fileInfo, err := os.Stat(toolFilePath)
		if err != nil {
			return tool, toolInfo, err
		}

		// If the downloaded content is a directory, search for a .json file in it
		if fileInfo.IsDir() {
			files, err := ioutil.ReadDir(toolFilePath)
			if err != nil {
				return tool, toolInfo, err
			}

			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".json") {
					toolFilePath = path.Join(toolFilePath, file.Name())
					break
				}
			}
		}
	} else {
		if _, err := os.Stat(toolPath); err == nil {
			toolFilePath = toolPath
			cid, err = ipfs.WrapAndPinFile(toolFilePath)
			if err != nil {
				return tool, toolInfo, fmt.Errorf("failed to pin tool file")
			}
			toolInfo.IPFS = cid
		} else {
			return tool, toolInfo, fmt.Errorf("tool not found")
		}
	}

	file, err := os.Open(toolFilePath)
	if err != nil {
		return tool, toolInfo, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return tool, toolInfo, fmt.Errorf("failed to read file: %w", err)
	}

	err = json.Unmarshal(bytes, &tool)
	if err != nil {
		return tool, toolInfo, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	toolInfo.Name = tool.Name

	return tool, toolInfo, nil
}

func toolToCmd(toolConfig Tool, ioEntry IO, ioGraph []IO) (string, error) {
	arguments := strings.Join(toolConfig.Arguments, " ")

	placeholderRegex := regexp.MustCompile(`\$\((inputs\..+?(\.filepath|\.basename|\.ext|\.default))\)`)
	fileMatches := placeholderRegex.FindAllStringSubmatch(arguments, -1)

	for _, match := range fileMatches {
		placeholder := match[0]
		key := strings.TrimSuffix(strings.TrimPrefix(match[1], "inputs."), ".filepath")
		key = strings.TrimSuffix(key, ".basename")
		key = strings.TrimSuffix(key, ".ext")
		key = strings.TrimSuffix(key, ".default")

		var replacement string
		input := ioEntry.Inputs[key]
		switch match[2] {
		case ".filepath":
			replacement = fmt.Sprintf("/inputs/%s", input.FilePath)
		case ".basename":
			replacement = strings.TrimSuffix(input.FilePath, filepath.Ext(input.FilePath))
		case ".ext":
			ext := filepath.Ext(input.FilePath)
			replacement = strings.TrimPrefix(ext, ".")
		case ".default":
			replacement = toolConfig.Inputs[key].Default
		}

		arguments = strings.Replace(arguments, placeholder, replacement, -1)
	}

	nonFilePlaceholderRegex := regexp.MustCompile(`\$\((inputs\..+?)\)`)
	nonFileMatches := nonFilePlaceholderRegex.FindAllStringSubmatch(arguments, -1)

	for _, match := range nonFileMatches {
		placeholder := match[0]
		key := strings.TrimPrefix(match[1], "inputs.")

		if input, ok := toolConfig.Inputs[key]; ok && input.Type != "File" {
			arguments = strings.Replace(arguments, placeholder, fmt.Sprintf("%v", input.Default), -1)
		}
	}

	cmd := fmt.Sprintf("%s \"%s\"", strings.Join(toolConfig.BaseCommand, " "), arguments)

	return cmd, nil
}

// You can use custom tools by passing the cid directly to plex -t arguments
var CORE_TOOLS = map[string]string{
	"equibind":             "QmZ2HarAgwZGjc3LBx9mWNwAQkPWiHMignqKup1ckp8NhB",
	"diffdock":             "QmSzetFkveiQYZ5FgpZdHHfsjMWYz5YzwMAvqUgUFhFPMM",
	"colabfold-mini":       "QmcRH74qfqDBJFku3mEDGxkAf6CSpaHTpdbe1pMkHnbcZD",
	"colabfold-standard":   "QmXnM1VpdGgX5huyU3zTjJovsu42KPfWhjxhZGkyvy9PVk",
	"colabfold-large":      "QmPYqMy19VFFuYztL6b5ruo4Kw4JWT583emStGrSYTH5Yi",
	"bam2fastq":            "QmbPUirWiWCv9sgdHLekf5AnoCdw4QPU2SyfGGKs9JRRbq",
	"oddt":                 "QmUx7NdxkXXZvbK1JXZVUYUBqsevWkbVxgTzpWJ4Xp4inf",
	"rfdiffusion":          "QmXnCBCtoYuPyGsEJVpjn5regHfFSYa8kx44e22XxDX2t2",
	"repeatmodeler":        "QmZdXxnUt1sFFR39CfkEUgiioUBf6qP5CUs8TCb7Wqn4MC",
	"gnina":                "QmYfGaWzxwi8HiWLdiX4iQXuuLXVKYrr6YC3DknEvZeSne",
	"batch-dlkcat":         "QmThdvypN8gDDwwyNnpSYsdwvyxCET8s1jym3HZCTaBzmD",
	"openbabel-pdb-to-sdf": "QmbbDSDZJp8G7EFaNKsT7Qe7S9iaaemZmyvS6XgZpdR5e3",
	"openbabel-rmsd":       "QmUxrKgAs5r42xVki4vtMskJa1Z7WA64wURkwywPMch7dA",
}
