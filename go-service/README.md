# Go PDF Report Service

A microservice built in Go that generates PDF reports for students by consuming the Node.js backend API.

## Overview

This service implements **Problem 4: Golang Developer Challenge** from the skill test requirements:
- Creates a standalone microservice in Go
- Consumes the Node.js backend `/api/v1/students/:id` endpoint 
- Generates downloadable PDF reports for students
- **Does NOT connect directly to the database** - only through Node.js API

## API Endpoints

### Student Report Generation
```
GET /api/v1/students/{id}/report
```
Generates and returns a PDF report for the specified student ID.

### Health Check
```
GET /health
```
Returns service health status.

## Project Structure

```
go-service/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── api/                 # HTTP handlers and routing
│   ├── client/              # Node.js API client
│   └── pdf/                 # PDF generation logic
├── pkg/
│   └── models/              # Data models
├── bin/                     # Compiled binaries
├── go.mod                   # Go module definition
└── README.md               # This file
```

## Prerequisites

- Go 1.21+
- Node.js backend running on http://localhost:5007
- PostgreSQL database seeded with student data

## Development

### Build
```bash
go build -o bin/go-service ./cmd/main.go
```

### Run
```bash
# Default port 8080
./bin/go-service

# Or specify port
PORT=8081 ./bin/go-service
```

### Test
```bash
# Health check
curl http://localhost:8080/health

# Generate student report (once implemented)
curl http://localhost:8080/api/v1/students/2/report
```

## Dependencies

- `github.com/gorilla/mux` - HTTP router
- Additional PDF generation libraries will be added

## Implementation Status

- [x] Project structure and basic setup
- [ ] Node.js API client implementation
- [ ] PDF generation functionality  
- [ ] Student report endpoint implementation
- [ ] Authentication handling
- [ ] Integration testing

## Next Steps

1. Implement HTTP client for Node.js API consumption
2. Add PDF generation capabilities
3. Complete the student report endpoint
4. Add authentication handling
5. Perform end-to-end testing 