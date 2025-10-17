package app

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"goida/internal/config"
	"goida/internal/database"
	"goida/internal/handlers"
	"goida/internal/middleware"
	"goida/internal/repository"
	"goida/internal/services"
)

type App struct {
	config *config.Config
	db     *database.Database
	router *mux.Router
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	db, err := database.New(cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)
	if err != nil {
		return nil, err
	}

	app := &App{
		config: cfg,
		db:     db,
		router: mux.NewRouter(),
	}

	if err := app.setup(); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) setup() error {
	userRepo := repository.NewUserRepository(a.db.DB)
	articleRepo := repository.NewArticleRepository(a.db.DB)
	roleRepo := repository.NewRoleRepository(a.db.DB)
	authCredentialsRepo := repository.NewAuthCredentialsRepository(a.db.DB)
	commentRepo := repository.NewCommentRepository(a.db.DB)

	userService := services.NewUserService(userRepo, roleRepo, authCredentialsRepo)
	authService := services.NewAuthService(userRepo, authCredentialsRepo, a.config.JWTSecret)
	articleService := services.NewArticleService(articleRepo, userRepo, commentRepo)
	commentService := services.NewCommentService(commentRepo, articleRepo)

	authMiddleware := middleware.NewAuthMiddleware(authService)
	validator := middleware.NewValidator()

	userHandler := handlers.NewUserHandler(userService, validator)
	authHandler := handlers.NewAuthHandler(authService, validator)
	articleHandler := handlers.NewArticleHandler(articleService, validator)
	commentHandler := handlers.NewCommentHandler(commentService, validator)
	roleHandler := handlers.NewRoleHandler(roleRepo)
	authCredentialsHandler := handlers.NewAuthCredentialsHandler(authCredentialsRepo, validator)

	a.setupRoutes(userHandler, authHandler, articleHandler, roleHandler, authCredentialsHandler, commentHandler, authMiddleware)

	return nil
}

func (a *App) Run() error {
	port := a.config.Server.Port
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	logrus.Infof("Server starting on port %s", port)
	return http.ListenAndServe(":"+port, a.router)
}

func (a *App) Close() error {
	return a.db.Close()
}
