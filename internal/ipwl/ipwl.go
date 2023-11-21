package ipwl

import (
	"fmt"
	"log"

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
		bacalhauJob, err := bacalhau.CreateBacalhauJob(ioEntry.Inputs, toolConfig.DockerPull, selector, maxTime, memory, cpu, toolConfig.GpuBool, toolConfig.NetworkBool, annotations)
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
