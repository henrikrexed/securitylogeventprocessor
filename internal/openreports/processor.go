package openreports

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

// Processor handles transformation of OpenReports logs into security events
type Processor struct {
	logger *zap.Logger
	config *Config
}

// NewProcessor creates a new OpenReports processor
func NewProcessor(logger *zap.Logger, config *Config) (*Processor, error) {
	return &Processor{
		logger: logger,
		config: config,
	}, nil
}

// ProcessLogRecord processes a single log record and transforms it into multiple security events
// Returns a slice of new log records (one per result) or nil if this is not an OpenReports log
func (p *Processor) ProcessLogRecord(ctx context.Context, logRecord *plog.LogRecord, resource pcommon.Resource, scopeLogs plog.ScopeLogs) ([]plog.LogRecord, error) {
	// Check if this is an OpenReports log by looking for the kind field
	attrs := logRecord.Attributes()
	kindVal, exists := attrs.Get("kind")
	if !exists || kindVal.AsString() != "Report" {
		// Not an OpenReports log, skip
		p.logger.Debug("Log record does not match OpenReports processor - kind field check",
			zap.Bool("kind_exists", exists),
			zap.String("kind_value", func() string {
				if exists {
					return kindVal.AsString()
				}
				return "missing"
			}()),
			zap.String("trace_id", logRecord.TraceID().String()))
		return nil, nil
	}

	apiVersionVal, exists := attrs.Get("apiVersion")
	if !exists || apiVersionVal.AsString() != "openreports.io/v1alpha1" {
		// Not an OpenReports log, skip
		p.logger.Debug("Log record does not match OpenReports processor - apiVersion check",
			zap.Bool("apiVersion_exists", exists),
			zap.String("apiVersion_value", func() string {
				if exists {
					return apiVersionVal.AsString()
				}
				return "missing"
			}()),
			zap.String("trace_id", logRecord.TraceID().String()))
		return nil, nil
	}

	// Log that we've identified an OpenReports log
	p.logger.Debug("OpenReports log identified - processing",
		zap.String("trace_id", logRecord.TraceID().String()),
		zap.String("span_id", logRecord.SpanID().String()),
		zap.String("timestamp", logRecord.Timestamp().String()))

	// Extract metadata for logging
	metadataName, metadataNameExists := attrs.Get("metadata.name")
	scopeName, scopeNameExists := attrs.Get("scope.name")
	scopeKind, scopeKindExists := attrs.Get("scope.kind")

	metadataNameStr := ""
	if metadataNameExists {
		metadataNameStr = metadataName.AsString()
	}
	scopeNameStr := ""
	if scopeNameExists {
		scopeNameStr = scopeName.AsString()
	}
	scopeKindStr := ""
	if scopeKindExists {
		scopeKindStr = scopeKind.AsString()
	}

	p.logger.Debug("OpenReports log metadata",
		zap.String("metadata.name", metadataNameStr),
		zap.String("scope.name", scopeNameStr),
		zap.String("scope.kind", scopeKindStr))

	// Extract the results array
	resultsVal, exists := attrs.Get("results")
	if !exists {
		p.logger.Warn("OpenReports log has no results field",
			zap.String("metadata.name", metadataNameStr))
		return nil, nil
	}

	p.logger.Debug("Parsing OpenReports results array",
		zap.String("results_type", resultsVal.Type().String()),
		zap.Bool("results_exists", exists))

	// Parse results - it's stored as an array/slice of JSON strings
	var resultsArray []string
	if resultsVal.Type() == pcommon.ValueTypeSlice {
		// If it's a slice, extract each element as a string
		slice := resultsVal.Slice()
		p.logger.Debug("Results is a slice type",
			zap.Int("slice_length", slice.Len()))
		for i := 0; i < slice.Len(); i++ {
			resultsArray = append(resultsArray, slice.At(i).AsString())
		}
	} else if resultsVal.Type() == pcommon.ValueTypeStr {
		// If it's a single JSON string containing an array, parse it
		resultStr := resultsVal.AsString()
		p.logger.Debug("Results is a string type, attempting JSON parse",
			zap.Int("string_length", len(resultStr)))
		var jsonArray []string
		if err := json.Unmarshal([]byte(resultStr), &jsonArray); err == nil {
			resultsArray = jsonArray
			p.logger.Debug("Successfully parsed JSON array from string",
				zap.Int("array_length", len(resultsArray)))
		} else {
			// Try as single string
			p.logger.Debug("JSON parse failed, treating as single string result",
				zap.Error(err))
			resultsArray = []string{resultStr}
		}
	} else {
		p.logger.Warn("OpenReports log results field has unexpected type",
			zap.String("type", resultsVal.Type().String()),
			zap.String("metadata.name", metadataNameStr))
		return nil, nil
	}

	p.logger.Debug("Parsed results array",
		zap.Int("total_results", len(resultsArray)))

	if len(resultsArray) == 0 {
		p.logger.Debug("OpenReports log has empty results array",
			zap.String("metadata.name", metadataNameStr))
		return nil, nil
	}

	// Extract remaining metadata from the original log
	metadataNamespaceVal, metadataNamespaceExists := attrs.Get("metadata.namespace")
	scopeNamespaceVal, scopeNamespaceExists := attrs.Get("scope.namespace")
	scopeUIDVal, scopeUIDExists := attrs.Get("scope.uid")
	scopeAPIVersionVal, scopeAPIVersionExists := attrs.Get("scope.apiVersion")
	timestamp := logRecord.Timestamp()

	metadataNamespaceStr := ""
	if metadataNamespaceExists {
		metadataNamespaceStr = metadataNamespaceVal.AsString()
	}
	scopeNamespaceStr := ""
	if scopeNamespaceExists {
		scopeNamespaceStr = scopeNamespaceVal.AsString()
	}
	scopeUIDStr := ""
	if scopeUIDExists {
		scopeUIDStr = scopeUIDVal.AsString()
	}
	scopeAPIVersionStr := ""
	if scopeAPIVersionExists {
		scopeAPIVersionStr = scopeAPIVersionVal.AsString()
	}

	// Extract workload information from owner references
	p.logger.Debug("Extracting workload information",
		zap.String("scope.name", scopeNameStr),
		zap.String("scope.namespace", scopeNamespaceStr))
	workloadInfo := extractWorkloadInfo(attrs, scopeNameStr, scopeNamespaceStr)

	if workloadInfo.name != "" {
		p.logger.Debug("Workload information extracted",
			zap.String("workload.name", workloadInfo.name),
			zap.String("workload.kind", workloadInfo.kind),
			zap.String("workload.namespace", workloadInfo.namespace),
			zap.String("workload.uid", workloadInfo.uid))
	} else {
		p.logger.Debug("No workload information found - will infer from pod name if applicable")
	}

	// Log status filter configuration
	if len(p.config.StatusFilter) > 0 {
		p.logger.Debug("Status filter active",
			zap.Strings("allowed_statuses", p.config.StatusFilter),
			zap.Int("total_results", len(resultsArray)))
	} else {
		p.logger.Debug("No status filter configured - processing all results",
			zap.Int("total_results", len(resultsArray)))
	}

	// Create a new log record for each result
	var newRecords []plog.LogRecord
	processedCount := 0
	filteredCount := 0

	for i := 0; i < len(resultsArray); i++ {
		resultJSONStr := resultsArray[i]

		p.logger.Debug("Parsing result",
			zap.Int("result_index", i),
			zap.Int("result_length", len(resultJSONStr)))

		// Parse the result JSON
		var result Result
		if err := json.Unmarshal([]byte(resultJSONStr), &result); err != nil {
			p.logger.Warn("Failed to parse result JSON",
				zap.Int("result_index", i),
				zap.String("result_preview", func() string {
					if len(resultJSONStr) > 200 {
						return resultJSONStr[:200] + "..."
					}
					return resultJSONStr
				}()),
				zap.Error(err))
			continue
		}

		p.logger.Debug("Parsed result successfully",
			zap.Int("result_index", i),
			zap.String("policy", result.Policy),
			zap.String("rule", result.Rule),
			zap.String("result", result.Result),
			zap.String("source", result.Source))

		// Filter by status if configured
		if len(p.config.StatusFilter) > 0 {
			if !p.isStatusAllowed(result.Result) {
				p.logger.Debug("Skipping result due to status filter",
					zap.Int("result_index", i),
					zap.String("status", result.Result),
					zap.String("policy", result.Policy),
					zap.String("rule", result.Rule),
					zap.Strings("allowed_statuses", p.config.StatusFilter))
				filteredCount++
				continue
			}
		}

		p.logger.Debug("Result passed status filter, creating security event",
			zap.Int("result_index", i),
			zap.String("status", result.Result),
			zap.String("policy", result.Policy),
			zap.String("rule", result.Rule))

		// Create a new log record for this result
		newRecord := plog.NewLogRecord()

		// Copy basic fields from original
		newRecord.SetTimestamp(timestamp)
		newRecord.SetObservedTimestamp(logRecord.ObservedTimestamp())
		newRecord.SetSeverityNumber(logRecord.SeverityNumber())
		newRecord.SetSeverityText(logRecord.SeverityText())
		newRecord.SetTraceID(logRecord.TraceID())
		newRecord.SetSpanID(logRecord.SpanID())
		newRecord.SetFlags(logRecord.Flags())

		// Transform the result into a security event
		p.transformToSecurityEvent(&newRecord, result, map[string]interface{}{
			"metadata.name":      metadataNameStr,
			"metadata.namespace": metadataNamespaceStr,
			"scope.name":         scopeNameStr,
			"scope.namespace":    scopeNamespaceStr,
			"scope.kind":         scopeKindStr,
			"scope.uid":          scopeUIDStr,
			"scope.apiVersion":   scopeAPIVersionStr,
			"workload.name":      workloadInfo.name,
			"workload.kind":      workloadInfo.kind,
			"workload.namespace": workloadInfo.namespace,
			"workload.uid":       workloadInfo.uid,
		}, attrs)

		newRecords = append(newRecords, newRecord)
		processedCount++
	}

	p.logger.Info("OpenReports log processing completed",
		zap.Int("original_logs", 1),
		zap.Int("total_results", len(resultsArray)),
		zap.Int("processed_results", processedCount),
		zap.Int("filtered_results", filteredCount),
		zap.Int("security_events_created", len(newRecords)),
		zap.String("metadata.name", metadataNameStr),
		zap.String("scope.name", scopeNameStr))

	p.logger.Debug("OpenReports log transformation summary",
		zap.Int("original_logs", 1),
		zap.Int("total_results", len(resultsArray)),
		zap.Int("processed_results", processedCount),
		zap.Int("filtered_results", filteredCount),
		zap.Int("security_events_created", len(newRecords)),
		zap.Bool("status_filter_enabled", len(p.config.StatusFilter) > 0),
		zap.String("metadata.name", metadataNameStr))

	return newRecords, nil
}

// Result represents a single result from the OpenReports results array
type Result struct {
	Source     string                 `json:"source"`
	Timestamp  Timestamp              `json:"timestamp"`
	Message    string                 `json:"message"`
	Policy     string                 `json:"policy"`
	Properties map[string]interface{} `json:"properties"`
	Result     string                 `json:"result"` // pass, fail, error, skip
	Rule       string                 `json:"rule"`
	Scored     bool                   `json:"scored"`
	Severity   string                 `json:"severity,omitempty"`
	Category   string                 `json:"category,omitempty"`
}

// Timestamp represents the timestamp in the result
type Timestamp struct {
	Seconds int64 `json:"seconds"`
	Nanos   int64 `json:"nanos"`
}

// transformToSecurityEvent transforms a result into a security event log record
func (p *Processor) transformToSecurityEvent(logRecord *plog.LogRecord, result Result, metadata map[string]interface{}, originalAttrs pcommon.Map) {
	attrs := logRecord.Attributes()

	// Generate event ID
	eventID := uuid.New().String()
	attrs.PutStr("event.id", eventID)

	// Hardcoded event fields
	attrs.PutStr("event.version", "1.309")
	attrs.PutStr("event.category", "COMPLIANCE")
	attrs.PutStr("event.name", "Compliance finding event")
	attrs.PutStr("event.type", "COMPLIANCE_FINDING")

	// Event description: "Policy violation on <pod> for rule <rule>" or appropriate message based on result
	scopeName := getString(metadata, "scope.name")
	rule := result.Rule
	if rule == "" {
		rule = "unknown"
	}

	var eventDescription string
	switch result.Result {
	case "fail":
		eventDescription = fmt.Sprintf("Policy violation on %s for rule %s", scopeName, rule)
	case "pass":
		eventDescription = fmt.Sprintf("Policy check passed on %s for rule %s", scopeName, rule)
	case "error":
		eventDescription = fmt.Sprintf("Policy check error on %s for rule %s", scopeName, rule)
	case "skip":
		eventDescription = fmt.Sprintf("Policy check skipped on %s for rule %s", scopeName, rule)
	default:
		eventDescription = fmt.Sprintf("Policy evaluation on %s for rule %s", scopeName, rule)
	}
	attrs.PutStr("event.description", eventDescription)

	// Product fields (empty for now)
	attrs.PutStr("product.name", "")
	attrs.PutStr("product.vendor", "")

	// Smartscape type - K8S_POD if scope.kind is Pod
	scopeKind := getString(metadata, "scope.kind")
	if scopeKind == "Pod" {
		attrs.PutStr("smartscape.type", "K8S_POD")
	}

	// Map finding.severity to dt.security.risk.level
	riskLevel := mapSeverityToRiskLevel(result.Severity)
	attrs.PutStr("dt.security.risk.level", riskLevel)

	// Calculate risk score based on level
	riskScore := calculateRiskScore(riskLevel)
	attrs.PutDouble("dt.security.risk.score", riskScore)

	// Object fields
	scopeUID := getString(metadata, "scope.uid")
	if scopeUID != "" {
		attrs.PutStr("object.id", scopeUID)
	}
	if scopeKind != "" {
		attrs.PutStr("object.type", scopeKind)
	}

	// Finding fields
	attrs.PutStr("finding.description", result.Message)
	findingID := uuid.New().String()
	attrs.PutStr("finding.id", findingID)

	if result.Severity != "" {
		attrs.PutStr("finding.severity", result.Severity)
	}

	// Finding time.created from result timestamp
	if result.Timestamp.Seconds > 0 {
		resultTime := time.Unix(result.Timestamp.Seconds, result.Timestamp.Nanos)
		logRecord.SetTimestamp(pcommon.NewTimestampFromTime(resultTime))
		// Also store as finding.time.created
		attrs.PutStr("finding.time.created", resultTime.Format(time.RFC3339Nano))
	}

	// Finding title: policy + rule
	findingTitle := result.Policy
	if result.Rule != "" {
		findingTitle = fmt.Sprintf("%s - %s", result.Policy, result.Rule)
	}
	attrs.PutStr("finding.title", findingTitle)

	// Finding type is the policy
	if result.Policy != "" {
		attrs.PutStr("finding.type", result.Policy)
	}

	// Finding URL (empty for now)
	attrs.PutStr("finding.url", "")

	// Compliance fields
	if result.Rule != "" {
		attrs.PutStr("compliance.control", result.Rule)
	}
	if result.Policy != "" {
		attrs.PutStr("compliance.requirements", result.Policy)
	}
	// compliance.standards can be omitted or hardcoded
	// For now, we'll omit it or use category if available
	if result.Category != "" {
		attrs.PutStr("compliance.standards", result.Category)
	}

	// Map result.result to compliance.status
	complianceStatus := mapResultToComplianceStatus(result.Result)
	attrs.PutStr("compliance.status", complianceStatus)

	// Copy all k8s.* fields from original log
	copyK8sFields(attrs, originalAttrs, metadata)

	// Set the log body/content to the security event message
	logRecord.Body().SetStr(result.Message)
}

// mapSeverityToRiskLevel maps finding severity to dt.security.risk.level
func mapSeverityToRiskLevel(severity string) string {
	switch severity {
	case "critical":
		return "CRITICAL"
	case "high":
		return "HIGH"
	case "medium":
		return "MEDIUM"
	case "low":
		return "LOW"
	default:
		// Default to MEDIUM if severity is not provided or unknown
		return "MEDIUM"
	}
}

// calculateRiskScore calculates the risk score based on risk level
func calculateRiskScore(riskLevel string) float64 {
	switch riskLevel {
	case "CRITICAL":
		return 10.0
	case "HIGH":
		return 8.9
	case "MEDIUM":
		return 6.9
	case "LOW":
		return 3.9
	default:
		return 0.0
	}
}

// mapResultToComplianceStatus maps result.result to compliance.status
func mapResultToComplianceStatus(result string) string {
	switch result {
	case "pass":
		return "PASSED"
	case "fail":
		return "FAILED"
	case "error":
		return "MANUAL" // errors may need manual review
	case "skip":
		return "NOT_RELEVANT"
	default:
		return "MANUAL" // default for unknown statuses
	}
}

// workloadInfo represents extracted workload information
type workloadInfo struct {
	name      string
	kind      string
	namespace string
	uid       string
}

// extractWorkloadInfo extracts workload information from owner references or pod name
func extractWorkloadInfo(attrs pcommon.Map, podName string, namespace string) workloadInfo {
	info := workloadInfo{}
	info.namespace = namespace // Workload namespace is the same as pod namespace

	// Try to extract from owner references first
	ownerRefsVal, exists := attrs.Get("metadata.ownerReferences")
	if exists {
		// ownerReferences is stored as an array of JSON strings
		var ownerRefs []string
		if ownerRefsVal.Type() == pcommon.ValueTypeSlice {
			slice := ownerRefsVal.Slice()
			for i := 0; i < slice.Len(); i++ {
				ownerRefs = append(ownerRefs, slice.At(i).AsString())
			}
		} else if ownerRefsVal.Type() == pcommon.ValueTypeStr {
			var jsonArray []string
			if err := json.Unmarshal([]byte(ownerRefsVal.AsString()), &jsonArray); err == nil {
				ownerRefs = jsonArray
			} else {
				ownerRefs = []string{ownerRefsVal.AsString()}
			}
		}

		// Parse owner references to find workload
		for _, ownerRefStr := range ownerRefs {
			var ownerRef map[string]interface{}
			if err := json.Unmarshal([]byte(ownerRefStr), &ownerRef); err == nil {
				kind, ok := ownerRef["kind"].(string)
				if ok && isWorkloadKind(kind) {
					info.kind = kind
					if name, ok := ownerRef["name"].(string); ok {
						info.name = name
					}
					if uid, ok := ownerRef["uid"].(string); ok {
						info.uid = uid
					}
					break // Take the first workload owner
				}
			}
		}
	}

	// If we couldn't find workload from owner references, try to infer from pod name
	if info.name == "" && podName != "" {
		// Pod names typically follow pattern: <workload-name>-<hash>-<random>
		// e.g., "cert-manager-cainjector-89fd4b8f9-t9xlf" -> "cert-manager-cainjector"
		// Extract workload name by removing hash and random suffix
		parts := splitPodName(podName)
		if len(parts) >= 2 {
			// Remove the last two parts (hash and random)
			workloadName := ""
			for i := 0; i < len(parts)-2; i++ {
				if i > 0 {
					workloadName += "-"
				}
				workloadName += parts[i]
			}
			if workloadName != "" {
				info.name = workloadName
				// Default to Deployment if kind is not known
				if info.kind == "" {
					info.kind = "Deployment"
				}
			}
		}
	}

	return info
}

// isWorkloadKind checks if a Kubernetes kind is a workload type
func isWorkloadKind(kind string) bool {
	workloadKinds := map[string]bool{
		"Deployment":  true,
		"StatefulSet": true,
		"DaemonSet":   true,
		"Job":         true,
		"CronJob":     true,
		"ReplicaSet":  true,
	}
	return workloadKinds[kind]
}

// splitPodName splits a pod name into its components
func splitPodName(podName string) []string {
	var parts []string
	current := ""
	for _, char := range podName {
		if char == '-' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

// copyK8sFields copies all k8s.* fields from the original attributes
func copyK8sFields(targetAttrs pcommon.Map, originalAttrs pcommon.Map, metadata map[string]interface{}) {
	// Copy k8s.* fields from original attributes
	originalAttrs.Range(func(key string, value pcommon.Value) bool {
		if len(key) > 4 && key[:4] == "k8s." {
			copyValue(targetAttrs, key, value)
		}
		return true
	})

	// Also add k8s fields from metadata if available
	if scopeName, ok := metadata["scope.name"]; ok {
		targetAttrs.PutStr("k8s.pod.name", fmt.Sprintf("%v", scopeName))
	}
	if scopeNamespace, ok := metadata["scope.namespace"]; ok {
		targetAttrs.PutStr("k8s.namespace.name", fmt.Sprintf("%v", scopeNamespace))
	}
	if scopeKind, ok := metadata["scope.kind"]; ok {
		kindStr := fmt.Sprintf("%v", scopeKind)
		targetAttrs.PutStr("k8s.resource.kind", kindStr)
		if kindStr == "Pod" {
			targetAttrs.PutStr("k8s.pod.name", getString(metadata, "scope.name"))
		}
	}
	if scopeUID, ok := metadata["scope.uid"]; ok {
		targetAttrs.PutStr("k8s.resource.uid", fmt.Sprintf("%v", scopeUID))
	}

	// Add workload fields
	if workloadName, ok := metadata["workload.name"]; ok && workloadName != "" {
		workloadKind := getString(metadata, "workload.kind")
		if workloadKind == "Deployment" {
			targetAttrs.PutStr("k8s.deployment.name", fmt.Sprintf("%v", workloadName))
		} else if workloadKind == "StatefulSet" {
			targetAttrs.PutStr("k8s.statefulset.name", fmt.Sprintf("%v", workloadName))
		} else if workloadKind == "DaemonSet" {
			targetAttrs.PutStr("k8s.daemonset.name", fmt.Sprintf("%v", workloadName))
		}
		targetAttrs.PutStr("k8s.workload.name", fmt.Sprintf("%v", workloadName))
		targetAttrs.PutStr("k8s.workload.kind", workloadKind)
	}
	if workloadNamespace, ok := metadata["workload.namespace"]; ok && workloadNamespace != "" {
		targetAttrs.PutStr("k8s.workload.namespace", fmt.Sprintf("%v", workloadNamespace))
	}
	if workloadUID, ok := metadata["workload.uid"]; ok && workloadUID != "" {
		targetAttrs.PutStr("k8s.workload.uid", fmt.Sprintf("%v", workloadUID))
	}
}

// copyValue copies a pcommon.Value to the target map
func copyValue(target pcommon.Map, key string, value pcommon.Value) {
	switch value.Type() {
	case pcommon.ValueTypeStr:
		target.PutStr(key, value.Str())
	case pcommon.ValueTypeInt:
		target.PutInt(key, value.Int())
	case pcommon.ValueTypeDouble:
		target.PutDouble(key, value.Double())
	case pcommon.ValueTypeBool:
		target.PutBool(key, value.Bool())
	case pcommon.ValueTypeSlice:
		slice := target.PutEmptySlice(key)
		value.Slice().CopyTo(slice)
	case pcommon.ValueTypeMap:
		m := target.PutEmptyMap(key)
		value.Map().CopyTo(m)
	}
}

// getString safely gets a string value from metadata
func getString(metadata map[string]interface{}, key string) string {
	if val, ok := metadata[key]; ok {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

// isStatusAllowed checks if a result status is in the allowed filter list
func (p *Processor) isStatusAllowed(status string) bool {
	// If no filter is configured, allow all statuses
	if len(p.config.StatusFilter) == 0 {
		return true
	}

	// Check if status is in the filter list
	for _, allowedStatus := range p.config.StatusFilter {
		if status == allowedStatus {
			return true
		}
	}

	return false
}

// putAttributeValue sets an attribute value based on the type
func putAttributeValue(attrs pcommon.Map, key string, value interface{}) {
	switch v := value.(type) {
	case string:
		attrs.PutStr(key, v)
	case int:
		attrs.PutInt(key, int64(v))
	case int64:
		attrs.PutInt(key, v)
	case float64:
		attrs.PutDouble(key, v)
	case bool:
		attrs.PutBool(key, v)
	default:
		// Convert to string for other types
		attrs.PutStr(key, fmt.Sprintf("%v", v))
	}
}
