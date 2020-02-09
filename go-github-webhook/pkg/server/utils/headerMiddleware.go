package utils

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type GithubHeader struct {
	Event     string
	Guid      string
	Signature string
}

func getGithubHeader(header *http.Header) (*GithubHeader, error) {
	event, ok := (*header)["X-Github-Event"]
	if !ok {
		return nil, errors.New("X-Github-Event header is missing")
	}
	guid, ok := (*header)["X-Github-Delivery"]
	if !ok {
		return nil, errors.New("X-Github-Delivery header is missing")
	}
	sign, ok := (*header)["X-Hub-Signature"]
	if !ok {
		return nil, errors.New("X-Hub-Signature header is missing")
	}
	return &GithubHeader{strings.Join(event, ""), strings.Join(guid, ""), strings.Join(sign, "")}, nil
}

func HeaderMiddle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(hrw http.ResponseWriter, req *http.Request) {
		header, err := getGithubHeader(&req.Header)
		if err != nil {
			SendErrorResponse(hrw, &ErrResp{http.StatusBadRequest, err.Error()})
			return
		}

		ctx := context.WithValue(req.Context(), "header", header)
		next.ServeHTTP(hrw, req.WithContext(ctx))
	})
}
