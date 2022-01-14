package main

import (
	"encoding/json"
	"github.com/fanfansun/goServer/api/db"
	"github.com/fanfansun/goServer/api/model"
	"github.com/fanfansun/goServer/api/session"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"net/http"
)

func CreateUser(w http.ResponseWriter, R *http.Request, _ httprouter.Params) {
	//_, err := io.WriteString(w, "Crete User Handler")
	//if err != nil {
	//	fmt.Printf("注册模块异常%s\n", err)
	//	return
	//}
	res, _ := ioutil.ReadAll(R.Body)
	ubody := &model.UserModelInterface{}

	if err := json.Unmarshal(res, ubody); err != nil {
		sendErrorResponse(w, model.ErrorRequestBodyParseFailed)
		return
	}
	if err := db.AddUserInformation(ubody.Username, ubody.Password); err != nil {
		sendErrorResponse(w, model.ErrorDBError)
	}

	id := session.GenerateNewSessionId(ubody.Username)
	su := &model.SignedUp{Success: true, SessionId: id}

	if resp, err := json.Marshal(su); err != nil {
		sendErrorResponse(w, model.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w, string(resp), 201)
	}

}

func Login(w http.ResponseWriter, _ *http.Request, param httprouter.Params) {
	//rqUsername := param.ByName("user_name")
	//_, err := io.WriteString(w, rqUsername)
	//if err != nil {
	//	fmt.Printf("登录模块异常%s\n", err)
	//	return
	//}
	uname := param.ByName("user_name")
	io.WriteString(w, uname)
}
