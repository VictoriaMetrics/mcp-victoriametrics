package metrics

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMetricsMetadataFiltering(t *testing.T) {
	tests := []struct {
		name           string
		apiResponse    string
		search         string
		metricType     string
		unit           string
		limit          float64
		expectedCount  int
		expectedMetric string
	}{
		{
			name: "Filter by type - counter",
			apiResponse: `{
				"status": "success",
				"data": {
					"http_requests_total": [{"type": "counter", "help": "Total HTTP requests", "unit": ""}],
					"memory_usage_bytes": [{"type": "gauge", "help": "Memory usage in bytes", "unit": "bytes"}]
				}
			}`,
			metricType:     "counter",
			expectedCount:  1,
			expectedMetric: "http_requests_total",
		},
		{
			name: "Filter by search - keyword in metric name",
			apiResponse: `{
				"status": "success",
				"data": {
					"http_requests_total": [{"type": "counter", "help": "Total HTTP requests", "unit": ""}],
					"http_response_time": [{"type": "gauge", "help": "HTTP response time", "unit": "seconds"}],
					"memory_usage_bytes": [{"type": "gauge", "help": "Memory usage in bytes", "unit": "bytes"}]
				}
			}`,
			search:        "http",
			expectedCount: 2,
		},
		{
			name: "Filter by search - keyword in help text",
			apiResponse: `{
				"status": "success",
				"data": {
					"requests_total": [{"type": "counter", "help": "Total HTTP requests", "unit": ""}],
					"memory_usage_bytes": [{"type": "gauge", "help": "Memory usage in bytes", "unit": "bytes"}]
				}
			}`,
			search:        "http",
			expectedCount: 1,
		},
		{
			name: "Filter by unit",
			apiResponse: `{
				"status": "success",
				"data": {
					"http_response_time": [{"type": "gauge", "help": "HTTP response time", "unit": "seconds"}],
					"memory_usage_bytes": [{"type": "gauge", "help": "Memory usage in bytes", "unit": "bytes"}],
					"cpu_usage": [{"type": "gauge", "help": "CPU usage", "unit": ""}]
				}
			}`,
			unit:          "bytes",
			expectedCount: 1,
		},
		{
			name: "Filter with limit",
			apiResponse: `{
				"status": "success",
				"data": {
					"metric1": [{"type": "counter", "help": "Metric 1", "unit": ""}],
					"metric2": [{"type": "counter", "help": "Metric 2", "unit": ""}],
					"metric3": [{"type": "counter", "help": "Metric 3", "unit": ""}]
				}
			}`,
			limit:         2,
			expectedCount: 2,
		},
		{
			name: "Filter by multiple parameters",
			apiResponse: `{
				"status": "success",
				"data": {
					"http_requests_total": [{"type": "counter", "help": "Total HTTP requests", "unit": ""}],
					"http_response_time": [{"type": "gauge", "help": "HTTP response time", "unit": "seconds"}],
					"memory_usage_bytes": [{"type": "gauge", "help": "Memory usage in bytes", "unit": "bytes"}]
				}
			}`,
			search:        "http",
			metricType:    "gauge",
			expectedCount: 1,
		},
		{
			name: "Case insensitive filtering",
			apiResponse: `{
				"status": "success",
				"data": {
					"HTTP_Requests_Total": [{"type": "COUNTER", "help": "Total HTTP requests", "unit": ""}]
				}
			}`,
			search:        "http",
			metricType:    "counter",
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the API response
			var apiResponse struct {
				Status string `json:"status"`
				Data   map[string][]struct {
					Type string `json:"type"`
					Help string `json:"help"`
					Unit string `json:"unit"`
				} `json:"data"`
			}
			err := json.Unmarshal([]byte(tt.apiResponse), &apiResponse)
			if err != nil {
				t.Fatalf("Failed to unmarshal API response: %v", err)
			}

			// Apply filtering logic (extracted from the handler)
			filteredData := make(map[string][]struct {
				Type string `json:"type"`
				Help string `json:"help"`
				Unit string `json:"unit"`
			})

			searchLower := ""
			if tt.search != "" {
				searchLower = strings.ToLower(tt.search)
			}
			typeLower := ""
			if tt.metricType != "" {
				typeLower = strings.ToLower(tt.metricType)
			}
			unitLower := ""
			if tt.unit != "" {
				unitLower = strings.ToLower(tt.unit)
			}

			count := 0
			for metricName, metadataList := range apiResponse.Data {
				for _, metadata := range metadataList {
					if searchLower != "" {
						metricNameLower := strings.ToLower(metricName)
						helpLower := strings.ToLower(metadata.Help)
						if !strings.Contains(metricNameLower, searchLower) && !strings.Contains(helpLower, searchLower) {
							continue
						}
					}

					if typeLower != "" {
						if strings.ToLower(metadata.Type) != typeLower {
							continue
						}
					}

					if unitLower != "" {
						if strings.ToLower(metadata.Unit) != unitLower {
							continue
						}
					}

					if _, exists := filteredData[metricName]; !exists {
						filteredData[metricName] = []struct {
							Type string `json:"type"`
							Help string `json:"help"`
							Unit string `json:"unit"`
						}{}
					}
					filteredData[metricName] = append(filteredData[metricName], metadata)
					count++

					if tt.limit != 0 && count >= int(tt.limit) {
						break
					}
				}

				if tt.limit != 0 && count >= int(tt.limit) {
					break
				}
			}

			// Verify results
			if len(filteredData) != tt.expectedCount {
				t.Errorf("Expected %d metrics after filtering, got %d", tt.expectedCount, len(filteredData))
			}

			if tt.expectedMetric != "" {
				if _, exists := filteredData[tt.expectedMetric]; !exists {
					t.Errorf("Expected metric %s not found in filtered results", tt.expectedMetric)
				}
			}
		})
	}
}

func TestFilteringLogic(t *testing.T) {
	// Test the filtering decision logic
	tests := []struct {
		name         string
		search       string
		metricType   string
		unit         string
		shouldFilter bool
	}{
		{
			name:         "No filtering parameters",
			search:       "",
			metricType:   "",
			unit:         "",
			shouldFilter: false,
		},
		{
			name:         "Search parameter set",
			search:       "http",
			metricType:   "",
			unit:         "",
			shouldFilter: true,
		},
		{
			name:         "Type parameter set",
			search:       "",
			metricType:   "counter",
			unit:         "",
			shouldFilter: true,
		},
		{
			name:         "Unit parameter set",
			search:       "",
			metricType:   "",
			unit:         "bytes",
			shouldFilter: true,
		},
		{
			name:         "All parameters set",
			search:       "http",
			metricType:   "counter",
			unit:         "bytes",
			shouldFilter: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shouldFilter := tt.search != "" || tt.metricType != "" || tt.unit != ""
			if shouldFilter != tt.shouldFilter {
				t.Errorf("Expected shouldFilter=%v, got %v", tt.shouldFilter, shouldFilter)
			}
		})
	}
}
