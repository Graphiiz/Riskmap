# riskmap-backtend

# Getting Started with Golang Backend

# Steps
## Start Database
- go inside db/ directory
- run command:
    docker compose up

## Start Backend
- run command for loading dependencies:
    go mod tidy
- run command:
    go run server.go
    
Backend will serve frontend via port 1234
