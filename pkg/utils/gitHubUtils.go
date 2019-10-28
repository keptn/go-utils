package utils

import (
	"fmt"
	"os"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// Checkout clons a GitHub repo and checks out the specified branch
func Checkout(gitHubOrg string, project string, branch string) (*git.Repository, error) {

	err := os.RemoveAll(project)
	if err != nil {
		return nil, err
	}

	var repo *git.Repository
	if os.Getenv("GITHUB_USERNAME") != "" && os.Getenv("GITHUB_TOKEN") != "" {
		// If credentials are available, use them
		repo, err = git.PlainClone(project, false, &git.CloneOptions{
			URL: "https://github.com/" + gitHubOrg + "/" + project + ".git",
			Auth: &http.BasicAuth{
				Username: os.Getenv("GITHUB_USERNAME"), // anything except an empty string
				Password: os.Getenv("GITHUB_TOKEN"),
			},
		})
	} else {
		repo, err = git.PlainClone(project, false, &git.CloneOptions{
			URL: "https://github.com/" + gitHubOrg + "/" + project + ".git",
		})
	}
	w, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	err = repo.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	})
	if err != nil {
		return nil, err
	}

	// Checking out branch
	return repo, w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		Force:  true,
	})
}

// CheckoutPrevCommit moves the HEAD pointer to the previous commit. It returns the original HEAD.
func CheckoutPrevCommit(repo *git.Repository) (*plumbing.Reference, error) {

	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}

	_, err = commitIter.Next()
	if err != nil {
		return nil, err
	}
	c, err := commitIter.Next()
	if err != nil {
		return nil, err
	}

	w, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	return ref, w.Checkout(&git.CheckoutOptions{
		Hash: c.Hash,
	})
}

// CheckoutReference moves the HEAD pointer to the specified reference.
func CheckoutReference(repo *git.Repository, ref *plumbing.Reference) error {

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	return w.Checkout(&git.CheckoutOptions{
		Hash: ref.Hash(),
	})
}
