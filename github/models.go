package github

type Author struct {
	Login string `json:"login"`
}

type PullRequest struct {
	Title     string  `json:"title"`
	Number    int     `json:"number"`
	Url       string  `json:"url"`
	State     string  `json:"state"`
	Author    *Author `json:"author"`
	CreatedAt string  `json:"createdAt"`
	MergedAt  string  `json:"mergedAt"`
	ClosedAt  string  `json:"closedAt"`
	// HeadRefName is dependabot's branch, e.g.
	// "dependabot/npm_and_yarn/ui/web/lodash-4.17.21". The second segment is
	// the ecosystem, which is the only reliable way to attribute a pull
	// request to a package manager without a second API call.
	HeadRefName string `json:"headRefName"`
}

// AuthorLogin is the pull request's author, or an empty string when GitHub
// reports none (a deleted account, for example).
func (pullRequest PullRequest) AuthorLogin() string {
	if pullRequest.Author == nil {
		return ""
	}

	return pullRequest.Author.Login
}

// GraphqlError is one entry from a GraphQL "errors" block. GitHub returns
// these with HTTP 200, so they must be inspected rather than assumed absent.
type GraphqlError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type searchNode struct {
	Node PullRequest `json:"node"`
}

type searchResult struct {
	IssueCount int          `json:"issueCount"`
	Edges      []searchNode `json:"edges"`
}

type searchResponseData struct {
	Search searchResult `json:"search"`
}

// SearchResponse is the raw shape returned by the pull request search.
type SearchResponse struct {
	Data   searchResponseData `json:"data"`
	Errors []GraphqlError     `json:"errors,omitempty"`
}

// SearchQuery is a GraphQL request body.
type SearchQuery struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables,omitempty"`
}
