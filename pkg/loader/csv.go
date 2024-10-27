package loader

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"

	"distributor/internal/domain"
	"distributor/internal/repository"
)

const chunkSize = 1000

func LoadLocationsFromCSV(filename string, repo repository.LocationRepository) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skip header
	if _, err := reader.Read(); err != nil {
		return err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	numGoroutines := (len(records) + chunkSize - 1) / chunkSize
	log.Printf("Reading records = %v using goroutines = %v\n", len(records), numGoroutines)
	errChan := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(startIdx int) {
			defer wg.Done()

			endIdx := startIdx + chunkSize
			if endIdx > len(records) {
				endIdx = len(records)
			}

			for _, record := range records[startIdx:endIdx] {
				location := &domain.Location{
					City:     record[0],
					Province: record[1],
					Country:  record[2],
				}
				locationString := fmt.Sprintf("%s-%s-%s", location.City, location.Province, location.Country)
				if err := repo.Store(location, locationString); err != nil {
					errChan <- err
					return
				}
			}
		}(i * chunkSize)
	}

	wg.Wait()
	close(errChan)

	// Check for any errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	log.Printf("Reading records completed\n")
	return nil
}
