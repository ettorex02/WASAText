package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	rt.router.POST("/session", rt.SessionHandler)
	rt.router.GET("/users/:userId", rt.GetUserHandler)
	rt.router.PATCH("/users/:userId", rt.SetMyUserName)
	rt.router.PATCH("/users/:userId/photo", rt.SetMyPhoto)
	rt.router.GET("/search/users", rt.SearchUsersHandler) // <--- cambiato qui
	rt.router.POST("/conversations", rt.CreateConversationHandler)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
