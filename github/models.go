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
}

type PullRequestNode struct {
	Node PullRequest `json:"node"`
}

type PullRequestEdges struct {
	Edges []PullRequestNode `json:"edges"`
}

type Repository struct {
	PullRequests PullRequestEdges `json:"pullRequests"`
}

type PullRequestResponseData struct {
	Repository Repository `json:"repository"`
}

// GraphqlError is one entry from a GraphQL "errors" block. GitHub returns
// these with HTTP 200, so they must be inspected rather than assumed absent.
type GraphqlError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type PullRequestResponse struct {
	Data   PullRequestResponseData `json:"data"`
	Errors []GraphqlError          `json:"errors,omitempty"`
}

type PullRequestQuery struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables,omitempty"`
}
