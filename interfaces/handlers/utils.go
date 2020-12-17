package handlers

import (
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func serveJsonUnmarshalErr(req *fasthttp.RequestCtx, err error) {
	logrus.WithField("type", "unmarshal json error").Error(err)
	req.Response.SetStatusCode(fasthttp.StatusBadRequest)
	req.Response.SetBodyString(message(err.Error()))
}

func serveJsonMarshalErr(req *fasthttp.RequestCtx, err error) {
	logrus.WithField("type", "marshal json error").Error(err)
	serveInternalErr(req, err)
}

func serveInternalErr(req *fasthttp.RequestCtx, err error) {
	req.SetStatusCode(fasthttp.StatusInternalServerError)
	req.SetBodyString(message(err.Error()))
}

func setBody(req *fasthttp.RequestCtx, any easyjson.Marshaler) error {
	body, err := easyjson.Marshal(any)
	if err != nil {
		serveJsonMarshalErr(req, err)
		return err
	}
	req.SetBody(body)
	return nil
}

func safeUnmarshal(req *fasthttp.RequestCtx, any easyjson.Unmarshaler) error {
	err := easyjson.Unmarshal(req.Request.Body(), any)
	if err != nil {
		serveJsonUnmarshalErr(req, err)
		return err
	}
	return nil
}

func setStatus(req *fasthttp.RequestCtx, status int) {
	req.SetStatusCode(status)
}

func message(message string) string {
	return fmt.Sprintf(`{"message":"%s"}`, message)
}