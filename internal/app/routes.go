package app

import (
	"goida/internal/handlers"
	"goida/internal/middleware"
)

func (a *App) setupRoutes(
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	articleHandler *handlers.ArticleHandler,
	roleHandler *handlers.RoleHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	a.router.Use(middleware.CORSMiddleware)
	a.router.Use(middleware.LoggingMiddleware)

	a.setupPublicRoutes(userHandler, authHandler, articleHandler, roleHandler)
	a.setupProtectedRoutes(authHandler, articleHandler, userHandler, authMiddleware)
	a.setupAdminRoutes(userHandler, roleHandler, authMiddleware)
}

func (a *App) setupPublicRoutes(
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	articleHandler *handlers.ArticleHandler,
	roleHandler *handlers.RoleHandler,
) {
	a.router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")
	a.router.HandleFunc("/api/users", userHandler.CreateUser).Methods("POST")
	a.router.HandleFunc("/api/articles", articleHandler.ListArticles).Methods("GET")
	a.router.HandleFunc("/api/articles/{id}", articleHandler.GetArticle).Methods("GET")
	a.router.HandleFunc("/api/users/{authorId}/articles", articleHandler.GetUserArticles).Methods("GET")
}

func (a *App) setupProtectedRoutes(
	authHandler *handlers.AuthHandler,
	articleHandler *handlers.ArticleHandler,
	userHandler *handlers.UserHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	authRouter := a.router.PathPrefix("/api").Subrouter()
	authRouter.Use(authMiddleware.RequireAuth)

	authRouter.HandleFunc("/auth/profile", authHandler.GetProfile).Methods("GET")
	authRouter.HandleFunc("/articles", articleHandler.CreateArticle).Methods("POST")
	authRouter.HandleFunc("/articles/{id}", articleHandler.UpdateArticle).Methods("PUT")
	authRouter.HandleFunc("/articles/{id}", articleHandler.DeleteArticle).Methods("DELETE")
	authRouter.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
}

func (a *App) setupAdminRoutes(
	userHandler *handlers.UserHandler,
	roleHandler *handlers.RoleHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	adminRouter := a.router.PathPrefix("/api/admin").Subrouter()
	adminRouter.Use(authMiddleware.RequireAdmin)

	adminRouter.HandleFunc("/users", userHandler.ListUsers).Methods("GET")
	adminRouter.HandleFunc("/roles", roleHandler.ListRoles).Methods("GET")
	adminRouter.HandleFunc("/roles/{id}", roleHandler.GetRole).Methods("GET")
}
