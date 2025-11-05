# Test Coverage

This document describes the unit test coverage for the OpenReports processor.

## Test Files

- `processor_test.go`: Tests for the main processor logic
- `config_test.go`: Tests for configuration validation

## Test Coverage

### Processor Tests (`processor_test.go`)

#### Basic Functionality
- ✅ `TestProcessLogRecord_NotOpenReportsLog`: Verifies non-OpenReports logs are skipped
- ✅ `TestProcessLogRecord_OpenReportsLog_NoResults`: Verifies handling of logs without results
- ✅ `TestProcessLogRecord_OpenReportsLog_EmptyResults`: Verifies handling of empty results array
- ✅ `TestProcessLogRecord_OpenReportsLog_SingleResult`: Verifies single result processing
- ✅ `TestProcessLogRecord_OpenReportsLog_MultipleResults`: Verifies multiple results expansion

#### Status Filtering
- ✅ `TestProcessLogRecord_StatusFilter_OnlyFailures`: Verifies filtering to only "fail" status
- ✅ `TestProcessLogRecord_StatusFilter_MultipleStatuses`: Verifies filtering with multiple allowed statuses
- ✅ `TestProcessLogRecord_StatusFilter_Empty`: Verifies empty filter processes all statuses

#### Field Mapping
- ✅ `TestTransformToSecurityEvent_FieldMapping`: Comprehensive test of all field mappings
- ✅ `TestProcessLogRecord_WithSeverityAndCategory`: Verifies severity and category handling

#### Risk Assessment
- ✅ `TestMapSeverityToRiskLevel`: Tests severity to risk level mapping
- ✅ `TestCalculateRiskScore`: Tests risk score calculation
- ✅ `TestMapResultToComplianceStatus`: Tests result to compliance status mapping

#### Workload Extraction
- ✅ `TestExtractWorkloadInfo_FromOwnerReferences`: Verifies workload extraction from owner references
- ✅ `TestExtractWorkloadInfo_FromPodName`: Verifies workload name inference from pod name
- ✅ `TestExtractWorkloadInfo_StatefulSet`: Verifies StatefulSet workload extraction
- ✅ `TestIsWorkloadKind`: Tests workload kind detection
- ✅ `TestSplitPodName`: Tests pod name parsing

#### Edge Cases
- ✅ `TestProcessLogRecord_InvalidJSON`: Verifies handling of invalid JSON in results
- ✅ `TestProcessLogRecord_TimestampMapping`: Verifies timestamp mapping from results

### Configuration Tests (`config_test.go`)

#### Validation Tests
- ✅ `TestConfig_Validate_ValidStatuses`: Tests valid status filter configurations
- ✅ `TestConfig_Validate_InvalidStatuses`: Tests invalid status filter configurations

## Running Tests

```bash
# Run all tests
go test ./processor/securityevent/internal/openreports/...

# Run with verbose output
go test -v ./processor/securityevent/internal/openreports/...

# Run with coverage
go test -cover ./processor/securityevent/internal/openreports/...

# Generate coverage report
go test -coverprofile=coverage.out ./processor/securityevent/internal/openreports/...
go tool cover -html=coverage.out
```

## Test Structure

Each test follows the pattern:
1. **Setup**: Create processor and test data
2. **Execute**: Call the method under test
3. **Assert**: Verify expected behavior and outputs

Tests use:
- `github.com/stretchr/testify/assert` for assertions
- `github.com/stretchr/testify/require` for critical assertions that stop execution
- `go.uber.org/zap/zaptest` for test logging

## Coverage Areas

- ✅ OpenReports log detection
- ✅ Result parsing and expansion
- ✅ Status filtering
- ✅ Security event field mapping
- ✅ Risk level and score calculation
- ✅ Compliance status mapping
- ✅ Workload information extraction
- ✅ Configuration validation
- ✅ Error handling (invalid JSON, missing fields)
- ✅ Edge cases (empty arrays, missing data)

## Notes

- Tests use realistic OpenReports log structures based on the example provided
- All valid status values are tested: "pass", "fail", "error", "skip"
- Workload extraction tests cover both owner references and pod name inference
- Field mapping tests verify all security event schema fields are populated correctly

