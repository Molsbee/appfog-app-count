package service

import (
	"encoding/json"
	"github.com/molsbee/appfog-app-count/model"
	"log"
	"os/exec"
	"strconv"
)

var DefaultCloudFoundryClient = &CloudFoundryClient{}

type CloudFoundryClient struct {
}

func (c *CloudFoundryClient) GetOrganizations() model.OrganizationResponse {
	organizationResponse := c.getOrganizationsByPageNumber(0)
	for i := 2; i <= organizationResponse.TotalPages; i++ {
		organizationResponse.Resources = append(organizationResponse.Resources, c.getOrganizationsByPageNumber(i).Resources...)
	}

	return organizationResponse
}

func (c *CloudFoundryClient) getOrganizationsByPageNumber(pageNumber int) model.OrganizationResponse {
	page := strconv.Itoa(pageNumber)

	url := "/v2/organizations?page=" + page + "&results-per-page=100"
	if pageNumber == 0 {
		url = "/v2/organizations?results-per-page=100"
	}

	out, err := exec.Command("cf", "curl", url).Output()
	if err != nil {
		log.Fatal("error retreiving organizations", err)
	}

	orgPageList := &model.OrganizationResponse{}
	json.Unmarshal(out, orgPageList)

	return *orgPageList
}

func (c *CloudFoundryClient) GetOrganizationDetails(orgID string) model.OrganizationResource {
	out, err := exec.Command("cf", "curl", "/v2/organizations/"+orgID).Output()
	if err != nil {
		log.Fatal("error retreiving organization details", err)
	}

	organizationDetails := &model.OrganizationResource{}
	json.Unmarshal(out, organizationDetails)

	return *organizationDetails
}

func (c *CloudFoundryClient) GetOrganizationApplications(orgID string) []model.ApplicationResource {
	appOut, err := exec.Command("cf", "curl", "/v2/apps?q=organization_guid:"+orgID+"&results-per-page=100").Output()
	if err != nil {
		log.Fatal("error retrieving apps for organization", err)
	}

	applicationResponse := &model.ApplicationsResponse{}
	json.Unmarshal(appOut, applicationResponse)

	var applications = applicationResponse.Resources
	for i := 2; i <= applicationResponse.TotalPages; i++ {
		page := strconv.Itoa(i)
		pageOut, err := exec.Command("cf", "curl", "/v2/apps?q=organization_guid:"+orgID+"&page="+page+"&results-per-page=100").Output()
		if err != nil {
			log.Fatal("error retrieving all applications entities", err)
		}

		pageResponse := &model.ApplicationsResponse{}
		json.Unmarshal(pageOut, pageResponse)
		applications = append(applications, pageResponse.Resources...)
	}

	return applications
}
