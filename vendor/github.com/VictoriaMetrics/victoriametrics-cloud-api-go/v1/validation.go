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

// validateCommonDeploymentParams validates parameters common to both create and update operations
func validateCommonDeploymentParams(
	name string,
	tier uint32,
	maintenanceWindow MaintenanceWindow,
	storageSize uint64,
	storageSizeUnit StorageUnit,
	retention uint32,
	retentionUnit DurationUnit,
	deduplicationUnit DurationUnit,
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
	if deduplicationUnit != DurationUnitSecond && deduplicationUnit != DurationUnitMillisecond {
		return fmt.Errorf("invalid deduplication unit: %s, only seconds and milliseconds are supported", deduplicationUnit)
	}
	return nil
}

// validateCreateDeploymentParams validates parameters specific to deployment creation
func validateCreateDeploymentParams(
	deploymentType DeploymentType,
	region string,
	provider DeploymentCloudProvider,
	deploymentStorageSize uint64,
	storageSizeUnit StorageUnit,
) error {
	if deploymentType != DeploymentTypeSingleNode && deploymentType != DeploymentTypeCluster {
		return fmt.Errorf("invalid deployment type: %s", deploymentType)
	}
	if region == "" {
		return fmt.Errorf("deployment region cannot be empty")
	}
	if provider != DeploymentCloudProviderAWS {
		return fmt.Errorf("unsupported deployment cloud provider: %s", provider)
	}
	if deploymentType == DeploymentTypeSingleNode &&
		storageSizeUnit == StorageUnitTB && deploymentStorageSize > 16 {
		return fmt.Errorf("single-node deployments cannot have more than 16 TB of storage")
	}
	return nil
}
