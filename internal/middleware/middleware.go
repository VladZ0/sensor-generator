package middleware

import (
	"sensors-generator/internal/apperror"

	"github.com/gin-gonic/gin"
)

func HandleErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				switch err.Err.(type) {
				case *apperror.AppError:
					appError := err.Err.(*apperror.AppError)
					c.JSON(appError.TransportCode, gin.H{"error": appError.Message})
					c.Abort()
					return
				}
			}

			c.JSON(apperror.ErrInternalSystem.TransportCode, gin.H{"error": apperror.ErrInternalSystem.Message})
			c.Abort()
		}
	}
}
