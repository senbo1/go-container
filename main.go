package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error starting container: %v\n", err)
		os.Exit(1)
	}
}

func runContainerProcess() {
	fmt.Printf("Inside container: %v\n", os.Args[2:])

	setupCgroups()

	if err := syscall.Sethostname([]byte("my-container")); err != nil {
		fmt.Printf("Error setting hostname: %v\n", err)
		os.Exit(1)
	}

	if err := syscall.Chroot("/home/user/containerfs"); err != nil {
		fmt.Printf("Error changing root: %v\n", err)
		os.Exit(1)
	}

	if err := os.Chdir("/"); err != nil {
		fmt.Printf("Error changing directory: %v\n", err)
		os.Exit(1)
	}

	if err := syscall.Mount("proc", "proc", "proc", 0, ""); err != nil {
		fmt.Printf("Error mounting proc: %v\n", err)
		os.Exit(1)
	}

	if err := syscall.Mount("tmpfs", "tmp", "tmpfs", 0, ""); err != nil {
		fmt.Printf("Error mounting tmpfs: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running container command: %v\n", err)
	}

	if err := syscall.Unmount("tmp", 0); err != nil {
		fmt.Printf("Error unmounting tmpfs: %v\n", err)
	}

	if err := syscall.Unmount("proc", 0); err != nil {
		fmt.Printf("Error unmounting proc: %v\n", err)
	}
}

func setupCgroups() {
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")
	containerGroup := filepath.Join(pids, "mycontainer")

	if err := os.Mkdir(containerGroup, 0755); err != nil && !os.IsExist(err) {
		fmt.Printf("Error creating cgroup: %v\n", err)
		os.Exit(1)
	}

	writeFile := func(path, content string) {
		if err := os.WriteFile(path, []byte(content), 0700); err != nil {
			fmt.Printf("Error writing to %s: %v\n", path, err)
			os.Exit(1)
		}
	}

	writeFile(filepath.Join(containerGroup, "pids.max"), "20")
	writeFile(filepath.Join(containerGroup, "notify_on_release"), "1")
	writeFile(filepath.Join(containerGroup, "cgroup.procs"), strconv.Itoa(os.Getpid()))
}
