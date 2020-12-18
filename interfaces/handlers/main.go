package handlers

import (
	"github.com/fasthttp/router"
	"github.com/mark-by/tp-db-bykhovets/application"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var app *application.App

func SetApp(implApp *application.App) {
	app = implApp
}

func ServeAPI() {
	r := router.New()

	r.POST("/api/forum/create", forumCreate)
	r.GET("/api/forum/{slug}/details", forumDetails)
	r.POST("/api/forum/{slug}/create", forumThreadCreate)
	r.GET("/api/forum/{slug}/users", forumUsers)
	r.GET("/api/forum/{slug}/threads", forumThreads)

	r.GET("/api/post/{id}/details", postDetail)
	r.POST("/api/post/{id}/details", postUpdate)

	r.POST("/api/service/clear", serviceClear)
	r.GET("/api/service/status", serviceInfo)

	r.POST("/api/thread/{slug_or_id}/create", threadCreatePosts)
	r.GET("/api/thread/{slug_or_id}/details", threadDetails)
	r.POST("/api/thread/{slug_or_id}/details", threadUpdate)
	r.GET("/api/thread/{slug_or_id}/posts", postsFromThread)
	r.POST("/api/thread/{slug_or_id}/vote", voteThread)

	r.POST("/api/user/{nickname}/create", createUser)
	r.GET("/api/user/{nickname}/profile", detailsUser)
	r.POST("/api/user/{nickname}/profile", updateUser)

	if err := fasthttp.ListenAndServe("0.0.0.0:5000", r.Handler); err != nil {
		logrus.Fatalf("Fail to start server: %s", err)
	}
}
