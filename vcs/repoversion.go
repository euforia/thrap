package vcs

import (
	"fmt"

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
	return fmt.Sprintf("%s-%d-%s", rver.Tag, rver.Count, rver.Hash.String()[:8])
}

// GetRepoVersion returns the latest tag, the count from that tag and the hash
// of the last commit
func GetRepoVersion(path string) RepoVersion {
	tagLabel := defaultVersionTag

	repo, err := git.PlainOpen(path)
	if err != nil {
		return RepoVersion{Tag: tagLabel}
	}

	tags, _ := repo.Tags()

	var rc *object.Commit

	var lastTag *plumbing.Reference
	for {
		t, err := tags.Next()
		if err != nil {
			// fmt.Println(err)
			// return RepoVersion{}
			break
		}

		lastTag = t
	}

	// lastTag, _ := tags.Next()
	if lastTag != nil {
		tagLabel = lastTag.Name().Short()
		to, _ := repo.TagObject(lastTag.Hash())
		if to != nil {
			rc, _ = to.Commit()
		}
	}

	cmtIter, _ := repo.CommitObjects()

	var c int

	if rc == nil {

		cmtIter.ForEach(func(arg1 *object.Commit) error {
			c++
			return nil
		})

	} else {
		// Calculate the total and mark where the tag was found
		var t int
		cmtIter.ForEach(func(arg1 *object.Commit) error {
			t++
			if arg1.ID().String() == rc.ID().String() {
				c = t
			}

			return nil
		})
		// Subtract tag count from total to get the count from the tag
		c = t - c
	}

	rp := RepoVersion{Tag: tagLabel, Count: c}

	head, err := repo.Head()
	if err == nil {
		rp.Hash = head.Hash()
	}

	return rp
}
