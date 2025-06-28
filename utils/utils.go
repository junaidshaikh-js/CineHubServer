package utils

import "strconv"

func ParseID(idStr string) (int, error) {
	return strconv.Atoi(idStr)
}
