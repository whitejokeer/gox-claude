package posts

import (
    "github.com/gin-gonic/gin"
)

// RegisterRoutes registers all routes for the posts service
func RegisterRoutes(router *gin.RouterGroup, handlers *Handlers) {
    // Register your routes
    router.GET("/posts", handlers.HandleList)
    router.POST("/posts", handlers.HandleCreate)
    router.GET("/posts/:id", handlers.HandleGet)
    router.PUT("/posts/:id", handlers.HandleUpdate)
    router.DELETE("/posts/:id", handlers.HandleDelete)
}