package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kenmobility/github-api-hex/common/message"
	"github.com/kenmobility/github-api-hex/common/response"
	"github.com/kenmobility/github-api-hex/internal/controllers"
)

type CommitHandlers struct {
	commitController controllers.CommitController
}

func NewCommitHandler(commitController controllers.CommitController) *CommitHandlers {
	return &CommitHandlers{
		commitController: commitController,
	}
}

func (ch CommitHandlers) GetCommitsByRepositoryId(ctx *gin.Context) {
	query := getPagingInfo(ctx)

	repositoryId := ctx.Param("repoId")

	if repositoryId == "" {
		response.Failure(ctx, http.StatusBadRequest, "repoId is required", nil)
		return
	}

	repoName, commits, err := ch.commitController.GetAllCommitsByRepository(ctx, repositoryId, query)
	if err == message.ErrInvalidRepositoryId {
		response.Failure(ctx, http.StatusBadRequest, "invalid repo Id", err)
		return
	}
	if err != nil {
		response.Failure(ctx, http.StatusInternalServerError, "invalid repo Id", err)
		return
	}

	msg := fmt.Sprintf("%s commits fetched successfully", repoName)

	response.Success(ctx, http.StatusOK, msg, commits)
}

func (ch CommitHandlers) GetTopCommitAuthors(ctx *gin.Context) {
	repositoryId := ctx.Param("repoId")

	if repositoryId == "" {
		response.Failure(ctx, http.StatusBadRequest, "repoId is required", nil)
		return
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))

	if err != nil {
		response.Failure(ctx, http.StatusBadRequest, "invalid limit", err)
		return
	}

	authors, err := ch.commitController.GetTopRepositoryCommitAuthors(ctx, repositoryId, limit)
	if err != nil {
		response.Failure(ctx, http.StatusInternalServerError, "error fetching top authors", err)
		return
	}

	msg := fmt.Sprintf("%v top commit authors fetched successfully", limit)

	response.Success(ctx, http.StatusOK, msg, authors)
}
