package vcs

import (
	"bytes"
	"fmt"
	"io"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

const (
	defaultVersionTag = "v0.0.0"
)

// RepoVersion is the calculate repo / project / stack version
type RepoVersion struct {
	// Last known tag.  Defaults to v0.0.0
	Tag string
	// Count from the tag
	Count int
	// Last commit hash
	Hash plumbing.Hash
}

func (rver RepoVersion) String() string {
	if rver.Count == 0 {
		return rver.Tag
	}
	return fmt.Sprintf("%s-%d-%s", rver.Tag, rver.Count, rver.Hash.String()[:8])
}

// returns nil if tags cannot be found
func getLatestTag(repo *git.Repository) (lastTag *plumbing.Reference) {
	tags, _ := repo.Tags()

	for {
		t, err := tags.Next()
		if err != nil {
			break
		}

		lastTag = t
	}

	return
}

// GetRepoVersion returns the latest tag, the count from that tag and the hash
// of the last commit
func GetRepoVersion(path string) RepoVersion {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return RepoVersion{Tag: defaultVersionTag}
	}
	ver, _ := getRepoVersion(repo)
	return ver
}

func getRepoVersion(repo *git.Repository) (RepoVersion, error) {
	rv := RepoVersion{Tag: defaultVersionTag}

	head, err := repo.Head()
	if err != nil {
		return rv, err
	}
	rv.Hash = head.Hash()

	iter, err := repo.Log(&git.LogOptions{From: rv.Hash})
	if err != nil {
		return rv, err
	}

	var (
		cmtHash   plumbing.Hash
		latestTag = getLatestTag(repo)
	)

	if latestTag != nil {
		rv.Tag = latestTag.Name().Short()

		to, _ := repo.TagObject(latestTag.Hash())
		if to != nil {
			// Annotated tag
			cmtHash = to.Target
		} else {
			// Tag with no annotation
			cmtHash = latestTag.Hash()
		}

	}

	iter.ForEach(func(c *object.Commit) error {
		if bytes.Compare(cmtHash[:], c.Hash[:]) == 0 {
			return io.EOF
		}
		rv.Count++
		return nil
	})

	return rv, nil
}
