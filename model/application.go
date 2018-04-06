package model

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