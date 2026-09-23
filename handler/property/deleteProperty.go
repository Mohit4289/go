package property

import (
	"gin-quickstart/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteProperty(propertyService *service.PropertyService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDVal, ok := ctx.Get("user_id")
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "userID is blank",
			})
			return
		}

		userID, ok := userIDVal.(int)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid user_id",
			})
			return
		}

		err := propertyService.DeleteProperty(ctx, userID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Deleted successfully",
		})
	}
}
