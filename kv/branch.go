package kv

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/vkorehov/go-git2consul/repository"
	"github.com/apex/log"
	"github.com/hashicorp/consul/api"
	"gopkg.in/libgit2/git2go.v24"
)

// Push a repository branch to the KV
// TODO: Optimize for PUT only on changes instead of the entire repo
func (h *KVHandler) putBranch(repo *repository.Repository, branch *git.Branch) error {
	// Checkout branch
	repo.CheckoutBranch(branch, &git.CheckoutOpts{
		Strategy: git.CheckoutForce,
	})

	// h, _ := repo.Head()
	// bn, _ := h.Branch().Name()
	// log.Debugf("(consul) pushBranch(): Branch: %s Head: %s", bn, h.Target().String())

	var pushFile = func(fullpath string, info os.FileInfo, err error) error {
		// Walk error
		if err != nil {
			return err
		}

		// Skip the .git directory
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		// Do not push directories
		if info.IsDir() {
			return nil
		}

		key := strings.TrimPrefix(fullpath, repo.Workdir())
		kvPath := key
		h.logger.Debugf("KV PUT changes: %s", kvPath)

		data, err := ioutil.ReadFile(fullpath)
		if err != nil {
                        h.logger.Debugf("KV FAILED changes: '%s' %s", fullpath, kvPath)
			return err
		}

		p := &api.KVPair{
			Key:   kvPath,
			Value: data,
		}

		_, err = h.Put(p, nil)
		if err != nil {
			return err
		}

		return nil
	}

	err := filepath.Walk(repo.Workdir(), pushFile)
	if err != nil {
		log.WithError(err).Debug("PUT branch error")
	}

	return nil
}
