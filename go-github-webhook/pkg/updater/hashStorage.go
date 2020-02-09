package updater

import (
	"encoding/json"
	"io/ioutil"
)

type HashStorage struct {
	Collection []hashContainer `json:"collection"`
}

type hashContainer struct {
	Branch string `json:"branch"`
	Arch   string `json:"arch"`
	Os     string `json:"os"`
	Hashes hash   `json:"hashes"`
}

type hash struct {
	Alpine     string `json:"alpine"`
	Dockerfile string `json:"dockerfile"`
	ReCommit   string `json:"re-commit"`
}

// Returns empty hash if no entry was found
func (storage *HashStorage) Get(branch string, arch string, os string) *hash {
	for _, hashC := range storage.Collection {
		if hashC.Branch == branch && hashC.Arch == arch && hashC.Os == os {
			return &hashC.Hashes
		}
	}
	return &hash{
		Alpine:     "",
		Dockerfile: "",
		ReCommit:   "",
	}
}

func (storage *HashStorage) Update(branch string, arch string, os string, newHash *hash) {
	curHash := storage.Get(branch, arch, os)
	if curHash != nil {
		curHash.Alpine = newHash.Alpine
		curHash.Dockerfile = newHash.Dockerfile
		curHash.ReCommit = newHash.ReCommit
		return
	}
	storage.Collection = append(storage.Collection, hashContainer{
		Branch: branch,
		Arch:   arch,
		Os:     os,
		Hashes: *newHash,
	})
}

func saveHashStorage(filename string, storage *HashStorage) error {
	file, _ := json.MarshalIndent(*storage, "", "    ")
	return ioutil.WriteFile(filename, file, 0644)
}

func loadHashStorage(file string) (*HashStorage, error) {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var storage HashStorage
	err = json.Unmarshal(raw, &storage)
	if err != nil {
		return nil, err
	}
	return &storage, nil
}
