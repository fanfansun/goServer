package main

import (
	"encoding/json"
	"gitee.com/fan_zi_xin/go_server/api/model"
	"io"
	"net/http"
)

func sendErrorResponse(w http.ResponseWriter, errResp model.ErrResponse) {
	w.WriteHeader(errResp.HttpSc)

	resStr, _ := json.Marshal(&errResp.Error)
	io.WriteString(w, string(resStr))
}

func sendNormalResponse(w http.ResponseWriter, resp string, sc int) {
	w.WriteHeader(sc)
	io.WriteString(w, resp)
}
