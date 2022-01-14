package main

import (
	"encoding/json"
	model2 "github.com/fanfansun/goServer/api/model"
	"io"
	"net/http"
)

func sendErrorResponse(w http.ResponseWriter, errResp model2.ErrResponse) {
	w.WriteHeader(errResp.HttpSc)

	resStr, _ := json.Marshal(&errResp.Error)
	io.WriteString(w, string(resStr))
}

func sendNormalResponse(w http.ResponseWriter, resp string, sc int) {
	w.WriteHeader(sc)
	io.WriteString(w, resp)
}
