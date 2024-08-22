package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kenmobility/github-api-hex/config"
	"github.com/kenmobility/github-api-hex/db"
	"github.com/kenmobility/github-api-hex/internal/controllers"
	"github.com/kenmobility/github-api-hex/internal/handlers"
	"github.com/kenmobility/github-api-hex/internal/integrations/api"
	"github.com/kenmobility/github-api-hex/internal/repositories"
	"github.com/kenmobility/github-api-hex/internal/routes"
	"github.com/kenmobility/github-api-hex/services"
)

func main() {
	// load env variables
	config := config.LoadConfig()

	// establish database connection
	db := db.NewDatabase(*config)

	// seed 'chromium/chromium' repo as default repository to repositories table
	db.SeedRepository(config)

	// Initialize repositories
	commitRepo := repositories.NewGormCommitRepository(db.Db)
	repoRepo := repositories.NewGormRepositoryRepository(db.Db)

	// Initialize controllers
	commitController := controllers.NewCommitController(commitRepo, repoRepo)
	repoController := controllers.NewRepositoryController(repoRepo, config)

	// Initialize handlers
	commitHandler := handlers.NewCommitHandler(commitController)
	repositoryHandler := handlers.NewRepositoryHandler(repoController)

	// Initialize API Client
	gitHubAPIClient := api.NewGitHubAPI(config.GitHubApiBaseURL, config.GitHubToken, commitRepo, repoRepo)

	// instantiate the GitHubAPI service
	githubService := services.NewGithubService(gitHubAPIClient, commitRepo,
		repoRepo, config)

	// start GitHub tracking service asynchronously
	go githubService.StartTracking()

	server := gin.New()

	// register routes
	routes.CommitRoutes(server, commitHandler)
	routes.RepositoryRoutes(server, repositoryHandler)

	//run server
	if err := server.Run(fmt.Sprintf("%s:%s", config.Address, config.Port)); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s", err)
	}
}
