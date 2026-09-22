package property

import (
	"fmt"
	"gin-quickstart/service"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func AddProperty(propertyService *service.PropertyService) gin.HandlerFunc {
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

		name := ctx.PostForm("name")
		if name == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "name is required",
			})
			return
		}

		file, err := ctx.FormFile("photo")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "photo is required",
			})
			return
		}

		uploadDir := "./uploads"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to create upload directory",
			})
			return
		}

		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
		filePath := filepath.Join(uploadDir, filename)

		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to save file",
			})
			return
		}

		photoID := 1
		data, err := propertyService.AddProperty(ctx, userID, name, photoID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			os.Remove(filePath)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Added successfully",
			"data":    data,
		})
	}
}
