package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/labdao/plex/internal/ipfs"
	"github.com/labdao/plex/internal/ipwl"
	"github.com/spf13/cobra"
)

var (
	ioPath  string
	toolCid string
)

var vectorizeCmd = &cobra.Command{
	Use:   "vectorize",
	Short: "Transform an IO JSON file into a list of outputs",
	Long:  `Transform an IO JSON file into a list of outputs.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := VectorizeOutputs(ioPath, toolCid, outputDir)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func VectorizeOutputs(ioPath string, toolCid string, outputDir string) (map[string]ipwl.OutputValues, error) {
	isCID := ipfs.IsValidCID(ioPath)
	id := uuid.New()
	workDirPath := ""

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if outputDir != "" {
		absPath, err := filepath.Abs(outputDir)
		if err != nil {
			return nil, err
		}
		cwd = absPath
		workDirPath = cwd
		err = os.MkdirAll(workDirPath, 0755)
		if err != nil {
			return nil, err
		}
	}

	var ioJsonFilPath string
	if isCID {
		workDirPath = path.Join(cwd, id.String())
		err = os.Mkdir(workDirPath, 0755)
		if err != nil {
			return nil, err
		}
		ioJsonFilPath = path.Join(workDirPath, "io.json")
		err = ipfs.DownloadFileContents(ioPath, ioJsonFilPath)
		if err != nil {
			return nil, err
		}
	} else {
		ioJsonFilPath, err = filepath.Abs(ioPath)
		if err != nil {
			return nil, err
		}
		if workDirPath == "" {
			workDirPath = filepath.Dir(ioJsonFilPath)
		}
	}

	file, err := os.Open(ioJsonFilPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var ios []ipwl.IO
	err = json.Unmarshal(bytes, &ios)
	if err != nil {
		return nil, err
	}

	outputMap := make(map[string]ipwl.OutputValues)
	for i, io := range ios {
		if io.Tool.IPFS == toolCid {
			for key, output := range io.Outputs {
				fileOutput, ok := output.(ipwl.FileOutput)
				if ok {
					ov := outputMap[key]

					filePath := fmt.Sprintf("entry-%d/outputs/%s", i, fileOutput.FilePath)
					absoluteFilePath := path.Join(workDirPath, filePath)

					// Check if the file is already downloaded
					if _, err := os.Stat(absoluteFilePath); os.IsNotExist(err) {
						// Download the file from IPFS to the local file path
						err = ipfs.DownloadToDirectory(fileOutput.IPFS, filepath.Dir(absoluteFilePath))
						if err != nil {
							return nil, err
						}
					}

					ov.FilePaths = append(ov.FilePaths, absoluteFilePath)
					ov.CIDs = append(ov.CIDs, fileOutput.IPFS)
					outputMap[key] = ov
				}
			}
		}
	}

	// Save the output map to a JSON file
	outputVectorsPath := path.Join(workDirPath, "output-vectors.json")
	outputVectorsFile, err := os.Create(outputVectorsPath)
	if err != nil {
		return nil, err
	}
	defer outputVectorsFile.Close()

	jsonData, err := json.MarshalIndent(outputMap, "", "  ")
	if err != nil {
		return nil, err
	}
	outputVectorsFile.Write(jsonData)

	fmt.Println("Output Vectors Were saved at:", outputVectorsPath)

	return outputMap, nil
}

func init() {
	vectorizeCmd.Flags().StringVarP(&ioPath, "ioPath", "i", "", "CID or file path of IO JsON")
	vectorizeCmd.Flags().StringVarP(&toolCid, "toolCid", "t", "", "Only vectorize output CIDs")
	vectorizeCmd.Flags().StringVarP(&outputDir, "outputDir", "o", "", "Only vectorize output CIDs")

	rootCmd.AddCommand(vectorizeCmd)
}
