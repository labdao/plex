package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/labdao/plex/internal/ipfs"
	"github.com/labdao/plex/internal/ipwl"
	"github.com/spf13/cobra"
)

var (
	inputs           string
	scatteringMethod string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initilizes an IO JSON from Tool config and inputs",
	Long:  `This command initlizes an IO JSON based on the provided Tool and inputs.`,
	Run: func(cmd *cobra.Command, args []string) {
		var kwargs map[string][]string
		err := json.Unmarshal([]byte(inputs), &kwargs)
		if err != nil {
			log.Fatal("Invalid inputs JSON:", err)
		}

		ioJson, err := InitilizeIo(toolPath, scatteringMethod, kwargs)
		if err != nil {
			log.Fatal(err)
		}

		// Convert the ioJson to bytes
		data, err := json.MarshalIndent(ioJson, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal the ioJson: %v", err)
		}

		// Create a temp file
		tempFile, err := ioutil.TempFile("", "io.json")
		if err != nil {
			log.Fatalf("Failed to create temporary file: %v", err)
		}

		// Write the data to the temp file
		_, err = tempFile.Write(data)
		if err != nil {
			log.Fatalf("Failed to write to temporary file: %v", err)
		}

		cid, err := ipfs.PinFile(tempFile.Name())
		if err != nil {
			log.Fatalf("Failed to pin the file: %v", err)
		}

		// Used by Python SDK do not change
		fmt.Println("Pinned IO JSON CID:", cid)

		// Remember to close the file
		if err := tempFile.Close(); err != nil {
			log.Fatalf("Failed to close the temporary file: %v", err)
		}
	},
}

func InitilizeIo(toolPath string, scatteringMethod string, inputVectors map[string][]string) ([]ipwl.IO, error) {
	// Open the file and load its content
	tool, toolInfo, err := ipwl.ReadToolConfig(toolPath)
	if err != nil {
		return nil, err
	}

	// Check if all kwargs are in the tool's inputs
	for inputKey := range inputVectors {
		if _, exists := tool.Inputs[inputKey]; !exists {
			log.Printf("The argument %s is not in the tool inputs.\n", inputKey)
			log.Printf("Available keys: %v\n", tool.Inputs)
			return nil, fmt.Errorf("the argument %s is not in the tool inputs", inputKey)
		}
	}

	// Handle scattering methods and create the inputsList
	var inputsList [][]string
	switch scatteringMethod {
	case "dotProduct":
		// check if all lists have the same length
		var vectorLength int
		for _, v := range inputVectors {
			if vectorLength == 0 {
				vectorLength = len(v)
				continue
			}
			if len(v) != vectorLength {
				return nil, fmt.Errorf("all input arguments must have the same length for dot_product scattering method")
			}
			vectorLength = len(v)
		}
		for i := 0; i < vectorLength; i++ {
			tmp := []string{}
			for _, v := range inputVectors {
				tmp = append(tmp, v[i])
			}
			inputsList = append(inputsList, tmp)
		}
	case "crossProduct":
		cartesian := func(arrs ...[]string) [][]string {
			result := [][]string{{}}

			for _, arr := range arrs {
				var temp [][]string

				for _, res := range result {
					for _, str := range arr {
						product := append([]string{}, res...)
						product = append(product, str)
						temp = append(temp, product)
					}
				}

				result = temp
			}

			return result
		}

		keys := make([]string, 0, len(inputVectors))
		for k := range inputVectors {
			keys = append(keys, k)
		}
		arrays := make([][]string, len(inputVectors))
		for i, k := range keys {
			arrays[i] = inputVectors[k]
		}
		inputsList = cartesian(arrays...)
	default:
		return nil, fmt.Errorf("invalid scattering method: %s", scatteringMethod)
	}

	// populate ioJSONGraph based on inputsList
	var ioJSONGraph []ipwl.IO
	for _, inputs := range inputsList {
		io := ipwl.IO{
			Tool:    toolInfo,
			Inputs:  make(map[string]ipwl.FileInput),
			Outputs: make(map[string]ipwl.Output),
			State:   "created",
			ErrMsg:  "",
		}

		inputKeys := make([]string, 0, len(inputVectors))
		for k := range inputVectors {
			inputKeys = append(inputKeys, k)
		}

		for i, inputValue := range inputs {
			inputKey := inputKeys[i]

			if strings.Count(inputValue, "/") == 1 {
				parts := strings.Split(inputValue, "/")
				cid := parts[0]
				fileName := parts[1]
				if !ipfs.IsValidCID(cid) {
					return nil, fmt.Errorf("invalid CID: %s", cid)
				}
				io.Inputs[inputKey] = ipwl.FileInput{
					Class:    tool.Inputs[inputKey].Type,
					FilePath: fileName,
					IPFS:     cid,
				}
			} else {
				cid, err := ipfs.WrapAndPinFile(inputValue) // Pin the file and get the CID
				if err != nil {
					return nil, err
				}
				io.Inputs[inputKey] = ipwl.FileInput{
					Class:    tool.Inputs[inputKey].Type,
					FilePath: filepath.Base(inputValue), // Use the respective input value from inputsList
					IPFS:     cid,                       // Use the CID returned by WrapAndPinFile
				}
			}
		}

		for outputKey := range tool.Outputs {
			io.Outputs[outputKey] = ipwl.FileOutput{
				Class:    tool.Outputs[outputKey].Type,
				FilePath: "", // Assuming filepath is empty, adapt as needed
				IPFS:     "", // Assuming IPFS is not provided, adapt as needed
			}
		}
		ioJSONGraph = append(ioJSONGraph, io)
	}

	return ioJSONGraph, nil
}

func init() {
	initCmd.Flags().StringVarP(&toolPath, "toolPath", "t", "", "Path of the Tool config (can be a local path or an IPFS CID)")
	initCmd.Flags().StringVarP(&inputs, "inputs", "i", "{}", "Inputs in JSON format")
	initCmd.Flags().StringVarP(&scatteringMethod, "scatteringMethod", "", "{}", "Inputs in JSON format")

	rootCmd.AddCommand(initCmd)
}
