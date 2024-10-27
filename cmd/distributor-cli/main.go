package main

import (
	"flag"
	"fmt"
	"log"

	"distributor/internal/config"
	"distributor/internal/service"
	"distributor/internal/storage/memory"
	"distributor/pkg/loader"
)

func main() {
	cfg := parseFlags()

	// Initialize repositories
	locationRepo := memory.NewLocationRepository()
	distributorRepo := memory.NewDistributorRepository()

	// Initialize service
	permissionService := service.NewPermissionService(locationRepo, distributorRepo)

	// Load data
	if err := loader.LoadLocationsFromCSV(cfg.CSVPath, locationRepo); err != nil {
		log.Fatalf("Failed to load locations: %v", err)
	}

	if err := loader.LoadPermissions(cfg.PermPath, distributorRepo); err != nil {
		log.Fatalf("Failed to load permissions: %v", err)
	}

	// Check permission if requested
	if cfg.CheckLocation != "" && cfg.DistributorName != "" {
		result, err := permissionService.CheckPermission(cfg.DistributorName, cfg.CheckLocation)
		if err != nil {
			log.Fatalf("Error checking permission: %v", err)
		}
		if result {
			fmt.Println("YES")
		} else {
			fmt.Println("NO")
		}
	}
}

func parseFlags() *config.Config {
	cfg := &config.Config{}
	flag.StringVar(&cfg.CSVPath, "csv", "cities.csv", "Path to the cities CSV file")
	flag.StringVar(&cfg.PermPath, "perm", "permissions.txt", "Path to the permissions file")
	flag.StringVar(&cfg.CheckLocation, "check", "", "Check permission for location (format: CITY-PROVINCE-COUNTRY)")
	flag.StringVar(&cfg.DistributorName, "dist", "", "Distributor name to check permission for")
	flag.Parse()
	return cfg
}
