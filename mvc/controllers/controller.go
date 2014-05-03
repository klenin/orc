package controllers

import (
	//"github.com/uas/session"
	"net/http"
)

type BaseController struct{}

type Controller struct {
	Request  *http.Request
	Response http.ResponseWriter
	//Session  *session.Manager
}
