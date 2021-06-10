package repo

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
)

func GetRunArtifacts(path, gitToken string, runID int64) (artifactIDs []*github.Artifact, err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	owner, name, err := GetOwnerName(path)
	if err != nil {
		return
	}

	artifactList, _, err := client.Actions.ListWorkflowRunArtifacts(ctx, owner, name, runID, nil)
	artifactIDs = artifactList.Artifacts
	return
}

func DownloadArtifact(path, gitToken, outputPath string, workflowID int64, artifact *github.Artifact) (err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	owner, name, err := GetOwnerName(path)
	if err != nil {
		return err
	}

	url, _, err := client.Actions.DownloadArtifact(ctx, owner, name, *artifact.ID, true)
	if err != nil {
		return err
	}

	resp, err := http.Get(url.String())
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Zip archive reader inspired by answer from
	// https://stackoverflow.com/questions/50539118/golang-unzip-response-body/50539327
	reader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through files in compressed archive.
	for _, zipFile := range reader.File {
		buf, err := readCompressedFile(zipFile)
		if err != nil {
			log.Println(err)
			continue
		}

		err = os.MkdirAll(filepath.Dir(filepath.Join(outputPath, fmt.Sprintf("%d", workflowID), artifact.GetName(), zipFile.Name)), os.ModePerm)
		if err != nil {
			return err
		}

		if zipFile.FileInfo().IsDir() {
			continue
		}

		err = os.WriteFile(
			filepath.Join(outputPath, fmt.Sprintf("%d", workflowID), artifact.GetName(), zipFile.Name),
			buf,
			os.ModePerm,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func readCompressedFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
