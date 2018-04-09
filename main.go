package main

import (
	"fmt"
	"github.com/molsbee/appfog-app-count/service"
)

func main() {
	organizations := service.DefaultCloudFoundryClient.GetOrganizations()
	guids := organizations.GetGUID()
	results := setupWorkers(organizations.GetGUID())

	totalApplications := 0
	totalRunningApplications := 0
	for i := 1; i <= len(guids); i++ {
		result := <-results
		totalApplications += result.TotalCount
		totalRunningApplications += result.RunningCount

		if result.TotalCount != 0 {
			fmt.Printf("%-20s Total Apps: %-4d Running Apps: %d\n", result.OrganizationName, result.TotalCount, result.RunningCount)
		}
	}
	close(results)

	fmt.Printf("ORG COUNT: %d Total Apps: %d Running Apps: %d", len(guids), totalApplications, totalRunningApplications)
}

func setupWorkers(guids []string) (chan WorkerResponse) {
	jobs := make(chan string, len(guids))
	results := make(chan WorkerResponse, len(guids))
	for id := 1; id <= 5; id++ {
		go func() {
			for orgID := range jobs {
				organizationDetails := service.DefaultCloudFoundryClient.GetOrganizationDetails(orgID)
				applications := service.DefaultCloudFoundryClient.GetOrganizationApplications(orgID)

				appCount := 0
				for _, app := range applications {
					if app.Entity.State != "STOPPED" {
						appCount++
					}
				}

				results <- WorkerResponse{OrganizationName: organizationDetails.Entity.Name, RunningCount: appCount, TotalCount: len(applications)}
			}
		}()
	}

	for _, orgID := range guids {
		jobs <- orgID
	}
	close(jobs)

	return results
}

type WorkerResponse struct {
	OrganizationName string
	RunningCount     int
	TotalCount       int
}
