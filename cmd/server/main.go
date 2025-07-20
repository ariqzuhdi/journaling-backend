package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cheeszy/journaling/controllers"
	"github.com/cheeszy/journaling/initializers"
	"github.com/cheeszy/journaling/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("FE_DOMAIN")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ===== Public Routes =====
	public := router.Group("/api")
	{
		public.POST("/register", controllers.Register)
		public.POST("/login", controllers.Login)
		public.POST("/resend-verification", controllers.ResendVerificationEmail)
		public.POST("/reset-password", controllers.ResetPasswordWithRecoveryKey)

		public.GET("/verify", controllers.VerifyEmail)
		public.GET("/monkeytype", controllers.MonkeyAPI)
		public.GET("/posts", controllers.PostsIndex)
		public.GET("/health", controllers.HealthHandler)

		// Optional/Commented routes
		// public.GET("/posts/:id", controllers.PostsShowById)
		// public.GET("/users", controllers.Users)
	}

	// ===== Protected Routes =====
	protected := router.Group("/api")
	protected.Use(middleware.RequireAuth)
	protected.Use(middleware.RequireRLS)
	{
		protected.GET("/posts/user/:username", controllers.PostsShowAllPosts)
		protected.POST("/posts", controllers.PostsCreate)
		protected.PUT("/posts/:id", controllers.PostsUpdate)
		protected.DELETE("/posts/:id", controllers.PostsDelete)

		protected.POST("/logout", controllers.Logout)
		protected.GET("/user", controllers.GetCurrentUser)
		protected.GET("/", controllers.HomeHandler)

		protected.PUT("/account/change-username", controllers.ChangeUsername)
		protected.PUT("/account/change-email", controllers.ChangeEmail)
	}

	router.NoRoute(controllers.NotFoundHandler)
	fmt.Println(os.Getenv("DOMAIN"))
	router.Run(":3000")
}
