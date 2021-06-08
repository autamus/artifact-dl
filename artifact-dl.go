package main

import (
	"encoding/json"
	"log"

	"github.com/autamus/artifact-dl/config"
	"github.com/autamus/artifact-dl/repo"
)

func main() {
	// Initialize & Parse Commits from Input JSON String
	commits := []string{}
	err := json.Unmarshal([]byte(config.Global.Input.Commits), &commits)
	if err != nil {
		log.Fatal(err)
	}

	// Get workflow run IDs from commits.
	runIDs, err := repo.GetWorkflowRuns(
		config.Global.Repo.Path,
		config.Global.Git.Token,
		commits,
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, run := range runIDs {
		artifactIDs, err := repo.GetRunArtifacts(
			config.Global.Repo.Path,
			config.Global.Git.Token,
			run,
		)
		if err != nil {
			log.Fatal(err)
		}

		for _, artifact := range artifactIDs {
			err = repo.DownloadArtifact(
				config.Global.Repo.Path,
				config.Global.Git.Token,
				config.Global.Output.Path,
				run,
				artifact,
			)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
