package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"net/http"
)

func validMAC(message []byte, messageMAC string, key []byte) bool {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return messageMAC == "sha1="+hex.EncodeToString(expectedMAC)
}

func SignatureMiddle(webhookSecret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(hrw http.ResponseWriter, req *http.Request) {
			body, err := ioutil.ReadAll(req.Body)
			if PrintError(err) {
				SendErrorResponse(hrw, GenericErrorResp)
				return
			}

			githubHeader := req.Context().Value("header").(*GithubHeader)
			if !validMAC(body, githubHeader.Signature, webhookSecret) {
				SendErrorResponse(hrw, &ErrResp{http.StatusBadRequest, "X-Hub-Signature does not match."})
				return
			}
			req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			next.ServeHTTP(hrw, req)
		})
	}
}
