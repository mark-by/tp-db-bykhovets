package handlers

import (
	"github.com/fasthttp/router"
	"github.com/mark-by/tp-db-bykhovets/domain/repository"
	"github.com/valyala/fasthttp"
	"log"
)

func ServeAPI(repositories *repository.Repositories) {
	r := router.New()

	r.POST("/forum/create", forumCreate)
	r.GET("/forum/{slug}/details", forumDetails)
	r.POST("/forum/{slug}/create", forumThreadCreate)
	r.GET("/forum/{slug}/users", forumUsers)
	r.GET("/forum/{slug}/threads", forumThreads)

	r.GET("/post/{id}/details", postDetail)
	r.POST("/post/{id}/details", postUpdate)

	r.POST("/service/clear", serviceClear)
	r.GET("/service/status", serviceInfo)

	r.POST("/thread/{slug_or_id}/create", threadCreatePosts)
	r.GET("/thread/{slug_or_id}/details", threadDetails)
	r.POST("/thread/{slug_or_id}/details", threadUpdate)
	r.GET("/thread/{slug_or_id}/posts", postsFromThread)
	r.POST("/thread/{slug_or_id}/vote", voteThread)

	r.POST("/user/{nickname}/create", createUser)
	r.GET("/user/{nickname}/profile", detailsUser)
	r.POST("/user/{nickname}/profile", updateUser)

	if err := fasthttp.ListenAndServe("0.0.0.0:8000", r.Handler); err != nil {
		log.Fatalf("Fail to start server: %s", err)
	}
}
