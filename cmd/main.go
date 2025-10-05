package main

import (
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/jexlor/votingapp/db/store"
	"github.com/jexlor/votingapp/internal/api"
	"github.com/jexlor/votingapp/internal/api/middleware"
	"github.com/jexlor/votingapp/internal/auth"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	middlewareCSRF "github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("failed to load .env file", "error", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port
	}

	dbUrl := os.Getenv("DB_URL")

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatal("failed to ping db:", err)
	}

	a := api.Config{
		DB: store.New(conn),
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	auth.SetJWTKey(jwtSecret)

	e := echo.New()

	e.Use(middlewareCSRF.CSRFWithConfig(middlewareCSRF.CSRFConfig{
		TokenLookup:    "form:csrf_token",
		CookieName:     "csrf_token",
		CookiePath:     "/",
		CookieMaxAge:   86400,
		CookieSecure:   false,
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteLaxMode,
	}))
	v1 := e.Group("/v1")

	// TODO: configure CORS
	// v1.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"https://labstack.com", "https://labstack.net"},
	// 	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	// }))

	v1.POST("/register", a.HandleRegister)
	v1.POST("/login", a.HandleLogin)
	v1.GET("/login", api.HandleLoginPage)
	v1.GET("/register", api.HandleRegisterPage)
	authGroup := v1.Group("", middleware.JWTAuthMiddleware())
	authGroup.GET("/comments", a.HandlerGetAllComments)
	authGroup.POST("/comments", a.HandlerCreateComment)
	authGroup.POST("/comments/:id/delete", a.HandlerDeleteComment)
	authGroup.POST("/comments/react", a.HandlerReactComment)

	if err := e.Start(":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}
