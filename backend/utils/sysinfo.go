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
		return 0, fmt.Errorf("Could not get total server memory")
	}

	return info.Totalram, nil
}

// Get how much much available RAM server has
func GetAvailableMemory() (uint64, error) {
	var info syscall.Sysinfo_t
	err := syscall.Sysinfo(&info)
	if err != nil {
		return 0, fmt.Errorf("Could not get total available memory on the server")
	}

	return info.Freeram, nil
}
