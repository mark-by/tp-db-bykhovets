package handlers

import (
	"github.com/mark-by/tp-db-bykhovets/domain/entity"
	"github.com/mark-by/tp-db-bykhovets/domain/entityErrors"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

func postDetail(req *fasthttp.RequestCtx) {
	contentTypeJson(req)
	idStr := req.UserValue("id").(string)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		req.SetStatusCode(fasthttp.StatusBadRequest)
		req.SetBodyString(message("id must be int"))
		return
	}
	relatedStr := string(req.QueryArgs().Peek("related"))
	related := strings.Split(relatedStr, ",")
	post, err := app.Post.Get(int64(id), related)
	if err != nil {
		req.SetBodyString(message(err.Error()))
		if err == entityErrors.PostNotFound {
			req.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		req.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	req.SetStatusCode(fasthttp.StatusOK)
	_ = setBody(req, post)
}

func postUpdate(req *fasthttp.RequestCtx) {
	contentTypeJson(req)
	idStr := req.UserValue("id").(string)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		req.SetStatusCode(fasthttp.StatusBadRequest)
		req.SetBodyString(message("id must be int"))
		return
	}
	post := &entity.Post{}
	err = safeUnmarshal(req, post)
	if err != nil {
		return
	}
	post.ID = int64(id)
	err = app.Post.Update(post)
	if err != nil {
		req.SetBodyString(message(err.Error()))
		if err == entityErrors.PostNotFound {
			req.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		req.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	req.SetStatusCode(fasthttp.StatusOK)
	_ = setBody(req, post)
}
