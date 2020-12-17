package handlers

import (
	"github.com/valyala/fasthttp"
)

func serviceClear(req *fasthttp.RequestCtx) {
	err := app.Service.Clear()
	if err != nil {
		serveInternalErr(req, err)
	}
	req.SetStatusCode(fasthttp.StatusOK)
}

func serviceInfo(req *fasthttp.RequestCtx) {
	contentTypeJson(req)
	status, err := app.Service.Status()
	if err != nil {
		serveInternalErr(req, err)
	}
	req.SetStatusCode(fasthttp.StatusOK)
	_ = setBody(req, status)
}
