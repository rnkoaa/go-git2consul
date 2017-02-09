package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vkorehov/go-git2consul/config/mock"
	"github.com/vkorehov/go-git2consul/testutil"
	"gopkg.in/libgit2/git2go.v24"
)

func TestDiffStatus(t *testing.T) {
	gitRepo, cleanup := testutil.GitInitTestRepo(t)
	defer cleanup()

	repoConfig := mock.RepoConfig(gitRepo.Workdir())
	dstPath := filepath.Join(os.TempDir(), repoConfig.Name)

	localRepo, err := git.Clone(repoConfig.Url, dstPath, &git.CloneOptions{})
	if err != nil {
		t.Fatal(err)
	}

	repo := &Repository{
		Repository: localRepo,
		Config:     repoConfig,
	}

	h, err := repo.Head()
	if err != nil {
		t.Fatal(err)
	}

	oldRef := h.Target().String()

	// Push a commit to the repository
	testutil.GitCommitTestRepo(t)

	_, err = repo.Pull("master")
	if err != nil {
		t.Fatal(err)
	}

	deltas, err := repo.DiffStatus(oldRef)
	if err != nil {
		t.Fatal(err)
	}

	if len(deltas) == 0 {
		t.Fatal("Expected deltas from pull changes")
	}

	if deltas[0].Status != git.DeltaModified {
		t.Fatalf("Expected DeltaModified on %s", deltas[0].OldFile.Path)
	}

	// Cleanup
	defer func() {
		err = os.RemoveAll(repo.Workdir())
		if err != nil {
			t.Fatal(err)
		}
	}()
}
