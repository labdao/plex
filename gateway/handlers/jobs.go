package handlers

import (
	"fmt"
)

func InitJobHandler() {
	// log that this function is being hit
	fmt.Print("InitJobHandler hit")
}

func RunJobHandler() {
	// log that this function is being hit
	fmt.Print("RunJobHandler hit")
}

func GetJobHandler() {
	// log that this function is being hit
	fmt.Print("GetJobHandler hit")
}

func GetJobsHandler() {
	// log that this function is being hit
	fmt.Print("GetJobsHandler hit")
}
