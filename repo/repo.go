package repo

import (
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

func GetURL(path string) (url string, err error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return url, err
	}
	remotes, err := r.Remotes()
	if err != nil {
		return url, err
	}
	return remotes[0].Config().URLs[0], nil
}

func GetOwnerName(path string) (repoOwner, repoName string, err error) {
	url, err := GetURL(path)
	if err != nil {
		return repoOwner, repoName, err
	}
	repoName = strings.TrimSuffix(filepath.Base(url), filepath.Ext(url))
	repoOwner = filepath.Base(filepath.Dir(url))
	return repoOwner, repoName, nil
}
