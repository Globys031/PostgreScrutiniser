// Code for PostgreSQL Configuration
package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	sysctl "github.com/lorenzosaino/go-sysctl"

	"github.com/Globys031/PostgreScrutiniser/backend/utils"
)

/*
TO DO:
Currently I've modified /var/lib/pgsql/15/data/pg_hba.conf
to allow all connections without any password on localhost.
This is bad, should reenable peer authentification later and have code
take that into account:

[root@localhost backend]# grep trust /var/lib/pgsql/15/data/pg_hba.conf
*/

// TO DO:
// also return documentation for every setting (probably best to do on the frontend side actually)

/* TO DO:
Tam tikri grazins tik general recommendations (ties kurie virs funkcijos) pamarkinti "GENERALREC"
Pasvarstyt kaip front'e idet... Gal padaryt tris boxes kur
auksciausiai yra: "found %{n} appliable suggestions"
zemiau: General recommendations %{n}
zemiausiai: Passed without suggestions %{n}
Prie "General recommendations" galima pridet ui-info, kad "this group contains recommendations
that could not be fully decided based on your system parameters (for example, because PostgreScrutiniser does not know what queries your applications are typically running). Suggestions here should be applied at your own discretion based on details provided"
*/

/* TO DO:
Add somewhere in the GUI a prompt that:
1. This app assumes that block_size is the default value of 8192 bytes as defined by `block_size` and described here: https://pgpedia.info/b/block_size.html#:~:text=The%20default%20value%20for%20block_size,(PostgreSQL%208.4%20and%20later).
2. This app assumes that all integer type settings have unit values specified
*/

type ResourceSetting struct {
	// TO DO:Consider having separate maps for differnent vartypes
	// TO DO: consider removing minVal, maxVal, vartype
	Name     string // name of the setting
	Value    string // value of the setting
	Vartype  string // what type value is (boolean, integer, enum, etc...)
	Unit     string // s, ms, kB, 8kB, etc...
	MinVal   string // Minimum allowed value (needed for validation)
	MaxVal   string // Maximum allowed value (needed for validation)
	EnumVals string // If an enumrator, this stores enum values. For internal use only. Not exposted by API

	SuggestedValue string // Value that will be suggested after running check
	Details        string // Details informing why a value was suggested
	// Specifies if applying suggestion requires rebooting postgresql server
	// If set to false, but a SuggestedValue has been applied, will do
	// `pg_ctl reload` automatically
	// https://www.postgresql.org/docs/15/app-pg-ctl.html
	// TO DO: check to confirm docs are referring to restarting postgresql, not the server itself
	RequiresRestart bool
	GotError        bool // specifies whether check got an error
}

type Configuration struct {
	// postgreUser string
	// postgreUserPass string
	path     string // Path to postgresql.conf
	settings map[string]ResourceSetting
	// settings [34]string
}

// Meant for initialising Configuration upon first api call so that same
// reference can be reused for later calls.
func InitChecks(logger *utils.Logger) *Configuration {
	// user, pass := findPostgresCredentials();

	// Find path to postgresql.conf and get postgresql settings
	// Error logging handled inside the functions
	configFile, _ := findConfigFile(logger)
	ResourceSettings, _ := getPGSettings(logger)
	conf := Configuration{path: configFile, settings: ResourceSettings}

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

// Function for finding main postgresql user credentials
// to later be used for peer authentification
func findPostgresCredentials() (string, error) {
	panic("findPostgresCredentials not implemented")
}

// Returns path to postgresql.conf if it exists.
func findConfigFile(logger *utils.Logger) (string, error) {
	output, err := exec.Command("psql", "-U", "postgres", "-c", "SHOW config_file").Output()
	if err != nil {
		logger.LogError("Failed finding postgresql.conf: " + err.Error())
		return "", err
	}

	var filePath string
	scanner := bufio.NewScanner(bytes.NewReader(output))
	// Return line that has postgresql.conf in it
	for scanner.Scan() {
		line := scanner.Text()
		if match, _ := regexp.MatchString("postgresql.conf", line); match {
			filePath = strings.TrimSpace(line)
			break
		}
	}

	// check to confirm postgresql.conf file exists
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		err = fmt.Errorf("postgresql.conf was found with `psql -U` but could not be opened: %v", err)
		logger.LogError("Failed finding postgresql.conf: " + err.Error())
		return "", err
	}

	if err := scanner.Err(); err != nil {
		logger.LogError("Failed finding postgresql.conf: " + err.Error())
		return "", err
	}

	return filePath, nil
}

// Stores data returned by `pg_settings` into a map
// only for settings we're interested in.
func getPGSettings(logger *utils.Logger) (map[string]ResourceSetting, error) {
	// Settings we're interested in
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
	// ResourceSettings := [33]string{"shared_buffers",
	// 	"huge_pages",
	// 	"huge_page_size",
	// 	"temp_buffers",
	// 	"max_prepared_transactions",
	// 	"work_mem",
	// 	"hash_mem_multiplier",
	// 	"maintenance_work_mem",
	// 	"autovacuum_work_mem",
	// 	"logical_decoding_work_mem",
	// 	"max_stack_depth",
	// 	"shared_memory_type",
	// 	"dynamic_shared_memory_type",
	// 	"min_dynamic_shared_memory",
	// 	"temp_file_limit",
	// 	"vacuum_cost_delay",
	// 	"vacuum_cost_page_hit",
	// 	"vacuum_cost_page_miss",
	// 	"vacuum_cost_page_dirty",
	// 	"vacuum_cost_limit",
	// 	"bgwriter_delay",
	// 	"bgwriter_lru_maxpages",
	// 	"bgwriter_lru_multiplier",
	// 	"bgwriter_flush_after",
	// 	"backend_flush_after",
	// 	"effective_io_concurrency",
	// 	"maintenance_io_concurrency",
	// 	"max_worker_processes",
	// 	"max_parallel_workers_per_gather",
	// 	"max_parallel_maintenance_workers",
	// 	"max_parallel_workers",
	// 	"parallel_leader_participation",
	// 	"old_snapshot_threshold"}

	// Use exec.Command to run the `psql` command and capture the output
	cmd := exec.Command("psql", "-U", "postgres", "-c", "select name,setting,vartype,unit,min_val,max_val,enumvals from pg_settings")
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.LogError("Failed fetching PostgreSql settings: " + err.Error())
		return nil, err
	}

	// Initialize the map to store the results
	settingsMap := make(map[string]ResourceSetting)

	// Use bufio.NewScanner to read the output line by line
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()

		// Use regexp.MatchString to check if the line contains one of the ResourceSettings
		for _, resource := range ResourceSettings {
			match, _ := regexp.MatchString(resource, line)
			if match {
				// If the line contains a ResourceSetting, split it by "|"
				// and store the name and setting in the map
				fields := strings.Split(line, "|")
				name := strings.TrimSpace(fields[0])
				value := strings.TrimSpace(fields[1])
				vartype := strings.TrimSpace(fields[2])
				unit := strings.TrimSpace(fields[3])
				minVal := strings.TrimSpace(fields[4])
				maxVal := strings.TrimSpace(fields[5])
				EnumVals := strings.TrimSpace(fields[6])

				resSetting := ResourceSetting{
					Name:     name,     // name of the setting
					Value:    value,    // value of the setting
					Vartype:  vartype,  // what type value is (boolean, integer, enum, etc...)
					Unit:     unit,     // s, ms, kB, 8kB, etc...
					MinVal:   minVal,   // Minimum allowed value (needed for validation)
					MaxVal:   maxVal,   // Maximum allowed value (needed for validation)
					EnumVals: EnumVals, // If an enumrator, this stores enum values
				}

				settingsMap[name] = resSetting
				break
			}
		}
	}

	// Return map of runtime config settings
	return settingsMap, nil
}

func getSpecificPGSetting(setting string) (*ResourceSetting, error) {
	formattedArg := fmt.Sprintf("select name,setting,vartype,unit,min_val,max_val,enumvals from pg_settings WHERE name = '%s'", setting)

	// Execute command to fetch postgre setting
	cmd := exec.Command("psql", "-U", "postgres", "-c", formattedArg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	// This will store the final result of the setting to be restored
	var returnSetting ResourceSetting

	// Use bufio.NewScanner to read the output line by line until we get the line with setting
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()

		match, _ := regexp.MatchString(setting, line)
		if match {
			fields := strings.Split(line, "|")
			name := strings.TrimSpace(fields[0])
			value := strings.TrimSpace(fields[1])
			vartype := strings.TrimSpace(fields[2])
			unit := strings.TrimSpace(fields[3])
			minVal := strings.TrimSpace(fields[4])
			maxVal := strings.TrimSpace(fields[5])
			EnumVals := strings.TrimSpace(fields[6])

			returnSetting = ResourceSetting{
				Name:     name,     // name of the setting
				Value:    value,    // value of the setting
				Vartype:  vartype,  // what type value is (boolean, integer, enum, etc...)
				Unit:     unit,     // s, ms, kB, 8kB, etc...
				MinVal:   minVal,   // Minimum allowed value (needed for validation)
				MaxVal:   maxVal,   // Maximum allowed value (needed for validation)
				EnumVals: EnumVals, // If an enumrator, this stores enum values
			}
		}
	}

	return &returnSetting, nil
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
		logger.LogError("Failed shared_buffers check because could not get total server memory: " + err.Error())
		sharedBuffers.GotError = true
		return nil, err
	}

	// 2. Convert total server memory to a unit that's used by shared_buffers
	totalMemoryAsString := utils.Uint64ToString(totalMemory)
	totalMemoryConverted, err := utils.ConvertBasedOnUnit(totalMemoryAsString, "B", sharedBuffers.Unit)
	if err != nil {
		logger.LogError("Failed shared_buffers check: " + err.Error())
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
			logger.LogError("Failed shared_buffers check because could not get total available memory: " + err.Error())
			sharedBuffers.GotError = true
			return nil, err
		}
		availableMemoryAsString := utils.Uint64ToString(availableMemory)
		availableMemoryConverted, err := utils.ConvertBasedOnUnit(availableMemoryAsString, "B", sharedBuffers.Unit)
		if err != nil {
			logger.LogError("Failed shared_buffers check: " + err.Error())
			sharedBuffers.GotError = true
			return nil, err
		}
		suggestion = availableMemoryConverted * 0.30
		sharedBuffers.Details = "Total server memory < 1GB. Suggest using 30% of available RAM"

		// Suggested value cannot be lower than 128MB
		lowestRecommendedValueAsString := utils.Float32ToString(lowestRecommendedValue)
		lowestRecommendedValue, err = utils.ConvertBasedOnUnit(lowestRecommendedValueAsString, "MB", sharedBuffers.Unit)
		if err != nil {
			logger.LogError("Failed shared_buffers check: " + err.Error())
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
	// https://www.postgresql.org/docs/current/runtime-config-resource.html
	// "This parameter can only be set at server start."
	sharedBuffers.RequiresRestart = true
	conf.settings["shared_buffers"] = sharedBuffers
	return &sharedBuffers, err
}

func (conf *Configuration) CheckHugePages(logger *utils.Logger) (*ResourceSetting, error) {
	hugePages := conf.settings["huge_pages"]

	// Get nr_hugepages value
	kernelPagesString, err := sysctl.Get("vm.nr_hugepages")
	if err != nil {
		logger.LogError("Failed huge_pages check: " + err.Error())
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
		logger.LogError("Failed huge_page_size check: " + err.Error())
		hugePageSize.GotError = true
		return nil, err
	}
	kernelNrHugePages, err := strconv.Atoi(kernelPagesString)
	if err != nil {
		logger.LogError("Failed huge_page_size check: " + err.Error())
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

	hugePageSize.RequiresRestart = true
	conf.settings["huge_page_size"] = hugePageSize
	return &hugePageSize, nil
}

// GENERALREC
func (conf *Configuration) CheckTempBuffers(logger *utils.Logger) (*ResourceSetting, error) {
	tempBuffers := conf.settings["temp_buffers"]

	currentValue, err := utils.StringToUint64(tempBuffers.Value)
	if err != nil {
		logger.LogError("Failed temp_buffers check: " + err.Error())
		tempBuffers.GotError = true
		return nil, err
	}
	currentValueAsString := utils.Uint64ToString(currentValue)
	currentValueConverted, err := utils.ConvertBasedOnUnit(currentValueAsString, "8kB", tempBuffers.Unit)
	if err != nil {
		logger.LogError("Failed temp_buffers check: " + err.Error())
		tempBuffers.GotError = true
		return nil, err
	}
	if currentValueConverted > 1024 { // 1024 8kB is equivalent to 8MB
		tempBuffers.Details = "Current temp_buffers value is set to more than 8MB. If there are multiple databases used by different application, consider changing this setting per database. It is recommended to increase this value only for applications that rely heavily on temporary tables"
	} else if currentValueConverted < 1024 {
		tempBuffers.Details = "Current temp_buffers value is set to less than 8MB. The cost of setting a large value in sessions that do not actually need many temporary buffers is only a buffer descriptor, or about 64 bytes. Consider increasing this to the recommended default value"

		convertedDefault, err := utils.ConvertBasedOnUnit("1024", "8kB", tempBuffers.Unit)
		if err != nil {
			logger.LogError("Failed temp_buffers check: " + err.Error())
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
	maxConnections, err := getSpecificPGSetting("max_connections")
	if err != nil {
		logger.LogError("Failed max_prepared_transactions check: " + err.Error())
		maxPreparedTransactions.GotError = true
		return nil, err
	}

	if maxPreparedTransactions.Value == "0" {
		maxPreparedTransactions.Details = " If you are using prepared transactions, you will probably want max_prepared_transactions to be at least as large as max_connections, so that every session can have a prepared transaction pending."
		maxPreparedTransactions.SuggestedValue = maxConnections.Value
	}

	maxPreparedTransactions.RequiresRestart = true
	conf.settings["max_prepared_transactions"] = maxPreparedTransactions
	return &maxPreparedTransactions, nil
}

func (conf *Configuration) CheckWorkMem(logger *utils.Logger) (*ResourceSetting, error) {
	workMem := conf.settings["work_mem"]

	// 1. Get amount of available memory and connections
	availableMemory, err := utils.GetAvailableMemory()
	if err != nil {
		logger.LogError("Failed work_mem check: " + err.Error())
		workMem.GotError = true
		return nil, err
	}
	maxConnections, err := getSpecificPGSetting("max_connections")
	if err != nil {
		logger.LogError("Failed work_mem check: " + err.Error())
		workMem.GotError = true
		return nil, err
	}

	// 2. suggestion = availablememory / max_connections
	maxConnectionsValue, err := utils.StringToUint64(maxConnections.Value)
	if err != nil {
		logger.LogError("Failed work_mem check: " + err.Error())
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
		logger.LogError("Failed hash_mem_multiplier check: " + err.Error())
		hashMemMultiplier.GotError = true
		return nil, err
	}

	// If more than 40MB
	if workMemValueAsMB > 40 {
		hashMemMultiplierAsFloat32, err := utils.StringToFloat32(hashMemMultiplier.Value)
		if err != nil {
			logger.LogError("Failed hash_mem_multiplier check: " + err.Error())
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
		fmt.Println("ieina")
		hashMemMultiplier.Details = "Generally default value works best. If your application uses hash-based operations and PostgreSQL often ends up spilling (creates workfiles on disk to compensate for lack of memory), consider increasing this after having increased work_mem above 40MB."
		hashMemMultiplier.SuggestedValue = "2"
	}

	conf.settings["hash_mem_multiplier"] = hashMemMultiplier
	return &hashMemMultiplier, nil
}

func (conf *Configuration) CheckMaintenanceWorkMem(logger *utils.Logger) (*ResourceSetting, error) {
	// 1. Get maintenance_work_mem, autovacumm_max_workers and available memory on server
	maintenanceWorkMem := conf.settings["maintenance_work_mem"]
	autovacuumMaxWorkers, availableMem, err := getWorkMemRelatedValues(logger, &maintenanceWorkMem)
	if err != nil {
		logger.LogError("failed maintenance_work_mem check: " + err.Error())
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
		logger.LogError("Failed maintenance_work_mem check: " + err.Error())
		maintenanceWorkMem.GotError = true
		return nil, err
	}

	if suggestionAsMB < 64 {
		maintenanceWorkMemAsMB, err := utils.ConvertBasedOnUnit(maintenanceWorkMem.Value, maintenanceWorkMem.Unit, "MB")
		if err != nil {
			logger.LogError("Failed maintenance_work_mem check: " + err.Error())
			maintenanceWorkMem.GotError = true
			return nil, err
		}
		if maintenanceWorkMemAsMB != 64 {
			maintenanceWorkMem.Details = "Currently there is not enough available memory on the server to go above default maintenance_work_mem value"
			defaultAsUnit, err := utils.ConvertBasedOnUnit("64", "MB", maintenanceWorkMem.Unit)
			if err != nil {
				logger.LogError("Failed maintenance_work_mem check: " + err.Error())
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
	autovacuumMaxWorkers, availableMem, err := getWorkMemRelatedValues(logger, &autovacuumWorkMem)
	if err != nil {
		logger.LogError("Failed autovacuum_work_mem check: " + err.Error())
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
		logger.LogError("Failed autovacuum_work_mem check: " + err.Error())
		autovacuumWorkMem.GotError = true
		return nil, err
	}

	if suggestionAsMB < 64 && autovacuumWorkMem.Value != "-1" {
		defaultAsUnit, err := utils.ConvertBasedOnUnit("64", "MB", autovacuumWorkMem.Unit)
		if err != nil {
			logger.LogError("Failed autovacuum_work_mem check: " + err.Error())
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
func getWorkMemRelatedValues(logger *utils.Logger, setting *ResourceSetting) (float32, float32, error) {
	// 1. Get autovacuum_max_workers
	tmpMaxWorkers, err := getSpecificPGSetting("autovacuum_max_workers")
	if err != nil {
		logger.LogError("failed getting autovacuum_max_workers: " + err.Error())
		tmpMaxWorkers.GotError = true
		return 0, 0, err
	}
	autovacuumMaxWorkers, err := utils.StringToFloat32(tmpMaxWorkers.Value)
	if err != nil {
		logger.LogError("failed getting autovacuum_max_workers: " + err.Error())
		tmpMaxWorkers.GotError = true
		return 0, 0, err
	}

	// 1. Get available memory on server
	availableMemory, err := utils.GetAvailableMemory()
	if err != nil {
		logger.LogError("failed getting available memory on server check: " + err.Error())
		tmpMaxWorkers.GotError = true
		return 0, 0, err
	}
	availableMemoryAsString := utils.Uint64ToString(availableMemory)
	availableMemoryConverted, err := utils.ConvertBasedOnUnit(availableMemoryAsString, "B", setting.Unit)
	if err != nil {
		logger.LogError("failed getting available memory on server check: " + err.Error())
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
		logger.LogError("failed logical_decoding_work_mem check: " + err.Error())
		logicalDecodingWorkMem.GotError = true
		return nil, err
	}
	availableMemoryAsString := utils.Uint64ToString(availableMemory)
	availableMemoryConverted, err := utils.ConvertBasedOnUnit(availableMemoryAsString, "B", logicalDecodingWorkMem.Unit)
	if err != nil {
		logger.LogError("failed logical_decoding_work_mem check: " + err.Error())
		logicalDecodingWorkMem.GotError = true
		return nil, err
	}

	// 2. make suggestion by dividing available memory by 8 and rounding to nearest power of 2
	suggestion := availableMemoryConverted / 8
	suggestionRounded := utils.RoundToPowerOf2(uint64(suggestion))

	// 3. If suggested value is less than 64MB, make no suggestion
	suggestionAsMB, err := utils.ConvertBasedOnUnit(utils.Uint64ToString(suggestionRounded), logicalDecodingWorkMem.Unit, "MB")
	if err != nil {
		logger.LogError("failed logicalDecodingWorkMem check: " + err.Error())
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
		logger.LogError("failed max_stack_depth check: " + err.Error())
		maxStackDepth.GotError = true
		return nil, err
	}

	// 2. System set stack depth is not equal to max_stack_depth, suggest system stack depth
	systemStackDepthAsUnit, err := utils.ConvertBasedOnUnit(utils.Uint64ToString(systemStackDepth), "B", maxStackDepth.Unit)
	if err != nil {
		logger.LogError("failed max_stack_depth check: " + err.Error())
		maxStackDepth.GotError = true
		return nil, err
	}

	maxStackDepthFloat32, err := utils.StringToFloat32(maxStackDepth.Value)
	if err != nil {
		logger.LogError("failed max_stack_depth check: " + err.Error())
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
		logger.LogError("failed shared_memory_type check: " + err.Error())
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
		logger.LogError("failed dynamic_shared_memory_type check: " + err.Error())
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
