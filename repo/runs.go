package repo

import (
	"context"
	"strings"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
)

func GetWorkflowRuns(path, gitToken string, commits []string) (result []int64, err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	owner, name, err := GetOwnerName(path)
	if err != nil {
		return result, err
	}

	// Setup Initial Values for Loop
	page := 1
	count := 0

	// Loop through pages of results.
	for count < len(commits) {
		runs, _, err := client.Actions.ListRepositoryWorkflowRuns(
			ctx,
			owner,
			name,
			&github.ListWorkflowRunsOptions{
				ListOptions: github.ListOptions{
					Page: page,
				},
			})
		if err != nil {
			return result, err
		}

		// Loop through list of workflow runs in page.
		for _, run := range runs.WorkflowRuns {
			for _, entry := range commits {
				if strings.HasPrefix(*run.HeadSHA, entry) {
					count++
					result = append(result, *run.ID)
				}
			}
		}
		page++
	}
	return result, nil
}
