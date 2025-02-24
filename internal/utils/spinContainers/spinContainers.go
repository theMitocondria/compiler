package spincontainers

import (
	"log"
	"os"
	"os/exec"
)

func SpinContainer() {
	cmd := exec.Command("docker-compose", "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to spin the containers: %v", err)
	}
}
func DownContainer() {
	cmd := exec.Command("docker-compose", "down")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to down the containers: %v", err)
	}
}
