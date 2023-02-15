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
	// Specifies if applying suggestion requires rebooting eerver
	// If set to false, but a suggestedValue has been applied, will do
	// `pg_ctl reload`
	// https://www.postgresql.org/docs/15/app-pg-ctl.html
	requiresReboot bool
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

	// setting := "huge_pages"
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
	totalMemoryAsString := strconv.FormatUint(totalMemory, 10)
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
		availableMemoryAsString := strconv.FormatUint(availableMemory, 10)
		availableMemoryConverted, err := utils.ConvertBasedOnUnit(availableMemoryAsString, "B", sharedBuffers.unit)
		if err != nil {
			return err
		}
		suggestion = availableMemoryConverted * 0.30
		sharedBuffers.details = "Total server memory < 1GB. Suggest using 30% of available RAM"

		// Suggested value cannot be lower than 128MB
		lowestRecommendedValueAsString := strconv.FormatFloat(float64(lowestRecommendedValue), 'f', -1, 32)
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
	roundedSuggestion := utils.RoundToPowerOf2(int(suggestion))
	sharedBuffers.suggestedValue = strconv.Itoa(roundedSuggestion)
	// https://www.postgresql.org/docs/current/runtime-config-resource.html
	// "This parameter can only be set at server start."
	sharedBuffers.requiresReboot = true
	conf.settings["shared_buffers"] = sharedBuffers
	return nil
}

// https://www.postgresql.org/docs/current/runtime-config-resource.html
// "This parameter can only be set at server start."

// Checks huge_pages
func (conf *configuration) checkHugePages() error {
	hugePages := conf.settings["huge_pages"]

	fmt.Println(hugePages.enumVals)

	// Get nr_hugepages value
	kernelPagesString, err := sysctl.Get("vm.nr_hugepages")
	if err != nil {
		return err
	}
	kernelNrHugePages, err := strconv.Atoi(kernelPagesString)

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
