package handlers

import (
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
	"github.com/mark-by/tp-db-bykhovets/domain/entityErrors"
	"github.com/valyala/fasthttp"
	"strconv"
)

func threadCreatePosts(req *fasthttp.RequestCtx) {
	contentTypeJson(req)
	slugOrId := req.UserValue("slug_or_id").(string)
	var posts entity.PostList
	err := safeUnmarshal(req, &posts)
	if err != nil {
		return
	}
	if err = app.Thread.CreatePosts(slugOrId, posts); err != nil {
		switch err {
		case entityErrors.ThreadNotFound:
			req.SetStatusCode(fasthttp.StatusNotFound)
		case entityErrors.ParentNotExist:
			req.SetStatusCode(fasthttp.StatusConflict)
		case entityErrors.UserNotFound:
			req.SetStatusCode(fasthttp.StatusBadRequest)
		default:
			req.SetStatusCode(fasthttp.StatusInternalServerError)
		}
		req.SetBodyString(message(err.Error()))
		return
	}

	req.SetStatusCode(fasthttp.StatusOK)
	_ = setBody(req, posts)
}

func threadDetails(req *fasthttp.RequestCtx) {
	contentTypeJson(req)
	slugOrId := req.UserValue("slug_or_id").(string)
	thread, err := app.Thread.Get(slugOrId)
	if err != nil {
		req.SetBodyString(message(err.Error()))
		if err == entityErrors.ThreadNotFound {
			req.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		req.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	req.SetStatusCode(fasthttp.StatusOK)
	_ = setBody(req, thread)
}

func threadUpdate(req *fasthttp.RequestCtx) {
	contentTypeJson(req)
	slugOrId := req.UserValue("slug_or_id").(string)
	thread := &entity.Thread{}
	err := safeUnmarshal(req, thread)
	if err != nil {
		return
	}
	err = app.Thread.Update(slugOrId, thread)
	if err != nil {
		req.SetBodyString(message(err.Error()))
		if err == entityErrors.ThreadNotFound {
			req.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		req.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	req.SetStatusCode(fasthttp.StatusOK)
	_ = setBody(req, thread)
}

func postsFromThread(req *fasthttp.RequestCtx) {
	contentTypeJson(req)
	slugOrId := req.UserValue("slug_or_id").(string)
	query := req.QueryArgs()
	limit, _ := strconv.Atoi(string(query.Peek("limit")))
	since, _ := strconv.Atoi(string(query.Peek("since")))
	desc := query.GetBool("desc")
	sort := string(query.Peek("sort"))

	posts, err := app.Thread.GetPosts(slugOrId, limit, since, desc, sort)
	if err != nil {
		req.SetBodyString(message(err.Error()))
		if err == entityErrors.ThreadNotFound {
			req.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		req.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	req.SetStatusCode(fasthttp.StatusOK)
	_ = setBody(req, posts)
}

func voteThread(req *fasthttp.RequestCtx) {
	contentTypeJson(req)
	slugOrId := req.UserValue("slug_or_id").(string)
	vote := &entity.Vote{}
	err := safeUnmarshal(req, vote)
	if err != nil {
		return
	}

	thread, err := app.Thread.Vote(slugOrId, vote)
	if err != nil {
		req.SetBodyString(message(err.Error()))
		if err == entityErrors.ThreadNotFound {
			req.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		if err == entityErrors.UserNotFound {
			req.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
		req.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	req.SetStatusCode(fasthttp.StatusOK)
	_ = setBody(req, thread)
}
