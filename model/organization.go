package model

type OrganizationResponse struct {
	TotalResults int                    `json:"total_results"`
	TotalPages   int                    `json:"total_pages"`
	PrevURL      string                 `json:"prev_url"`
	NextURL      string                 `json:"next_url"`
	Resources    []OrganizationResource `json:"resources"`
}

func (orgs *OrganizationResponse) GetGUID() []string {
	var guids []string
	for _, resource := range orgs.Resources {
		guids = append(guids, resource.Metadata.GUID)
	}

	return guids
}

type OrganizationResource struct {
	Metadata Metadata           `json:"metadata"`
	Entity   OrganizationEntity `json:"entity"`
}

type OrganizationEntity struct {
	Name                     string `json:"name"`
	BillingEnabled           bool   `json:"billing_enabled"`
	QuotaDefinitionGUID      string `json:"quota_definition_guid"`
	Status                   string `json:"status"`
	QuotaDefinitionURL       string `json:"quota_definition_url"`
	SpacesURL                string `json:"spaces_url"`
	DomainsURL               string `json:"domains_url"`
	PrivateDomainsURL        string `json:"private_domains_url"`
	UsersURL                 string `json:"users_url"`
	ManagersURL              string `json:"managers_url"`
	BillingManagersURL       string `json:"billing_managers_url"`
	AuditorsURL              string `json:"auditors_url"`
	AppEventsURL             string `json:"app_events_url"`
	SpaceQuotaDefinitionsURL string `json:"space_quota_definitions_url"`
}
