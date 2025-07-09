# Go PDF Service - Automated Testing Guide

This document provides comprehensive information about the automated testing suite for the Go PDF microservice.

## Overview

The testing suite includes:
- **Unit Tests** - Individual component testing
- **Integration Tests** - Service integration with mock and real backends  
- **End-to-End Tests** - Complete workflow testing
- **Performance Tests** - Load and stress testing
- **Benchmarks** - Performance benchmarking

## Quick Start

### Run All Basic Tests
```bash
./run_tests.sh
```

### Run All Tests (Including Performance)
```bash
./run_tests.sh --all
```

### Run with Real Backend
```bash
./run_tests.sh --real-backend
```

## Test Files

| File | Purpose |
|------|---------|
| `test_helpers.go` | Test utilities and mock server setup |
| `integration_test.go` | Integration tests with authentication and PDF generation |
| `e2e_test.go` | End-to-end workflow tests |
| `performance_test.go` | Performance and load tests |
| `run_tests.sh` | Test runner script with various options |

## Test Types

### 1. Unit Tests

Located in individual package directories (`internal/api/`, `internal/client/`, `internal/pdf/`).

**Run unit tests only:**
```bash
./run_tests.sh --unit-only
```

**Manual execution:**
```bash
go test -v ./internal/api/...
go test -v ./internal/client/...
go test -v ./internal/pdf/...
```

### 2. Integration Tests

Tests service integration with mock Node.js backend.

**Run integration tests only:**
```bash
./run_tests.sh --integration-only
```

**Test scenarios:**
- Health check with/without authentication
- PDF generation with various authentication methods
- Error handling for invalid requests
- Different authentication mechanisms (cookies, headers)

### 3. End-to-End Tests

Tests complete user workflows.

**Run E2E tests only:**
```bash
./run_tests.sh --e2e-only
```

**Test scenarios:**
- Complete PDF generation workflow
- Error recovery scenarios
- Concurrent request handling
- Real-world usage patterns

### 4. Performance Tests

Load and performance testing.

**Run performance tests:**
```bash
./run_tests.sh --with-performance
```

**Test scenarios:**
- Light load (5 users, 50 requests)
- Medium load (10 users, 100 requests)  
- Heavy load (20 users, 200 requests)
- Memory usage under sustained load
- Response time consistency

### 5. Benchmarks

Performance benchmarking using Go's built-in benchmark framework.

**Run benchmarks:**
```bash
./run_tests.sh --with-benchmarks
```

**Available benchmarks:**
- PDF generation performance
- Health check performance

## Test Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `RUN_REAL_BACKEND_TESTS` | Enable tests with real Node.js backend | `false` |
| `SAVE_TEST_PDFS` | Save generated PDFs for manual verification | `false` |
| `AUTH_MODE` | Set to `test` for test mode | Set automatically |
| `NODEJS_API_URL` | Node.js backend URL | `http://localhost:5007` |

### Command Line Options

| Option | Description |
|--------|-------------|
| `--unit-only` | Run only unit tests |
| `--integration-only` | Run only integration tests |
| `--e2e-only` | Run only end-to-end tests |
| `--with-performance` | Include performance tests |
| `--with-benchmarks` | Include benchmark tests |
| `--all` | Run all test types |
| `--real-backend` | Include real backend tests |
| `--help` | Show help message |

## Test Scenarios

### Scenario 1: Health Check Testing

Tests service health endpoint:
- With authentication (should be healthy)
- Without authentication (should be unhealthy)
- Response format validation

### Scenario 2: Authentication Testing

Tests various authentication methods:
- Cookie-based authentication
- Authorization header authentication  
- Custom header authentication
- Mixed authentication methods
- Invalid/missing tokens

### Scenario 3: PDF Generation Testing

Tests PDF report generation:
- Successful generation with valid student ID
- Error handling for invalid student ID
- Authentication requirement validation
- PDF content validation

### Scenario 4: Error Handling Testing

Tests error scenarios:
- Student not found (404)
- Missing student ID (400)
- Authentication failures (401/500)
- Invalid request formats

### Scenario 5: Concurrent Testing

Tests system under load:
- Multiple simultaneous requests
- Performance metrics collection
- Error rate monitoring
- Response time consistency

## Mock Backend

The test suite includes a comprehensive mock Node.js server that simulates:

### Student Data
- Student ID 2: Alice Johnson (detailed test data)
- Student ID 1: Generic test student  
- Student ID 999: Returns 404 (student not found)
- Other IDs: Return generic test data

### Authentication
- Validates access tokens in cookies/headers
- Requires CSRF tokens for security
- Returns appropriate error codes for missing auth

### Endpoints
- `GET /api/v1/students/{id}` - Student data retrieval
- `GET /api/v1/dashboard` - Health check endpoint

## Real Backend Testing

### Prerequisites

1. Node.js backend running on `http://localhost:5007`
2. PostgreSQL database with seeded student data
3. Valid admin credentials

### Configuration

Enable real backend tests:
```bash
export RUN_REAL_BACKEND_TESTS=true
./run_tests.sh
```

Or use the flag:
```bash
./run_tests.sh --real-backend
```

### Authentication

Real backend tests automatically:
1. Login with admin credentials (`admin@school-admin.com`)
2. Extract fresh authentication tokens
3. Use tokens for subsequent requests

## Performance Metrics

The performance tests collect and validate:

### Response Time Metrics
- Average response time
- Minimum/maximum response times
- Response time consistency
- Standard deviation

### Throughput Metrics  
- Requests per second
- Concurrent request handling
- Success/failure rates

### Load Testing Thresholds
- **Light Load:** 5 users, 50 requests, max 10s
- **Medium Load:** 10 users, 100 requests, max 20s  
- **Heavy Load:** 20 users, 200 requests, max 45s
- **Failure Rate:** Max 5% allowed
- **Response Time:** Max 2s average

## Test Output

### Files Generated

All test outputs are saved to `test_output/` directory:

| File | Description |
|------|-------------|
| `test_results.log` | Detailed test execution logs |
| `test_report.md` | Summary report with statistics |
| `*.pdf` | Generated PDF files (if `SAVE_TEST_PDFS=true`) |

### Report Format

The test report includes:
- Test configuration summary
- Pass/fail statistics  
- Performance metrics
- Generated file listings

## Troubleshooting

### Common Issues

**Tests fail with "connection refused":**
- Ensure Go service builds successfully
- Check if ports 8080 is available
- Verify mock server starts correctly

**Real backend tests fail:**
- Ensure Node.js backend is running on port 5007
- Verify database is seeded with test data
- Check admin credentials are correct

**Performance tests timeout:**
- Increase timeout values in test runner
- Reduce load test parameters
- Check system resources

**PDF validation fails:**
- Verify gofpdf library is installed
- Check student data completeness
- Review PDF generation logic

### Debug Mode

For detailed debugging:
```bash
go test -v -timeout=120s .
```

View detailed logs:
```bash
tail -f test_output/test_results.log
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Go PDF Service Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v3
      with:
        go-version: 1.19
    
    - name: Run Tests
      run: |
        cd go-service
        ./run_tests.sh --all
        
    - name: Upload Test Results
      uses: actions/upload-artifact@v3
      with:
        name: test-results
        path: go-service/test_output/
```

### Docker Testing

```dockerfile
FROM golang:1.19-alpine
WORKDIR /app
COPY go-service/ .
RUN go mod download
RUN chmod +x run_tests.sh
CMD ["./run_tests.sh", "--all"]
```

## Best Practices

### Test Development
1. Write tests for new features before implementation
2. Use descriptive test names and documentation
3. Include both positive and negative test cases
4. Mock external dependencies appropriately

### Test Execution
1. Run tests frequently during development  
2. Include performance tests in CI/CD pipeline
3. Monitor test execution times
4. Review test reports regularly

### Test Maintenance
1. Update tests when API changes
2. Review and update performance thresholds
3. Clean up test artifacts regularly
4. Document test changes in commit messages

## Contributing

When adding new tests:

1. Follow existing test patterns and naming conventions
2. Add tests to appropriate test files
3. Update this documentation
4. Ensure tests pass in both mock and real backend modes
5. Include performance considerations for new features

## Examples

### Running Specific Test Patterns

```bash
# Run only health check tests
go test -v -run TestHealthCheck .

# Run only PDF generation tests  
go test -v -run TestStudentReportGeneration .

# Run only authentication tests
go test -v -run TestAuthenticationMethods .

# Run performance tests with custom timeout
go test -v -timeout=300s -run TestPerformanceUnderLoad .
```

### Custom Test Configuration

```bash
# Test with custom backend URL
export NODEJS_API_URL=http://localhost:3000
./run_tests.sh

# Save test PDFs for manual review
export SAVE_TEST_PDFS=true
./run_tests.sh --with-performance

# Run with real backend and save outputs
RUN_REAL_BACKEND_TESTS=true SAVE_TEST_PDFS=true ./run_tests.sh --all
```

This comprehensive testing suite ensures the Go PDF service is robust, performant, and reliable across various scenarios and load conditions. 