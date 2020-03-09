package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func ApiKeyMiddleware(validKeys []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(hrw http.ResponseWriter, req *http.Request) {
			body, err := ioutil.ReadAll(req.Body)
			if PrintError(err) {
				SendErrorResponse(hrw, GenericErrorResp)
				return
			}

			key, ok := req.Header["X-Iceflower-Apikey"]
			if !ok {
				SendErrorResponse(hrw, &ErrResp{400, "X-Iceflower-Apikey header is missing"})
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
				SendErrorResponse(hrw, &ErrResp{401, "invalid api key"})
				return
			}

			req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			next.ServeHTTP(hrw, req)
		})
	}
}
