package model

type OrganizationResponse struct {
	TotalResults int                    `json:"total_results"`
	TotalPages   int                    `json:"total_pages"`
	PrevURL      string                 `json:"prev_url"`
	NextURL      string                 `json:"next_url"`
	Resources    []OrganizationResource `json:"resources"`
}

type OrganizationResource struct {
	Metadata Metadata           `json:"metadata"`
	Entity   OrganizationEntity `json:"entity"`
}

type Metadata struct {
	GUID      string `json:"guid"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
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

type ApplicationsResponse struct {
	TotalResults int                   `json:"total_results"`
	TotalPages   int                   `json:"total_pages"`
	PrevURL      string                `json:"prev_url"`
	NextURL      string                `json:"next_url"`
	Resources    []ApplicationResource `json:"resources"`
}

type ApplicationResource struct {
	Metadata Metadata          `json:"metadata"`
	Entity   ApplicationEntity `json:"entity"`
}

type ApplicationEntity struct {
	Name                     string            `json:"name"`
	Production               string            `json:"production"`
	SpaceGUID                string            `json:"space_guid"`
	StackGUID                string            `json:"stack_guid"`
	Buildpack                string            `json:"buildpack"`
	DetectedBuildpack        string            `json:"detected_buildpack"`
	DetectedBuildpackGUID    string            `json:"detected_buildpack_guid"`
	EnvironmentJSON          string            `json:"environment_json"`
	Memory                   int               `json:"memory"`
	Instances                int               `json:"instances"`
	DiskQuota                int               `json:"disk_quota"`
	State                    string            `json:"state"`
	Version                  string            `json:"version"`
	Command                  string            `json:"command"`
	Console                  bool              `json:"console"`
	Debug                    string            `json:"debug"`
	StagingTaskID            string            `json:"staging_task_id"`
	PackageState             string            `json:"package_state"`
	HealthCheckType          string            `json:"health_check_type"`
	HealthCheckTimeout       string            `json:"health_check_timeout"`
	StagingFailedReason      string            `json:"staging_failed_reason"`
	StagingFailedDescription string            `json:"staging_failed_description"`
	Diego                    bool              `json:"diego"`
	DockerImage              string            `json:"docker_image"`
	PackageUpdatedAt         string            `json:"package_updated_at"`
	DetectedStartCommand     string            `json:"detected_start_command"`
	EnableSSH                bool              `json:"enable_ssh"`
	DockerCredentialsJSON    map[string]string `json:"docker_credentials_json"`
	Ports                    string            `json:"ports"`
	SpaceURL                 string            `json:"space_url"`
	StackURL                 string            `json:"stack_url"`
	RoutesURL                string            `json:"routes_url"`
	EventsURL                string            `json:"events_url"`
	ServiceBindingsURL       string            `json:"service_bindings_url"`
	RouteMappingsURL         string            `json:"route_mappings_url"`
}
