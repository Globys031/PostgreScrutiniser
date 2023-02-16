// This file is meant for converting various units

package utils

import (
	"fmt"
	"strconv"
)

// Convert B, kB, 8kB, to MB
//
// I am assuming that block_size is the default value
// of 8192 bytes as defined by `block_size` and described here:
// https://pgpedia.info/b/block_size.html#:~:text=The%20default%20value%20for%20block_size,(PostgreSQL%208.4%20and%20later).
func ConvertToMB(value string, unit string) (float32, error) {
	size, err := strconv.ParseFloat(value, 32)

	if err != nil {
		return 0, fmt.Errorf("Could not convert value: %s, to MB", value)
	}

	switch unit {
	case "MB": // Do nothing if unit is already MB
	case "kB":
		size /= 1024
	case "8kB":
		size = size * 8 / 1024
	case "B":
		size /= 1024 * 1024
	default:
		return 0, fmt.Errorf("unsupported unit: %s", unit)
	}

	return float32(size), nil
}

// Convert B, 8kB, MB to KB
func ConvertToKB(value string, unit string) (float32, error) {
	size, err := strconv.ParseFloat(value, 32)

	if err != nil {
		return 0, fmt.Errorf("Could not convert value: %s, to KB", value)
	}

	switch unit {
	case "MB":
		size *= 1024
	case "kB": // Do nothing if unit is already kB
	case "8kB":
		size *= 8
	case "B":
		size /= 1024
	default:
		return 0, fmt.Errorf("unsupported unit: %s", unit)
	}

	return float32(size), nil
}

func ConvertToB(value string, unit string) (float32, error) {
	size, err := strconv.ParseFloat(value, 32)

	if err != nil {
		return 0, fmt.Errorf("Could not convert value: %s, to KB", value)
	}

	switch unit {
	case "MB":
		size *= 1024 * 1024
	case "kB":
		size *= 1024
	case "8kB":
		size *= 8 * 1024
	case "B": // Do nothing if unit is already B
	default:
		return 0, fmt.Errorf("unrecognized unit: %s", unit)
	}

	return float32(size), nil
}

func ConvertTo8KB(value string, unit string) (float32, error) {
	size, err := strconv.ParseFloat(value, 32)

	if err != nil {
		return 0, fmt.Errorf("Could not convert value: %s, to KB", value)
	}

	switch unit {
	case "MB":
		size *= 1024 / 8
	case "kB":
		size /= 8
	case "8kB": // Do nothing if unit is already 8kB
	case "B":
		size /= 8192
	default:
		return 0, fmt.Errorf("unsupported unit: %s", unit)
	}

	return float32(size), nil
}

func ConvertBasedOnUnit(value string, unitToConvertFrom string, unitToConvertTo string) (float32, error) {
	var convertedValue float32
	var err error

	switch unitToConvertTo {
	case "MB":
		convertedValue, err = ConvertToMB(value, unitToConvertFrom)
	case "kB":
		convertedValue, err = ConvertToKB(value, unitToConvertFrom)
	case "8kB":
		convertedValue, err = ConvertTo8KB(value, unitToConvertFrom)
	case "B":
		convertedValue, err = ConvertToB(value, unitToConvertFrom)
	default:
		return 0, fmt.Errorf("unsupported unit: %s", unitToConvertFrom)
	}

	return convertedValue, err
}

func RoundToPowerOf2(n uint64) uint64 {
	// This if statement does not round to power of 2
	// However, due to nature of the app, if we do pass a value of 0, we want 0 returned
	if n == 0 {
		return 0
	}

	power := uint(0)
	// While (1*2)^power is less or equal to n
	for 1<<power <= n {
		power++
	}
	// If n = n^0 = 1
	if power == 0 {
		return 1
	}
	if n-(1<<(power-1)) < 1<<power-n {
		return 1 << (power - 1)
	}
	return 1 << power
}

//////////////////////////////////////
// Below functions came as a necessity for easier conversion management

func StringToFloat32(value string) (float32, error) {
	result, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return 0, err
	}
	return float32(result), nil
}

func StringToFloat64(value string) (float64, error) {
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func StringToUint64(value string) (uint64, error) {
	result, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func StringToInt(value string) (int, error) {
	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func Float32ToString(value float32) string {
	return strconv.FormatFloat(float64(value), 'f', -1, 32)
}

func Float64ToString(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func Uint64ToString(value uint64) string {
	return strconv.FormatUint(value, 10)
}

//////////////////////////////////////
