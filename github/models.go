package github

type Author struct {
	Login string `json:"login"`
}

type PullRequest struct {
	Title  string  `json:"title"`
	Number int     `json:"number"`
	Url    string  `json:"url"`
	State  string  `json:"state"`
	Author *Author `json:"user"`
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

type PullRequestResponse struct {
	Data PullRequestResponseData `json:"data"`
}

type PullRequestQuery struct {
	Query string `json:"query"`
}
