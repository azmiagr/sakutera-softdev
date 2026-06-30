package rest

import (
	"fmt"
	"os"

	"github.com/azmiagr/sakutera-softdev/internal/service"
	"github.com/azmiagr/sakutera-softdev/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type Rest struct {
	router     *gin.Engine
	service    *service.Service
	middleware middleware.Interface
}

func NewRest(service *service.Service, middleware middleware.Interface) *Rest {
	return &Rest{
		router:     gin.Default(),
		service:    service,
		middleware: middleware,
	}
}

func (r *Rest) MountEndpoint() {
	r.router.Use(r.middleware.Cors())
	baseURL := r.router.Group("/api/v1")

	auth := baseURL.Group("/auth")
	auth.POST("/register", r.Register)
	auth.POST("/verify-otp", r.VerifyOTP)

	onboarding := baseURL.Group("/onboarding")
	onboarding.Use(r.middleware.AuthenticateUser())
	onboarding.GET("/work-categories", r.GetWorkCategories)
	onboarding.POST("/work-platform", r.SelectWorkPlatform)
	onboarding.GET("/income-sources", r.GetIncomeSources)
	onboarding.POST("/income-source", r.SelectIncomeSource)

	protected := baseURL.Group("")
	protected.Use(r.middleware.AuthenticateUser())
	protected.GET("/dashboard", r.GetDashboard)

	tx := protected.Group("/transactions")
	tx.GET("/sources", r.GetTransactionSources)
	tx.POST("", r.CreateTransaction)
	tx.GET("", r.GetTransactions)
}

func (r *Rest) Run() {
	addr := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")

	r.router.Run(fmt.Sprintf("%s:%s", addr, port))
}
