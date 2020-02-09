package server

import (
	"net/http"

	"github.com/IceflowRE/redeclipse-server-docker/pkg/server/utils"
	"github.com/IceflowRE/redeclipse-server-docker/pkg/updater"
)

type postResp struct {
	Message string `json:"message"`
	Guid    string `json:"guid"`
}

func EventHandler(updaterConfig *updater.AppConfig, storage *updater.HashStorage, buildCtx *updater.BuildContext) func(hrw http.ResponseWriter, req *http.Request) {
	return http.HandlerFunc(func(hrw http.ResponseWriter, req *http.Request) {
		githubHeader := req.Context().Value("header").(*utils.GithubHeader)
		switch githubHeader.Event {
		case "ping":
			pingEvent(hrw, req)
		case "push":
			pushEvent(hrw, req, updaterConfig, storage, buildCtx)
		default:
			utils.SendErrorResponse(hrw, &utils.ErrResp{Status: http.StatusOK, Message: "'" + githubHeader.Event + "' is not supported, aborting"})
			return
		}
	})
}
