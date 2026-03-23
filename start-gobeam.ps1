Write-Host "Starting GoBeam and backend servers..."

# Start the 4 backend servers in separate windows so you can see their logs
Start-Process powershell -ArgumentList "-NoExit", "-Command", "go run servers/server1.go"
Start-Process powershell -ArgumentList "-NoExit", "-Command", "go run servers/server2.go"
Start-Process powershell -ArgumentList "-NoExit", "-Command", "go run servers/server3.go"
Start-Process powershell -ArgumentList "-NoExit", "-Command", "go run servers/server4.go"

Start-Sleep -Seconds 2

Write-Host "Backend servers started in new windows."
Write-Host "Starting the GoBeam Load Balancer on port 8080..."
Write-Host "Press Ctrl+C to stop the Load Balancer."

# Start the load balancer in this window
go run main.go
