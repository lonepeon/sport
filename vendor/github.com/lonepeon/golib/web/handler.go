package web

import (
	"net/http"
)

type Handler interface {
	Handle(Context, http.ResponseWriter, *http.Request) Response
}

type HandlerFunc func(Context, http.ResponseWriter, *http.Request) Response
