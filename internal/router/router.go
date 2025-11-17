package router

import (
	"database/sql"

	_ "github.com/giang19062001/gin-golang-standard/docs"
	"github.com/giang19062001/gin-golang-standard/internal/config"
	"github.com/giang19062001/gin-golang-standard/internal/controllers"
	middlewares "github.com/giang19062001/gin-golang-standard/internal/middleware"
	"github.com/giang19062001/gin-golang-standard/internal/repositories"
	"github.com/giang19062001/gin-golang-standard/internal/services"
	rabbitmq "github.com/giang19062001/gin-golang-standard/pkg/rabbit"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// * CHUẨN Clean Architecture trong Go
// 1. REPOSITORY, SERVERICE TRẢ VỀ INTERFACE
// 2. CONTROLLER TRẢ VỀ STRUCT

func SetupRouter(g *gin.Engine, cfg *config.Config, db *sql.DB, mq *rabbitmq.MqService) {

	// SWAGGER
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//email - rabbitmq
	emailService := services.NewEmailService(mq)

	// user
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo, cfg.JwtSecret)
	userController := controllers.NewUserController(userService)

	// event
	eventRepo := repositories.NewEventRepository(db)
	eventService := services.NewEventService(eventRepo, userService)
	eventController := controllers.NewEventController(eventService)

	// attendee
	attendeeRepo := repositories.NewAttendeRepository(db)
	attendeeService := services.NewAttendeeService(attendeeRepo, userService, eventService, emailService)
	attendeeController := controllers.NewAttendeeController(attendeeService)

	v1 := g.Group("/api/v1")
	{
		// Public routes
		public := v1.Group("")
		{
			// AUTH
			public.POST("/auth/register", userController.RegisterUser)
			public.POST("/auth/login", userController.LoginUser)

			// USER
			public.GET("/users", userController.GetAllUser)

			// EVENT
			public.GET("/events", eventController.GetAllEvents)
			public.POST("/events", eventController.CreateEvent)
			public.GET("/events/:id", eventController.GetEvent)
			public.PUT("/events/:id", eventController.UpdateEvent)
			public.DELETE("/events/:id", eventController.DeleteEvent)
			public.POST("/events/images", eventController.UploadImages)

			// ATTENDE
			public.POST("/attendees/register", attendeeController.RegisterAttendee)
			public.DELETE("/attendees/:eventId/:userId", attendeeController.DeleteAttendee)
			public.GET("/attendees/events/:userId", attendeeController.GetEventsByUser)
			public.GET("/attendees/users/:eventId", attendeeController.GetUsersByEvent)

			// USER
			public.PUT("/users/avatar", userController.UpdateAvatar)

		}

		// Protected routes (cần Bearer token)
		protected := v1.Group("")
		protected.Use(middlewares.AuthMiddleware(cfg.JwtSecret, userService))
		{
			// USER
			protected.GET("/users/profile", userController.GetProfile)
		}
	}
}
