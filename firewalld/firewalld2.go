package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

// logInfo prints informational messages
func logInfo(message string) {
	log.Printf("[INFO] %s\n", message)
}

// logWarning prints warning messages
func logWarning(message string) {
	log.Printf("[WARNING] %s\n", message)
}

// logError prints error messages and exits
func logError(message string) {
	log.Fatalf("[ERROR] %s\n", message)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	logInfo("Starting firewalld port checker...")

	// Check if firewall-cmd exists
	_, err := exec.LookPath("firewall-cmd")
	if err != nil {
		logError("firewall-cmd not found. Is firewalld installed and running?")
	}

	logInfo("Found firewall-cmd binary.")

	// Check if firewalld service is running
	checkStatus := exec.Command("firewall-cmd", "--state")
	var statusOut bytes.Buffer
	var statusErr bytes.Buffer
	checkStatus.Stdout = &statusOut
	checkStatus.Stderr = &statusErr

	err = checkStatus.Run()
	if err != nil {
		logError(fmt.Sprintf("Failed to query firewalld status: %v - %s", err, statusErr.String()))
	}

	state := strings.TrimSpace(statusOut.String())
	if state != "running" {
		logWarning(fmt.Sprintf("firewalld is not running (state: %s)", state))
	} else {
		logInfo("firewalld is active and running.")
	}

	// Pause briefly
	time.Sleep(500 * time.Millisecond)

	// Now list the open ports
	logInfo("Querying open ports...")
	cmd := exec.Command("firewall-cmd", "--list-ports")

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		logError(fmt.Sprintf("Failed to list open ports: %v - %s", err, stderr.String()))
	}

	portsRaw := out.String()
	ports := strings.Fields(portsRaw)

	fmt.Println("\n========== FIREWALLD OPEN PORTS ==========")
	if len(ports) == 0 {
		fmt.Println("No open ports found in firewalld.")
		logWarning("No open ports detected.")
	} else {
		for idx, port := range ports {
			fmt.Printf("%d. %s\n", idx+1, port)
		}
		logInfo(fmt.Sprintf("Total open ports found: %d", len(ports)))
	}
	fmt.Println("===========================================")

	logInfo("Firewalld port check completed successfully.")
}
