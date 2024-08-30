package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kenmobility/github-api-hex/common/helpers"
	"github.com/kenmobility/github-api-hex/common/message"
	"github.com/kenmobility/github-api-hex/common/response"
	"github.com/kenmobility/github-api-hex/dtos"
	"github.com/kenmobility/github-api-hex/internal/controllers"
)

type RepositoryHandlers struct {
	repositoryController controllers.RepositoryController
}

func NewRepositoryHandler(repositoryController controllers.RepositoryController) *RepositoryHandlers {
	return &RepositoryHandlers{
		repositoryController: repositoryController,
	}
}

func (rh RepositoryHandlers) AddRepository(ctx *gin.Context) {
	var input dtos.AddRepositoryRequestDto

	err := ctx.BindJSON(&input)
	if err != nil {
		response.Failure(ctx, http.StatusBadRequest, "invalid input", err)
		return
	}

	inputErrors := helpers.ValidateInput(input)
	if inputErrors != nil {
		response.Failure(ctx, http.StatusBadRequest, message.ErrInvalidInput.Error(), inputErrors)
		return
	}

	repo, err := rh.repositoryController.AddRepository(ctx, input)
	if err != nil {
		response.Failure(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}

	response.Success(ctx, http.StatusCreated, "New Repository added successfully", repo)
}

func (rh RepositoryHandlers) TrackRepository(ctx *gin.Context) {
	var input dtos.TrackRepositoryRequestDto

	err := ctx.BindJSON(&input)
	if err != nil {
		response.Failure(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}

	inputErrors := helpers.ValidateInput(input)
	if inputErrors != nil {
		response.Failure(ctx, http.StatusBadRequest, message.ErrInvalidInput.Error(), inputErrors)
		return
	}

	repo, err := rh.repositoryController.TrackRepository(ctx, input)
	if err != nil {
		response.Failure(ctx, http.StatusBadRequest, err.Error(), err)
		return
	}

	response.Success(ctx, http.StatusOK, "new repository tracking successful", repo)
}

func (rh RepositoryHandlers) FetchAllRepositories(ctx *gin.Context) {
	repos, err := rh.repositoryController.GetAllRepositories(ctx)
	if err != nil {
		response.Failure(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}

	response.Success(ctx, http.StatusOK, "successfully fetched all repos", repos)
}
