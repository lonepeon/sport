package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Context interface {
	StdCtx() context.Context
	AddData(string, interface{})
	AddFlash(f FlashMessage)
	Response(httpCode int, template string, data map[string]interface{}) Response
	Redirect(w http.ResponseWriter, httpCode int, target string) Response
	NotFoundResponse(format string, vars ...interface{}) Response
	InternalServerErrorResponse(format string, vars ...interface{}) Response
	Vars(r *http.Request) map[string]string
}

type ContextImpl struct {
	context.Context

	data              map[string]interface{}
	session           *sessions.Session
	tmplConfiguration TmplConfiguration
}

func (c *ContextImpl) StdCtx() context.Context {
	return c
}

func (c *ContextImpl) AddFlash(f FlashMessage) {
	c.session.AddFlash(&f)
}

func (c *ContextImpl) AddData(key string, data interface{}) {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}

	c.data[key] = data
}

func (c *ContextImpl) Response(httpCode int, template string, data map[string]interface{}) Response {
	d := c.data
	if d == nil {
		d = make(map[string]interface{})
	}

	for k, v := range data {
		d[k] = v
	}

	return Response{
		HTTPCode:   httpCode,
		LogMessage: "response sent",
		Layout:     c.tmplConfiguration.Layout,
		Template:   template,
		Data:       d,
	}
}

func (c *ContextImpl) Redirect(w http.ResponseWriter, httpCode int, target string) Response {
	w.Header().Add("Location", target)

	return Response{
		HTTPCode:   httpCode,
		Layout:     "",
		LogMessage: "redirected to " + target,
		Template:   c.tmplConfiguration.RedirectionTemplate,
		Data:       map[string]interface{}{"Target": target},
	}
}

func (c *ContextImpl) NotFoundResponse(format string, vars ...interface{}) Response {
	return Response{
		HTTPCode:   http.StatusNotFound,
		LogMessage: fmt.Sprintf(format, vars...),
		Layout:     c.tmplConfiguration.ErrorLayout,
		Template:   c.tmplConfiguration.NotFoundTemplate,
	}
}

func (c *ContextImpl) InternalServerErrorResponse(format string, vars ...interface{}) Response {
	return Response{
		HTTPCode:   http.StatusInternalServerError,
		LogMessage: fmt.Sprintf(format, vars...),
		Layout:     c.tmplConfiguration.ErrorLayout,
		Template:   c.tmplConfiguration.InternalServerErrorTemplate,
	}
}

func (c *ContextImpl) Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}
