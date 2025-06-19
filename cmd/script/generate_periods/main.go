package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

type PayrollPeriod struct {
	ID          int
	PeriodStart string
	PeriodEnd   string
	Status      string
}

func main() {
	startDate := time.Date(2023, 10, 10, 0, 0, 0, 0, time.UTC)
	cutoffDate := time.Date(2024, 5, 10, 0, 0, 0, 0, time.UTC)
	periods := []PayrollPeriod{}

	// Create 20 periods
	for i := 1; i <= 20; i++ {
		periodStart := startDate
		periodEnd := startDate.AddDate(0, 1, -1) // 9th of next month

		status := "open"
		if !periodEnd.After(cutoffDate) {
			status = "closed"
		}

		periods = append(periods, PayrollPeriod{
			ID:          i,
			PeriodStart: periodStart.Format("2006-01-02"),
			PeriodEnd:   periodEnd.Format("2006-01-02"),
			Status:      status,
		})

		// Advance to 10th of next month
		startDate = startDate.AddDate(0, 1, 0)
	}

	// Create CSV file
	file, err := os.Create("./cmd/script/payroll_periods.csv")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"id", "period_start", "period_end", "status"})

	// Write rows
	for _, p := range periods {
		writer.Write([]string{
			fmt.Sprintf("%d", p.ID),
			p.PeriodStart,
			p.PeriodEnd,
			p.Status,
		})
	}

	fmt.Println("CSV file 'payroll_periods.csv' created successfully.")
}
