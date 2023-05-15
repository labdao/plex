// package docker

// import (
// 	"bufio"
// 	"fmt"
// 	"os/exec"
// 	"path/filepath"
// )

// func InstructionToDockerCmd(container, cmd, jobDir string, gpu bool) string {
// 	outputsDir := filepath.Join(jobDir, "outputs")
// 	fmt.Println(outputsDir)
// 	gpuFlag := ""
// 	if gpu {
// 		gpuFlag = "--gpus"
// 	}
// 	return fmt.Sprintf("docker run %s -v %s:/inputs -v %s:/outputs %s /bin/bash -c '%s'", gpuFlag, jobDir, outputsDir, container, cmd)
// }

// func RunDockerJob(container, cmd, jobDir string, gpu bool) error {
// 	fmt.Println(InstructionToDockerCmd(container, cmd, jobDir, gpu))
// 	dockerCmd := InstructionToDockerCmd(container, cmd, jobDir, gpu)
// 	cmdExec := exec.Command("/bin/sh", "-c", dockerCmd)

// 	// Set up a pipe to read the command's stdout and stderr
// 	stdout, err := cmdExec.StdoutPipe()
// 	if err != nil {
// 		return err
// 	}
// 	stderr, err := cmdExec.StderrPipe()
// 	if err != nil {
// 		return err
// 	}

// 	// Start the command
// 	if err := cmdExec.Start(); err != nil {
// 		fmt.Println("Error starting Docker command:", err)
// 		return err
// 	}

// 	// Use a scanner to read the output line by line
// 	go func() {
// 		scanner := bufio.NewScanner(stdout)
// 		for scanner.Scan() {
// 			fmt.Println(scanner.Text())
// 		}
// 	}()

// 	// Use a scanner to read the error output line by line
// 	go func() {
// 		scanner := bufio.NewScanner(stderr)
// 		for scanner.Scan() {
// 			fmt.Println(scanner.Text())
// 		}
// 	}()

// 	// Wait for the command to finish
// 	if err := cmdExec.Wait(); err != nil {
// 		return err
// 	}
// 	return nil
// }
