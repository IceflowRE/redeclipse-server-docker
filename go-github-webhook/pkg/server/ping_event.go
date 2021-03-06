package server

import (
	"net/http"

	"github.com/IceflowRE/redeclipse-server-docker/pkg/server/utils"
	"github.com/google/go-github/v33/github"
)

func pingEvent(hrw http.ResponseWriter, req *http.Request) {
	utils.ResponseJSON(hrw, http.StatusOK, postResp{"Ping accepted", github.DeliveryID(req)})
}
