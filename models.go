package githubrepodetails

import (
	"time"
)

// RepoDetailsRaw defines the raw repo details.
type RepoDetailsRaw struct {
	DatabaseID int `json:"databaseId"`
	// Now in unix epochs in milliseconds
	Timestamp        int64                    `json:"timestamp"`
	FullName         string                   `json:"fullName"`
	CreatedAt        time.Time                `json:"createdAt"`
	UpdatedAt        *time.Time               `json:"updatedAt,omitempty"`
	IsArchived       bool                     `json:"isArchived"`
	IsFork           bool                     `json:"isFork"`
	IsLocked         bool                     `json:"isLocked"`
	IsDisabled       bool                     `json:"isDisabled"`
	DiskUsage        int                      `json:"diskUsage"`
	TotalForks       int                      `json:"totalForks"`
	TotalWatchers    int                      `json:"totalWatchers"`
	TotalStargazers  int                      `json:"totalStargazers"`
	TotalReleases    int                      `json:"totalReleases"`
	TotalIssuesOpen  int                      `json:"totalIssuesOpen"`
	TotalIssuesClose int                      `json:"totalIssuesClose"`
	Documents        *repoDetailsRawDocuments `json:"documents,omitempty"`
	Dynamics         *repoDetailsRawDynamics  `json:"dynamics,omitempty"`
}

type repoDetailsRawDocuments struct {
	HealthPercentage       int        `json:"healthPercentage"`
	Description            string     `json:"description"`
	Documentation          string     `json:"documentation"`
	ReadmeURL              string     `json:"readmeURL"`
	LicenseType            string     `json:"licenseType"`
	LicenseName            string     `json:"licenseName"`
	LicenseURL             string     `json:"licenseURL"`
	CodeOfConductName      string     `json:"codeOfConductName"`
	CodeOfConductKey       string     `json:"codeOfConductKey"`
	CodeOfConductURL       string     `json:"codeOfConductURL"`
	ContributingURL        string     `json:"contributingURL"`
	IssueTemplateURL       string     `json:"issueTemplateURL"`
	PullRequestTemplateURL string     `json:"pullRequestTemplateURL"`
	UpdatedAt              *time.Time `json:"updatedAt,omitempty"`
}

type repoDetailsRawDynamics struct {
	Forks      *repoDetailsRawDynamicsForks       `json:"forks,omitempty"`
	Stargazers []*repoDetailsRawDynamicStargazers `json:"stargazers,omitempty"`
}

type repoDetailsRawDynamicsForks struct {
	TotalDiskUsage int          `json:"totalDiskUsage"`
	CreatedAt      []*time.Time `json:"createdAt,omitempty"`
}

type repoDetailsRawDynamicStargazers struct {
	Login     string     `json:"login"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

// fetchRepoDetails defines the features about the repo.
type fetchRepoDetails struct {
	GraphQLOverall      *respGraphQLOverallDataRepository
	CommunityScore      *respCommunityScore
	ContributorsCommits []*respContributorsCommits
}

func (d *fetchRepoDetails) ToOutput() *RepoDetailsRaw {
	now := time.Now().UTC().UnixMilli()
	return &RepoDetailsRaw{
		DatabaseID:       d.GraphQLOverall.DatabaseID,
		Timestamp:        now,
		FullName:         d.GraphQLOverall.FullName,
		CreatedAt:        d.GraphQLOverall.CreatedAt,
		UpdatedAt:        d.GraphQLOverall.UpdatedAt,
		IsArchived:       d.GraphQLOverall.IsArchived,
		IsFork:           d.GraphQLOverall.IsFork,
		IsLocked:         d.GraphQLOverall.IsLocked,
		IsDisabled:       d.GraphQLOverall.IsDisabled,
		DiskUsage:        d.GraphQLOverall.DiskUsage,
		TotalForks:       d.GraphQLOverall.Forks.TotalCount,
		TotalWatchers:    d.GraphQLOverall.Watchers.TotalCount,
		TotalStargazers:  d.GraphQLOverall.Stargazers.TotalCount,
		TotalReleases:    d.GraphQLOverall.Releases.TotalCount,
		TotalIssuesOpen:  d.GraphQLOverall.IssuesOpen.TotalCount,
		TotalIssuesClose: d.GraphQLOverall.IssuesClosed.TotalCount,
		Documents: &repoDetailsRawDocuments{
			HealthPercentage:       d.CommunityScore.HealthPercentage,
			Description:            d.CommunityScore.Description,
			Documentation:          d.CommunityScore.Documentation,
			ReadmeURL:              d.CommunityScore.Files.Readme.HtmlURL,
			LicenseType:            d.CommunityScore.Files.License.SpdxID,
			LicenseName:            d.CommunityScore.Files.License.Name,
			LicenseURL:             d.CommunityScore.Files.License.HtmlURL,
			CodeOfConductName:      d.CommunityScore.Files.CodeOfConduct.Name,
			CodeOfConductKey:       d.CommunityScore.Files.CodeOfConduct.Key,
			CodeOfConductURL:       d.CommunityScore.Files.CodeOfConductFile.HtmlURL,
			ContributingURL:        d.CommunityScore.Files.Contributing.HtmlURL,
			IssueTemplateURL:       d.CommunityScore.Files.IssueTemplate.HtmlURL,
			PullRequestTemplateURL: d.CommunityScore.Files.PullRequestTemplate.HtmlURL,
			UpdatedAt:              d.CommunityScore.UpdatedAt,
		},
		Dynamics: &repoDetailsRawDynamics{
			Forks: &repoDetailsRawDynamicsForks{
				TotalDiskUsage: d.GraphQLOverall.Forks.TotalDiskUsage,
			},
		},
	}
}

type detailsGraphQL struct {
	Details *respGraphQLOverallDataRepository
	Error   error
}

type detailsGraphQLForks struct {
	Details *respGraphQLForksData
	Error   error
}

type detailsContributors struct {
	Details []*respContributorsCommits
	Error   error
}

type detailsCommunityScore struct {
	Details *respCommunityScore
	Error   error
}

// respGraphQL defines the object response from github graphQL.
// see https://docs.github.com/en/graphql/overview/explorer
type respGraphQLOverall struct {
	Data   *respGraphQLOverallData `json:"data,omitempty"`
	Errors *respGraphQLErrors      `json:"errors,omitempty"`
}

type respGraphQLForksData struct {
	Repository respGraphQLForksDataRepository `json:"repository"`
}

type respGraphQLForksDataRepository struct {
	Forks struct {
		Nodes    []*respGraphQLForksDataRepositoryNode `json:"nodes"`
		PageInfo struct {
			EndCursor string `json:"endCursor"`
		} `json:"pageInfo"`
	} `json:"forks"`
}

type respGraphQLForksDataRepositoryNode struct {
	CreatedAt time.Time `json:"createdAt"`
	PushedAt  time.Time `json:"pushedAt"`
}

type respGraphQLOverallData struct {
	Repository respGraphQLOverallDataRepository `json:"repository"`
}

type respGraphQLOverallDataRepository struct {
	DatabaseID  int        `json:"databaseId"`
	FullName    string     `json:"fullName"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
	Description string     `json:"description"`
	IsArchived  bool       `json:"isArchived"`
	IsFork      bool       `json:"isFork"`
	IsLocked    bool       `json:"isLocked"`
	IsDisabled  bool       `json:"isDisabled"`
	// Disk usage in kB
	DiskUsage int `json:"diskUsage"`
	Forks     struct {
		TotalCount     int `json:"totalCount"`
		TotalDiskUsage int `json:"totalDiskUsage"`
	} `json:"forks"`
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
		Nodes []dynamicsNode `json:"nodes"`
	} `json:"releasesDynamics"`

	IssuesOpen struct {
		TotalCount int `json:"totalCount"`
	} `json:"issuesOpen"`

	IssuesClosed struct {
		TotalCount int `json:"totalCount"`
	} `json:"issuesClosed"`

	IssuesClosingDynamics struct {
		Nodes []dynamicsNode `json:"nodes"`
	} `json:"issuesClosingDynamics"`
	IssuesOpeningDynamics struct {
		Nodes []dynamicsNode `json:"nodes"`
	} `json:"issuesOpeningDynamics"`
}

type dynamicsNode struct {
	CreatedAt *time.Time   `json:"createdAt"`
	UpdatedAt *time.Time   `json:"updatedAt,omitempty"`
	Author    *issueAuthor `json:"author,omitempty"`
	Labels    *issueLabels `json:"labels,omitempty"`
}

type issueAuthor struct {
	Login string `json:"login"`
}

type issueLabels struct {
	Nodes []labelNode `json:"nodes"`
}

type labelNode struct {
	Name string `json:"name"`
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
	FieldName string `json:"fieldName"`
}

type respGraphQLErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// see https://docs.github.com/en/rest/reference/commits
type respContributorsCommits struct {
	URL     string         `json:"url"`
	SHA     string         `json:"sha"`
	HtmlURL string         `json:"html_url"`
	Commit  *commitDetails `json:"commit"`
}

type commitDetails struct {
	Author       *commiterDetails `json:"author"`
	Commiter     *commiterDetails `json:"commiter"`
	Message      string           `json:"message"`
	CommentCnt   int              `json:"comment_count"`
	Verification struct {
		Verified  bool   `json:"verified"`
		Reason    string `json:"reason"`
		Signature string `json:"signature"`
		Payload   string `json:"Payload"`
	} `json:"verification"`
}

type commiterDetails struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

// see https://docs.github.com/en/rest/reference/repository-metrics
type respCommunityScore struct {
	HealthPercentage int    `json:"health_percentage"`
	Description      string `json:"description"`
	Documentation    string `json:"documentation"`
	Files            struct {
		CodeOfConduct struct {
			Key  string `json:"key"`
			Name string `json:"name"`
			respCommunityScoreFile
		} `json:"code_of_conduct,omitempty"`
		CodeOfConductFile   respCommunityScoreFile `json:"code_of_conduct_file,omitempty"`
		Contributing        respCommunityScoreFile `json:"contributing,omitempty"`
		IssueTemplate       respCommunityScoreFile `json:"issue_template,omitempty"`
		PullRequestTemplate respCommunityScoreFile `json:"pull_request_template,omitempty"`
		Readme              respCommunityScoreFile `json:"readme,omitempty"`
		License             struct {
			Key    string `json:"key"`
			Name   string `json:"name"`
			SpdxID string `json:"spdx_id"`
			NodeID string `json:"node_id"`
			respCommunityScoreFile
		} `json:"license"`
	} `json:"files"`
	UpdatedAt             *time.Time `json:"updated_at"`
	ContentReportsEnabled bool       `json:"content_reports_enabled"`
}

type respCommunityScoreFile struct {
	URL     string `json:"url"`
	HtmlURL string `json:"html_url"`
}
