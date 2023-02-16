// File for system info

package utils

import (
	"fmt"
	"syscall"
)

// Get how much total RAM server has
func GetTotalMemory() (uint64, error) {
	var info syscall.Sysinfo_t
	err := syscall.Sysinfo(&info)
	if err != nil {
		return 0, fmt.Errorf("Could not get total server memory: %v", err)
	}

	return info.Totalram, nil
}

// Get how much much available RAM server has
func GetAvailableMemory() (uint64, error) {
	var info syscall.Sysinfo_t
	err := syscall.Sysinfo(&info)
	if err != nil {
		return 0, fmt.Errorf("Could not get total available memory on the server: %v", err)
	}

	return info.Freeram, nil
}

// Get maximum stack depth set on system (equivalent would be `ulimit -s`)
func GetStackSize() (uint64, error) {
	var rlimit syscall.Rlimit

	err := syscall.Getrlimit(syscall.RLIMIT_STACK, &rlimit)
	if err != nil {
		return 0, fmt.Errorf("Could not get stack size: %v", err)
	}

	// Returns stack size in bytes
	return rlimit.Cur, nil
}
