package spincontainers

import (
	"log"
	"os/exec"
)

func SpinContainer() {
	cmd := exec.Command("docker-compose", "up")

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to spin the containers: %v", err)
	}
}
