package server

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/IceflowRE/redeclipse-server-docker/pkg/server/utils"
	"github.com/google/go-github/v30/github"
)

func githubMiddleware(secretToken []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(hrw http.ResponseWriter, req *http.Request) {
			payload, err := github.ValidatePayload(req, secretToken)
			if err != nil {
				utils.SendErrorResponse(hrw, &utils.ErrResp{Status: http.StatusBadRequest, Message: err.Error()})
				return
			}
			sign, ok := req.Header["X-Hub-Signature"]
			if !ok {
				utils.SendErrorResponse(hrw, &utils.ErrResp{Status: http.StatusBadRequest, Message: "X-Hub-Signature header is missing"})
				return
			}
			signature := strings.Join(sign, "")
			if err = github.ValidateSignature(signature, payload, secretToken); err != nil {
				utils.SendErrorResponse(hrw, &utils.ErrResp{Status: http.StatusBadRequest, Message: err.Error()})
				return
			}
			event := github.WebHookType(req)
			accepted := false
			for _, curEvent := range supportedEvents {
				if event == curEvent {
					accepted = true
					break
				}
			}
			if !accepted {
				utils.ResponseJSON(hrw, http.StatusOK,
					postResp{"'" + event + "' is not supported, aborting", github.DeliveryID(req)},
				)
				return
			}

			parsedPayload, err := github.ParseWebHook(github.WebHookType(req), payload)
			if err != nil {
				log.Println("ERROR: could not parse: " + event)
				utils.SendErrorResponse(hrw, &utils.ErrResp{Status: http.StatusBadRequest, Message: err.Error()})
				return
			}

			ctx := context.WithValue(req.Context(), "payload", parsedPayload)
			next.ServeHTTP(hrw, req.WithContext(ctx))
		})
	}
}
