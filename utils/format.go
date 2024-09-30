package utils

import "strconv"

// StringToUint64 converts a string to a uint64
func StringToUint64(s string) uint64 {
    id, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
    return id
}
// StringToInt64 converts a string to a int64
func StringToFloat64(s string) float64 {
	id, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return id
}
// StringToInt converts a string to a int
func StringToInt(s string) int {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return id
}