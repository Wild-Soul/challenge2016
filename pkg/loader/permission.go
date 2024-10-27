package loader

import (
	"bufio"
	"os"
	"strings"

	"distributor/internal/domain"
	"distributor/internal/repository"
)

func LoadPermissions(filename string, repo repository.DistributorRepository) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentDist *domain.Distributor

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "Permissions for ") {
			parts := strings.Fields(line)
			distName := parts[2]
			var parentName string

			if len(parts) > 3 && parts[3] == "<" {
				parentName = parts[4]
				// Validate parent exists
				if _, err := repo.Find(parentName); err != nil {
					return domain.ErrInvalidParent
				}
			}

			currentDist = domain.NewDistributor(distName, parentName)
			if err := repo.Store(currentDist); err != nil {
				return err
			}
		} else if strings.HasPrefix(line, "INCLUDE: ") || strings.HasPrefix(line, "EXCLUDE: ") {
			isInclude := strings.HasPrefix(line, "INCLUDE: ")
			location := strings.TrimPrefix(line, map[bool]string{
				true:  "INCLUDE: ",
				false: "EXCLUDE: ",
			}[isInclude])

			if err := currentDist.UpdatePermission(location, isInclude); err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}
