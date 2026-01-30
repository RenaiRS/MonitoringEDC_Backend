package dto

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"log"
	"github.com/xuri/excelize/v2"
)

type ImportMesinDTO struct {
	TerminalID    string
	MID           string
	NamaNasabah   string
	Kota          string
	Cabang        string
	TipeEDC       string
	Vendor        string
	StatusData    string
	TanggalPasang *time.Time
	BiayaBulanan  int
}

func GetRawCell(xl *excelize.File, sheet string, row, col int) string {
	cell, _ := excelize.CoordinatesToCellName(col+1, row+1)
	val, _ := xl.GetCellValue(sheet, cell, excelize.Options{
		RawCellValue: true,
	})
	return strings.TrimSpace(val)
}

func GetExcelDate(xl *excelize.File, sheet string, row, col int) (*time.Time, error) {
    cell, _ := excelize.CoordinatesToCellName(col+1, row+1)

    // Ambil nilai raw dari Excel
    raw, err := xl.GetCellValue(sheet, cell, excelize.Options{
        RawCellValue: true,
    })
    if err != nil || raw == "" {
        return nil, errors.New("empty date")
    }

    log.Printf("[DEBUG] row=%d col=%d raw='%s'", row+1, col+1, raw)

    if num, err := strconv.ParseFloat(raw, 64); err == nil {
        if num < 1 || num > 2958465 { // 2958465 ≈ 9999-12-31
            return nil, errors.New("excel date out of range")
        }
        t, err := excelize.ExcelDateToTime(num, false)
        if err == nil {
            return &t, nil
        }
    }

    layouts := []string{"2006-01-02", "1/2/2006", "01/02/2006"}
    for _, layout := range layouts {
        if t, err := time.Parse(layout, raw); err == nil {
            return &t, nil
        }
    }

    return nil, errors.New("invalid date format")
}
