package ipwl

import (
	"fmt"
	"log"
	"strings"

	"github.com/labdao/plex/internal/bacalhau"
)

func SubmitIoList(ioList []IO, selector string, maxTime int, annotations []string) []IO {
	submittedIOList := make([]IO, len(ioList))
	for i, ioEntry := range ioList {
		log.Printf("Submitting IO entry %d \n", i)
		submittedIOList[i] = ioEntry
		log.Println("Reading tool config")
		toolConfig, _, err := ReadToolConfig(ioEntry.Tool.IPFS)
		if err != nil {
			submittedIOList[i].State = "failed"
			submittedIOList[i].ErrMsg = fmt.Sprintf("error reading tool config: %v", err)
			continue
		}
		log.Println("Creating cmd")
		cmd, err := toolToCmd(toolConfig, ioEntry)
		if err != nil {
			submittedIOList[i].State = "failed"
			submittedIOList[i].ErrMsg = fmt.Sprintf("error reading tool config: %v", err)
			continue
		}
		log.Printf("cmd: %s \n", cmd)
		log.Println("mapping inputs")
		fileInputs := make(map[string]string)
		fileArrayInputs := make(map[string][]string)

		for key, input := range ioEntry.Inputs {
			fmt.Println("going thru inputs with new code")
			fmt.Println(input)
			switch v := input.(type) {
			case string:
				if strings.HasPrefix(v, "Qm") {
					fmt.Println("found file input")
					fmt.Println(v)
					fileInputs[key] = v
				} else {
					fmt.Println("input is a string but does not have 'Qm' prefix")
				}
			case []interface{}: // Changed from []string to []interface{}
				fmt.Println("found slice, checking each for 'Qm' prefix")
				var stringArray []string
				allValid := true
				for _, elem := range v {
					str, ok := elem.(string)
					if !ok || !strings.HasPrefix(str, "Qm") {
						allValid = false
						fmt.Println("element is not a string or does not have 'Qm' prefix:", elem)
						break
					}
					stringArray = append(stringArray, str)
				}
				if allValid && len(stringArray) > 0 {
					fmt.Println("found file array")
					fmt.Println(stringArray)
					fileArrayInputs[key] = stringArray
				}
			default:
				fmt.Println("input is neither a string nor a slice of strings")
			}
		}

		log.Println("creating bacalhau job")
		// this memory type conversion is for backwards compatibility with the -app flag
		var memory int
		if toolConfig.MemoryGB == nil {
			memory = 0
		} else {
			memory = *toolConfig.MemoryGB
		}

		var cpu float64
		if toolConfig.Cpu == nil {
			cpu = 0
		} else {
			cpu = *toolConfig.Cpu
		}

		log.Println("creating bacalhau job")
		bacalhauJob, err := bacalhau.CreateBacalhauJob(fileInputs, fileArrayInputs, toolConfig.DockerPull, selector, cmd, maxTime, memory, cpu, toolConfig.GpuBool, toolConfig.NetworkBool, annotations)
		if err != nil {
			submittedIOList[i].State = "failed"
			submittedIOList[i].ErrMsg = fmt.Sprintf("error creating Bacalhau job: %v", err)
			continue
		}

		log.Println("submitting bacalhau job")
		submittedJob, err := bacalhau.SubmitBacalhauJob(bacalhauJob)
		if err != nil {
			submittedIOList[i].State = "failed"
			submittedIOList[i].ErrMsg = fmt.Sprintf("error submitting Bacalhau job: %v", err)
			continue
		}
		submittedIOList[i].State = "new"
		submittedIOList[i].BacalhauJobId = submittedJob.Metadata.ID
	}
	log.Println("returning io submited list")
	return submittedIOList
}
