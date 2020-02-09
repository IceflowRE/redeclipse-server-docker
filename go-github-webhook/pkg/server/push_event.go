package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/IceflowRE/redeclipse-server-docker/pkg/server/utils"
	"github.com/IceflowRE/redeclipse-server-docker/pkg/updater"
	"github.com/google/go-github/github"
)


func bindPushEvent(req *http.Request) (*github.PushEvent, *utils.ErrResp) {
	body, err := ioutil.ReadAll(req.Body)
	if utils.PrintError(err) {
		return nil, utils.GenericErrorResp
	}

	var obj github.PushEvent
	err = json.Unmarshal(body, &obj)
	if utils.PrintError(err) {
		return nil, &utils.ErrResp{Status: http.StatusBadRequest, Message: "Malformed body."}
	}
	switch {
	case obj.Ref == nil:
		return nil, &utils.ErrResp{Status: http.StatusBadRequest, Message: "Ref is empty."}
	}
	return &obj, nil
}

var branchRx = regexp.MustCompile(`^refs/heads/(\w+)$`)

func pushEvent(hrw http.ResponseWriter, req *http.Request, updaterConfig *updater.AppConfig, storage *updater.HashStorage, buildCtx *updater.BuildContext) {
	payload, err := bindPushEvent(req)
	if err != nil {
		utils.SendErrorResponse(hrw, err)
		return
	}
	githubHeader := req.Context().Value("header").(*utils.GithubHeader)
	branchMatch := branchRx.FindStringSubmatch(*payload.Ref)

	if len(branchMatch) == 2 && update(updaterConfig, storage, buildCtx, branchMatch[1]) {
		utils.ResponseJSON(hrw, http.StatusCreated,
			postResp{"Update started for '" + branchMatch[1] + "'", githubHeader.Guid},
		)
	} else {
		utils.ResponseJSON(hrw, http.StatusOK,
			postResp{"Payload was not for an accepted branch, aborting.", githubHeader.Guid},
		)
	}
}
