# artifact-dl
A simple GitHub Action to grab all of the artifacts from a set of workflow runs by commit.

## Introduction
Sometimes in a workflow it'd be really handy to be able pull a bunch of build artifacts from previous builds for comparison. This action aims to make that as easy as providing a JSON list of the commits for which workflow runs you'd like to download artifacts for.

## Usage
Artifact-dl can either be run as a GitHub Action or a Docker Container.

```yaml
name: "Pull Some Artifacts"

on:
  push:
    branches: main

jobs:
  auto-scan:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v2
        with:
          fetch-depth: '0'

      - name: Get All Commits In Branch
        id: merge-commits
        uses: autamus/merge-commits@main
        
      - name: Pull All Artifacts In Branch Runs
        uses: autamus/artifact-dl@main
        with:
          git_token: ${{ secrets.GIT_TOKEN }}
          input_commits: '${{ steps.merge-commits.outputs.commits }}'
          # For Docker Container Actions Your Repository is Mounted at '/github/workspace'
          # So here we'll save our output to the 'artifacts' directory.
          output_path: '/github/workspace/artifacts'
```