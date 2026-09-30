package v1

import (
	"encoding/json"
	"maps"
	"time"
)

// DeploymentType - type of the deployment (single_node / cluster / vlogs_single / vtraces_single)
type DeploymentType string

const (
	// DeploymentTypeSingleNode - single node VictoriaMetrics deployment
	DeploymentTypeSingleNode DeploymentType = "single_node"
	// DeploymentTypeCluster - cluster VictoriaMetrics deployment
	DeploymentTypeCluster DeploymentType = "cluster"
	// DeploymentTypeVLogs - single node VictoriaLogs deployment
	DeploymentTypeVLogs DeploymentType = "vlogs_single"
	// DeploymentTypeVTraces - single node VictoriaTraces deployment.
	// Available only for accounts where VictoriaTraces is enabled.
	DeploymentTypeVTraces DeploymentType = "vtraces_single"
)

// SupportsDeduplication reports whether the deployment type has a deduplication window.
// VictoriaLogs and VictoriaTraces deployments do not, and the API ignores the
// deduplication fields of create and update requests for them.
func (t DeploymentType) SupportsDeduplication() bool {
	return t == DeploymentTypeSingleNode || t == DeploymentTypeCluster
}

// SupportsAlertingRules reports whether the deployment type runs vmalert and therefore
// serves the rule file endpoints. VictoriaLogs and VictoriaTraces deployments do not.
func (t DeploymentType) SupportsAlertingRules() bool {
	return t == DeploymentTypeSingleNode || t == DeploymentTypeCluster
}

func (t DeploymentType) String() string {
	return string(t)
}

// DeploymentCloudProvider - cloud provider type for deployment
type DeploymentCloudProvider string

const (
	// DeploymentCloudProviderAWS - AWS cloud provider
	DeploymentCloudProviderAWS DeploymentCloudProvider = "aws"
)

func (p DeploymentCloudProvider) String() string {
	return string(p)
}

// DeploymentStatus - current status of the deployment
type DeploymentStatus string

const (
	DeploymentStatusProvisioning DeploymentStatus = "PROVISIONING"
	DeploymentStatusRunning      DeploymentStatus = "RUNNING"
	DeploymentStatusError        DeploymentStatus = "ERROR"
	DeploymentStatusStopped      DeploymentStatus = "STOPPED"
)

func (s DeploymentStatus) String() string {
	return string(s)
}

// DurationUnit represents a unit of time
type DurationUnit string

const (
	// DurationUnitMillisecond - milliseconds
	DurationUnitMillisecond DurationUnit = "ms"
	// DurationUnitSecond - seconds
	DurationUnitSecond DurationUnit = "s"
	// DurationUnitDay - days
	DurationUnitDay DurationUnit = "d"
	// DurationUnitMonth - months
	DurationUnitMonth DurationUnit = "m"
)

// AccessMode defines token access mode
type AccessMode string

const (
	// AccessModeRead - read-only access mode
	AccessModeRead AccessMode = "r"
	// AccessModeWrite - write-only access mode
	AccessModeWrite AccessMode = "w"
	// AccessModeReadWrite - read+write access mode
	AccessModeReadWrite AccessMode = "rw"
)

func (a AccessMode) String() string {
	return string(a)
}

// StorageUnit - storage unit type (GB / TB)
type StorageUnit string

const (
	// StorageUnitGB - gigabyte storage unit
	StorageUnitGB StorageUnit = "GB"
	// StorageUnitTB - terabyte storage unit
	StorageUnitTB StorageUnit = "TB"
)

// MaintenanceWindow - maintenance window for the deployment
type MaintenanceWindow string

const (
	// MaintenanceWindowWeekendDays - maintenance window on weekdays
	MaintenanceWindowWeekendDays MaintenanceWindow = "Sat-Sun 3-4am"
	// MaintenanceWindowBusinessDays - maintenance window on business days
	MaintenanceWindowBusinessDays MaintenanceWindow = "Mon-Fri 4-5am"
)

// FlagList - list of command-line flags
type FlagList []string

// DeploymentFlags - Customized command-line flags for the deployment
type DeploymentFlags struct {
	// CommonFlags - Common command-line flags for vmsingle component
	SingleFlags FlagList `json:"single_flags"`
	// SelectFlags - Customized command-line flags for the vmselect component
	SelectFlags FlagList `json:"select_flags"`
	// StorageFlags - Customized command-line flags for the vmstorage component
	StorageFlags FlagList `json:"storage_flags"`
	// InsertFlags - Customized command-line flags for the vminsert component
	InsertFlags FlagList `json:"insert_flags"`
}

func (u StorageUnit) String() string {
	return string(u)
}

// TierInfoCommon holds the fields every tier reports, regardless of its deployment type.
type TierInfoCommon struct {
	// ID is the unique identifier of the tier of given type
	ID uint32 `json:"id"`
	// Type of the deployment (single_node / cluster / vlogs_single / vtraces_single)
	Type DeploymentType `json:"type"`
	// CloudProvider is the name of Cloud provider of the deployment (aws)
	CloudProvider DeploymentCloudProvider `json:"cloud_provider"`
	// Name is the name of the tier
	Name string `json:"name"`
	// ComputeCostPerHour is the cost of the deployment per hour
	ComputeCostPerHour float64 `json:"compute_cost_per_hour"`
	// AccessTokenConcurrentRequests is the maximum number of concurrent requests for each access token
	AccessTokenConcurrentRequests int `json:"access_token_concurrent_requests"`
	// AccessTokenLimit is the maximum number of access tokens for deployments of this tier
	AccessTokenLimit int `json:"access_token_limit"`
}

// MetricsTierInfo holds the limits reported by single_node and cluster tiers.
type MetricsTierInfo struct {
	// IngestionRate is the maximum ingestion rate of the tier
	IngestionRate int `json:"ingestion_rate"`
	// ActiveTimeSeries is the maximum number of active time series of the tier
	ActiveTimeSeries int `json:"active_time_series"`
	// NewSeriesOver24h is the maximum number of new series over 24 hours of the tier
	NewSeriesOver24h int `json:"new_series_over_24h"`
	// DatapointsReadRate is the maximum read rate of the tier
	DatapointsReadRate int `json:"datapoints_read_rate"`
	// SeriesReadPerQuery is the maximum number of series read per query of the tier
	SeriesReadPerQuery int `json:"series_read_per_query"`
}

// LogsTierInfo holds the limits reported by vlogs_single tiers.
//
// The byte-valued limits are int64 rather than int: a tier can report more than 2 GB,
// which does not fit in a 32-bit int and would fail to decode on 32-bit targets.
type LogsTierInfo struct {
	// IngestionRateBytes is the maximum ingestion rate in bytes per second
	IngestionRateBytes int64 `json:"ingestion_rate_bytes"`
	// ActiveLogStreams is the maximum number of active log streams of the tier
	ActiveLogStreams int `json:"active_log_streams"`
	// NewStreamsOver24h is the maximum number of new log streams over 24 hours of the tier
	NewStreamsOver24h int `json:"new_streams_over_24h"`
	// DataReadRate is the maximum read rate in bytes per second of the tier
	DataReadRate int64 `json:"data_read_rate"`
	// BytesPerQuery is the maximum number of bytes scanned per query
	BytesPerQuery int64 `json:"bytes_per_query"`
}

// TracesTierInfo holds the limits reported by vtraces_single tiers.
//
// The JSON keys match the vlogs_single ones, since VictoriaTraces reports its limits
// in terms of log streams as well. The type is kept separate so that the two can
// diverge without breaking callers.
//
// The byte-valued limits are int64, for the same reason as LogsTierInfo.
type TracesTierInfo struct {
	// IngestionRateBytes is the maximum ingestion rate in bytes per second
	IngestionRateBytes int64 `json:"ingestion_rate_bytes"`
	// ActiveLogStreams is the maximum number of active log streams of the tier
	ActiveLogStreams int `json:"active_log_streams"`
	// NewStreamsOver24h is the maximum number of new log streams over 24 hours of the tier
	NewStreamsOver24h int `json:"new_streams_over_24h"`
	// DataReadRate is the maximum read rate in bytes per second of the tier
	DataReadRate int64 `json:"data_read_rate"`
	// BytesPerQuery is the maximum number of bytes scanned per query
	BytesPerQuery int64 `json:"bytes_per_query"`
}

// TierInfo represents the information about the tier in public VMCloud API.
//
// The limits a tier reports depend on its Type, so they live in a dedicated struct
// per deployment type. At most one of Metrics, Logs and Traces is set: single_node
// and cluster tiers fill Metrics, vlogs_single tiers fill Logs and vtraces_single
// tiers fill Traces. The other two stay nil, and a tier of a type this version of
// the SDK does not know leaves all three nil rather than guessing.
//
// Type is the single discriminator: both UnmarshalJSON and MarshalJSON pick the
// limit struct from it, so a limit struct that does not match Type is not written.
type TierInfo struct {
	TierInfoCommon
	// Metrics holds the limits of single_node and cluster tiers, nil for other types
	Metrics *MetricsTierInfo `json:"-"`
	// Logs holds the limits of vlogs_single tiers, nil for other types
	Logs *LogsTierInfo `json:"-"`
	// Traces holds the limits of vtraces_single tiers, nil for other types
	Traces *TracesTierInfo `json:"-"`
}

// UnmarshalJSON decodes a tier into its common fields plus the limit struct its Type
// calls for. A type this version of the SDK does not know leaves all three limit
// structs nil, so that a new deployment type is not mistaken for a metrics tier.
func (t *TierInfo) UnmarshalJSON(data []byte) error {
	t.TierInfoCommon = TierInfoCommon{}
	t.Metrics, t.Logs, t.Traces = nil, nil, nil

	if err := json.Unmarshal(data, &t.TierInfoCommon); err != nil {
		return err
	}
	switch t.Type {
	case DeploymentTypeSingleNode, DeploymentTypeCluster:
		var limits MetricsTierInfo
		if err := json.Unmarshal(data, &limits); err != nil {
			return err
		}
		t.Metrics = &limits
	case DeploymentTypeVLogs:
		var limits LogsTierInfo
		if err := json.Unmarshal(data, &limits); err != nil {
			return err
		}
		t.Logs = &limits
	case DeploymentTypeVTraces:
		var limits TracesTierInfo
		if err := json.Unmarshal(data, &limits); err != nil {
			return err
		}
		t.Traces = &limits
	}
	return nil
}

// MarshalJSON encodes a tier back into the flat object the API returns. It selects the
// limit struct by Type, mirroring UnmarshalJSON, so that the two stay symmetric even
// when a caller has filled in a limit struct that does not match Type.
//
// The receiver is deliberately a value, mirroring time.Time and net.IP: encoding/json
// only reaches a pointer-receiver MarshalJSON when the tier is addressable, so a
// pointer receiver here would make json.Marshal(tier) drop every limit without an
// error. TestTierInfoMarshalValue guards that.
func (t TierInfo) MarshalJSON() ([]byte, error) {
	parts := []any{t.TierInfoCommon}
	switch t.Type {
	case DeploymentTypeSingleNode, DeploymentTypeCluster:
		if t.Metrics != nil {
			parts = append(parts, t.Metrics)
		}
	case DeploymentTypeVLogs:
		if t.Logs != nil {
			parts = append(parts, t.Logs)
		}
	case DeploymentTypeVTraces:
		if t.Traces != nil {
			parts = append(parts, t.Traces)
		}
	}

	fields := map[string]json.RawMessage{}
	for _, part := range parts {
		encoded, err := json.Marshal(part)
		if err != nil {
			return nil, err
		}
		var partFields map[string]json.RawMessage
		if err := json.Unmarshal(encoded, &partFields); err != nil {
			return nil, err
		}
		maps.Copy(fields, partFields)
	}
	return json.Marshal(fields)
}

// TierInfoList represents the list of TierInfo
type TierInfoList []TierInfo

// DeploymentSummary - simplified representation of deployment for list view
type DeploymentSummary struct {
	// ID - unique identifier of the deployment
	ID string `json:"id"`
	// Name - human-readable name of the deployment
	Name string `json:"name"`
	// Type of the deployment (single_node / cluster / vlogs_single / vtraces_single)
	Type DeploymentType `json:"type"`
	// Tier - tier identifier of the deployment
	Tier uint32 `json:"tier"`
	// Version of VictoriaMetrics used in the deployment
	Version string `json:"version"`
	// CloudProvider - ID of then cloud provider of the deployment
	CloudProvider DeploymentCloudProvider `json:"cloud_provider"`
	// Region of the deployment in specified cloud provider
	Region string `json:"region"`
	// CreatedAt - timestamp of deployment creation
	CreatedAt time.Time `json:"created_at"`
	// Status - current status of the deployment
	Status DeploymentStatus `json:"status"`
}

// DeploymentSummaryList represents the list of DeploymentSummary
type DeploymentSummaryList []DeploymentSummary

// DeploymentPrice - price of the deployment (USD per month)
type DeploymentPrice struct {
	// ComputeCost - cost of the compute resources
	ComputeCost float64 `json:"compute_cost"`
	// StorageCost - cost of the storage resources
	StorageCost float64 `json:"storage_cost"`
	// TotalCost - total cost of the deployment (without network costs)
	TotalCost float64 `json:"total_cost"`
}

// DeploymentInfo - deployment details for the public API
type DeploymentInfo struct {
	// ID - unique identifier of the deployment
	ID string `json:"id"`
	// Name - human-readable name of the deployment
	Name string `json:"name"`
	// Type of the deployment (single_node / cluster / vlogs_single / vtraces_single)
	Type DeploymentType `json:"type"`
	// Tier - tier identifier of the deployment
	Tier uint32 `json:"tier"`
	// Version of VictoriaMetrics used in the deployment
	Version string `json:"version"`
	// CloudProvider - ID of then cloud provider of the deployment
	CloudProvider DeploymentCloudProvider `json:"cloud_provider"`
	// Region of the deployment in specified cloud provider
	Region string `json:"region"`
	// CreatedAt - timestamp of deployment creation
	CreatedAt time.Time `json:"created_at"`
	// Status - current status of the deployment
	Status DeploymentStatus `json:"status"`
	// RetentionValue - retention period of the deployment
	RetentionValue uint32 `json:"retention_value"`
	// RetentionUnit - retention period unit of the deployment
	RetentionUnit DurationUnit `json:"retention_unit"`
	// DeduplicationValue - deduplication period of the deployment. Absent for
	// vlogs_single and vtraces_single deployments, which have no deduplication window.
	// A pointer, so that a metrics deployment with a zero window stays distinguishable
	// from a deployment that has none. Use Deduplication for a safe read.
	DeduplicationValue *uint32 `json:"deduplication_value,omitempty"`
	// DeduplicationUnit - deduplication period unit of the deployment. Absent for
	// vlogs_single and vtraces_single deployments.
	DeduplicationUnit *DurationUnit `json:"deduplication_unit,omitempty"`
	// StorageSizeGb - storage size of the deployment
	StorageSizeGb uint64 `json:"storage_size_gb"`
	// MaintenanceWindow - maintenance window of the deployment
	MaintenanceWindow MaintenanceWindow `json:"maintenance_window"`
	// Price - price of the deployment
	Price DeploymentPrice `json:"price"`
	// VMSingleSettings - settings of the single component. Used by single_node,
	// vlogs_single and vtraces_single deployments.
	VMSingleSettings []string `json:"vmsingle_settings,omitempty"`
	// VMStorageSettings - VMStorage settings for cluster deployment
	VMStorageSettings []string `json:"vmstorage_settings,omitempty"`
	// VMSelectSettings - VMSelect settings for cluster deployment
	VMSelectSettings []string `json:"vmselect_settings,omitempty"`
	// VMInsertSettings - VMInsert settings for cluster deployment
	VMInsertSettings []string `json:"vminsert_settings,omitempty"`
	// AccessEndpoint - endpoint of the deployment (URL entrypoint to API of the deployment)
	AccessEndpoint string `json:"access_endpoint"`
}

// Deduplication reports the deduplication window of the deployment. ok is false for
// vlogs_single and vtraces_single deployments, which have no deduplication window, and for
// a response that does not carry both fields.
//
// Type decides whether a window applies at all, rather than the two fields being present:
// a response that spells the absent window out as an explicit zero and an empty unit still
// reports ok false for a type that has no window.
func (d DeploymentInfo) Deduplication() (value uint32, unit DurationUnit, ok bool) {
	switch d.Type {
	case DeploymentTypeVLogs, DeploymentTypeVTraces:
		return 0, "", false
	}
	if d.DeduplicationValue == nil || d.DeduplicationUnit == nil {
		return 0, "", false
	}
	return *d.DeduplicationValue, *d.DeduplicationUnit, true
}

// DeploymentInfoList represents the list of DeploymentInfo
type DeploymentInfoList []DeploymentInfo

// DeploymentCreationRequest represents the request for creating a deployment
type DeploymentCreationRequest struct {
	// Name - human-readable name of the deployment
	Name string `json:"name"`
	// Type of the deployment (single_node / cluster / vlogs_single / vtraces_single)
	Type DeploymentType `json:"type"`
	// Provider - cloud provider of the deployment
	Provider DeploymentCloudProvider `json:"provider"`
	// Region of the deployment in specified cloud provider
	Region string `json:"region"`
	// Tier - tier identifier of the deployment
	Tier uint32 `json:"tier"`
	// StorageSize - storage size in units specified in StorageSizeUnit
	StorageSize uint64 `json:"storage_size"`
	// StorageSizeUnit - storage size unit (GB / TB)
	StorageSizeUnit StorageUnit `json:"storage_size_unit"`
	// Deduplication window for the deployment in units specified in DeduplicationUnit.
	// Required for single_node and cluster deployments, where a zero window is a valid
	// setting distinct from having none: write new(uint32(0)) to send one. Leave it nil
	// for vlogs_single and vtraces_single, which have no deduplication window; setting it
	// for those types is rejected rather than quietly dropped.
	Deduplication *uint32 `json:"deduplication,omitempty"`
	// DeduplicationUnit - deduplication window unit for the deployment. Set it together
	// with Deduplication; nil for the types that have no deduplication window.
	DeduplicationUnit *DurationUnit `json:"deduplication_unit,omitempty"`
	// Retention period for the deployment in units specified in RetentionUnit
	Retention uint32 `json:"retention"`
	// RetentionUnit - retention period unit for the deployment
	RetentionUnit DurationUnit `json:"retention_unit"`
	// MaintenanceWindow - maintenance window for the deployment
	MaintenanceWindow MaintenanceWindow `json:"maintenance_window"`
}

// DeploymentUpdateRequest represents the request for updating a deployment
type DeploymentUpdateRequest struct {
	// Name - human-readable name of the deployment
	Name string `json:"name"`
	// Tier - tier identifier of the deployment
	Tier uint32 `json:"tier"`
	// StorageSize - storage size in units specified in StorageSizeUnit
	StorageSize uint64 `json:"storage_size"`
	// StorageSizeUnit - storage size unit (GB / TB)
	StorageSizeUnit StorageUnit `json:"storage_size_unit"`
	// Deduplication window for the deployment in units specified in DeduplicationUnit.
	// An update body carries no deployment type, so nil is what carries the meaning here:
	// it leaves the deployment's current window untouched, which is what a vlogs_single or
	// vtraces_single update needs - those types have no window - and what a metrics update
	// that is not changing deduplication needs too. Set it together with DeduplicationUnit
	// to change the window; a zero window is a valid setting and is sent as such.
	// Use new(uint32(10)) to take the address of a literal.
	Deduplication *uint32 `json:"deduplication,omitempty"`
	// DeduplicationUnit - deduplication window unit for the deployment. Set it together
	// with Deduplication; nil leaves the current unit untouched.
	// Use new(DurationUnitSecond) to take the address of a literal.
	DeduplicationUnit *DurationUnit `json:"deduplication_unit,omitempty"`
	// Retention period for the deployment in units specified in RetentionUnit
	Retention uint32 `json:"retention"`
	// RetentionUnit - retention period unit for the deployment
	RetentionUnit DurationUnit `json:"retention_unit"`
	// MaintenanceWindow - maintenance window for the deployment
	MaintenanceWindow MaintenanceWindow `json:"maintenance_window"`
	// Flags - customized command-line flags for the deployment
	Flags DeploymentFlags `json:"flags"`
}

// AccessTokensList represents the list of AccessToken
type AccessTokensList []AccessToken

// AccessToken represents the access token in public VMCloud API
type AccessToken struct {
	// ID is the unique identifier of the access token
	ID string `json:"id"`
	// Secret is the secret value of the access token (only 4 symbols are returned in access tokens list, for full secret use reveal endpoint)
	Secret string `json:"value"`
	// Type is the access mode of the token (read-only, write-only, read+write)
	Type AccessMode `json:"type"`
	// Description is the human-readable description of the access token
	Description string `json:"description"`
	// CreatedBy is the user who created the access token
	CreatedBy string `json:"created_by"`
	// CreatedAt is the timestamp of the access token creation
	CreatedAt time.Time `json:"created_at"`
	// TenantID represents the unique identifier of the tenant associated with this access token (optional)
	TenantID string `json:"tenant_id,omitempty"`
	// Timestamp of the last usage of the access token (within the last 7 days)
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
}

// AccessTokenCreateRequest represents the request for creating an access token
type AccessTokenCreateRequest struct {
	// Type is the access mode of the token (read-only, write-only, read+write)
	Type AccessMode `json:"type"`
	// Description is the human-readable description of the access token
	Description string `json:"description"`
	// TenantID represents the unique identifier of the tenant associated with this access token (optional)
	TenantID string `json:"tenant_id,omitempty"`
}

// RegionInfo represents the information about the region in public VMCloud API
type RegionInfo struct {
	// Provider is the name of the cloud provider
	CloudProvider DeploymentCloudProvider `json:"cloud_provider"`
	// Name is the name of the region
	Name string `json:"name"`
}

// RegionInfoList represents the list of RegionInfo
type RegionInfoList []RegionInfo

// CloudProviderInfo represents the information about the cloud provider in public VMCloud API
type CloudProviderInfo struct {
	// ID is the unique identifier of the cloud provider
	ID DeploymentCloudProvider `json:"id"`
	// URL is the URL of the cloud provider
	URL string `json:"url"`
}

// CloudProviderInfoList represents the list of CloudProviderInfo
type CloudProviderInfoList []CloudProviderInfo
