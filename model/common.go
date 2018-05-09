package model

type Metadata struct {
	GUID      string `json:"guid"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type WorkerResponse struct {
	OrganizationName string
	RunningCount     int
	TotalCount       int
}
