package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/IceflowRE/redeclipse-server-docker/pkg/server/utils"
	"github.com/IceflowRE/redeclipse-server-docker/pkg/structs"
)

func bindHashPost(req *http.Request) (*structs.HashContainer, *utils.ErrResp) {
	body, err := ioutil.ReadAll(req.Body)
	if utils.PrintError(err) {
		return nil, utils.GenericErrorResp
	}

	var obj structs.HashContainer
	err = json.Unmarshal(body, &obj)
	if utils.PrintError(err) {
		return nil, &utils.ErrResp{Status: http.StatusBadRequest, Message: "Malformed body."}
	}
	return &obj, nil
}

func HashPost(storage *structs.HashStorage) func(hrw http.ResponseWriter, req *http.Request) {
	return func(hrw http.ResponseWriter, req *http.Request) {
		payload, err := bindHashPost(req)
		if err != nil {
			utils.SendErrorResponse(hrw, err)
			return
		}

		// if update
		if payload.Hashes != nil {
			if err := storage.UpdateLocal(payload.Ref, payload.Arch, payload.Os, payload.Hashes); err != nil {
				utils.SendErrorResponse(hrw, &utils.ErrResp{Status: 500, Message: err.Error()})
			}
			if err := storage.SaveToFile(); err != nil {
				utils.SendErrorResponse(hrw, &utils.ErrResp{Status: 500, Message: err.Error()})
			}
		}
		hashes, _ := storage.GetLocal(payload.Ref, payload.Arch, payload.Os)
		utils.ResponseJSON(hrw, http.StatusOK,
			structs.HashContainer{
				Ref:    payload.Ref,
				Arch:   payload.Arch,
				Os:     payload.Os,
				Hashes: hashes,
			},
		)
	}
}
