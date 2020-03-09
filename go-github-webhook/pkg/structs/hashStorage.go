package structs

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
)

type HashStorage struct {
	Collection []HashContainer `json:"collection"`
	mux        sync.RWMutex    `json:"-"`
	file       string
	url        string
	apiKey     string
}

type HashContainer struct {
	Ref    string `json:"ref"`
	Arch   string `json:"arch"`
	Os     string `json:"os"`
	Hashes *Hash  `json:"hashes"`
}

type Hash struct {
	Alpine     string `json:"alpine"`
	Dockerfile string `json:"dockerfile"`
	ReCommit   string `json:"re-commit"`
}

func NewLocalStorage(file string) *HashStorage {
	return &HashStorage{
		Collection: make([]HashContainer, 0),
		file:       file,
	}
}

func NewOnlineStorage(url string, apiKey string) *HashStorage {
	return &HashStorage{
		Collection: make([]HashContainer, 0),
		url:        url,
		apiKey:     apiKey,
	}
}

func (storage *HashStorage) IsLocal() bool {
	return storage.file != ""
}

// creates a new entry, if none exists, else returns existing
func (storage *HashStorage) createLocal(ref string, arch string, os string) (*Hash, error) {
	if curHash, _ := storage.GetLocal(ref, arch, os); curHash != nil {
		return curHash, nil
	}
	newHashCont := HashContainer{
		Ref:  ref,
		Arch: arch,
		Os:   os,
		Hashes: &Hash{
			Alpine:     "",
			Dockerfile: "",
			ReCommit:   "",
		},
	}
	storage.mux.Lock()
	defer storage.mux.Unlock()
	storage.Collection = append(storage.Collection, newHashCont)
	return newHashCont.Hashes, nil
}

func (storage *HashStorage) Get(ref string, arch string, os string) (*Hash, error) {
	if storage.IsLocal() {
		return storage.GetLocal(ref, arch, os)
	} else {
		return storage.getOnline(ref, arch, os)
	}
}

func (storage *HashStorage) GetLocal(ref string, arch string, os string) (*Hash, error) {
	storage.mux.RLock()
	defer storage.mux.RUnlock()
	for _, hashC := range storage.Collection {
		if hashC.Ref == ref && hashC.Arch == arch && hashC.Os == os {
			return hashC.Hashes, nil
		}
	}
	return nil, nil
}

func (storage *HashStorage) getOnline(ref string, arch string, os string) (*Hash, error) {
	resp, err := hashRequest(&HashContainer{ref, arch, os, nil}, storage.url, storage.apiKey)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read body: " + err.Error())
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("failed to get a good response: " + string(body))
	}
	var hashCont HashContainer
	if err = json.Unmarshal(body, &hashCont); err != nil {
		return nil, errors.New("failed to unmarshal request")
	}
	return hashCont.Hashes, nil
}

func (storage *HashStorage) Update(ref string, arch string, os string, newHash *Hash) error {
	if storage.IsLocal() {
		return storage.UpdateLocal(ref, arch, os, newHash)
	} else {
		return storage.updateOnline(ref, arch, os, newHash)
	}
}

func (storage *HashStorage) UpdateLocal(ref string, arch string, os string, newHash *Hash) error {
	curHash, _ := storage.createLocal(ref, arch, os)
	storage.mux.Lock()
	defer storage.mux.Unlock()
	*curHash = *newHash
	return nil
}

func (storage *HashStorage) updateOnline(ref string, arch string, os string, newHash *Hash) error {
	resp, err := hashRequest(&HashContainer{ref, arch, os, newHash}, storage.url, storage.apiKey)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("failed to get a good response")
	}
	return nil
}

// skip if online hash storage
func (storage *HashStorage) SaveToFile() error {
	if !storage.IsLocal() {
		return nil
	}
	storage.mux.Lock()
	defer storage.mux.Unlock()
	file, _ := json.MarshalIndent(*storage, "", "    ")
	return ioutil.WriteFile(storage.file, file, 0644)
}

// skip if not local
func (storage *HashStorage) LoadFromFile() (err error) {
	if !storage.IsLocal() || !FileExists(storage.file) {
		return nil
	}
	var raw []byte
	if raw, err = ioutil.ReadFile(storage.file); err != nil {
		return err
	}
	err = json.Unmarshal(raw, storage)
	return err
}

func hashRequest(hashCont *HashContainer, url string, apiKey string) (*http.Response, error) {
	data, err := json.Marshal(hashCont)
	if err != nil {
		return nil, errors.New("failed to marshal request")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, errors.New("failed to create new request")
	}
	req.Header.Set("X-Iceflower-Apikey", apiKey)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("failed to do request: " + err.Error())
	}
	return resp, nil
}
