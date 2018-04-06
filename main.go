package main

import (
	"encoding/json"
	"fmt"
	"github.com/molsbee/appfog-app-count/model"
	"log"
	"os/exec"
	"strconv"
)

func main() {
	out, err := exec.Command("cf", "curl", "/v2/organizations?results-per-page=100").Output()
	if err != nil {
		log.Fatal("error retreiving organizations", err)
	}

	organizationResponse := &model.OrganizationResponse{}
	json.Unmarshal(out, organizationResponse)

	guids := getOrganizationID(*organizationResponse)
	for i := 2; i <= organizationResponse.TotalPages; i++ {
		page := strconv.Itoa(i)
		out, err := exec.Command("cf", "curl", "/v2/organizations?page="+page+"&results-per-page=100").Output()
		if err != nil {
			log.Fatal("error retreiving organizations", err)
		}

		orgPageList := &model.OrganizationResponse{}
		json.Unmarshal(out, orgPageList)

		guids = append(guids, getOrganizationID(*orgPageList)...)
	}

	totalApplications := 0
	totalRunningApplications := 0
	for _, orgID := range guids {
		organizationDetails := getOrganizationDetails(orgID)
		applications := getOrganizationApplications(orgID)

		appCount := 0
		for _, app := range applications {
			if app.Entity.State != "STOPPED" {
				appCount++
				totalRunningApplications++
			}
		}

		orgName := organizationDetails.Entity.Name

		orgApplicationCount := len(applications)
		totalApplications += orgApplicationCount
		if orgApplicationCount != 0 {
			fmt.Printf("%-20s Total Apps: %-4d Running Apps: %d\n", orgName, orgApplicationCount, appCount)
		}
	}

	fmt.Printf("ORG COUNT: %d Total Apps: %d Running Apps: %d", len(guids), totalApplications, totalRunningApplications)
}

func getOrganizationID(org model.OrganizationResponse) []string {
	var guid []string
	for _, resource := range org.Resources {
		guid = append(guid, resource.Metadata.GUID)
	}
	return guid
}

func getOrganizationDetails(orgID string) model.OrganizationResource {
	out, err := exec.Command("cf", "curl", "/v2/organizations/"+orgID).Output()
	if err != nil {
		log.Fatal("error retreiving organization details", err)
	}

	organizationDetails := &model.OrganizationResource{}
	json.Unmarshal(out, organizationDetails)

	return *organizationDetails
}

func getOrganizationApplications(orgID string) []model.ApplicationResource {
	appOut, err := exec.Command("cf", "curl", "/v2/apps?q=organization_guid:"+orgID+"&results-per-page=100").Output()
	if err != nil {
		log.Fatal("error retrieving apps for organization", err)
	}

	applicationResponse := &model.ApplicationsResponse{}
	json.Unmarshal(appOut, applicationResponse)

	var applications = applicationResponse.Resources
	for i := 2; i <= applicationResponse.TotalPages; i++ {
		page := strconv.Itoa(i)
		pageOut, err := exec.Command("cf", "curl", "/v2/apps?q=organization_guid"+orgID+"&page="+page+"&results-per-page=100").Output()
		if err != nil {
			log.Fatal("error retrieving all applications entities", err)
		}

		pageResponse := &model.ApplicationsResponse{}
		json.Unmarshal(pageOut, pageResponse)

		applications = append(applications, pageResponse.Resources...)
	}

	return applications
}
