package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <command> [args...]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		runContainer()
	case "child":
		runContainerProcess()
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}

func runContainer() {
	fmt.Printf("Starting container process: %v\n", os.Args[2:])
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error starting container: %v\n", err)
		os.Exit(1)
	}
}

func runContainerProcess() {
	fmt.Printf("Inside container: %v\n", os.Args[2:])
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := syscall.Sethostname([]byte("my-container")); err != nil {
		fmt.Printf("Error setting hostname: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running container command: %v\n", err)
		os.Exit(1)
	}
}
