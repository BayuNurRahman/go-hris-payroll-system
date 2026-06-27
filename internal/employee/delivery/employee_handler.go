package delivery

import "github.com/gin-gonic/gin"

func RegisterEmployeeRoutes(r *gin.RouterGroup) {
	r.GET("/employees", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "list employees"})
	})
	r.POST("/employees", func(c *gin.Context) {
		c.JSON(201, gin.H{"message": "create employee"})
	})
}
