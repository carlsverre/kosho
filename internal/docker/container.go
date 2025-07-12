package docker

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"kosho/internal/worktree"
)

func StartInteractiveShell(kw *worktree.KoshoWorktree) error {
	containerName := kw.ContainerName()

	fmt.Printf("Starting interactive shell in container for worktree '%s'\n", kw.WorktreeName)

	// Check if container exists and is running
	running, err := IsContainerRunning(containerName)
	if err != nil {
		return fmt.Errorf("failed to check container status: %w", err)
	}

	if !running {
		// Create and start the container
		if err := CreateContainer(kw); err != nil {
			return fmt.Errorf("failed to create container: %w", err)
		}
		if err := StartContainer(containerName); err != nil {
			return fmt.Errorf("failed to start container: %w", err)
		}
	}

	// Attach to the container
	cmd := exec.Command("docker", "exec", "-it", containerName, "/bin/zsh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		}
		return fmt.Errorf("failed to attach to container: %w", err)
	}

	return nil
}

func CreateContainer(kw *worktree.KoshoWorktree) error {
	containerName := kw.ContainerName()
	configVolume := kw.ConfigVolumeName()
	historyVolume := kw.HistoryVolumeName()
	worktreePath := kw.WorktreePath()

	// Ensure the kosho-runtime image exists
	if err := EnsureKoshoImage(); err != nil {
		return fmt.Errorf("failed to ensure kosho-runtime image: %w", err)
	}

	// Check if container already exists
	if exists, err := ContainerExists(containerName); err != nil {
		return fmt.Errorf("failed to check if container exists: %w", err)
	} else if exists {
		fmt.Printf("Container '%s' already exists\n", containerName)
		return nil
	}

	// Create named volumes if they don't exist
	if err := CreateNamedVolume(configVolume); err != nil {
		return fmt.Errorf("failed to create config volume: %w", err)
	}
	if err := CreateNamedVolume(historyVolume); err != nil {
		return fmt.Errorf("failed to create history volume: %w", err)
	}

	// Build the Docker run command
	args := []string{
		"run",
		"-d", // detached
		"--name", containerName,
		"--cap-add", "NET_ADMIN",
		"--cap-add", "NET_RAW",
		"-v", fmt.Sprintf("%s:/workspace", worktreePath),
		"-v", fmt.Sprintf("%s:/home/ubuntu/.claude", configVolume),
		"-v", fmt.Sprintf("%s:/commandhistory", historyVolume),
		"--workdir", "/workspace",
		"kosho-runtime",
		"sleep", "infinity", // Keep container running
	}

	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create container: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("Created container '%s'\n", containerName)
	return nil
}

func StartContainer(containerName string) error {
	// Check if container is already running
	running, err := IsContainerRunning(containerName)
	if err != nil {
		return fmt.Errorf("failed to check container status: %w", err)
	}
	if running {
		fmt.Printf("Container '%s' is already running\n", containerName)
		return nil
	}

	cmd := exec.Command("docker", "start", containerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start container: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("Started container '%s'\n", containerName)
	return nil
}

func StopContainer(containerName string) error {
	// Check if container is running
	running, err := IsContainerRunning(containerName)
	if err != nil {
		return fmt.Errorf("failed to check container status: %w", err)
	}
	if !running {
		fmt.Printf("Container '%s' is not running\n", containerName)
		return nil
	}

	cmd := exec.Command("docker", "stop", containerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop container: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("Stopped container '%s'\n", containerName)
	return nil
}

func CreateNamedVolume(volumeName string) error {
	// Check if volume already exists
	cmd := exec.Command("docker", "volume", "inspect", volumeName)
	if err := cmd.Run(); err == nil {
		// Volume already exists
		return nil
	}

	// Create the volume
	cmd = exec.Command("docker", "volume", "create", volumeName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create volume '%s': %w\nOutput: %s", volumeName, err, string(output))
	}

	fmt.Printf("Created volume '%s'\n", volumeName)
	return nil
}

func ContainerExists(containerName string) (bool, error) {
	cmd := exec.Command("docker", "inspect", containerName)
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 1 {
				// Container doesn't exist
				return false, nil
			}
		}
		return false, fmt.Errorf("failed to inspect container: %w", err)
	}
	return true, nil
}

func IsContainerRunning(containerName string) (bool, error) {
	cmd := exec.Command("docker", "inspect", "-f", "{{.State.Running}}", containerName)
	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 1 {
				// Container doesn't exist, so it's not running
				return false, nil
			}
		}
		return false, fmt.Errorf("failed to check container status: %w", err)
	}

	running := strings.TrimSpace(string(output)) == "true"
	return running, nil
}

func RemoveContainer(containerName string) error {
	// Stop the container first if it's running
	if running, err := IsContainerRunning(containerName); err != nil {
		return fmt.Errorf("failed to check container status: %w", err)
	} else if running {
		if err := StopContainer(containerName); err != nil {
			return fmt.Errorf("failed to stop container before removal: %w", err)
		}
	}

	// Remove the container
	cmd := exec.Command("docker", "rm", containerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove container: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("Removed container '%s'\n", containerName)
	return nil
}

func ListContainers() ([]string, error) {
	cmd := exec.Command("docker", "ps", "-a", "--filter", "name=kosho-", "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	var containers []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			containers = append(containers, line)
		}
	}

	return containers, nil
}

func BuildKoshoImage() error {
	// Get the directory where the Dockerfile is located
	cmd := exec.Command("docker", "build", "-t", "kosho-runtime", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Building kosho-runtime Docker image...")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to build kosho-runtime image: %w", err)
	}

	fmt.Println("Successfully built kosho-runtime image")
	return nil
}

func EnsureKoshoImage() error {
	// Check if the kosho-runtime image exists
	cmd := exec.Command("docker", "inspect", "kosho-runtime")
	if err := cmd.Run(); err == nil {
		// Image exists
		return nil
	}

	// Image doesn't exist, build it
	return BuildKoshoImage()
}
