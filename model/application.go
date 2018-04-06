package model

import "encoding/json"

type ApplicationsResponse struct {
	TotalResults int                   `json:"total_results"`
	TotalPages   int                   `json:"total_pages"`
	PrevURL      json.RawMessage       `json:"prev_url"`
	NextURL      json.RawMessage       `json:"next_url"`
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
	Buildpack                json.RawMessage   `json:"buildpack"`
	DetectedBuildpack        json.RawMessage   `json:"detected_buildpack"`
	DetectedBuildpackGUID    json.RawMessage   `json:"detected_buildpack_guid"`
	EnvironmentJSON          json.RawMessage   `json:"environment_json"`
	Memory                   int               `json:"memory"`
	Instances                int               `json:"instances"`
	DiskQuota                int               `json:"disk_quota"`
	State                    string            `json:"state"`
	Version                  string            `json:"version"`
	Command                  json.RawMessage   `json:"command"`
	Console                  bool              `json:"console"`
	Debug                    json.RawMessage   `json:"debug"`
	StagingTaskID            json.RawMessage   `json:"staging_task_id"`
	PackageState             string            `json:"package_state"`
	HealthCheckType          string            `json:"health_check_type"`
	HealthCheckTimeout       json.RawMessage   `json:"health_check_timeout"`
	StagingFailedReason      json.RawMessage   `json:"staging_failed_reason"`
	StagingFailedDescription json.RawMessage   `json:"staging_failed_description"`
	Diego                    bool              `json:"diego"`
	DockerImage              json.RawMessage   `json:"docker_image"`
	PackageUpdatedAt         string            `json:"package_updated_at"`
	DetectedStartCommand     string            `json:"detected_start_command"`
	EnableSSH                bool              `json:"enable_ssh"`
	DockerCredentialsJSON    map[string]string `json:"docker_credentials_json"`
	Ports                    json.RawMessage   `json:"ports"`
	SpaceURL                 string            `json:"space_url"`
	StackURL                 string            `json:"stack_url"`
	RoutesURL                string            `json:"routes_url"`
	EventsURL                string            `json:"events_url"`
	ServiceBindingsURL       string            `json:"service_bindings_url"`
	RouteMappingsURL         string            `json:"route_mappings_url"`
}
