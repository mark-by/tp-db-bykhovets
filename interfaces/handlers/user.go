package handlers

import (
	"github.com/mailru/easyjson"
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
	"github.com/mark-by/tp-db-bykhovets/domain/entityErrors"
	"github.com/valyala/fasthttp"
)

func createUser(req *fasthttp.RequestCtx) {
	nickname := req.UserValue("nickname").(string)
	user := &entity.User{}
	err := easyjson.Unmarshal(req.Request.Body(), user)
	if err != nil {
		serveJsonUnmarshalErr(req, err)
		return
	}

	user.NickName = nickname
	conflictUsers, err := app.User.Create(user)
	if err != nil {
		serveInternalErr(req, err)
		return
	}
	if conflictUsers != nil {
		req.SetStatusCode(fasthttp.StatusConflict)
		_ = setBody(req, conflictUsers)
		return
	}

	req.SetStatusCode(fasthttp.StatusCreated)
	_ = setBody(req, user)
}

func updateUser(req *fasthttp.RequestCtx) {
	nickname := req.UserValue("nickname").(string)
	user := &entity.User{}
	user.NickName = nickname
	err := safeUnmarshal(req, user)
	if err != nil {
		return
	}

	err = app.User.Update(user)
	if err != nil {
		switch err {
		case entityErrors.UserNotFound:
			req.SetStatusCode(fasthttp.StatusNotFound)
		case entityErrors.UserConflict:
			req.SetStatusCode(fasthttp.StatusConflict)
		default:
			req.SetStatusCode(fasthttp.StatusInternalServerError)
		}
		req.SetBodyString(message(err.Error()))
		return
	}

	req.SetStatusCode(fasthttp.StatusOK)
	_ = setBody(req, user)
}

func detailsUser(req *fasthttp.RequestCtx) {
	user, err := app.User.Get(req.UserValue("nickname").(string))
	if err != nil {
		if err == entityErrors.UserNotFound {
			req.SetStatusCode(fasthttp.StatusNotFound)
			req.SetBodyString(message(err.Error()))
			return
		}
		serveInternalErr(req, err)
		return
	}

	req.SetStatusCode(fasthttp.StatusOK)
	_ = setBody(req, user)
}
