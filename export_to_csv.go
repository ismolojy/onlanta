package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// Экспорт данных в CSV. Проблемы разделены через "; ".
func exportToCsv(filename string, channel <-chan hostsWithProblems) (int, error) {
	count := 0
	file, err := os.Create(filename)
	if err != nil {
		return count, fmt.Errorf("File don`t created: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"hostid", "host", "name", "problems"}
	if err := writer.Write(headers); err != nil {
		return count, fmt.Errorf("Export CSV: headers error: %v", err)
	}

	for hostWithProblems := range channel {
		problems := strings.Join(hostWithProblems.Problems, "; ")
		record := []string{
			hostWithProblems.Host.HostId,
			hostWithProblems.Host.Host,
			hostWithProblems.Host.Name,
			problems,
		}
		if err := writer.Write(record); err != nil {
			return count, fmt.Errorf("Record Error: %v", err)
		}
		count++
	}
	return count, nil
}
