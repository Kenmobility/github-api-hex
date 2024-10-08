package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kenmobility/github-api-hex/internal/handlers"
)

func RepositoryRoutes(r *gin.Engine, rh *handlers.RepositoryHandlers) {
	r.POST("/repository", rh.AddRepository)
	r.GET("/repositories", rh.FetchAllRepositories)
	r.POST("/repository/track", rh.TrackRepository)
}
