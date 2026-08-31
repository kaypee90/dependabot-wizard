package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	defaultGraphqlEndpoint = "https://api.github.com/graphql"
	requestTimeout         = 30 * time.Second
	maxResults             = 100
)

const (
	// DefaultAuthor identifies dependabot in GitHub's search syntax. Filtering
	// on the author is what actually distinguishes dependabot's pull requests:
	// the "dependencies" label is applied by default but is not guaranteed —
	// a `labels:` key in dependabot.yml replaces the defaults, and label
	// creation can be blocked entirely, leaving the pull requests unlabelled.
	DefaultAuthor = "app/dependabot"

	// DefaultLabel is dependabot's default label. It is no longer applied
	// automatically; pass it explicitly to narrow a search further.
	DefaultLabel = "dependencies"
)

// Filter narrows which pull requests are counted. The zero value means "every
// pull request dependabot opened", which is what the dashboard and the report
// want. An empty Author disables author filtering entirely.
type Filter struct {
	Author string
	Label  string
}

// DependabotFilter is the default: authored by dependabot, any label.
func DependabotFilter() Filter {
	return Filter{Author: DefaultAuthor}
}

// endpoint is a variable so tests can point the client at a stub server.
var endpoint = defaultGraphqlEndpoint

const searchQueryTemplate = `query($search: String!) {
  search(query: $search, type: ISSUE, first: %d) {
    issueCount
    edges {
      node {
        ... on PullRequest {
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

// quoteSearchTerm wraps a value so a label containing spaces stays one term.
func quoteSearchTerm(value string) string {
	return `"` + strings.ReplaceAll(value, `"`, ``) + `"`
}

// BuildSearchQuery assembles the GitHub search expression. Search is used
// rather than repository.pullRequests because it filters by author server
// side; filtering client side would mean the first 100 pull requests in a
// busy repository could contain no dependabot ones at all.
func BuildSearchQuery(owner string, repo string, filter Filter) string {
	terms := []string{
		"repo:" + owner + "/" + repo,
		"is:pr",
		"sort:created-desc",
	}

	if author := strings.TrimSpace(filter.Author); author != "" {
		terms = append(terms, "author:"+author)
	}

	if label := strings.TrimSpace(filter.Label); label != "" {
		terms = append(terms, "label:"+quoteSearchTerm(label))
	}

	return strings.Join(terms, " ")
}

// Describe renders the filter for a human, e.g. for an empty-result message.
func (filter Filter) Describe() string {
	author := strings.TrimSpace(filter.Author)
	label := strings.TrimSpace(filter.Label)

	switch {
	case author != "" && label != "":
		return fmt.Sprintf("opened by %s and labelled %q", author, label)
	case author != "":
		return fmt.Sprintf("opened by %s", author)
	case label != "":
		return fmt.Sprintf("labelled %q", label)
	default:
		return "in this repository"
	}
}

// GetPullRequests returns the matching pull requests, newest first.
func GetPullRequests(ctx context.Context, owner string, repo string, filter Filter, token string) ([]PullRequest, error) {
	query := SearchQuery{
		Query: fmt.Sprintf(searchQueryTemplate, maxResults),
		Variables: map[string]string{
			"search": BuildSearchQuery(owner, repo, filter),
		},
	}

	return GetPullRequestsFromQuery(ctx, query, token)
}

func GetPullRequestsFromQuery(ctx context.Context, query SearchQuery, token string) ([]PullRequest, error) {
	jsonData, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("encoding query: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling GitHub: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("GitHub rejected the token (401); check that it is valid and has repo scope")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub returned %s", resp.Status)
	}

	var result SearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	// A GraphQL error arrives with HTTP 200 and an empty data block, so it has
	// to be checked explicitly or the caller sees "no pull requests".
	if len(result.Errors) > 0 {
		messages := make([]string, 0, len(result.Errors))
		for _, item := range result.Errors {
			messages = append(messages, item.Message)
		}
		return nil, fmt.Errorf("GitHub reported: %s", strings.Join(messages, "; "))
	}

	pullRequests := make([]PullRequest, 0, len(result.Data.Search.Edges))
	for _, edge := range result.Data.Search.Edges {
		// Search can return issues too; a pull request always has a number.
		if edge.Node.Number == 0 {
			continue
		}
		pullRequests = append(pullRequests, edge.Node)
	}

	sort.SliceStable(pullRequests, func(i, j int) bool {
		return pullRequests[i].CreatedAt > pullRequests[j].CreatedAt
	})

	return pullRequests, nil
}
