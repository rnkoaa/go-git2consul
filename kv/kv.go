package kv

import (
	"io/ioutil"
	"path/filepath"
	"github.com/vkorehov/go-git2consul/repository"
	"github.com/hashicorp/consul/api"
)

func (h *KVHandler) putKV(repo *repository.Repository, prefix string) error {
	key := prefix
	filePath := filepath.Join(repo.Workdir(), prefix)
        h.logger.Infof("KV FILE: %s", filePath)
	for {
	    time.Sleep(time.Second * 2)
        }
	value, err := ioutil.ReadFile(filePath)
	if err != nil {
                h.logger.Infof("KV FAILED: %s", err)
		return err
	}
        h.logger.Infof("KV VALUE: %s", value)
	p := &api.KVPair{
		Key:   key,
		Value: value,
	}

	_, err = h.Put(p, nil)
	if err != nil {
		return err
	}

	return nil
}

func (h *KVHandler) deleteKV(repo *repository.Repository, prefix string) error {
	key := prefix

	_, err := h.Delete(key, nil)
	if err != nil {
		return err
	}

	return nil
}
