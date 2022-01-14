package main

import (
	"github.com/fanfansun/goServer/api/model"
	"github.com/fanfansun/goServer/api/session"
	"net/http"
)

var HEADER_FIELD_SESSION = "FF-Session-Id"
var HEADER_FIELD_UNAME = "FF-User-Name"

// 校验Session,没有Session 加入 session
func validateUserSession(r *http.Request) bool {
	// GET 返回与给定键值关联的第一个值,没有值，get返回 “”
	sid := r.Header.Get(HEADER_FIELD_SESSION)
	if len(sid) == 0 {
		return false
	}
	uname, ok := session.IsSessionExpired(sid)
	if ok {
		return false
	}

	r.Header.Add(HEADER_FIELD_UNAME, uname)
	return true
}

func ValidateUser(w http.ResponseWriter, r *http.Request) bool {
	uname := r.Header.Get(HEADER_FIELD_UNAME)
	if len(uname) == 0 {
		sendErrorResponse(w, model.ErrorNotAuthUser)
		return false
	}

	return true
}
