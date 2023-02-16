// Code for PostgreSQL configuration
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	sysctl "github.com/lorenzosaino/go-sysctl"

	"github.com/Globys031/PostgreScrutiniser/backend/utils"
)

// TO DO:
// Currently I've modified /var/lib/pgsql/15/data/pg_hba.conf
// to allow all connections without any password on localhost.
// This is bad, should reenable peer authentification later and have code
// take that into account:
//
// [root@localhost backend]# grep trust /var/lib/pgsql/15/data/pg_hba.conf

// TO DO:
// also return documentation for every setting

// TO DO:
// Tam tikri grazins tik general recommendations (ties kurie virs funkcijos) pamarkinti "GENERALREC"
// Pasvarstyt kaip front'e idet... Gal padaryt tris boxes kur
// auksciausiai yra: "found %{n} appliable suggestions"
// zemiau: General recommendations %{n}
// zemiausiai: Passed without suggestions %{n}
// Prie "General recommendations" galima pridet ui-info, kad "this group contains recommendations
// that could not be fully decided based on your system parameters (for example, because PostgreScrutiniser does not know what queries your applications are typically running). Suggestions here should be applied at your own discretion based on details provided"

// TO DO:
// Add somewhere in the GUI a prompt that:
// 1. This app assumes that block_size is the default value of 8192 bytes as defined by `block_size` and described here: https://pgpedia.info/b/block_size.html#:~:text=The%20default%20value%20for%20block_size,(PostgreSQL%208.4%20and%20later).
// 2. This app assumes that all integer type settings have unit values specified

type resourceSetting struct {
	// TO DO: change minVal, maxVal, enumVals to something that's not string
	// Consider having separate maps for differnent vartypes

	name     string // name of the setting
	value    string // value of the setting
	vartype  string // what type value is (boolean, integer, enum, etc...)
	unit     string // s, ms, kB, 8kB, etc...
	minVal   string // Minimum allowed value (needed for validation)
	maxVal   string // Maximum allowed value (needed for validation)
	enumVals string // If an enumrator, this stores enum values

	suggestedValue string // Value that will be suggested after running check
	details        string // Details informing why a value was suggested
	// Specifies if applying suggestion requires rebooting postgresql server
	// If set to false, but a suggestedValue has been applied, will do
	// `pg_ctl reload` automatically
	// https://www.postgresql.org/docs/15/app-pg-ctl.html
	// TO DO: check to confirm docs are referring to restarting postgresql, not the server itself
	requiresRestart bool
}

type configuration struct {
	// postgreUser string
	// postgreUserPass string
	path     string // Path to postgresql.conf
	settings map[string]resourceSetting
	// settings [34]string
}

func RunChecks() {
	// user, pass := findPostgresCredentials();

	// Find path to postgresql.conf
	configFile, err := findConfigFile()
	if err != nil {
		// TO DO: use https://pkg.go.dev/log instead
		fmt.Printf("Unable to find postgresql.conf:\n%v\n\n", err)
	}

	resourceSettings, err := getPGSettings()
	conf := configuration{path: configFile, settings: resourceSettings}

	///////////////////////////////
	// Run checks
	conf.checkSharedBuffers()
	conf.checkHugePages()
	conf.checkHugePageSize()
	conf.checkTempBuffers()
	conf.checkMaxPreparedTransactions()
	conf.checkWorkMem()
	conf.checkHashMemMultiplier()
	conf.checkMaintenanceWorkMem()

	// setting := "maintenance_work_mem"
	// fmt.Printf("RunChecks() %s.value: %s\n", setting, conf.settings[setting].value)
	// fmt.Printf("RunChecks() %s.suggestedValue: %s\n", setting, conf.settings[setting].suggestedValue)

}

// Function for finding main postgresql user credentials
// to later be used for peer authentification
func findPostgresCredentials() (string, error) {
	panic("findPostgresCredentials not implemented")
}

// Returns path to postgresql.conf if it exists.
func findConfigFile() (string, error) {
	output, err := exec.Command("psql", "-U", "postgres", "-c", "SHOW config_file").Output()
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	// Return line that has postgresql.conf in it
	for scanner.Scan() {
		line := scanner.Text()
		if match, _ := regexp.MatchString("postgresql.conf", line); match {
			return strings.TrimSpace(line), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("something unexpected happened")
}

// Stores data returned by `pg_settings` into a map
// only for settings we're interested in.
func getPGSettings() (map[string]resourceSetting, error) {
	// Settings we're interested in
	resourceSettings := [34]string{"shared_buffers",
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
		"dynamic_shared_memory_type",
		"min_dynamic_shared_memory",
		"temp_file_limit",
		"max_files_per_process",
		"vacuum_cost_delay",
		"vacuum_cost_page_hit",
		"vacuum_cost_page_miss",
		"vacuum_cost_page_dirty",
		"vacuum_cost_limit",
		"bgwriter_delay",
		"bgwriter_lru_maxpages",
		"bgwriter_lru_multiplier",
		"bgwriter_flush_after",
		"backend_flush_after",
		"effective_io_concurrency",
		"maintenance_io_concurrency",
		"max_worker_processes",
		"max_parallel_workers_per_gather",
		"max_parallel_maintenance_workers",
		"max_parallel_workers",
		"parallel_leader_participation",
		"old_snapshot_threshold"}

	// Use exec.Command to run the `psql` command and capture the output
	cmd := exec.Command("psql", "-U", "postgres", "-c", "select name,setting,vartype,unit,min_val,max_val,enumvals from pg_settings")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	// Initialize the map to store the results
	settingsMap := make(map[string]resourceSetting)

	// Use bufio.NewScanner to read the output line by line
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()

		// Use regexp.MatchString to check if the line contains one of the resourceSettings
		for _, resource := range resourceSettings {
			match, _ := regexp.MatchString(resource, line)
			if match {
				// If the line contains a resourceSetting, split it by "|"
				// and store the name and setting in the map
				fields := strings.Split(line, "|")
				name := strings.TrimSpace(fields[0])
				value := strings.TrimSpace(fields[1])
				vartype := strings.TrimSpace(fields[2])
				unit := strings.TrimSpace(fields[3])
				minVal := strings.TrimSpace(fields[4])
				maxVal := strings.TrimSpace(fields[5])
				enumVals := strings.TrimSpace(fields[6])

				resSetting := resourceSetting{
					name:     name,     // name of the setting
					value:    value,    // value of the setting
					vartype:  vartype,  // what type value is (boolean, integer, enum, etc...)
					unit:     unit,     // s, ms, kB, 8kB, etc...
					minVal:   minVal,   // Minimum allowed value (needed for validation)
					maxVal:   maxVal,   // Maximum allowed value (needed for validation)
					enumVals: enumVals, // If an enumrator, this stores enum values
				}

				settingsMap[name] = resSetting
				break
			}
		}
	}

	// Return map of runtime config settings
	return settingsMap, nil
}

func getSpecificPGSetting(setting string) (*resourceSetting, error) {
	formattedArg := fmt.Sprintf("select name,setting,vartype,unit,min_val,max_val,enumvals from pg_settings WHERE name = '%s'", setting)

	// Execute command to fetch postgre setting
	cmd := exec.Command("psql", "-U", "postgres", "-c", formattedArg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	// This will store the final result of the setting to be restored
	var returnSetting resourceSetting

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
			enumVals := strings.TrimSpace(fields[6])

			returnSetting = resourceSetting{
				name:     name,     // name of the setting
				value:    value,    // value of the setting
				vartype:  vartype,  // what type value is (boolean, integer, enum, etc...)
				unit:     unit,     // s, ms, kB, 8kB, etc...
				minVal:   minVal,   // Minimum allowed value (needed for validation)
				maxVal:   maxVal,   // Maximum allowed value (needed for validation)
				enumVals: enumVals, // If an enumrator, this stores enum values
			}
		}
	}

	return &returnSetting, nil
}

// Function used to ensure that parameter value being set
// is amongst options returned by `postgres -c "select enumvals...`
// Helps avoid exposing incorrect setting suggestions to users.
func setEnumTypeSuggestedValue(setting *resourceSetting, valueToSet string) error {
	// 1. Extract enumerator values from something like `{off,on,try}`
	start := strings.Index(setting.enumVals, "{")
	end := strings.Index(setting.enumVals, "}")
	if start == -1 || end == -1 {
		return fmt.Errorf("Could not extract enumerator values from enumvals")
	}

	enumVals := strings.Split(setting.enumVals[start+1:end], ",")

	// 2. Set value we're trying to set if it exists
	for _, enumValue := range enumVals {
		if enumValue == valueToSet {
			setting.suggestedValue = valueToSet
			return nil
		}
	}
	return fmt.Errorf("There is no %s in setting's %s enumerator %s", valueToSet, setting.name, setting.enumVals)
}

////////////////////////////////////////////////////////
// Configuration functions
////////////////////////////////////////////////////////

// Checks what unit is used for shared_buffers, applies necessary conversions
// and sets final suggestion as a shared_buffers unit to closest value that's power of 2
func (conf *configuration) checkSharedBuffers() error {
	var suggestion float32
	var lowestRecommendedValue float32 = 128 // Cannot suggest value that is lower than 128MB
	var GigabyteInBytes uint64 = 1073741824  // 1GB
	sharedBuffers := conf.settings["shared_buffers"]

	// 1. Get total server memory
	totalMemory, err := utils.GetTotalMemory()
	if err != nil {
		return err
	}

	// 2. Convert total server memory to a unit that's used by shared_buffers
	totalMemoryAsString := utils.Uint64ToString(totalMemory)
	totalMemoryConverted, err := utils.ConvertBasedOnUnit(totalMemoryAsString, "B", sharedBuffers.unit)
	if err != nil {
		return err
	}

	// 3. Suggest value that's n% memory depending on how much total and free memory we have
	//
	// If total server memory > 1GB, suggest 25% of total server RAM
	if totalMemory > GigabyteInBytes {
		suggestion = totalMemoryConverted * 0.25
		sharedBuffers.details = "Total server memory > 1GB. Suggest using 25% of total server RAM"
	} else { // Else suggest 30% of what memory is currently available on the server
		availableMemory, err := utils.GetAvailableMemory()
		if err != nil {
			return err
		}
		availableMemoryAsString := utils.Uint64ToString(availableMemory)
		availableMemoryConverted, err := utils.ConvertBasedOnUnit(availableMemoryAsString, "B", sharedBuffers.unit)
		if err != nil {
			return err
		}
		suggestion = availableMemoryConverted * 0.30
		sharedBuffers.details = "Total server memory < 1GB. Suggest using 30% of available RAM"

		// Suggested value cannot be lower than 128MB
		lowestRecommendedValueAsString := utils.Float32ToString(lowestRecommendedValue)
		lowestRecommendedValue, err = utils.ConvertBasedOnUnit(lowestRecommendedValueAsString, "MB", sharedBuffers.unit)
		if err != nil {
			return err
		}
		if suggestion < lowestRecommendedValue {
			suggestion = lowestRecommendedValue
			sharedBuffers.details = "Server lacks available memory. Lowest recommended value is 128MB"
		}
	}

	// Round suggestion to power of 2 and make sure there's no decimal point
	roundedSuggestion := utils.RoundToPowerOf2(uint64(suggestion))
	sharedBuffers.suggestedValue = utils.Uint64ToString(roundedSuggestion)
	// https://www.postgresql.org/docs/current/runtime-config-resource.html
	// "This parameter can only be set at server start."
	sharedBuffers.requiresRestart = true
	conf.settings["shared_buffers"] = sharedBuffers
	return nil
}

func (conf *configuration) checkHugePages() error {
	hugePages := conf.settings["huge_pages"]

	// Get nr_hugepages value
	kernelPagesString, err := sysctl.Get("vm.nr_hugepages")
	if err != nil {
		return err
	}
	kernelNrHugePages, err := utils.StringToInt(kernelPagesString)

	// if it's not set in kernel, no point in having hugePages set to on/try
	if kernelNrHugePages == 0 {
		hugePages.details = "Kernel parameter nr_hugepages is set to 0. Because of that, PostgreSQL cannot request huge pages"
		setEnumTypeSuggestedValue(&hugePages, "off")

		conf.settings["huge_pages"] = hugePages
		return nil
	}

	if hugePages.value == "on" {
		hugePages.details = "huge_pages current value is equal to 'on'. Failure to request huge pages will prevent the server from starting up"
		setEnumTypeSuggestedValue(&hugePages, "try")
	} else if hugePages.value == "off" {
		hugePages.details = "huge_pages current value is equal to 'off'. Use of huge pages results in smaller page tables and less CPU time spent on memory management, increasing performance"
		setEnumTypeSuggestedValue(&hugePages, "try")
	}

	// If kernelNrHugePages > 0 and hugePages.suggestedValue = "try", make no suggestions

	conf.settings["huge_pages"] = hugePages
	return nil
}

func (conf *configuration) checkHugePageSize() error {
	hugePageSize := conf.settings["huge_page_size"]

	// Get nr_hugepages value
	kernelPagesString, err := sysctl.Get("vm.nr_hugepages")
	if err != nil {
		return err
	}
	kernelNrHugePages, err := strconv.Atoi(kernelPagesString)

	// If vm.nr_hugepages is 0, then huge_page_size cannot be set
	if kernelNrHugePages == 0 {
		return nil
	}

	if hugePageSize.value != "0" {
		hugePageSize.suggestedValue = "0"
		hugePageSize.details = "Current huge_page_size value is set to a non 0 value. To prevent fragmentation, the same huge page size as the one set in your Linux kernel should be used. When set to 0, the default huge page size on the system will be used."
	}

	hugePageSize.requiresRestart = true
	conf.settings["huge_page_size"] = hugePageSize
	return nil
}

// GENERALREC
func (conf *configuration) checkTempBuffers() error {
	tempBuffers := conf.settings["temp_buffers"]

	currentValue, err := utils.StringToUint64(tempBuffers.value)
	if err != nil {
		return err
	}
	currentValueAsString := utils.Uint64ToString(currentValue)
	currentValueConverted, err := utils.ConvertBasedOnUnit(currentValueAsString, "8kB", tempBuffers.unit)
	if err != nil {
		return err
	}
	if currentValueConverted > 1024 { // 1024 8kB is equivalent to 8MB
		tempBuffers.details = "Current temp_buffers value is set to more than 8MB. If there are multiple databases used by different application, consider changing this setting per database. It is recommended to increase this value only for applications that rely heavily on temporary tables"
	} else if currentValueConverted < 1024 {
		tempBuffers.details = "Current temp_buffers value is set to less than 8MB. The cost of setting a large value in sessions that do not actually need many temporary buffers is only a buffer descriptor, or about 64 bytes. Consider increasing this to the recommended default value"

		convertedDefault, err := utils.ConvertBasedOnUnit("1024", "8kB", tempBuffers.unit)
		if err != nil {
			return err
		}
		tempBuffers.suggestedValue = utils.Float32ToString(convertedDefault)
	}

	conf.settings["temp_buffers"] = tempBuffers
	return nil
}

// GENERALREC
func (conf *configuration) checkMaxPreparedTransactions() error {
	maxPreparedTransactions := conf.settings["max_prepared_transactions"]

	// !!! consider passing this as an argument because there are currently two
	// separate functions that ask for max_connections
	maxConnections, err := getSpecificPGSetting("max_connections")
	if err != nil {
		return err
	}

	if maxPreparedTransactions.value == "0" {
		maxPreparedTransactions.details = " If you are using prepared transactions, you will probably want max_prepared_transactions to be at least as large as max_connections, so that every session can have a prepared transaction pending."
		maxPreparedTransactions.suggestedValue = maxConnections.value
	}

	maxPreparedTransactions.requiresRestart = true
	conf.settings["max_prepared_transactions"] = maxPreparedTransactions
	return nil
}

func (conf *configuration) checkWorkMem() error {
	workMem := conf.settings["work_mem"]

	// 1. Get amount of available memory and connections
	availableMemory, err := utils.GetAvailableMemory()
	if err != nil {
		return err
	}
	maxConnections, err := getSpecificPGSetting("max_connections")
	if err != nil {
		return err
	}

	// 2. suggestion = availablememory / max_connections
	maxConnectionsValue, err := utils.StringToUint64(maxConnections.value)
	if err != nil {
		return err
	}
	suggestion := utils.Uint64ToString(utils.RoundToPowerOf2(availableMemory / maxConnectionsValue))
	suggestionAsWorkMemUnit, err := utils.ConvertBasedOnUnit(suggestion, "B", workMem.unit)
	workMem.suggestedValue = utils.Float32ToString(suggestionAsWorkMemUnit)

	// 3. Add details for decision
	workMem.details = "Suggested value is based on currently available memory on the server divided by max_connections. If using complex queries that involve sorts or hash tables, consider using double this value. It can also be set higher if this server is a dedicated database server and there is no concern that other software will run out of memory."

	conf.settings["work_mem"] = workMem
	return nil
}

// GENERALREC
func (conf *configuration) checkHashMemMultiplier() error {
	hashMemMultiplier := conf.settings["hash_mem_multiplier"]
	workMem := conf.settings["work_mem"]
	workMemValueAsMB, err := utils.ConvertBasedOnUnit(workMem.value, workMem.unit, "MB")
	if err != nil {
		return err
	}

	// If more than 40MB
	if workMemValueAsMB > 40 {
		hashMemMultiplierAsFloat32, err := utils.StringToFloat32(hashMemMultiplier.value)
		if err != nil {
			return err
		}
		suggestion := (hashMemMultiplierAsFloat32 + workMemValueAsMB*0.01)
		if suggestion > 8 {
			suggestion = 8
		}
		hashMemMultiplier.details = "Generally default value works best. If your application uses hash-based operations and PostgreSQL often ends up spilling (creates workfiles on disk to compensate for lack of memory), consider increasing this further. Suggested value is based on how much working memory is currently set."
		hashMemMultiplier.suggestedValue = utils.Float32ToString(suggestion)
	} else if hashMemMultiplier.value != "2" {
		fmt.Println("ieina")
		hashMemMultiplier.details = "Generally default value works best. If your application uses hash-based operations and PostgreSQL often ends up spilling (creates workfiles on disk to compensate for lack of memory), consider increasing this after having increased work_mem above 40MB."
		hashMemMultiplier.suggestedValue = "2"
	}

	conf.settings["hash_mem_multiplier"] = hashMemMultiplier
	return nil
}

func (conf *configuration) checkMaintenanceWorkMem() error {
	// 0. Fetch maintenance_work_mem and autovacumm_max_workers
	maintenanceWorkMem := conf.settings["maintenance_work_mem"]
	// tmpMaxWorkers only used to fetch autovacuum_max_workers
	tmpMaxWorkers, err := getSpecificPGSetting("autovacuum_max_workers")
	if err != nil {
		return err
	}
	autovacuumMaxWorkers, err := utils.StringToFloat32(tmpMaxWorkers.value)
	if err != nil {
		return err
	}

	// 1. Get available memory on server
	availableMemory, err := utils.GetAvailableMemory()
	if err != nil {
		return err
	}
	availableMemoryAsString := utils.Uint64ToString(availableMemory)
	availableMemoryConverted, err := utils.ConvertBasedOnUnit(availableMemoryAsString, "B", maintenanceWorkMem.unit)
	if err != nil {
		return err
	}

	// 2. Divide available memory by 8 * autovacuum_max_workers and round to nearest power of 2
	suggestion := availableMemoryConverted / 8 / autovacuumMaxWorkers
	suggestionRounded := utils.RoundToPowerOf2(uint64(suggestion))

	maintenanceWorkMem.details = fmt.Sprintf("This suggestion was made by dividing current available memory(%.2f%s) by 8 and multiplying by how many autovacuum_max_workers(%.0f) are set. Applications that heavily rely on maintenance operations, such as VACUUM, CREATE INDEX, and ALTER TABLE ADD FOREIGN KEY may want to increase this further by multiplying the suggested value by 2", availableMemoryConverted, maintenanceWorkMem.unit, autovacuumMaxWorkers)
	if err != nil {
		return err
	}
	maintenanceWorkMem.suggestedValue = utils.Uint64ToString(suggestionRounded)

	// 3. If suggestion is below default value and current value is not equal to default,
	// suggest default value instead.

	suggestionAsMB, err := utils.ConvertBasedOnUnit(utils.Uint64ToString(suggestionRounded), maintenanceWorkMem.unit, "MB")
	if err != nil {
		return err
	}

	if suggestionAsMB < 64 {
		maintenanceWorkMemAsMB, err := utils.ConvertBasedOnUnit(maintenanceWorkMem.value, maintenanceWorkMem.unit, "MB")
		if err != nil {
			return err
		}
		if maintenanceWorkMemAsMB != 64 {
			maintenanceWorkMem.details = "Currently there is not enough available memory on the server to go above default maintenance_work_mem value"
			defaultAsUnit, err := utils.ConvertBasedOnUnit("64", "MB", maintenanceWorkMem.unit)
			if err != nil {
				return err
			}
			maintenanceWorkMem.suggestedValue = utils.Float32ToString(defaultAsUnit)
		} else { // if suggestion was below default and current value is already set to default
			maintenanceWorkMem.details = ""
			maintenanceWorkMem.suggestedValue = ""
		}
	}

	conf.settings["maintenance_work_mem"] = maintenanceWorkMem
	return nil
}

///////////////////////////////////////////////////
///////////////////////////////////////////////////
///////////////////////////////////////////////////
///////////////////////////////////////////////////
