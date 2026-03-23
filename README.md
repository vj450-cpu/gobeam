# GoBeam

GoBeam is a simplified, modern load balancer written in Go, featuring a built-in health checker, metrics collection, and a web dashboard.

## Running the Project

1. **Start the backend servers** (in separate terminals or a script):
   ```bash
   go run servers/server1.go
   go run servers/server2.go
   go run servers/server3.go
   go run servers/server4.go
   ```

2. **Start the load balancer**:
   ```bash
   go run main.go
   ```

3. **Open the Dashboard**:
   Navigate to http://localhost:9000 in your web browser.

4. **Test the Load Balancer**:
   Send requests to the load balancer:
   ```bash
   curl http://localhost:8080
   ```
   You will see responses alternating between the backend servers.

## Features
- **Round Robin & Least Connections**: Load balancing algorithms.
- **Health Checking**: Automatically detects downed servers and stops sending traffic to them.
- **Metrics Dashboard**: A real-time web UI showing server status and traffic stats.
