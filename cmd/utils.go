package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// 验证年月参数格式
func validateMonthYear(monthYear string) error {
	parts := strings.Split(monthYear, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid format for month-year: %s. Expected format: YYYY-MM", monthYear)
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid year in month-year: %s. Expected format: YYYY-MM", monthYear)
	}

	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid month in month-year: %s. Expected format: YYYY-MM", monthYear)
	}

	if month < 1 || month > 12 {
		return fmt.Errorf("invalid month in month-year: %s. Month should be between 1 and 12", monthYear)
	}

	// 验证年份是否合理
	if year < 0 || year > time.Now().Year()+100 {
		return fmt.Errorf("invalid year in month-year: %s. Year should be within a reasonable range", monthYear)
	}

	return nil
}

func getLastMonth() string {
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)
	return lastMonth.Format("2006-01")
}

func getMonthsBetween(start, end string) ([]string, error) {
	var months []string

	startTime, err := time.Parse("2006-01", start)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}
	endTime, err := time.Parse("2006-01", end)

	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}

	for startTime.Before(endTime) || startTime.Equal(endTime) {
		months = append(months, startTime.Format("2006-01"))
		startTime = startTime.AddDate(0, 1, 0)
	}
	return months, nil

}
