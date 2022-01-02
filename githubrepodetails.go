/* Module to fetch details of a github repo using github API and github graphQL. */
package githubrepodetails

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/goccy/go-json"
	httpClient "github.com/kislerdm/github-repo-details/internal/http/client"
)

// Client defines the client to fetch data.
type Client struct {
	// http headers
	headers    map[string]string
	httpClient *httpClient.Client
}

// ClientOption defines the option to modify Client's behavior.
type ClientOption func(*Client)

const (
	defaultTimeout    = 2 * time.Second
	defaultPagination = 100
	apiBaseURL        = "https://api.github.com"
)

// NewClient defines the function to init a Client.
func NewClient(opts ...ClientOption) (*Client, error) {
	httpClient, err := httpClient.NewClient(httpClient.WithTimeout(defaultTimeout))
	if err != nil {
		return nil, err
	}
	c := &Client{
		headers: map[string]string{
			"Content-Type": "application/json",
		},
		httpClient: httpClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// WithAPIKey defines the option to set github API key.
func WithAPIKey(key string) ClientOption {
	return func(c *Client) {
		c.headers["Authorization"] = fmt.Sprintf("bearer %s", key)
	}
}

// FetchDetails defines the function to fetch repo details.
// Requires
// owner - the repo owner
// repo - the repo name
func (c *Client) FetchDetails(owner, repo string) (*RepoDetailsRaw, error) {
	resp := &fetchRepoDetails{
		GraphQLOverall:      &respGraphQLOverallDataRepository{},
		CommunityScore:      &respCommunityScore{},
		ContributorsCommits: []*respContributorsCommits{},
	}
	var wg sync.WaitGroup
	ch := make(chan error, 3)

	wg.Add(1)
	go func(wg *sync.WaitGroup, resp *fetchRepoDetails, ch chan error) {
		r, err := c.getGraphQLOverall(owner, repo)
		if err != nil {
			ch <- err
		} else {
			if r.Errors != nil {
				errs := ""
				for _, e := range r.Errors.Errors {
					errs = fmt.Sprintf("%s%s\n", errs, e.Message)
				}
				ch <- errors.New(errs)
			} else {
				resp.GraphQLOverall = &r.Data.Repository
			}
		}
		wg.Done()
	}(&wg, resp, ch)

	wg.Add(1)
	go func(wg *sync.WaitGroup, res *fetchRepoDetails, ch chan error) {
		r, err := c.getCommunityProfileMetrics(owner, repo)
		if err != nil {
			ch <- err
		} else {
			res.CommunityScore = r
		}
		wg.Done()
	}(&wg, resp, ch)

	wg.Add(1)
	go func(wg *sync.WaitGroup, res *fetchRepoDetails, ch chan error) {
		r, err := c.getCommits(owner, repo)
		if err != nil {
			ch <- err
		} else {
			res.ContributorsCommits = r
		}
		wg.Done()
	}(&wg, resp, ch)

	wg.Wait()

	eOutStr := ""
	for e := range ch {
		if e != nil {
			eOutStr = fmt.Sprintf("%s%s\n", eOutStr, e.Error())
		}
	}
	if eOutStr != "" {
		return nil, errors.New(eOutStr)
	}
	return resp.ToOutput(), nil
}

func getGraphQLQuery(owner, repo, queryArgs string) ([]byte, error) {
	q := map[string]string{
		"query": fmt.Sprintf(`{repository(name: "%s", owner: "%s") {%s}}`, repo, owner, queryArgs),
	}
	return json.Marshal(q)
}

func (c *Client) fetchGraphQL(query []byte) ([]byte, error) {
	headers := c.headers
	headers["X-REQUEST-TYPE"] = "GraphQL"

	url := fmt.Sprintf("%s/graphql", apiBaseURL)
	resp, err := c.httpClient.POST(url, headers, query)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (c *Client) getGraphQLOverall(owner, repo string) (result *respGraphQLOverall, err error) {
	queryArgs := `databaseId
    url
    fullName: nameWithOwner
    createdAt
    updatedAt
    description
    isArchived
    isFork
    isLocked
    isDisabled
    diskUsage
    forks {
      totalCount
      totalDiskUsage
    }
    stargazers {
      totalCount
    }
    watchers {
      totalCount
    }
    releases {
      totalCount
    }
    releasesDynamics: releases(first: 20, orderBy: {field: CREATED_AT, direction: DESC}) {
      nodes {
        createdAt
      }
    }
    issuesOpen: issues(states: OPEN) {
      totalCount
    }
    issuesClosed: issues(states: CLOSED) {
      totalCount
    }
    issuesClosingDynamics: issues(states: CLOSED, first: 100, orderBy: {direction: DESC, field: CREATED_AT}) {
      nodes {
        createdAt
        closedAt
        author{
          login
        }
        labels(first: 10){
          nodes{
            name
          }
        }
      }
    }
    issuesOpeningDynamics: issues(states: OPEN, first: 100, orderBy: {direction: DESC, field: CREATED_AT}) {
      nodes {
        createdAt
        author{
          login
        }
        labels(first: 10){
          nodes{
            name
          }
        }
      }
    }`
	query, err := getGraphQLQuery(owner, repo, queryArgs)
	if err != nil {
		return nil, err
	}

	resp, err := c.fetchGraphQL(query)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &result)
	return
}

func (c *Client) getCommunityProfileMetrics(owner, repo string) (result *respCommunityScore, err error) {
	url := fmt.Sprintf("%s/repos/%s/%s/community/profile", apiBaseURL, owner, repo)
	resp, err := c.httpClient.GET(url, c.headers)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp.Data, &result)
	return
}

func (c *Client) getCommits(owner, repo string) (result []*respContributorsCommits, err error) {
	pageIndex := 1
	for {
		res, err := c.getCommitsPage(owner, repo, pageIndex)
		if err != nil {
			return nil, err
		}
		result = append(result, res...)
		if len(res) < defaultPagination {
			break
		}
		pageIndex++
	}
	return
}

func (c *Client) getCommitsPage(owner, repo string, page int) (result []*respContributorsCommits, err error) {
	url := fmt.Sprintf("%s/repos/%s/%s/commits?per_page=%d&page=%d", apiBaseURL, owner, repo, defaultPagination, page)
	resp, err := c.httpClient.GET(url, c.headers)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp.Data, &result)
	return
}
