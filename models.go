package githubrepodetails

import "time"

// RepoDetails defines the features about the repo.
type RepoDetails struct {
	// URL contains the repo url
	URL string `json:"url"`
	// Timestamp contains response timestamp in unix epochs
	Timestamp int64 `json:"timestamp"`
	// CreatedAt contains the repo creation timestamp
	CreatedAt string `json:"created_at"`
	// LastUpdatedAt contains the timestamp of repo last update
	LastUpdatedAt string `json:"last_updated_at"`
	// StargazerCount contains the count of start for the repo
	StargazerCnt int `json:"stargazer_count"`
	// ForksCnt contains the count of the repo forks
	ForksCnt int `json:"forks_count"`
}

// respGraphQL defines the object response from github graphQL.
type respGraphQL struct {
	Data   *respGraphQLData   `json:"data,omitempty"`
	Errors *respGraphQLErrors `json:"errors,omitempty"`
}

type respGraphQLData struct {
	Repository respGraphQLDataRepository `json:"repository"`
}

type respGraphQLDataRepository struct {
	URL         string `json:"url"`
	LicenseInfo struct {
		Name string `json:"name"`
	} `json:"licenseInfo"`
	DatabaseID  int       `json:"databaseId"`
	FullName    string    `json:"fullName"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Description string    `json:"description"`
	IsArchived  bool      `json:"isArchived"`
	IsFork      bool      `json:"isFork"`
	IsLocked    bool      `json:"isLocked"`
	IsDisabled  bool      `json:"isDisabled"`
	// Disk usage in kB
	DiskUsage  int `json:"diskUsage"`
	ForkCount  int `json:"forkCount"`
	Stargazers struct {
		TotalCount int `json:"totalCount"`
	} `json:"stargazers"`
	Watchers struct {
		TotalCount int `json:"totalCount"`
	} `json:"watchers"`

	Releases struct {
		TotalCount int `json:"totalCount"`
	} `json:"releases"`
	ReleasesDynamics struct {
		Nodes []*dynamicsNode `json:"nodes"`
	} `json:"releasesDynamics"`

	IssuesOpen struct {
		TotalCount int `json:"totalCount"`
	} `json:"issuesOpen"`

	IssuesClosed struct {
		TotalCount int `json:"totalCount"`
	} `json:"issuesClosed"`

	IssuesClosingDynamics struct {
		Nodes []*dynamicsNode `json:"nodes"`
	} `json:"issuesClosingDynamics"`
	IssuesOpeningDynamics struct {
		Nodes []*dynamicsNode `json:"nodes"`
	} `json:"issuesOpeningDynamics"`
}

type dynamicsNode struct {
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type respGraphQLErrors struct {
	Errors []*respGraphQLError
}

type respGraphQLError struct {
	Path       []string                    `json:"path"`
	Extensions *respGraphQLErrorExtensions `json:"extensions"`
	Locations  []*respGraphQLErrorLocation `json:"locations"`
	Message    string                      `json:"message"`
}

type respGraphQLErrorExtensions struct {
	Code      string `json:"code"`
	TypeName  string `json:"typeName"`
	fieldName string `json:"fieldName"`
}

type respGraphQLErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}
