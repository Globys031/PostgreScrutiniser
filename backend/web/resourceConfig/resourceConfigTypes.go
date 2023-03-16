// Package resourceConfig provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package resourceConfig

// Defines values for GetResourceConfigByIdParamsConfig.
const (
	AutovacuumWorkMem       GetResourceConfigByIdParamsConfig = "autovacuum_work_mem"
	DynamicSharedMemoryType GetResourceConfigByIdParamsConfig = "dynamic_shared_memory_type"
	HashMemMultiplier       GetResourceConfigByIdParamsConfig = "hash_mem_multiplier"
	HugePageSize            GetResourceConfigByIdParamsConfig = "huge_page_size"
	HugePages               GetResourceConfigByIdParamsConfig = "huge_pages"
	LogicalDecodingWorkMem  GetResourceConfigByIdParamsConfig = "logical_decoding_work_mem"
	MaintenanceWorkMem      GetResourceConfigByIdParamsConfig = "maintenance_work_mem"
	MaxPreparedTransactions GetResourceConfigByIdParamsConfig = "max_prepared_transactions"
	MaxStackDepth           GetResourceConfigByIdParamsConfig = "max_stack_depth"
	SharedBuffers           GetResourceConfigByIdParamsConfig = "shared_buffers"
	SharedMemoryType        GetResourceConfigByIdParamsConfig = "shared_memory_type"
	TempBuffers             GetResourceConfigByIdParamsConfig = "temp_buffers"
	WorkMem                 GetResourceConfigByIdParamsConfig = "work_mem"
)

// Error defines model for Error.
type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// ResourceConfig defines model for resourceConfig.
type ResourceConfig struct {
	// Details Details informing why a value was suggested
	Details *string `json:"details,omitempty"`

	// EnumVals specifies what type of values setting could have
	EnumVals *string `json:"enum_vals,omitempty"`

	// GotError specifies whether check got an error
	GotError *bool `json:"got_error,omitempty"`

	// MaxVal Maximum allowed value (needed for validation)
	MaxVal *string `json:"max_val,omitempty"`

	// MinVal Minimum allowed value (needed for validation)
	MinVal *string `json:"min_val,omitempty"`

	// Name Name of the setting
	Name *string `json:"name,omitempty"`

	// RequiresRestart Whether a restart is required after changing the value
	RequiresRestart *bool `json:"requires_restart,omitempty"`

	// SuggestedValue Value that will be suggested after running check
	SuggestedValue *string `json:"suggested_value,omitempty"`

	// Unit Unit of measurement (s, ms, kB, 8kB, etc.)
	Unit *string `json:"unit,omitempty"`

	// Value Value of the setting
	Value *string `json:"value,omitempty"`

	// Vartype Type of the value (boolean, integer, enum, etc.)
	Vartype *string `json:"vartype,omitempty"`
}

// ResourceConfigPatchSchema defines model for resourceConfigPatchSchema.
type ResourceConfigPatchSchema struct {
	// Name Name of the setting
	Name string `json:"name"`

	// SuggestedValue Value that will be suggested after running check
	SuggestedValue string `json:"suggested_value"`
}

// PatchResourceConfigsJSONBody defines parameters for PatchResourceConfigs.
type PatchResourceConfigsJSONBody = []ResourceConfigPatchSchema

// GetResourceConfigByIdParamsConfig defines parameters for GetResourceConfigById.
type GetResourceConfigByIdParamsConfig string

// PatchResourceConfigsJSONRequestBody defines body for PatchResourceConfigs for application/json ContentType.
type PatchResourceConfigsJSONRequestBody = PatchResourceConfigsJSONBody
