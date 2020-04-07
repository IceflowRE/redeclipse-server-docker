package server

import (
	"net/http"

	"github.com/IceflowRE/redeclipse-server-docker/pkg/server/utils"
	"github.com/IceflowRE/redeclipse-server-docker/pkg/structs"
	"github.com/IceflowRE/redeclipse-server-docker/pkg/updater"
	"github.com/google/go-github/v30/github"
)

type postResp struct {
	Message string `json:"message"`
	Guid    string `json:"guid"`
}

var supportedEvents = [...]string{"ping", "push"}

func EventHandler(updaterConfig *updater.Config, storage *structs.HashStorage, workDir string) func(hrw http.ResponseWriter, req *http.Request) {
	return func(hrw http.ResponseWriter, req *http.Request) {
		switch payload := req.Context().Value("payload").(type) {
		case *github.PingEvent:
			pingEvent(hrw, req)
		case *github.CreateEvent:
			pushEvent(hrw, req, updaterConfig, storage, workDir, payload)
		default:
			utils.ResponseJSON(hrw, http.StatusOK, postResp{"'" + github.WebHookType(req) + "' is not supported,aborting", github.DeliveryID(req)})
			return
		}
	}
}
