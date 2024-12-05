package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	graphqlEndpoint = "https://api.github.com/graphql"
)

func GetPullRequests(owner string, repo string, token string) PullRequestResponse {
	queryContent := fmt.Sprintf("query {\n  repository(owner: \"%s\", name: \"%s\") {\n    pullRequests(states: [OPEN, CLOSED, MERGED], first: 100, labels: \"dependencies\") {\n      edges {\n        node {\n          title\n          number\n          url\n     state\n      author {\n            login\n          }\n          createdAt\n        }\n      }\n    }\n  }\n}", owner, repo)
	query := PullRequestQuery{
		Query: queryContent,
	}
	return GetPullRequestsFromQuery(query, token)
}

func GetPullRequestsFromQuery(query PullRequestQuery, token string) PullRequestResponse {
	jsonData, _ := json.Marshal(query)
	req, _ := http.NewRequest("POST", graphqlEndpoint, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result PullRequestResponse
	json.Unmarshal(body, &result)

	return result
}
