package server

import (
	"net/http"

	"github.com/IceflowRE/redeclipse-server-docker/pkg/server/utils"
	"github.com/IceflowRE/redeclipse-server-docker/pkg/structs"
	"github.com/IceflowRE/redeclipse-server-docker/pkg/updater"
)

type postResp struct {
	Message string `json:"message"`
	Guid    string `json:"guid"`
}

func EventHandler(updaterConfig *updater.Config, storage *structs.HashStorage, workDir string) func(hrw http.ResponseWriter, req *http.Request) {
	return func(hrw http.ResponseWriter, req *http.Request) {
		githubHeader := req.Context().Value("header").(*utils.GithubHeader)
		switch githubHeader.Event {
		case "ping":
			pingEvent(hrw, req)
		case "push":
			pushEvent(hrw, req, updaterConfig, storage, workDir)
		default:
			utils.SendErrorResponse(hrw, &utils.ErrResp{Status: http.StatusOK, Message: "'" + githubHeader.Event + "' is not supported, aborting"})
			return
		}
	}
}
