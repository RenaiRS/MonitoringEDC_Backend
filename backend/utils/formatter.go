package utils

import (	
	"time"
)

func GetNamaNasabah(nama, sumberData string) string {
	if sumberData == "vendor" {
		return "Belum Terdata"
	}
	return nama
}

// FormatDateOnly formats time.Time to YYYY-MM-DD string
// Returns empty string if time is zero
func FormatDateOnly(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

// FormatDateOnlyPtr formats *time.Time to YYYY-MM-DD string
// Returns empty string if pointer is nil
func FormatDateOnlyPtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

// FormatDateTime formats time.Time to YYYY-MM-DD HH:MM:SS string
// Returns empty string if time is zero
func FormatDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// FormatDateTimePtr formats *time.Time to YYYY-MM-DD HH:MM:SS string
// Returns empty string if pointer is nil
func FormatDateTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}