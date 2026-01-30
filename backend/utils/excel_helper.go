package utils

import (
	"strings"
)

// IsExcelFile checks if filename has excel extension
func IsExcelFile(filename string) bool {
	filename = strings.ToLower(filename)
	return strings.HasSuffix(filename, ".xlsx") || strings.HasSuffix(filename, ".xls")
}

// IsVendorHeaderRow checks if row is vendor excel header
func IsVendorHeaderRow(row []string) bool {
	if len(row) < 2 {
		return false
	}
	firstCol := strings.ToUpper(strings.TrimSpace(row[0]))
	secondCol := strings.ToUpper(strings.TrimSpace(row[1]))
	
	return firstCol == "NO" && secondCol == "TID"
}

// IsBankHeaderRow checks if row is bank excel header
func IsBankHeaderRow(row []string) bool {
	if len(row) < 2 {
		return false
	}
	
	firstCol := strings.ToUpper(strings.TrimSpace(row[0]))
	secondCol := strings.ToUpper(strings.TrimSpace(row[1]))
	
	return firstCol == "NO" && strings.Contains(secondCol, "TERMINAL")
}

// IsEmptyRow checks if all cells in row are empty
func IsEmptyRow(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

// GetCell safely retrieves cell value from row at given index
func GetCell(row []string, index int) string {
	if len(row) > index {
		return strings.TrimSpace(row[index])
	}
	return ""
}