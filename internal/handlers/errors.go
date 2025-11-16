package handlers

import (
	"net/http"

	"github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/errors"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			switch {
			case errors.IsTeamExists(err):
				c.JSON(http.StatusBadRequest, toErrorResponse(err))
			case errors.IsPRMerged(err) || errors.IsNotAssigned(err) || errors.IsNoCandidate(err) || errors.IsPRExists(err):
				c.JSON(http.StatusConflict, toErrorResponse(err))
			case errors.IsNotFound(err):
				c.JSON(http.StatusNotFound, toErrorResponse(err))
			default:
				if appErr, ok := err.(*errors.AppError); ok && appErr.Code == errors.CodeInvalidInput {
					c.JSON(http.StatusBadRequest, toErrorResponse(err))
				} else {
					c.JSON(http.StatusInternalServerError, toErrorResponse(
						errors.NewInvalidInput("Internal server error"),
					))
				}
			}
			c.Abort()
		}
	}
}

func toErrorResponse(err error) gin.H {
	if appErr, ok := err.(*errors.AppError); ok {
		return gin.H{
			"error": gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
			},
		}
	}

	return gin.H{
		"error": gin.H{
			"code":    "NOT_FOUND",
			"message": err.Error(),
		},
	}
}
