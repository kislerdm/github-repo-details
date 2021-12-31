/* Module to fetch details of a github repo using github API and github graphQL. */
package githubrepodetails

import (
	"encoding/json"
	"fmt"
	"time"

	httpClient "github.com/kislerdm/github-repo-details/internal/http"
)

// Client defines the client to fetch data.
type Client struct {
	// github API key
	apikey     string
	httpClient *httpClient.Client
}

// ClientOption defines the option to modify Client's behavior.
type ClientOption func(*Client)

const defaultTimeout = 2 * time.Second

// NewClient defines the function to init a Client.
func NewClient(opts ...ClientOption) (*Client, error) {
	httpClient, err := httpClient.NewClient(httpClient.WithTimeout(defaultTimeout))
	if err != nil {
		return nil, err
	}
	c := &Client{
		apikey:     "",
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
		c.apikey = key
	}
}

// Response defines the response object.
type Response struct {
	// Status response http status code.
	Status int8
	// Details
	Details *RepoDetails
}

// Run defines the function to fetch repo details.
// Requires
// owner - the repo owner
// repo - the repo name
func (c *Client) Run(owner, repo string) (Response, error) {
	return Response{}, nil
}

func getGraphQLQuery(owner, repo string) ([]byte, error) {
	q := map[string]string{
		"query": fmt.Sprintf(`{
			repository(name: "%s", owner: "%s") {
    url
    licenseInfo {
      name
    }
    databaseId
    fullName: nameWithOwner
    createdAt
    updatedAt
    description
    isArchived
    isFork
    isLocked
    isDisabled
    diskUsage
    forkCount
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
      }
    }
    issuesOpeningDynamics: issues(states: OPEN, first: 100, orderBy: {direction: DESC, field: CREATED_AT}) {
      nodes {
        createdAt
      }
    }
  }
}`, repo, owner),
	}
	return json.Marshal(q)
}

func (c *Client) GetGraphQLData(owner, repo string) (result *respGraphQL, err error) {
	headers := map[string]string{
		"Content-Type":   "application/json",
		"X-REQUEST-TYPE": "GraphQL",
	}
	if c.apikey != "" {
		headers["Authorization"] = fmt.Sprintf("bearer %s", c.apikey)
	}

	query, err := getGraphQLQuery(owner, repo)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.POST("https://api.github.com/graphql", headers, query)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp.Data, &result)
	return
}
