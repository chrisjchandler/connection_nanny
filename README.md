
# TCP Connection Nanny

## Overview

The **TCP Connection Nanny** is a Go-based monitoring application designed to automatically sever incoming TCP connections that exceed a specified time limit. The application continuously monitors all established TCP connections and utilizes the `tcpkill` tool to terminate connections that have been alive longer than the configured threshold. It helps prevent socket exhaustion and ensures efficient use of network resources in high-traffic environments.

## Purpose

This application is designed to manage long-lived TCP connections by monitoring their duration and automatically terminating those that exceed a predefined limit. It ensures that stale or misbehaving connections do not linger, which can lead to resource exhaustion and potential performance degradation in systems handling a large number of concurrent connections.

## Components and Requirements

### Requisite Components:
1. **Go Compiler**: You need Go to compile the nanny application. However, once compiled, you can run it without further Go dependencies.
   - [Download Go](https://golang.org/dl/) (if necessary).
  
2. **`tcpkill` from the `dsniff` Package**: This tool is responsible for forcibly terminating TCP connections by sending TCP RST packets. You need to have this package installed on the system where the nanny application is run.
   - Install `dsniff` via:
     ```bash
     sudo apt-get update
     sudo apt-get install dsniff
     ```

3. **Root Privileges**: The application requires root access because `tcpkill` needs elevated privileges to manipulate network connections.

### Supported Platforms:
- Linux-based systems with `tcpkill` available and root access.
- The system should have `ss` or `netstat` installed to gather connection details.

## Configuration

### Application Parameters:

- **Monitoring Interval** (`monitorInterval`):
  - The frequency at which the application checks for long-lived TCP connections.
  - Default: 30 seconds.
  
- **Maximum Connection Time** (`maxConnectionTime`):
  - The maximum time a connection is allowed to remain active before it is automatically severed.
  - Default: 5 minutes.

These values can be modified in the Go source code if needed:

```go
const monitorInterval = 30 * time.Second
const maxConnectionTime = 300 * time.Second

Logging
The application logs to the console (stdout) and provides details on the connections it severs, including:

The source and destination IP addresses and ports.
Any errors encountered when trying to sever a connection.
Example Log Output:

Monitoring established connections...
Killing connection: 10.0.15.73:18443 -> 174.201.120.36:64119
Connection killed successfully: 10.0.15.73:18443 -> 174.201.120.36:64119

Additional Features
Automated Monitoring: The application runs continuously in the background, requiring no user input after starting.
Customizable Connection Lifespan: Adjust the maxConnectionTime to suit your environment and traffic patterns.
Robustness: Uses tcpkill to reliably sever connections, ensuring long-lived connections do not exhaust system resources.
Notes
The tcpkill command must be installed and operational for the nanny application to function.
Ensure that the application has root access to manage network connections.
This tool is designed to mitigate issues in environments where connection mismanagement could cause system instability.

License
This project is licensed under the MIT License. See the LICENSE file for detail


