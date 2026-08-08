package main

import (
	"fmt"
	"time"

	"github.com/OpsNexusHQ/opsnexus-common/models"
)

func main() {
	agent := models.Agent{
		ID:       "local-agent",
		Name:     "OpsNexus Local Agent",
		Hostname: "localhost",
		OS:       "linux",
		Arch:     "amd64",
		Version:  "0.1.0",
		Status:   "online",
		LastSeen: time.Now(),
	}

	fmt.Printf("OpsNexus Agent\n")
	fmt.Printf("ID: %s\n", agent.ID)
	fmt.Printf("Hostname: %s\n", agent.Hostname)
	fmt.Printf("Status: %s\n", agent.Status)
}
