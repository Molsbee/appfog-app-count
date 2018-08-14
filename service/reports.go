package service

import (
	"encoding/csv"
	"fmt"
	"github.com/molsbee/appfog-app-count/model"
	"github.com/molsbee/go-cfclient"
	"github.com/pkg/errors"
	"net/url"
	"os"
	"strconv"
	"time"
)

func GenerateReport(username, password, region, endpoint string) error {
	fmt.Printf("collecting application data for region %s\n", region)
	client, err := cfclient.NewClient(&cfclient.Config{
		ApiAddress: endpoint,
		Username:   username,
		Password:   password,
	})

	if err != nil {
		return errors.New(fmt.Sprintf("unable to login with cloud foundry api at endpoint %s - error %s", endpoint, err.Error()))
	}

	orgs, err := client.ListOrgs()
	if err != nil {
		return errors.New("unable to collect a list of all orgs in region")
	}

	results := setupWorkers(client, orgs)

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
	for i := 1; i <= len(orgs); i++ {
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

	fmt.Printf("ORG COUNT: %d Total Apps: %d Running Apps: %d\n", len(orgs), totalApplications, totalRunningApplications)

	return nil
}

func setupWorkers(client *cfclient.Client, orgs []cfclient.Org) chan model.WorkerResponse {
	orgJobs := make(chan cfclient.Org, len(orgs))
	for _, org := range orgs {
		orgJobs <- org
	}

	results := make(chan model.WorkerResponse, len(orgs))

	for id := 1; id <= 5; id++ {
		go func() {
			for org := range orgJobs {
				apps, _ := client.ListAppsByQuery(url.Values{
					"q": {fmt.Sprintf("organization_guid:%s", org.Guid)},
				})

				appCount := 0
				for _, app := range apps {
					if app.State != "STOPPED" {
						appCount++
					}
				}

				results <- model.WorkerResponse{OrganizationName: org.Name, RunningCount: appCount, TotalCount: len(apps)}
			}
		}()
	}

	return results
}
