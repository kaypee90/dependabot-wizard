package github

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestGetPullRequests(t *testing.T) {
	token := ""
	if token == "" {
		t.Skip("Skipping explicit test. Set gitub token to run this test.")
	}
	response := GetPullRequests("kaypee90", "zebrok", token)
	assert.Equal(t, len(response.Data.Repository.PullRequests.Edges) > 0, true)
}
