// Code for PostgreSQL Configuration
package resourceConfig

import (
	"bufio"
	"database/sql"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	sysctl "github.com/lorenzosaino/go-sysctl"

	"github.com/Globys031/PostgreScrutiniser/backend/utils"
)

type ResourceSetting struct {
	Name     string // name of the setting
	Value    string // value of the setting
	Unit     string // s, ms, kB, 8kB, etc...
	EnumVals string // If an enumrator, this stores enum values. For internal use only. Not exposted by API

	SuggestedValue string // Value that will be suggested after running check
	Details        string // Details informing why a value was suggested
	GotError       bool   // specifies whether check got an error
}

type Configuration struct {
	dbHandler    *sql.DB
	path         string // Path to postgresql.conf
	autoConfPath string // Path to postgresql.auto.conf
	backupDir    string // directory to where postgresql.auto.conf will be backed up
	settings     map[string]ResourceSetting
	appUser      *utils.User // postgrescrutiniser user
	postgresUser *utils.User // postgresql user
}

////////////////////////////////////////////////////////////////////
// Below are functions used to check and make configuration suggestions
////////////////////////////////////////////////////////////////////

// Meant for initialising Configuration upon first api call so that same
// reference can be reused for later calls.
func InitChecks(dbHandler *sql.DB, appUser *utils.User, postgresUser *utils.User, logger *utils.Logger) *Configuration {
	configFilePath, _ := utils.FindConfigFile(logger)
	ResourceSettings, _ := getPGSettings(dbHandler, logger)
	autoConfPath := filepath.Dir(configFilePath) + "/postgresql.auto.conf"
	backupDir := "/usr/local/postgrescrutiniser/backups"

	conf := Configuration{dbHandler: dbHandler, path: configFilePath, autoConfPath: autoConfPath, backupDir: backupDir, settings: ResourceSettings, appUser: appUser, postgresUser: postgresUser}

	return &conf
}

func RunChecks(conf *Configuration, logger *utils.Logger) *map[string]ResourceSetting {
	// Error logging is handled inside the functions
	conf.CheckSharedBuffers(logger)
	conf.CheckHugePages(logger)
	conf.CheckHugePageSize(logger)
	conf.CheckTempBuffers(logger)
	conf.CheckMaxPreparedTransactions(logger)
	conf.CheckWorkMem(logger)
	conf.CheckHashMemMultiplier(logger)
	conf.CheckMaintenanceWorkMem(logger)
	conf.CheckAutovacuumWorkMem(logger)
	conf.ChecklogicalDecodingWorkMem(logger)
	conf.CheckMaxStackDepth(logger)
	conf.CheckSharedMemoryType(logger)
	conf.CheckDynamicSharedMemoryType(logger)

	return &conf.settings
}

// Stores data returned by `pg_settings` into a map. Only stores settings we're interested in.
func getPGSettings(dbHandler *sql.DB, logger *utils.Logger) (map[string]ResourceSetting, error) {
	ResourceSettings := [13]string{"shared_buffers",
		"huge_pages",
		"huge_page_size",
		"temp_buffers",
		"max_prepared_transactions",
		"work_mem",
		"hash_mem_multiplier",
		"maintenance_work_mem",
		"autovacuum_work_mem",
		"logical_decoding_work_mem",
		"max_stack_depth",
		"shared_memory_type",
		"dynamic_shared_memory_type"}

	// Prepare the SQL statement
	stmt, err := dbHandler.Prepare("SELECT name,setting,unit,enumvals FROM pg_settings")
	if err != nil {
		logger.LogError(fmt.Errorf("Failed preparing SQL statement: %v", err))
		return nil, err
	}
	defer stmt.Close()

	// Query the database
	rows, err := stmt.Query()
	if err != nil {
		logger.LogError(fmt.Errorf("Failed querying Postgres: %v", err))
		return nil, err
	}
	defer rows.Close()

	// Initialize the map to store the results
	settingsMap := make(map[string]ResourceSetting)

	// Use rows.Next() to read the output row by row
	for rows.Next() {
		var name, setting, unit, EnumVals sql.NullString
		if err := rows.Scan(&name, &setting, &unit, &EnumVals); err != nil {
			logger.LogError(fmt.Errorf("Failed scanning row: %v", err))
			return nil, err
		}

		// Use regexp.MatchString to check if row contains one of the ResourceSettings
		for _, resource := range ResourceSettings {
			match, _ := regexp.MatchString(resource, name.String)
			if match {
				// If the line contains a ResourceSetting, store the name and setting in the map
				resSetting := ResourceSetting{
					Name:     name.String,     // name of the setting
					Value:    setting.String,  // value of the setting
					Unit:     unit.String,     // s, ms, kB, 8kB, etc...
					EnumVals: EnumVals.String, // If an enumrator, this stores enum values
				}

				settingsMap[name.String] = resSetting
				break
			}
		}
	}
	// Return map of runtime config settings
	return settingsMap, nil
}

func (conf *Configuration) getSpecificPGSetting(settingName string, logger *utils.Logger) (*ResourceSetting, error) {
	// Prepare the SQL statement
	formattedArg := fmt.Sprintf("SELECT name,setting,unit,enumvals FROM pg_settings WHERE name = '%s'", settingName)
	stmt, err := conf.dbHandler.Prepare(formattedArg)
	if err != nil {
		logger.LogError(fmt.Errorf("Failed preparing SQL statement: %v", err))
		return nil, err
	}
	defer stmt.Close()

	// Query the database
	rows, err := stmt.Query()
	if err != nil {
		logger.LogError(fmt.Errorf("Failed querying Postgres: %v", err))
		return nil, err
	}
	defer rows.Close()

	// Read and return first resulting row
	rows.Next()
	var name, setting, unit, EnumVals sql.NullString
	if err := rows.Scan(&name, &setting, &unit, &EnumVals); err != nil {
		logger.LogError(fmt.Errorf("Failed scanning row: %v", err))
		return nil, err
	}
	resSetting := ResourceSetting{
		Name:     name.String,     // name of the setting
		Value:    setting.String,  // value of the setting
		Unit:     unit.String,     // s, ms, kB, 8kB, etc...
		EnumVals: EnumVals.String, // If an enumrator, this stores enum values
	}

	return &resSetting, nil
}

// Function used to ensure that parameter value being set
// is amongst options returned by `postgres -c "select enumvals...`
// Helps avoid exposing incorrect setting suggestions to users.
func setEnumTypeSuggestedValue(setting *ResourceSetting, valueToSet string) error {
	// 1. Extract enumerator values from something like `{off,on,try}`
	start := strings.Index(setting.EnumVals, "{")
	end := strings.Index(setting.EnumVals, "}")
	if start == -1 || end == -1 {
		return fmt.Errorf("Could not extract enumerator values from enumvals")
	}

	EnumVals := strings.Split(setting.EnumVals[start+1:end], ",")

	// 2. Set value we're trying to set if it exists
	for _, enumValue := range EnumVals {
		if enumValue == valueToSet {
			setting.SuggestedValue = valueToSet
			return nil
		}
	}
	return fmt.Errorf("There is no %s in setting's %s enumerator %s", valueToSet, setting.Name, setting.EnumVals)
}

////////////////////////////////////////////////////////
// Configuration functions
////////////////////////////////////////////////////////

// checks what unit is used for shared_buffers, applies necessary conversions
// and sets final suggestion as a shared_buffers unit to closest value that's power of 2
func (conf *Configuration) CheckSharedBuffers(logger *utils.Logger) (*ResourceSetting, error) {
	var suggestion float32
	var lowestRecommendedValue float32 = 128 // Cannot suggest value that is lower than 128MB
	var GigabyteInBytes uint64 = 1073741824  // 1GB
	sharedBuffers := conf.settings["shared_buffers"]

	// 1. Get total server memory
	totalMemory, err := utils.GetTotalMemory()
	if err != nil {
		logger.LogError(fmt.Errorf("Failed shared_buffers check because could not get total server memory: %v", err))
		sharedBuffers.GotError = true
		return nil, err
	}

	// 2. Convert total server memory to a unit that's used by shared_buffers
	totalMemoryAsString := utils.Uint64ToString(totalMemory)
	totalMemoryConverted, err := utils.ConvertBasedOnUnit(totalMemoryAsString, "B", sharedBuffers.Unit)
	if err != nil {
		logger.LogError(fmt.Errorf("Failed shared_buffers check: %v", err))
		sharedBuffers.GotError = true
		return nil, err
	}

	// 3. Suggest value that's n% memory depending on how much total and free memory we have
	//
	// If total server memory > 1GB, suggest 25% of total server RAM
	if totalMemory > GigabyteInBytes {
		suggestion = totalMemoryConverted * 0.25
		sharedBuffers.Details = "Total server memory > 1GB. Suggest using 25% of total server RAM"
	} else { // Else suggest 30% of what memory is currently available on the server
		availableMemory, err := utils.GetAvailableMemory()
		if err != nil {
			logger.LogError(fmt.Errorf("Failed shared_buffers check because could not get total available memory: %v", err))
			sharedBuffers.GotError = true
			return nil, err
		}
		availableMemoryAsString := utils.Uint64ToString(availableMemory)
		availableMemoryConverted, err := utils.ConvertBasedOnUnit(availableMemoryAsString, "B", sharedBuffers.Unit)
		if err != nil {
			logger.LogError(fmt.Errorf("Failed shared_buffers check: %v", err))
			sharedBuffers.GotError = true
			return nil, err
		}
		suggestion = availableMemoryConverted * 0.30
		sharedBuffers.Details = "Total server memory < 1GB. Suggest using 30% of available RAM"

		// Suggested value cannot be lower than 128MB
		lowestRecommendedValueAsString := utils.Float32ToString(lowestRecommendedValue)
		lowestRecommendedValue, err = utils.ConvertBasedOnUnit(lowestRecommendedValueAsString, "MB", sharedBuffers.Unit)
		if err != nil {
			logger.LogError(fmt.Errorf("Failed shared_buffers check: %v", err))
			sharedBuffers.GotError = true
			return nil, err
		}
		if suggestion < lowestRecommendedValue {
			suggestion = lowestRecommendedValue
			sharedBuffers.Details = "Server lacks available memory. Lowest recommended value is 128MB"
		}
	}

	// Round suggestion to power of 2 and make sure there's no decimal point
	roundedSuggestion := utils.RoundToPowerOf2(uint64(suggestion))
	sharedBuffers.SuggestedValue = utils.Uint64ToString(roundedSuggestion)
	conf.settings["shared_buffers"] = sharedBuffers
	return &sharedBuffers, err
}

func (conf *Configuration) CheckHugePages(logger *utils.Logger) (*ResourceSetting, error) {
	hugePages := conf.settings["huge_pages"]

	// Get nr_hugepages value
	kernelPagesString, err := sysctl.Get("vm.nr_hugepages")
	if err != nil {
		logger.LogError(fmt.Errorf("Failed huge_pages check: %v", err))
		hugePages.GotError = true
		return nil, err
	}
	kernelNrHugePages, err := utils.StringToInt(kernelPagesString)

	// if it's not set in kernel, no point in having hugePages set to on/try
	if kernelNrHugePages == 0 {
		hugePages.Details = "Kernel parameter nr_hugepages is set to 0. Because of that, PostgreSQL cannot request huge pages"
		setEnumTypeSuggestedValue(&hugePages, "off")

		conf.settings["huge_pages"] = hugePages
		return &hugePages, nil
	}

	if hugePages.Value == "on" {
		hugePages.Details = "huge_pages current value is equal to 'on'. Failure to request huge pages will prevent the server from starting up"
		setEnumTypeSuggestedValue(&hugePages, "try")
	} else if hugePages.Value == "off" {
		hugePages.Details = "huge_pages current value is equal to 'off'. Use of huge pages results in smaller page tables and less CPU time spent on memory management, increasing performance"
		setEnumTypeSuggestedValue(&hugePages, "try")
	}

	// If kernelNrHugePages > 0 and hugePages.SuggestedValue = "try", make no suggestions

	conf.settings["huge_pages"] = hugePages
	return &hugePages, nil
}

func (conf *Configuration) CheckHugePageSize(logger *utils.Logger) (*ResourceSetting, error) {
	hugePageSize := conf.settings["huge_page_size"]

	// Get nr_hugepages value
	kernelPagesString, err := sysctl.Get("vm.nr_hugepages")
	if err != nil {
		logger.LogError(fmt.Errorf("Failed huge_page_size check: %v", err))
		hugePageSize.GotError = true
		return nil, err
	}
	kernelNrHugePages, err := strconv.Atoi(kernelPagesString)
	if err != nil {
		logger.LogError(fmt.Errorf("Failed huge_page_size check: %v", err))
		hugePageSize.GotError = true
		return nil, err
	}

	// If vm.nr_hugepages is 0, then huge_page_size cannot be set
	if kernelNrHugePages == 0 {
		return &hugePageSize, nil
	}

	if hugePageSize.Value != "0" {
		hugePageSize.SuggestedValue = "0"
		hugePageSize.Details = "Current huge_page_size value is set to a non 0 value. To prevent fragmentation, the same huge page size as the one set in your Linux kernel should be used. When set to 0, the default huge page size on the system will be used."
	}

	conf.settings["huge_page_size"] = hugePageSize
	return &hugePageSize, nil
}

// GENERALREC
func (conf *Configuration) CheckTempBuffers(logger *utils.Logger) (*ResourceSetting, error) {
	tempBuffers := conf.settings["temp_buffers"]

	currentValue, err := utils.StringToUint64(tempBuffers.Value)
	if err != nil {
		logger.LogError(fmt.Errorf("Failed temp_buffers check: %v", err))
		tempBuffers.GotError = true
		return nil, err
	}
	currentValueAsString := utils.Uint64ToString(currentValue)
	currentValueConverted, err := utils.ConvertBasedOnUnit(currentValueAsString, "8kB", tempBuffers.Unit)
	if err != nil {
		logger.LogError(fmt.Errorf("Failed temp_buffers check: %v", err))
		tempBuffers.GotError = true
		return nil, err
	}
	if currentValueConverted > 1024 { // 1024 8kB is equivalent to 8MB
		tempBuffers.Details = "Current temp_buffers value is set to more than 8MB. If there are multiple databases used by different application, consider changing this setting per database. It is recommended to increase this value only for applications that rely heavily on temporary tables"
	} else if currentValueConverted < 1024 {
		tempBuffers.Details = "Current temp_buffers value is set to less than 8MB. The cost of setting a large value in sessions that do not actually need many temporary buffers is only a buffer descriptor, or about 64 bytes. Consider increasing this to the recommended default value"

		convertedDefault, err := utils.ConvertBasedOnUnit("1024", "8kB", tempBuffers.Unit)
		if err != nil {
			logger.LogError(fmt.Errorf("Failed temp_buffers check: %v", err))
			tempBuffers.GotError = true
			return nil, err
		}
		tempBuffers.SuggestedValue = utils.Float32ToString(convertedDefault)
	}

	conf.settings["temp_buffers"] = tempBuffers
	return &tempBuffers, nil
}

// GENERALREC
func (conf *Configuration) CheckMaxPreparedTransactions(logger *utils.Logger) (*ResourceSetting, error) {
	maxPreparedTransactions := conf.settings["max_prepared_transactions"]

	// !!! consider passing this as an argument because there are currently two
	// separate functions that ask for max_connections
	maxConnections, err := conf.getSpecificPGSetting("max_connections", logger)
	if err != nil {
		logger.LogError(fmt.Errorf("Failed max_prepared_transactions check: %v", err))
		maxPreparedTransactions.GotError = true
		return nil, err
	}

	if maxPreparedTransactions.Value == "0" {
		maxPreparedTransactions.Details = " If you are using prepared transactions, you will probably want max_prepared_transactions to be at least as large as max_connections, so that every session can have a prepared transaction pending."
		maxPreparedTransactions.SuggestedValue = maxConnections.Value
	}

	conf.settings["max_prepared_transactions"] = maxPreparedTransactions
	return &maxPreparedTransactions, nil
}

func (conf *Configuration) CheckWorkMem(logger *utils.Logger) (*ResourceSetting, error) {
	workMem := conf.settings["work_mem"]

	// 1. Get amount of available memory and connections
	availableMemory, err := utils.GetAvailableMemory()
	if err != nil {
		logger.LogError(fmt.Errorf("Failed work_mem check: %v", err))
		workMem.GotError = true
		return nil, err
	}
	maxConnections, err := conf.getSpecificPGSetting("max_connections", logger)
	if err != nil {
		logger.LogError(fmt.Errorf("Failed work_mem check: %v", err))
		workMem.GotError = true
		return nil, err
	}

	// 2. suggestion = availablememory / max_connections
	maxConnectionsValue, err := utils.StringToUint64(maxConnections.Value)
	if err != nil {
		logger.LogError(fmt.Errorf("Failed work_mem check: %v", err))
		workMem.GotError = true
		return nil, err
	}
	suggestion := utils.Uint64ToString(utils.RoundToPowerOf2(availableMemory / maxConnectionsValue))
	suggestionAsWorkMemUnit, err := utils.ConvertBasedOnUnit(suggestion, "B", workMem.Unit)
	workMem.SuggestedValue = utils.Float32ToString(suggestionAsWorkMemUnit)

	// 3. Add details for decision
	workMem.Details = "Suggested value is based on currently available memory on the server divided by max_connections. If using complex queries that involve sorts or hash tables, consider using double this value. It can also be set higher if this server is a dedicated database server and there is no concern that other software will run out of memory."

	conf.settings["work_mem"] = workMem
	return &workMem, nil
}

// GENERALREC
func (conf *Configuration) CheckHashMemMultiplier(logger *utils.Logger) (*ResourceSetting, error) {
	hashMemMultiplier := conf.settings["hash_mem_multiplier"]
	workMem := conf.settings["work_mem"]
	workMemValueAsMB, err := utils.ConvertBasedOnUnit(workMem.Value, workMem.Unit, "MB")
	if err != nil {
		logger.LogError(fmt.Errorf("Failed hash_mem_multiplier check: %v", err))
		hashMemMultiplier.GotError = true
		return nil, err
	}

	// If more than 40MB
	if workMemValueAsMB > 40 {
		hashMemMultiplierAsFloat32, err := utils.StringToFloat32(hashMemMultiplier.Value)
		if err != nil {
			logger.LogError(fmt.Errorf("Failed hash_mem_multiplier check: %v", err))
			hashMemMultiplier.GotError = true
			return nil, err
		}
		suggestion := (hashMemMultiplierAsFloat32 + workMemValueAsMB*0.01)
		if suggestion > 8 {
			suggestion = 8
		}
		hashMemMultiplier.Details = "Generally default value works best. If your application uses hash-based operations and PostgreSQL often ends up spilling (creates workfiles on disk to compensate for lack of memory), consider increasing this further. Suggested value is based on how much working memory is currently set."
		hashMemMultiplier.SuggestedValue = utils.Float32ToString(suggestion)
	} else if hashMemMultiplier.Value != "2" {
		hashMemMultiplier.Details = "Generally default value works best. If your application uses hash-based operations and PostgreSQL often ends up spilling (creates workfiles on disk to compensate for lack of memory), consider increasing this after having increased work_mem above 40MB."
		hashMemMultiplier.SuggestedValue = "2"
	}

	conf.settings["hash_mem_multiplier"] = hashMemMultiplier
	return &hashMemMultiplier, nil
}

func (conf *Configuration) CheckMaintenanceWorkMem(logger *utils.Logger) (*ResourceSetting, error) {
	// 1. Get maintenance_work_mem, autovacumm_max_workers and available memory on server
	maintenanceWorkMem := conf.settings["maintenance_work_mem"]
	autovacuumMaxWorkers, availableMem, err := conf.getWorkMemRelatedValues(logger, &maintenanceWorkMem)
	if err != nil {
		logger.LogError(fmt.Errorf("failed maintenance_work_mem check: %v", err))
		maintenanceWorkMem.GotError = true
		return nil, err
	}

	// 2. Divide available memory by 8 * autovacuum_max_workers and round to nearest power of 2
	suggestion := availableMem / 8 / autovacuumMaxWorkers
	suggestionRounded := utils.RoundToPowerOf2(uint64(suggestion))

	maintenanceWorkMem.Details = fmt.Sprintf("This suggestion was made by dividing current available memory(%.2f%s) by 8 and by how many autovacuum_max_workers(%.0f) are set. Applications that heavily rely on maintenance operations, such as VACUUM, CREATE INDEX, and ALTER TABLE ADD FOREIGN KEY may want to increase this further by multiplying the suggested value by 2", availableMem, maintenanceWorkMem.Unit, autovacuumMaxWorkers)
	maintenanceWorkMem.SuggestedValue = utils.Uint64ToString(suggestionRounded)

	// 3. If suggestion is below default value and current value is not equal to default,
	// suggest default value instead.

	suggestionAsMB, err := utils.ConvertBasedOnUnit(utils.Uint64ToString(suggestionRounded), maintenanceWorkMem.Unit, "MB")
	if err != nil {
		logger.LogError(fmt.Errorf("Failed maintenance_work_mem check: %v", err))
		maintenanceWorkMem.GotError = true
		return nil, err
	}

	if suggestionAsMB < 64 {
		maintenanceWorkMemAsMB, err := utils.ConvertBasedOnUnit(maintenanceWorkMem.Value, maintenanceWorkMem.Unit, "MB")
		if err != nil {
			logger.LogError(fmt.Errorf("Failed maintenance_work_mem check: %v", err))
			maintenanceWorkMem.GotError = true
			return nil, err
		}
		if maintenanceWorkMemAsMB != 64 {
			maintenanceWorkMem.Details = "Currently there is not enough available memory on the server to go above default maintenance_work_mem value"
			defaultAsUnit, err := utils.ConvertBasedOnUnit("64", "MB", maintenanceWorkMem.Unit)
			if err != nil {
				logger.LogError(fmt.Errorf("Failed maintenance_work_mem check: %v", err))
				maintenanceWorkMem.GotError = true
				return nil, err
			}
			maintenanceWorkMem.SuggestedValue = utils.Float32ToString(defaultAsUnit)
		} else { // if suggestion was below default and current value is already set to default
			maintenanceWorkMem.Details = ""
			maintenanceWorkMem.SuggestedValue = ""
		}
	}

	conf.settings["maintenance_work_mem"] = maintenanceWorkMem
	return &maintenanceWorkMem, nil
}

func (conf *Configuration) CheckAutovacuumWorkMem(logger *utils.Logger) (*ResourceSetting, error) {
	// 1. Get autovacuum_work_mem, autovacumm_max_workers and available memory on server
	autovacuumWorkMem := conf.settings["autovacuum_work_mem"]
	autovacuumMaxWorkers, availableMem, err := conf.getWorkMemRelatedValues(logger, &autovacuumWorkMem)
	if err != nil {
		logger.LogError(fmt.Errorf("Failed autovacuum_work_mem check: %v", err))
		autovacuumWorkMem.GotError = true
		return nil, err
	}

	// 2. Divide available memory by 8 * autovacuum_max_workers and round to nearest power of 2
	suggestion := availableMem / 4 / autovacuumMaxWorkers
	suggestionRounded := utils.RoundToPowerOf2(uint64(suggestion))

	autovacuumWorkMem.Details = fmt.Sprintf("This suggestion was made by dividing current available memory(%.2f%s) by 4 and by how many autovacuum_max_workers(%.0f) are set", availableMem, autovacuumWorkMem.Unit, autovacuumMaxWorkers)
	autovacuumWorkMem.SuggestedValue = utils.Uint64ToString(suggestionRounded)

	// 3. If suggestion is below default of maintenance_work_mem value
	// and current value is not already -1, suggest -1 instead.

	suggestionAsMB, err := utils.ConvertBasedOnUnit(utils.Uint64ToString(suggestionRounded), autovacuumWorkMem.Unit, "MB")
	if err != nil {
		logger.LogError(fmt.Errorf("Failed autovacuum_work_mem check: %v", err))
		autovacuumWorkMem.GotError = true
		return nil, err
	}

	if suggestionAsMB < 64 && autovacuumWorkMem.Value != "-1" {
		defaultAsUnit, err := utils.ConvertBasedOnUnit("64", "MB", autovacuumWorkMem.Unit)
		if err != nil {
			logger.LogError(fmt.Errorf("Failed autovacuum_work_mem check: %v", err))
			autovacuumWorkMem.GotError = true
			return nil, err
		}
		autovacuumWorkMem.Details = "Currently there is not enough available memory on the server to go above default autovacuum_work_mem value. Suggesting setting to -1 to rely on maintenance_work_mem instead"
		autovacuumWorkMem.SuggestedValue = utils.Float32ToString(defaultAsUnit)
	} else if suggestionAsMB < 64 { // if already set to -1, don't suggest anything
		autovacuumWorkMem.Details = ""
		autovacuumWorkMem.SuggestedValue = ""
	}

	conf.settings["autovacuum_work_mem"] = autovacuumWorkMem
	return &autovacuumWorkMem, nil
}

// Function for getting `autovacuum_max_workers` and available memory on server
// to later be used in checks for `maintenance_work_mem` and `autovacuum_work_mem`
func (conf *Configuration) getWorkMemRelatedValues(logger *utils.Logger, setting *ResourceSetting) (float32, float32, error) {
	// 1. Get autovacuum_max_workers
	tmpMaxWorkers, err := conf.getSpecificPGSetting("autovacuum_max_workers", logger)
	if err != nil {
		logger.LogError(fmt.Errorf("failed getting autovacuum_max_workers: %v", err))
		tmpMaxWorkers.GotError = true
		return 0, 0, err
	}
	autovacuumMaxWorkers, err := utils.StringToFloat32(tmpMaxWorkers.Value)
	if err != nil {
		logger.LogError(fmt.Errorf("failed getting autovacuum_max_workers: %v", err))
		tmpMaxWorkers.GotError = true
		return 0, 0, err
	}

	// 1. Get available memory on server
	availableMemory, err := utils.GetAvailableMemory()
	if err != nil {
		logger.LogError(fmt.Errorf("failed getting available memory on server check: %v", err))
		tmpMaxWorkers.GotError = true
		return 0, 0, err
	}
	availableMemoryAsString := utils.Uint64ToString(availableMemory)
	availableMemoryConverted, err := utils.ConvertBasedOnUnit(availableMemoryAsString, "B", setting.Unit)
	if err != nil {
		logger.LogError(fmt.Errorf("failed getting available memory on server check: %v", err))
		tmpMaxWorkers.GotError = true
		return 0, 0, err
	}

	return autovacuumMaxWorkers, availableMemoryConverted, nil
}

func (conf *Configuration) ChecklogicalDecodingWorkMem(logger *utils.Logger) (*ResourceSetting, error) {
	logicalDecodingWorkMem := conf.settings["logical_decoding_work_mem"]

	// 1. Get available memory on server
	availableMemory, err := utils.GetAvailableMemory()
	if err != nil {
		logger.LogError(fmt.Errorf("failed logical_decoding_work_mem check: %v", err))
		logicalDecodingWorkMem.GotError = true
		return nil, err
	}
	availableMemoryAsString := utils.Uint64ToString(availableMemory)
	availableMemoryConverted, err := utils.ConvertBasedOnUnit(availableMemoryAsString, "B", logicalDecodingWorkMem.Unit)
	if err != nil {
		logger.LogError(fmt.Errorf("failed logical_decoding_work_mem check: %v", err))
		logicalDecodingWorkMem.GotError = true
		return nil, err
	}

	// 2. make suggestion by dividing available memory by 8 and rounding to nearest power of 2
	suggestion := availableMemoryConverted / 8
	suggestionRounded := utils.RoundToPowerOf2(uint64(suggestion))

	// 3. If suggested value is less than 64MB, make no suggestion
	suggestionAsMB, err := utils.ConvertBasedOnUnit(utils.Uint64ToString(suggestionRounded), logicalDecodingWorkMem.Unit, "MB")
	if err != nil {
		logger.LogError(fmt.Errorf("failed logicalDecodingWorkMem check: %v", err))
		logicalDecodingWorkMem.GotError = true
		return nil, err
	}

	if !(suggestionAsMB < 64) {
		logicalDecodingWorkMem.Details = fmt.Sprintf("This suggestion was made by dividing current available memory(%.2f) by 8", availableMemoryConverted)
		logicalDecodingWorkMem.SuggestedValue = utils.Uint64ToString(suggestionRounded)
	}

	conf.settings["logical_decoding_work_mem"] = logicalDecodingWorkMem
	return &logicalDecodingWorkMem, nil
}

func (conf *Configuration) CheckMaxStackDepth(logger *utils.Logger) (*ResourceSetting, error) {
	maxStackDepth := conf.settings["max_stack_depth"]

	// 1. Get system stack depth
	systemStackDepth, err := utils.GetStackSize()
	if err != nil {
		logger.LogError(fmt.Errorf("failed max_stack_depth check: %v", err))
		maxStackDepth.GotError = true
		return nil, err
	}

	// 2. System set stack depth is not equal to max_stack_depth, suggest system stack depth
	systemStackDepthAsUnit, err := utils.ConvertBasedOnUnit(utils.Uint64ToString(systemStackDepth), "B", maxStackDepth.Unit)
	if err != nil {
		logger.LogError(fmt.Errorf("failed max_stack_depth check: %v", err))
		maxStackDepth.GotError = true
		return nil, err
	}

	maxStackDepthFloat32, err := utils.StringToFloat32(maxStackDepth.Value)
	if err != nil {
		logger.LogError(fmt.Errorf("failed max_stack_depth check: %v", err))
		maxStackDepth.GotError = true
		return nil, err
	}

	if systemStackDepthAsUnit != maxStackDepthFloat32 {
		maxStackDepth.SuggestedValue = utils.Float32ToString(systemStackDepthAsUnit)
		maxStackDepth.Details = "The ideal setting for this parameter is the actual stack size limit enforced by the kernel (as set by ulimit -s or local equivalent)."
	}

	conf.settings["max_stack_depth"] = maxStackDepth
	return &maxStackDepth, nil
}

func (conf *Configuration) CheckSharedMemoryType(logger *utils.Logger) (*ResourceSetting, error) {
	sharedMemoryType := conf.settings["shared_memory_type"]

	// Suggest boot_val (default) value

	details := "sysv option is discouraged because it typically requires non-default kernel settings to allow for large allocations"
	err := suggestDefault(&sharedMemoryType, details)
	if err != nil {
		logger.LogError(fmt.Errorf("failed shared_memory_type check: %v", err))
		sharedMemoryType.GotError = true
		return nil, err
	}

	conf.settings["shared_memory_type"] = sharedMemoryType
	return &sharedMemoryType, nil
}

func (conf *Configuration) CheckDynamicSharedMemoryType(logger *utils.Logger) (*ResourceSetting, error) {
	dynamicSharedMemoryType := conf.settings["dynamic_shared_memory_type"]

	// Suggest boot_val (default) value
	details := "Typically default value is best for this option"
	err := suggestDefault(&dynamicSharedMemoryType, details)
	if err != nil {
		logger.LogError(fmt.Errorf("failed dynamic_shared_memory_type check: %v", err))
		dynamicSharedMemoryType.GotError = true
		return nil, err
	}

	conf.settings["dynamic_shared_memory_type"] = dynamicSharedMemoryType
	return &dynamicSharedMemoryType, nil
}

func suggestDefault(setting *ResourceSetting, details string) error {
	// 1. Get boot_val (default value) for setting
	formattedArg := fmt.Sprintf("select name,boot_val from pg_settings WHERE name = '%s'", setting.Name)
	cmd := exec.Command("psql", "-U", "postgres", "-c", formattedArg)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	var bootVal string
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		match, _ := regexp.MatchString(setting.Name, line)
		if match {
			fields := strings.Split(line, "|")
			bootVal = strings.TrimSpace(fields[1])
		}
	}

	// 2. Suggest boot_val (default) if not already set to that
	if setting.Value != bootVal {
		setting.SuggestedValue = bootVal
		setting.Details = details
	}

	return nil
}

////////////////////////////////////////////////////////////////////
// Below are functions used or reset suggestions
////////////////////////////////////////////////////////////////////

// Set suggested parameter in postgresql.auto.conf. ALTER SYSTEM SET cannot specify a unit.
func (conf *Configuration) setSuggestion(db *sql.DB, paramName string, paramValue string, logger *utils.Logger) error {
	_, err := db.Exec(fmt.Sprintf("ALTER SYSTEM SET %s = '%s'", paramName, paramValue))
	if err != nil {
		logger.LogError(fmt.Errorf("failed to apply suggestion for %s: %v", paramName, err))
	}
	return err
}

/*
Meant for applying a single suggestion.
@path - full path to postgresql.auto.conf
@suggestion - setting to apply suggestion on
*/
func (conf *Configuration) ApplySuggestions(suggestions *PatchResourceConfigsJSONBody, logger *utils.Logger) error {
	// 1. Create a backup of postgresql.auto.conf
	if err := utils.BackupFile(conf.autoConfPath, conf.backupDir, conf.appUser, logger); err != nil {
		return err
	}

	// 2. Execute ALTER SYSTEM to apply all suggestions
	gotError := false
	for _, suggestion := range *suggestions {
		err := conf.setSuggestion(conf.dbHandler, suggestion.Name, suggestion.SuggestedValue, logger)
		if err != nil {
			gotError = true
		}
	}
	// 3. Reload the configuration file to apply the changes
	utils.ReloadConfiguration(conf.dbHandler, logger)

	if gotError {
		return fmt.Errorf("One or more suggestion could not be applied")
	}
	return nil
}

// Removes all content inside postgresql.auto.conf and reloads configuration
func (conf *Configuration) DiscardConfigs(logger *utils.Logger) error {
	// 1. Wipe postgresql.auto.conf content
	_, err := conf.dbHandler.Exec("ALTER SYSTEM RESET ALL")
	if err != nil {
		logger.LogError(fmt.Errorf("failed to reset postgresql.auto.conf content: %v", err))
	}

	// 2. Reload configuration files
	err = utils.ReloadConfiguration(conf.dbHandler, logger)
	return err
}
