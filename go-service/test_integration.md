# Integration Test Guide

This document provides step-by-step instructions for testing the complete PDF report generation flow from the Go microservice.

## Prerequisites

1. **PostgreSQL Database**: Running with seeded student data
2. **Node.js Backend**: Running on http://localhost:5007
3. **Go Service**: Built and ready to run on http://localhost:8080

## Test Scenarios

### Scenario 1: Health Check
Tests basic connectivity and API status.

```bash
# Start the Go service
./bin/go-service &

# Test health endpoint
curl http://localhost:8080/health

# Expected response (Node.js backend not authenticated):
# {"status":"unhealthy","service":"go-pdf-service","error":"Node.js API unavailable"}

# Or (Node.js backend accessible):
# {"status":"healthy","service":"go-pdf-service","nodejs_api":"connected"}
```

### Scenario 2: Authentication Token Testing
Tests authentication token handling.

```bash
# Start Go service with test tokens
AUTH_MODE=test ./bin/go-service &

# Test with valid authentication (if tokens haven't expired)
curl -o student_2_report.pdf http://localhost:8080/api/v1/students/2/report

# Test with custom tokens via headers
curl -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
     -H "X-CSRF-Token: YOUR_CSRF_TOKEN" \
     -o student_3_report.pdf \
     http://localhost:8080/api/v1/students/3/report

# Test with cookies (mimicking browser request)
curl -b "accessToken=YOUR_ACCESS_TOKEN;csrfToken=YOUR_CSRF_TOKEN" \
     -o student_4_report.pdf \
     http://localhost:8080/api/v1/students/4/report
```

### Scenario 3: End-to-End PDF Generation
Complete flow from student ID to PDF download.

```bash
# Step 1: Get fresh authentication tokens
# Login to Node.js backend first:
curl -X POST http://localhost:5007/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"admin@school-admin.com","password":"3OU4zn3q6Zh9"}' \
     -c cookies.txt

# Step 2: Extract tokens from cookies.txt and set environment variables
# (This step requires manual extraction of tokens)

# Step 3: Test with fresh tokens
curl -b cookies.txt \
     -o alice_johnson_report.pdf \
     http://localhost:8080/api/v1/students/2/report

# Step 4: Verify PDF was generated
file alice_johnson_report.pdf
# Expected: alice_johnson_report.pdf: PDF document

# Step 5: Open PDF to verify content
open alice_johnson_report.pdf  # On macOS
# Or: xdg-open alice_johnson_report.pdf  # On Linux
```

## Error Cases to Test

### Invalid Student ID
```bash
curl http://localhost:8080/api/v1/students/999/report
# Expected: {"error":"Student not found"} with 404 status
```

### No Authentication
```bash
curl http://localhost:8080/api/v1/students/2/report
# Expected: {"error":"Failed to fetch student data"} with 500 status
```

### Expired Tokens
```bash
# Use expired tokens from login_cookies.txt
curl -b "accessToken=EXPIRED_TOKEN;csrfToken=CSRF_TOKEN" \
     http://localhost:8080/api/v1/students/2/report
# Expected: {"error":"Failed to fetch student data"} with 500 status
```

## Expected Results

### Successful PDF Generation
When authentication is valid and student exists:

1. **HTTP Response**:
   - Status: 200 OK
   - Content-Type: application/pdf
   - Content-Disposition: attachment; filename=student_2_report.pdf

2. **PDF Content** should include:
   - Student Information: ID, Name, Email, Phone, Gender, DOB
   - Academic Information: Class, Section, Roll Number, Admission Date
   - Family Information: Father, Mother, Guardian details
   - Address Information: Current and Permanent addresses
   - Generated timestamp and service attribution

3. **Service Logs** should show:
   - "Successfully generated PDF report for student 2"

## Test Data

The following students are available for testing (seeded in database):

| ID | Name | Class | Section | Roll |
|----|------|-------|---------|------|
| 1 | John Smith | Grade 10 | A | 1 |
| 2 | Alice Johnson | Grade 10 | A | 2 |
| 3 | Bob Wilson | Grade 10 | B | 1 |
| 4 | Carol Davis | Grade 10 | B | 2 |
| 5 | David Brown | Grade 11 | A | 1 |
| 6 | Emma Taylor | Grade 11 | A | 2 |
| 7 | Frank Miller | Grade 11 | B | 1 |
| 8 | Grace Anderson | Grade 11 | B | 2 |
| 9 | Henry Thompson | Grade 12 | A | 1 |
| 10 | Ivy Garcia | Grade 12 | A | 2 |

## Troubleshooting

### Authentication Issues
- Verify Node.js backend is running and accessible
- Check token expiration dates
- Ensure tokens are properly formatted and valid

### PDF Generation Issues  
- Check Go service logs for detailed error messages
- Verify gofpdf library is properly installed
- Ensure student data is complete in database

### Connection Issues
- Verify both services are running on correct ports
- Check firewall settings
- Ensure database connectivity

## Performance Testing

```bash
# Generate multiple reports quickly
for i in {1..10}; do
  curl -b cookies.txt \
       -o "student_${i}_report.pdf" \
       http://localhost:8080/api/v1/students/${i}/report &
done
wait

# Check all files were generated
ls -la *.pdf
```

## Cleanup

```bash
# Stop services
pkill -f go-service
pkill -f "node.*server.js"

# Remove test files
rm -f *.pdf cookies.txt
``` 