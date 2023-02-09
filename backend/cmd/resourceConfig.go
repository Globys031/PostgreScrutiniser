// Code for PostgreSQL configuration
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// TO DO:
// Currently I've modified /var/lib/pgsql/15/data/pg_hba.conf
// to allow all connections without any password on localhost.
// This is bad, should reenable peer authentification later and have code
// take that into account:
//
// [root@localhost backend]# grep trust /var/lib/pgsql/15/data/pg_hba.conf

type resourceSetting struct {
	name     string // name of the setting
	value    string // value of the setting
	vartype  string // what type value is (boolean, integer, enum, etc...)
	unit     string // s, ms, kB, 8kB, etc...
	minVal   string // Minimum allowed value (needed for validation)
	maxVal   string // Maximum allowed value (needed for validation)
	enumVals string // If an enumrator, this stores enum values

	suggestedValue string // Value that will be suggested after running check
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
	conf.checkSharedBuffers()

	// // Run checks
	// conf.checkSharedBuffers()

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

	// Return the map of settings
	return settingsMap, nil
}

////////////////////////////
// Configuration functions

// func (conf *configuration) checkSharedBuffers(string, error) {

// }

// func (conf *configuration) setSharedBuffers(filepath string) error {
// 	f, err := os.OpenFile(filepath, os.O_RDONLY, 0)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	lines := []string{}
// 	scanner := bufio.NewScanner(f)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		if strings.HasPrefix(line, "shared_buffers") {
// 			line = "shared_buffers = 10MB"
// 		}
// 		lines = append(lines, line)
// 	}

// 	if err := scanner.Err(); err != nil {
// 		return err
// 	}

// 	if err := f.Truncate(0); err != nil {
// 		return err
// 	}
// 	if _, err := f.Seek(0, 0); err != nil {
// 		return err
// 	}

// 	w := bufio.NewWriter(f)
// 	for _, line := range lines {
// 		if _, err := w.WriteString(line + "\n"); err != nil {
// 			return err
// 		}
// 	}
// 	if err := w.Flush(); err != nil {
// 		return err
// 	}

// 	return nil
// }
