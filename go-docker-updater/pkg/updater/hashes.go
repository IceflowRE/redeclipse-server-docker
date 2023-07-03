package updater

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/felicianotech/sonar/sonar/docker"
)

type Hash struct {
	Alpine     string
	Dockerfile string
	ReCommit   string
}

func getFileHash(filename string) *string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		fmt.Println(err.Error())
		return nil
	}
	res := "sha256:" + hex.EncodeToString(hasher.Sum(nil))
	return &res
}

type baseResp struct {
	Images []image `json:"images"`
}

type image struct {
	Os     string `json:"os"`
	Arch   string `json:"architecture"`
	Digest string `json:"digest"`
}

func getAlpineHash(arch string, os string) *string {
	resp, err := http.Get("https://registry.hub.docker.com/v2/repositories/library/alpine/tags/latest")
	if err != nil {
		return nil
	}
	if resp.StatusCode != 200 {
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	var respUn baseResp
	err = json.Unmarshal(body, &respUn)
	for _, image := range respUn.Images {
		if image.Arch == arch && image.Os == os {
			return &image.Digest
		}
	}
	fmt.Println("base image for "+arch+", ", os+" not found")
	return nil
}

func getCommitHash(ref string) *string {
	out, err := exec.Command("git", "ls-remote", "https://github.com/redeclipse/base.git", ref).Output()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	hash := string(bytes.Split(out, []byte("\t"))[0])
	if hash == "" {
		fmt.Println("reference '" + ref + "' not found.")
		return nil
	}
	return &hash
}

func GetNewHashes(dockerfile string, ref string, arch string, os string) *Hash {
	reCommit := getCommitHash(ref)
	if reCommit == nil {
		fmt.Println("failed to get git commit hash")
		return nil
	}
	dockerHash := getFileHash(dockerfile)
	if dockerHash == nil {
		fmt.Println("dockerfile hash failed " + dockerfile)
		return nil
	}
	alpineHash := getAlpineHash(arch, os)
	if alpineHash == nil {
		fmt.Println("alpine hash failed")
		return nil
	}
	return &Hash{
		Alpine:     *alpineHash,
		Dockerfile: *dockerHash,
		ReCommit:   *reCommit,
	}
}

func GetCurrentHashes(dockerRepo string, ref string, arch string) (*Hash, error) {
	labels, err := docker.GetLabels(dockerRepo + ":" + arch + "-" + DockerTagFromRef(ref))
	if err != nil {
		return nil, err
	}

	return &Hash{
		ReCommit:   labels["re-commit"],
		Alpine:     labels["alpine-sha"],
		Dockerfile: labels["dockerfile-sha"],
	}, nil
}
