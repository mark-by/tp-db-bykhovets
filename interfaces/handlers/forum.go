package handlers

import (
	"github.com/mailru/easyjson"
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
	"github.com/mark-by/tp-db-bykhovets/domain/entityErrors"
	"github.com/valyala/fasthttp"
	"strconv"
)

func forumCreate(req *fasthttp.RequestCtx) {
	contentTypeJson(req)
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
			body, err := easyjson.Marshal(forum)
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
	contentTypeJson(req)
	forumSlug := req.UserValue("slug").(string)
	forum, err := app.Forum.Get(forumSlug)
	if err != nil {
		switch err {
		case entityErrors.ForumNotFound:
			req.SetStatusCode(fasthttp.StatusNotFound)
		default:
			req.SetStatusCode(fasthttp.StatusInternalServerError)
		}
		req.SetBodyString(message(err.Error()))
		return
	}
	req.SetStatusCode(fasthttp.StatusOK)
	_ = setBody(req, forum)
}

func forumThreadCreate(req *fasthttp.RequestCtx) {
	contentTypeJson(req)
	forumSlug := req.UserValue("slug").(string)
	thread := &entity.Thread{}
	err := safeUnmarshal(req, thread)
	if err != nil {
		return
	}

	err = app.Forum.CreateThread(forumSlug, thread)
	if err != nil {
		switch err {
		case entityErrors.ForumNotFound:
			fallthrough
		case entityErrors.UserNotFound:
			req.SetStatusCode(fasthttp.StatusNotFound)
		case entityErrors.ThreadAlreadyExist:
			req.SetStatusCode(fasthttp.StatusConflict)
			_ = setBody(req, thread)
		default:
			req.SetStatusCode(fasthttp.StatusInternalServerError)
		}
		if err != entityErrors.ThreadAlreadyExist {
			req.SetBodyString(message(err.Error()))
		}
		return
	}

	req.SetStatusCode(fasthttp.StatusCreated)
	_ = setBody(req, thread)
}

func forumUsers(req *fasthttp.RequestCtx) {
	contentTypeJson(req)
	forumSlug := req.UserValue("slug").(string)
	query := req.QueryArgs()
	limit, _ := strconv.Atoi(string(query.Peek("limit")))
	since := string(query.Peek("since"))
	desc := query.GetBool("desc")
	users, err := app.Forum.GetUsers(forumSlug, limit, since, desc)
	if err != nil {
		if err == entityErrors.ForumNotFound {
			req.SetStatusCode(fasthttp.StatusNotFound)
		} else {
			req.SetStatusCode(fasthttp.StatusInternalServerError)
		}
		req.SetBodyString(message(err.Error()))
		return
	}

	req.SetStatusCode(fasthttp.StatusOK)
	_ = setBody(req, users)
}

func forumThreads(req *fasthttp.RequestCtx) {
	contentTypeJson(req)
	forumSlug := req.UserValue("slug").(string)
	query := req.QueryArgs()
	limit, _ := strconv.Atoi(string(query.Peek("limit")))
	since := string(query.Peek("since"))
	desc := query.GetBool("desc")
	threads, err := app.Forum.GetThreads(forumSlug, limit, since, desc)
	if err != nil {
		if err == entityErrors.ForumNotFound {
			req.SetStatusCode(fasthttp.StatusNotFound)
			req.SetBodyString(message(err.Error()))
			return
		}
		serveInternalErr(req, err)
		return
	}

	req.SetStatusCode(fasthttp.StatusOK)
	_ = setBody(req, threads)
}
