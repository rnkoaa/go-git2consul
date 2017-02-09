package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vkorehov/go-git2consul/config/mock"
	"github.com/vkorehov/go-git2consul/testutil"
	"gopkg.in/libgit2/git2go.v24"
)

func TestCheckRef(t *testing.T) {
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

	ref := h.Target().String()

	// Push a commit to the repository
	testutil.GitCommitTestRepo(t)

	_, err = repo.Pull("master")
	if err != nil {
		t.Fatal(err)
	}

	err = repo.CheckRef(ref)
	if err != nil {
		t.Fatal(err)
	}

	// Cleanup
	defer func() {
		err = os.RemoveAll(repo.Workdir())
		if err != nil {
			t.Fatal(err)
		}
	}()
}
