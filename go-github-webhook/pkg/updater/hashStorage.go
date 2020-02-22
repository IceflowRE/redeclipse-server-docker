package updater

import (
	"encoding/json"
	"io/ioutil"
)

type HashStorage struct {
	Collection []hashContainer `json:"collection"`
}

type hashContainer struct {
	Ref    string `json:"ref"`
	Arch   string `json:"arch"`
	Os     string `json:"os"`
	Hashes *hash  `json:"hashes"`
}

type hash struct {
	Alpine     string `json:"alpine"`
	Dockerfile string `json:"dockerfile"`
	ReCommit   string `json:"re-commit"`
}

func NewHashStorage() *HashStorage {
	return &HashStorage{
		Collection: make([]hashContainer, 0),
	}
}

// if not found it will create a new default value
// never returns nil
func (storage *HashStorage) Get(ref string, arch string, os string) *hash {
	for _, hashC := range storage.Collection {
		if hashC.Ref == ref && hashC.Arch == arch && hashC.Os == os {
			return hashC.Hashes
		}
	}
	newHash := &hash{
		Alpine:     "",
		Dockerfile: "",
		ReCommit:   "",
	}
	storage.Collection = append(storage.Collection, hashContainer{
		Ref:    ref,
		Arch:   arch,
		Os:     os,
		Hashes: newHash,
	})
	return newHash
}

func (storage *HashStorage) Update(ref string, arch string, os string, newHash *hash) {
	curHash := storage.Get(ref, arch, os)
	*curHash = *newHash
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
	storage := NewHashStorage()
	err = json.Unmarshal(raw, storage)
	if err != nil {
		return nil, err
	}
	return storage, nil
}
