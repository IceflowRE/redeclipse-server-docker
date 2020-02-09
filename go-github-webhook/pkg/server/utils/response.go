package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
)

/*
value must be a VALUE NOT a pointer
 */
func ResponseJSON(hrw http.ResponseWriter, status int, value interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	if err := enc.Encode(value); PrintError(err) {
		SendErrorResponse(hrw, GenericErrorResp)
		return
	}

	content := buf.Bytes()
	SendJSON(hrw, status, &content)
}

/*
Sends the json data in content.
 */
func SendJSON(hrw http.ResponseWriter, status int, content *[]byte) {
	hrw.Header().Set("Content-Type", "application/json; charset=utf-8")
	hrw.WriteHeader(status)
	hrw.Write(*content)
}

var (
	GenericErrorResp = &ErrResp{http.StatusInternalServerError, "Internal Server Error"}
)

type ErrResp struct {
	Status  int
	Message string
}

func SendErrorResponse(hrw http.ResponseWriter, resp *ErrResp) {
	content := []byte("{\"error\":{\"code\":" + strconv.Itoa(resp.Status) + ",\"message\":\"" + resp.Message + "\"}}")
	SendJSON(hrw, resp.Status, &content)
}
