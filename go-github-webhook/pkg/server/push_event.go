package server

import (
	"net/http"

	"github.com/IceflowRE/redeclipse-server-docker/pkg/server/utils"
	"github.com/IceflowRE/redeclipse-server-docker/pkg/structs"
	"github.com/IceflowRE/redeclipse-server-docker/pkg/updater"
	"github.com/google/go-github/v30/github"
)

var updateManager = utils.NewUpdateCtx()

func pushEvent(hrw http.ResponseWriter, req *http.Request, updaterConfig *updater.Config, storage *structs.HashStorage, workDir string, payload *github.CreateEvent) {
	if payload.Ref == nil || *payload.Ref == "" {
		utils.SendErrorResponse(hrw, &utils.ErrResp{Status: http.StatusBadRequest, Message: "Ref is empty"})
		return
	}

	buildConfig := updaterConfig.Get(*payload.Ref)
	if buildConfig == nil {
		utils.ResponseJSON(hrw, http.StatusOK,
			postResp{"Payload was not for an accepted reference, aborting.", github.DeliveryID(req)},
		)
		return
	}
	utils.ResponseJSON(hrw, http.StatusCreated,
		postResp{"Update started for '" + *payload.Ref + "'", github.DeliveryID(req)},
	)

	if start := updateManager.Add(*payload.Ref); !start {
		return
	}
	go func() {
		for {
			updater.Build(updaterConfig, storage, buildConfig, workDir)
			if start := updateManager.Remove(*payload.Ref); !start {
				break
			}
		}
	}()
	return
}
