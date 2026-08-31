package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultGraphqlEndpoint = "https://api.github.com/graphql"
	requestTimeout         = 30 * time.Second
)

const pullRequestQueryTemplate = `query($owner: String!, $repo: String!, $label: String!) {
  repository(owner: $owner, name: $repo) {
    pullRequests(states: [OPEN, CLOSED, MERGED], first: 100, labels: [$label]) {
      edges {
        node {
          title
          number
          url
          state
          author {
            login
          }
          createdAt
        }
      }
    }
  }
}`

// endpoint is a variable so tests can point the client at a stub server.
var endpoint = defaultGraphqlEndpoint

// DefaultLabel is the label dependabot applies to the pull requests it opens.
const DefaultLabel = "dependencies"

// GetPullRequests reads labelled pull requests for a repository. Owner, repo
// and label travel as GraphQL variables rather than being interpolated into
// the query string, so a value containing a quote cannot alter the query.
func GetPullRequests(ctx context.Context, owner string, repo string, label string, token string) (PullRequestResponse, error) {
	if label == "" {
		label = DefaultLabel
	}

	query := PullRequestQuery{
		Query: pullRequestQueryTemplate,
		Variables: map[string]string{
			"owner": owner,
			"repo":  repo,
			"label": label,
		},
	}

	return GetPullRequestsFromQuery(ctx, query, token)
}

func GetPullRequestsFromQuery(ctx context.Context, query PullRequestQuery, token string) (PullRequestResponse, error) {
	var result PullRequestResponse

	jsonData, err := json.Marshal(query)
	if err != nil {
		return result, fmt.Errorf("encoding query: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return result, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("calling GitHub: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return result, fmt.Errorf("GitHub rejected the token (401); check that it is valid and has repo scope")
	}

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("GitHub returned %s", resp.Status)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return result, fmt.Errorf("decoding response: %w", err)
	}

	// A GraphQL error arrives with HTTP 200 and an empty data block, so it has
	// to be checked explicitly or the caller sees "no pull requests" instead.
	if len(result.Errors) > 0 {
		messages := make([]string, 0, len(result.Errors))
		for _, item := range result.Errors {
			messages = append(messages, item.Message)
		}
		return result, fmt.Errorf("GitHub reported: %s", strings.Join(messages, "; "))
	}

	return result, nil
}
