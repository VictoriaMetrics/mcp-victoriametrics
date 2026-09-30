package v1

import (
	"fmt"
	"regexp"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func isValidUUID(uuid string) bool {
	return uuidRegex.MatchString(uuid)
}

func checkDeploymentID(deploymentID string) error {
	if deploymentID == "" {
		return fmt.Errorf("deployment ID cannot be empty")
	}
	if !isValidUUID(deploymentID) {
		return fmt.Errorf("invalid deployment ID format: %s", deploymentID)
	}
	return nil
}

var tenantIDRegex = regexp.MustCompile(`^(\d+)(:\d+)?$`)

func isValidTenantID(tenantID string) bool {
	return tenantIDRegex.MatchString(tenantID)
}

// validateCommonDeploymentParams validates parameters common to both create and update operations.
//
// Storage size is only checked for a unit and a non-zero value. The valid range and the
// step it must sit on depend on the storage type, the deployment topology and the
// per-installation configuration, none of which the client can see, so the API is the
// authority on them.
func validateCommonDeploymentParams(
	name string,
	tier uint32,
	maintenanceWindow MaintenanceWindow,
	storageSize uint64,
	storageSizeUnit StorageUnit,
	retention uint32,
	retentionUnit DurationUnit,
) error {
	if name == "" {
		return fmt.Errorf("deployment name cannot be empty")
	}
	if tier == 0 {
		return fmt.Errorf("deployment tier cannot be empty")
	}
	if maintenanceWindow != MaintenanceWindowWeekendDays &&
		maintenanceWindow != MaintenanceWindowBusinessDays {
		return fmt.Errorf("invalid maintenance window: %s", maintenanceWindow)
	}
	if storageSize == 0 {
		return fmt.Errorf("deployment storage size cannot be zero")
	}
	if storageSizeUnit != StorageUnitGB && storageSizeUnit != StorageUnitTB {
		return fmt.Errorf("invalid storage size unit: %s", storageSizeUnit)
	}
	if retention == 0 {
		return fmt.Errorf("deployment retention cannot be zero")
	}
	if retentionUnit != DurationUnitDay && retentionUnit != DurationUnitMonth {
		return fmt.Errorf("invalid retention unit: %s, only days and months are supported", retentionUnit)
	}
	return nil
}

// isValidDeduplicationUnit reports whether unit is a unit the API accepts for a
// deduplication window.
func isValidDeduplicationUnit(unit DurationUnit) bool {
	return unit == DurationUnitSecond || unit == DurationUnitMillisecond
}

// validateDeduplicationForCreate checks the deduplication window of a create request.
// The deployment type is known here, so the rule is exact: metrics deployments must
// carry a unit, and for the other types the API ignores both fields.
func validateDeduplicationForCreate(deploymentType DeploymentType, deduplication *uint32, deduplicationUnit *DurationUnit) error {
	if !deploymentType.SupportsDeduplication() {
		if deduplication != nil || deduplicationUnit != nil {
			return fmt.Errorf("deduplication window is not supported for %s deployments", deploymentType)
		}
		return nil
	}
	if deduplication == nil || deduplicationUnit == nil {
		return fmt.Errorf("deduplication and deduplication unit must be set together")
	}
	if !isValidDeduplicationUnit(*deduplicationUnit) {
		return fmt.Errorf("invalid deduplication unit: %s, only seconds and milliseconds are supported", *deduplicationUnit)
	}
	return nil
}

// validateDeduplicationForUpdate checks the deduplication window of an update request.
//
// An update body carries no deployment type, so the request states its intent through the
// two pointers instead of the client having to guess it: both nil leaves the deployment's
// window untouched, which is what a vlogs_single or vtraces_single update needs and what a
// metrics update that is not changing deduplication needs too. Both set changes the
// window, a zero window included. Exactly one set is always a mistake.
func validateDeduplicationForUpdate(deduplication *uint32, deduplicationUnit *DurationUnit) error {
	if deduplication == nil && deduplicationUnit == nil {
		return nil
	}
	if deduplication == nil || deduplicationUnit == nil {
		return fmt.Errorf("deduplication and deduplication unit must be set together")
	}
	if !isValidDeduplicationUnit(*deduplicationUnit) {
		return fmt.Errorf("invalid deduplication unit: %s, only seconds and milliseconds are supported", *deduplicationUnit)
	}
	return nil
}

// isCreatableDeploymentType reports whether deploymentType is a type the API can create.
//
// This is a create-time question only. Do not reuse it to gate a list filter: the API
// decides which types it can filter on, and an SDK build that predates a new type must
// not stop a caller from asking for it.
func isCreatableDeploymentType(deploymentType DeploymentType) bool {
	switch deploymentType {
	case DeploymentTypeSingleNode, DeploymentTypeCluster, DeploymentTypeVLogs, DeploymentTypeVTraces:
		return true
	default:
		return false
	}
}

// validateCreateDeploymentParams validates parameters specific to deployment creation
func validateCreateDeploymentParams(
	deploymentType DeploymentType,
	region string,
	provider DeploymentCloudProvider,
) error {
	if !isCreatableDeploymentType(deploymentType) {
		return fmt.Errorf("invalid deployment type: %s", deploymentType)
	}
	if region == "" {
		return fmt.Errorf("deployment region cannot be empty")
	}
	if provider != DeploymentCloudProviderAWS {
		return fmt.Errorf("unsupported deployment cloud provider: %s", provider)
	}
	return nil
}
