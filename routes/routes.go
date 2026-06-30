package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	authDelivery "hris-payroll/internal/auth/delivery"
	leaveDelivery "hris-payroll/internal/leave/delivery"
	payrollDelivery "hris-payroll/internal/payroll/delivery"
	"hris-payroll/middleware"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	router *gin.Engine,
	db *gorm.DB,
	authHandler *authDelivery.AuthHandler,
	payrollHandler *payrollDelivery.PayrollHandler,
	leaveHandler *leaveDelivery.LeaveHandler,
) {
	// API v1 group
	api := router.Group("/api/v1")
	{
		// --- ENDPOINT PUBLIC / AUTHENTICATION ---
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", authHandler.LoginHandler)
			authGroup.POST("/register", middleware.JWTAuth(db), middleware.RequireRole("HRD"), authHandler.RegisterHandler)
			authGroup.POST("/logout", middleware.JWTAuth(db), authHandler.LogoutHandler)
		}

		// Payroll Routes (Protected by JWT & Only for HRD)
		payrollGroup := api.Group("/payroll")
		payrollGroup.Use(middleware.JWTAuth(db), middleware.RequireRole("HRD"))
		{
			payrollGroup.POST("/process", payrollHandler.ProcessPayrollHandler)
		}

		// Leave Routes (Requires Login & Open for Employee & HRD)
		leaveGroup := api.Group("/leaves")
		leaveGroup.Use(middleware.JWTAuth(db))
		{
			leaveGroup.POST("", leaveHandler.CreateLeaveHandler)
			leaveGroup.GET("/:id", leaveHandler.GetLeaveDetailHandler)
		}

		// Ping endpoint for testing
		api.GET("/ping", middleware.JWTAuth(db), func(c *gin.Context) {
			role, _ := c.Get("user_role")
			c.JSON(200, gin.H{
				"message": "Anda berhasil menembus middleware keamanan!",
				"role":    role,
			})
		})
	}
}
