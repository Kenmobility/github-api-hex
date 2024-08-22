package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kenmobility/github-api-hex/internal/handlers"
)

func CommitRoutes(r *gin.Engine, ch *handlers.CommitHandlers) {
	r.GET("/commits/:repoId", ch.GetCommitsByRepositoryId)
	r.GET("/top-authors", ch.GetTopCommitAuthors)
}
