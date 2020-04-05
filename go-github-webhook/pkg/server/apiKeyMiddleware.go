package server

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/IceflowRE/redeclipse-server-docker/pkg/server/utils"
)

func apiKeyMiddleware(validKeys []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(hrw http.ResponseWriter, req *http.Request) {
			body, err := ioutil.ReadAll(req.Body)
			if utils.PrintError(err) {
				utils.SendErrorResponse(hrw, utils.GenericErrorResp)
				return
			}

			key, ok := req.Header["X-Iceflower-Apikey"]
			if !ok {
				utils.SendErrorResponse(hrw, &utils.ErrResp{400, "X-Iceflower-Apikey header is missing"})
				return
			}
			isValid := false
			for _, curKey := range validKeys {
				if curKey == key[0] {
					isValid = true
					break
				}
			}
			if !isValid {
				utils.SendErrorResponse(hrw, &utils.ErrResp{401, "invalid api key"})
				return
			}

			req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			next.ServeHTTP(hrw, req)
		})
	}
}
