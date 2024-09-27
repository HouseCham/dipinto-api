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