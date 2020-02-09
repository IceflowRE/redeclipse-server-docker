package server

import (
	"net/http"
	"time"

	"github.com/IceflowRE/redeclipse-server-docker/pkg/server/utils"
)

func PingGet(hrw http.ResponseWriter, req *http.Request) {
	utils.ResponseJSON(hrw, http.StatusOK,
		map[string]string{
			"timestamp": time.Now().Format("2006-01-02T15:04:05.000000000Z07:00"),
		},
	)
}
