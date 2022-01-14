package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type middleWareHandler struct {
	r *httprouter.Router
}

func NewMiddleWareHandler(r *httprouter.Router) http.Handler {
	middleHandre := middleWareHandler{}
	middleHandre.r = r
	return middleHandre
}

func (middleHandre middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//check session
	validateUserSession(r)

	middleHandre.r.ServeHTTP(w, r)
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/user", CreateUser)

	router.POST("/user/:user_name", Login)

	return router
}

func main() {
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r)
	http.ListenAndServe(":8089", mh)
}
