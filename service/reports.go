package service

import (
	"github.com/molsbee/appfog-app-count/model"
	"fmt"
	"github.com/pkg/errors"
	"time"
	"os"
	"encoding/csv"
	"strconv"
)

func GenerateReport(username, password, region, endpoint string) error {
	fmt.Printf("collecting application data for region %s\n", region)
	if err := DefaultCloudFoundryClient.Login(endpoint, username, password); err != nil {
		return errors.New(fmt.Sprintf("unable to login with cloud foundry api at endpoint %s", endpoint))
	}

	organizations := DefaultCloudFoundryClient.GetOrganizations()
	guids := organizations.GetGUID()
	results := setupWorkers(organizations.GetGUID())

	now := time.Now()
	timeStamp := fmt.Sprintf("%d-%02d-%02d_%02d-%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute())
	fileName := fmt.Sprintf("appfog-app-count-%s-%s.csv", region, timeStamp)
	file, err := os.Create(fileName)
	if err != nil {
		return errors.New("unable to create csv file for application information")
	}

	writer := csv.NewWriter(file)
	writer.Write([]string{"org", "total_app_count", "running_app_count"})

	totalApplications := 0
	totalRunningApplications := 0
	for i := 1; i <= len(guids); i++ {
		result := <-results
		totalApplications += result.TotalCount
		totalRunningApplications += result.RunningCount

		if result.TotalCount != 0 {
			writer.Write([]string{result.OrganizationName, strconv.Itoa(result.TotalCount), strconv.Itoa(result.RunningCount)})
			fmt.Printf("%-20s Total Apps: %-4d Running Apps: %d\n", result.OrganizationName, result.TotalCount, result.RunningCount)
		}
	}
	close(results)

	writer.Flush()
	if err := writer.Error(); err != nil {
		return errors.New("error writing csv content into file")
	}

	fmt.Printf("ORG COUNT: %d Total Apps: %d Running Apps: %d\n", len(guids), totalApplications, totalRunningApplications)

	return nil
}

func setupWorkers(guids []string) (chan model.WorkerResponse) {
	jobs := make(chan string, len(guids))
	results := make(chan model.WorkerResponse, len(guids))
	for id := 1; id <= 5; id++ {
		go func() {
			for orgID := range jobs {
				organizationDetails := DefaultCloudFoundryClient.GetOrganizationDetails(orgID)
				applications := DefaultCloudFoundryClient.GetOrganizationApplications(orgID)

				appCount := 0
				for _, app := range applications {
					if app.Entity.State != "STOPPED" {
						appCount++
					}
				}

				results <- model.WorkerResponse{OrganizationName: organizationDetails.Entity.Name, RunningCount: appCount, TotalCount: len(applications)}
			}
		}()
	}

	for _, orgID := range guids {
		jobs <- orgID
	}
	close(jobs)

	return results
}
