package webtest

import (
	"net/http"

	"github.com/lonepeon/golib/web"
)

func MockedResponse(name string) web.Response {
	return web.Response{
		HTTPCode:   http.StatusNotImplemented,
		Layout:     "mockLayout",
		LogMessage: name,
		Data:       nil,
		Template:   "mockTpl",
	}
}
