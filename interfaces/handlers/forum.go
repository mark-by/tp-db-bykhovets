package handlers

import (
	"github.com/mailru/easyjson"
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
	"github.com/mark-by/tp-db-bykhovets/domain/entityErrors"
	"github.com/valyala/fasthttp"
)

func forumCreate(req *fasthttp.RequestCtx) {
	forum := &entity.Forum{}
	err := easyjson.Unmarshal(req.Request.Body(), forum)
	if err != nil {
		serveJsonUnmarshalErr(req, err)
		return
	}
	err = app.Forum.Create(forum)
	if err != nil {
		switch err {
		case entityErrors.UserNotFound:
			req.Response.SetStatusCode(fasthttp.StatusNotFound)
			req.Response.SetBodyString(message(err.Error()))
			return
		case entityErrors.ForumAlreadyExist:
			body , err := easyjson.Marshal(forum)
			if err != nil {
				serveJsonMarshalErr(req, err)
				return
			}
			req.Response.SetStatusCode(fasthttp.StatusConflict)
			req.Response.SetBody(body)
			return
		default:
			req.Response.SetStatusCode(fasthttp.StatusInternalServerError)
			req.Response.SetBodyString(message(err.Error()))
		}
	}

	body, err := easyjson.Marshal(forum)
	if err != nil {
		serveJsonMarshalErr(req, err)
		return
	}
	req.Response.SetStatusCode(fasthttp.StatusOK)
	req.Response.SetBody(body)
}

func forumDetails(req *fasthttp.RequestCtx) {

}

func forumThreadCreate(req *fasthttp.RequestCtx) {

}

func forumUsers(req *fasthttp.RequestCtx) {

}

func forumThreads(req *fasthttp.RequestCtx) {

}
