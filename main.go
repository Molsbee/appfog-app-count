package main

import (
	"fmt"
	"github.com/molsbee/appfog-app-count/service"
	"encoding/csv"
	"os"
	"strconv"
	"log"
	"flag"
	"github.com/howeyc/gopass"
	"time"
)

var (
	endpoints = map[string]string{"useast": "https://api.useast.appfog.ctl.io", "uswest": "https://api.uswest.appfog.ctl.io"}
)

func main() {
	username := flag.String("username", "", "string t3n account")
	flag.Parse()

	if *username == "" {
		fmt.Println("username is required to execute application - example: -username=t3n.account")
		return
	}

	fmt.Printf("Password: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		log.Fatal("unable to capture password", err)
	}

	for region, endpoint := range endpoints {
		fmt.Printf("collecting application data for region %s\n", region)

		if err := service.DefaultCloudFoundryClient.Login(endpoint, *username, string(pass)); err != nil {
			log.Fatalf("unable to login with cloud found api at endpoint %s", endpoint)
		}

		organizations := service.DefaultCloudFoundryClient.GetOrganizations()
		guids := organizations.GetGUID()
		results := setupWorkers(organizations.GetGUID())

		now := time.Now()
		timeStamp := fmt.Sprintf("%d-%02d-%02d_%02d-%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute())
		fileName := fmt.Sprintf("appfog-app-count-%s-%s.csv", region, timeStamp)
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatal("unable to create file", err)
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
			log.Fatal("error writing csv content")
		}

		fmt.Printf("ORG COUNT: %d Total Apps: %d Running Apps: %d", len(guids), totalApplications, totalRunningApplications)
	}
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
