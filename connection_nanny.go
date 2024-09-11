package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

// Configuration constants
const (
	monitorInterval   = 30 * time.Second  // How often to monitor connections
	maxConnectionTime = 300 * time.Second // Maximum allowed connection time (e.g., 5 minutes)
)

// Connection structure to hold relevant connection details
type Connection struct {
	srcIP        string
	srcPort      string
	dstIP        string
	dstPort      string
	established  time.Time
	connectionID string
}

// ExecuteCommand executes a shell command and returns its output
func ExecuteCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// ParseConnections parses the output of the ss/netstat command to find established connections
func ParseConnections(output string) ([]Connection, error) {
	var connections []Connection

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()

		// Ignore the header or empty lines
		if strings.Contains(line, "ESTAB") {
			fields := strings.Fields(line)
			// Example: tcp ESTAB 0 0 10.0.15.73:18443 174.201.120.36:64119
			if len(fields) >= 4 {
				src := fields[3]
				dst := fields[4]

				srcParts := strings.Split(src, ":")
				dstParts := strings.Split(dst, ":")

				if len(srcParts) == 2 && len(dstParts) == 2 {
					srcIP, srcPort := srcParts[0], srcParts[1]
					dstIP, dstPort := dstParts[0], dstParts[1]

					// For this example, we assume established time starts when the connection was seen
					conn := Connection{
						srcIP:       srcIP,
						srcPort:     srcPort,
						dstIP:       dstIP,
						dstPort:     dstPort,
						established: time.Now(),
						connectionID: fmt.Sprintf("%s:%s -> %s:%s",
							srcIP, srcPort, dstIP, dstPort),
					}
					connections = append(connections, conn)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return connections, nil
}

// DropConnection uses tcpkill to sever a TCP connection
func DropConnection(conn Connection) error {
	fmt.Printf("Killing connection: %s:%s -> %s:%s\n", conn.srcIP, conn.srcPort, conn.dstIP, conn.dstPort)

	// Construct the tcpkill command to match the connection
	// -9 sends a RST to kill the connection immediately
	cmd := exec.Command("tcpkill", "-9", fmt.Sprintf("host %s and port %s", conn.srcIP, conn.srcPort))
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to kill connection: %s", err)
	}

	fmt.Printf("Connection killed successfully: %s\n", conn.connectionID)
	return nil
}

// MonitorConnections monitors and severs stale connections
func MonitorConnections() {
	for {
		fmt.Println("Monitoring established connections...")

		// Fetch active connections using ss (or netstat if ss is not available)
		output, err := ExecuteCommand("ss", "-tanp")
		if err != nil {
			log.Fatalf("Failed to execute ss command: %v", err)
		}

		// Parse the output to find established connections
		connections, err := ParseConnections(output)
		if err != nil {
			log.Fatalf("Failed to parse connections: %v", err)
		}

		// Check for long-lived connections and drop them
		now := time.Now()
		for _, conn := range connections {
			if now.Sub(conn.established) > maxConnectionTime {
				// Drop connections older than the maximum allowed time
				err := DropConnection(conn)
				if err != nil {
					log.Printf("Failed to drop connection: %s", err)
				}
			}
		}

		// Sleep until the next monitoring interval
		time.Sleep(monitorInterval)
	}
}

func main() {
	fmt.Println("Starting connection nanny...")
	MonitorConnections()
}
