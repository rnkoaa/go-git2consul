package repository

import (
	"fmt"
	"gopkg.in/libgit2/git2go.v24"
)

// Clone the repository. Cloning will only checkout tracked branches.
// A destination path to clone to needs to be provided
func (r *Repository) Clone(path string) error {
	r.Lock()
	defer r.Unlock()

	// Clone the first tracked branch instead of the default branch
	if len(r.Config.Branches) == 0 {
		return fmt.Errorf("No tracked branches specified")
	}
	checkoutBranch := r.Config.Branches[0]

	repo, err := git.Clone(r.Config.Url, path, &git.CloneOptions{
		CheckoutOpts: &git.CheckoutOpts{
			Strategy: git.CheckoutNone,
		},
		CheckoutBranch: checkoutBranch,
	})
	if err != nil {
		return err
	}

	r.Repository = repo

	err = r.checkoutConfigBranches()
	if err != nil {
		return err
	}

        itr, err := repo.NewBranchIterator(git.BranchLocal)
        if err != nil {
                return err
        }
        defer itr.Free()

        var checkoutBranchFn = func(b *git.Branch, _ git.BranchType) error {
                branchName, err := b.Name()
                if err != nil {
                        return err
                }
                analysis, err := r.Pull(branchName)
                if err != nil {
                        return err
                }

                // If there is a change, send the repo RepoChangeCh
                switch {
                case analysis&git.MergeAnalysisUpToDate != 0:
                        fmt.Printf("Up to date: %s/%s\n", r.Name(), branchName)
                case analysis&git.MergeAnalysisNormal != 0, analysis&git.MergeAnalysisFastForward != 0:
                        fmt.Printf("Changed: %s/%s\n", r.Name(), branchName)
                }
                return nil
        }

        err = itr.ForEach(checkoutBranchFn)
        if err != nil && !git.IsErrorCode(err, git.ErrIterOver) {
                return err
        }
	return nil
}
