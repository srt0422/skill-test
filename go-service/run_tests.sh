#!/bin/bash

# Go Service Test Runner
# This script runs various types of automated tests for the Go PDF service

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test configuration
GO_SERVICE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_OUTPUT_DIR="$GO_SERVICE_DIR/test_output"
LOG_FILE="$TEST_OUTPUT_DIR/test_results.log"

# Create test output directory
mkdir -p "$TEST_OUTPUT_DIR"

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to run a test command
run_test() {
    local test_name="$1"
    local test_cmd="$2"
    local description="$3"
    
    print_status "Running $test_name: $description"
    echo "==================== $test_name ====================" >> "$LOG_FILE"
    echo "Description: $description" >> "$LOG_FILE"
    echo "Command: $test_cmd" >> "$LOG_FILE"
    echo "Started at: $(date)" >> "$LOG_FILE"
    echo "" >> "$LOG_FILE"
    
    if eval "$test_cmd" >> "$LOG_FILE" 2>&1; then
        print_success "$test_name passed"
        echo "Result: PASSED" >> "$LOG_FILE"
    else
        print_error "$test_name failed"
        echo "Result: FAILED" >> "$LOG_FILE"
        return 1
    fi
    
    echo "" >> "$LOG_FILE"
    echo "" >> "$LOG_FILE"
}

# Function to check dependencies
check_dependencies() {
    print_status "Checking dependencies..."
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    # Check Go version
    GO_VERSION=$(go version | grep -o 'go[0-9]*\.[0-9]*')
    print_status "Using Go version: $GO_VERSION"
    
    print_success "Dependencies check passed"
}

# Function to build the service
build_service() {
    print_status "Building Go service..."
    
    cd "$GO_SERVICE_DIR"
    
    if go build -o bin/go-service cmd/main.go; then
        print_success "Go service built successfully"
    else
        print_error "Failed to build Go service"
        exit 1
    fi
}

# Function to run unit tests
run_unit_tests() {
    print_status "Running unit tests..."
    
    cd "$GO_SERVICE_DIR"
    
    # Run tests for each package
    run_test "API Unit Tests" \
        "go test -v ./internal/api/..." \
        "Unit tests for API handlers and authentication"
    
    run_test "Client Unit Tests" \
        "go test -v ./internal/client/..." \
        "Unit tests for Node.js client"
    
    run_test "PDF Generator Unit Tests" \
        "go test -v ./internal/pdf/..." \
        "Unit tests for PDF generation"
}

# Function to run integration tests
run_integration_tests() {
    print_status "Running integration tests..."
    
    cd "$GO_SERVICE_DIR"
    
    run_test "Integration Tests" \
        "go test -v -timeout=60s -run='TestHealthCheck|TestStudentReportGeneration|TestAuthenticationMethods|TestErrorHandling' ." \
        "Integration tests with mock backend"
    
    # Run real backend tests if available
    if [ "$RUN_REAL_BACKEND_TESTS" = "true" ]; then
        print_warning "Running tests with real backend (requires Node.js backend to be running)"
        run_test "Real Backend Integration Tests" \
            "go test -v -timeout=120s -run='TestWithRealBackend' ." \
            "Integration tests with real Node.js backend"
    else
        print_warning "Skipping real backend tests (set RUN_REAL_BACKEND_TESTS=true to enable)"
    fi
}

# Function to run end-to-end tests
run_e2e_tests() {
    print_status "Running end-to-end tests..."
    
    cd "$GO_SERVICE_DIR"
    
    run_test "E2E PDF Generation Workflow" \
        "go test -v -timeout=90s -run='TestEndToEndPDFGenerationWorkflow' ." \
        "Complete PDF generation workflow"
    
    run_test "E2E Error Recovery" \
        "go test -v -timeout=60s -run='TestEndToEndErrorRecovery' ." \
        "Error handling and recovery workflow"
    
    run_test "E2E Concurrent Requests" \
        "go test -v -timeout=120s -run='TestEndToEndConcurrentRequests' ." \
        "Concurrent request handling"
    
    run_test "E2E Real World Scenario" \
        "go test -v -timeout=90s -run='TestEndToEndRealWorldScenario' ." \
        "Realistic usage scenarios"
    
    # Run real backend E2E tests if available
    if [ "$RUN_REAL_BACKEND_TESTS" = "true" ]; then
        run_test "E2E Real Backend Tests" \
            "go test -v -timeout=120s -run='TestEndToEndWithRealBackend' ." \
            "End-to-end tests with real backend"
    fi
}

# Function to run performance tests
run_performance_tests() {
    print_status "Running performance tests..."
    
    cd "$GO_SERVICE_DIR"
    
    # Set environment variable to save test PDFs for manual verification
    export SAVE_TEST_PDFS=true
    
    run_test "Performance Under Load" \
        "go test -v -timeout=300s -run='TestPerformanceUnderLoad' ." \
        "Performance testing under various load conditions"
    
    run_test "Memory Usage Under Load" \
        "go test -v -timeout=120s -run='TestMemoryUsageUnderLoad' ." \
        "Memory usage during sustained load"
    
    run_test "Response Time Consistency" \
        "go test -v -timeout=120s -run='TestResponseTimeConsistency' ." \
        "Response time consistency testing"
    
    unset SAVE_TEST_PDFS
}

# Function to run benchmarks
run_benchmarks() {
    print_status "Running benchmarks..."
    
    cd "$GO_SERVICE_DIR"
    
    run_test "PDF Generation Benchmark" \
        "go test -bench=BenchmarkPDFGeneration -benchtime=10s -timeout=120s ." \
        "Benchmark for PDF generation performance"
    
    run_test "Health Check Benchmark" \
        "go test -bench=BenchmarkHealthCheck -benchtime=10s -timeout=60s ." \
        "Benchmark for health check performance"
}

# Function to generate test report
generate_report() {
    print_status "Generating test report..."
    
    local report_file="$TEST_OUTPUT_DIR/test_report.md"
    
    cat > "$report_file" << EOF
# Go PDF Service Test Report

**Generated:** $(date)  
**Go Version:** $(go version)  
**Test Directory:** $GO_SERVICE_DIR  

## Test Configuration

- **Mock Backend:** Used for most tests
- **Real Backend Tests:** ${RUN_REAL_BACKEND_TESTS:-false}
- **Performance Tests:** Included
- **Benchmarks:** Included

## Test Results Summary

EOF
    
    # Count passed/failed tests from log
    local total_tests=$(grep -c "Result: " "$LOG_FILE" || echo "0")
    local passed_tests=$(grep -c "Result: PASSED" "$LOG_FILE" || echo "0")
    local failed_tests=$(grep -c "Result: FAILED" "$LOG_FILE" || echo "0")
    
    cat >> "$report_file" << EOF
- **Total Tests:** $total_tests
- **Passed:** $passed_tests
- **Failed:** $failed_tests
- **Success Rate:** $(echo "scale=2; $passed_tests * 100 / $total_tests" | bc 2>/dev/null || echo "N/A")%

## Detailed Results

See \`test_results.log\` for detailed test output.

## Test Files Generated

EOF
    
    # List any generated test files
    if [ -d "$TEST_OUTPUT_DIR" ]; then
        echo "- Test logs and outputs in: \`$TEST_OUTPUT_DIR\`" >> "$report_file"
        if ls "$TEST_OUTPUT_DIR"/*.pdf >/dev/null 2>&1; then
            echo "- Generated PDF files:" >> "$report_file"
            for pdf in "$TEST_OUTPUT_DIR"/*.pdf; do
                echo "  - \`$(basename "$pdf")\`" >> "$report_file"
            done
        fi
    fi
    
    print_success "Test report generated: $report_file"
}

# Function to cleanup
cleanup() {
    print_status "Cleaning up..."
    
    # Remove any temporary files if needed
    # (Currently nothing to clean up)
    
    print_success "Cleanup completed"
}

# Main function
main() {
    echo "=================================================="
    echo "Go PDF Service Automated Test Suite"
    echo "=================================================="
    echo ""
    
    # Initialize log file
    echo "Go PDF Service Test Results" > "$LOG_FILE"
    echo "Generated: $(date)" >> "$LOG_FILE"
    echo "Go Version: $(go version)" >> "$LOG_FILE"
    echo "" >> "$LOG_FILE"
    
    # Parse command line arguments
    local run_unit=true
    local run_integration=true
    local run_e2e=true
    local run_performance=false
    local run_benchmarks=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --unit-only)
                run_integration=false
                run_e2e=false
                shift
                ;;
            --integration-only)
                run_unit=false
                run_e2e=false
                shift
                ;;
            --e2e-only)
                run_unit=false
                run_integration=false
                shift
                ;;
            --with-performance)
                run_performance=true
                shift
                ;;
            --with-benchmarks)
                run_benchmarks=true
                shift
                ;;
            --all)
                run_performance=true
                run_benchmarks=true
                shift
                ;;
            --real-backend)
                export RUN_REAL_BACKEND_TESTS=true
                shift
                ;;
            --help)
                echo "Usage: $0 [options]"
                echo ""
                echo "Options:"
                echo "  --unit-only         Run only unit tests"
                echo "  --integration-only  Run only integration tests"
                echo "  --e2e-only         Run only end-to-end tests"
                echo "  --with-performance Run performance tests"
                echo "  --with-benchmarks  Run benchmark tests"
                echo "  --all              Run all test types including performance and benchmarks"
                echo "  --real-backend     Include tests that require real Node.js backend"
                echo "  --help             Show this help message"
                echo ""
                echo "Environment variables:"
                echo "  RUN_REAL_BACKEND_TESTS=true  Enable real backend tests"
                echo "  SAVE_TEST_PDFS=true         Save generated PDFs for verification"
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                echo "Use --help for usage information"
                exit 1
                ;;
        esac
    done
    
    # Run tests
    check_dependencies
    build_service
    
    local failed_tests=0
    
    if [ "$run_unit" = true ]; then
        run_unit_tests || ((failed_tests++))
    fi
    
    if [ "$run_integration" = true ]; then
        run_integration_tests || ((failed_tests++))
    fi
    
    if [ "$run_e2e" = true ]; then
        run_e2e_tests || ((failed_tests++))
    fi
    
    if [ "$run_performance" = true ]; then
        run_performance_tests || ((failed_tests++))
    fi
    
    if [ "$run_benchmarks" = true ]; then
        run_benchmarks || ((failed_tests++))
    fi
    
    # Generate report
    generate_report
    cleanup
    
    # Final status
    echo ""
    echo "=================================================="
    if [ $failed_tests -eq 0 ]; then
        print_success "All tests completed successfully!"
        echo "Test report: $TEST_OUTPUT_DIR/test_report.md"
        echo "Detailed logs: $LOG_FILE"
    else
        print_error "$failed_tests test suite(s) failed"
        echo "Check logs for details: $LOG_FILE"
        exit 1
    fi
    echo "=================================================="
}

# Run main function with all arguments
main "$@" 