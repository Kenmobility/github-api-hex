package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kenmobility/github-api-hex/config"
	"github.com/kenmobility/github-api-hex/db"
	"github.com/kenmobility/github-api-hex/internal/controllers"
	"github.com/kenmobility/github-api-hex/internal/handlers"
	"github.com/kenmobility/github-api-hex/internal/repositories"
	"github.com/kenmobility/github-api-hex/internal/routes"
	"github.com/kenmobility/github-api-hex/services/github"
)

func main() {
	// load env variables
	config, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// establish database connection
	database, err := db.NewDatabase(*config)
	if err != nil {
		log.Fatalf("failed to establish database connection: %v", err)
	}

	// seed and set 'chromium/chromium' repo as default repository to track
	err = database.SeedRepository(config)
	if err != nil {
		log.Fatalf("failed to seed default repository to database: %v", err)
	}

	// Initialize repositories
	commitRepo := repositories.NewGormCommitRepository(database.Db)
	repoRepo := repositories.NewGormRepositoryRepository(database.Db)

	// Initialize controllers
	commitController := controllers.NewCommitController(commitRepo, repoRepo)
	repoController := controllers.NewRepositoryController(repoRepo, config)

	// Initialize handlers
	commitHandler := handlers.NewCommitHandler(commitController)
	repositoryHandler := handlers.NewRepositoryHandler(repoController)

	// Initialize Github Tracker service
	trackerService := github.NewGitHubAPIClient(config.GitHubApiBaseURL, config.GitHubToken, config.FetchInterval, commitRepo, repoRepo)

	ginEngine := gin.Default()

	// register routes
	routes.CommitRoutes(ginEngine, commitHandler)
	routes.RepositoryRoutes(ginEngine, repositoryHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Address, config.Port),
		Handler: ginEngine,
	}

	// start web server
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v\n", err)
		}
		log.Printf("GitHub Service is listening on address %s", server.Addr)
	}()

	// create a context with cancellation to gracefully shut down GitHub tracker service if server shuts down
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown: err=%v", err)
	}

	//start GitHub tracking service asynchronously
	go log.Println(trackerService.StartTracking(ctx, config.FetchInterval))
}
