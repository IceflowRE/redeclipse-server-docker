package server

import (
	"net/http"

	"github.com/IceflowRE/redeclipse-server-docker/pkg/server/utils"
)

func pingEvent(hrw http.ResponseWriter, req *http.Request) {
	githubHeader := req.Context().Value("header").(*utils.GithubHeader)
	utils.ResponseJSON(hrw, http.StatusOK,
		postResp{"Ping accepted", githubHeader.Guid},
	)
}
