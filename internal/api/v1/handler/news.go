package v1handler

import (
	"fmt"
	"hoc-thuat-toan/utils"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type NewsHandler struct {
}

type PostNewsV1Param struct {
	Title  string `form:"title" binding:"required"`
	Status string `form:"status" binding:"required,oneof=1 2"`
}

func NewNewsHandler() *NewsHandler {
	return &NewsHandler{}
}

func (n *NewsHandler) GetNewsV1(ctx *gin.Context) {
	slug := ctx.Param("slug")

	if slug == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Get news (V1)",
			"slug":    "xin-chao",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Get news (V1)",
			"slug":    slug,
		})
	}
}

func (n *NewsHandler) PostNewsV1(ctx *gin.Context) {
	var params PostNewsV1Param

	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	image, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	if image.Size > 5<<20 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File too large > (5 MB)"})
		return
	}

	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create upload folder"})
		return
	}

	dst := fmt.Sprintf("./uploads/%s", filepath.Base(image.Filename))

	if err := ctx.SaveUploadedFile(image, dst); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot save file"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Post category (V1)",
		"title":   params.Title,
		"status":  params.Status,
		"image":   image.Filename,
		"path":    dst,
	})
}

func (n *NewsHandler) PostUploadFileNewsV1(ctx *gin.Context) {
	var params PostNewsV1Param

	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	image, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	filename, err := utils.ValidateAndSaveFile(image, "./uploads")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Post category (V1)",
		"title":   params.Title,
		"status":  params.Status,
		"image":   image.Filename,
		"path":    "./upload/" + filename,
	})
}

func (n *NewsHandler) PostUploadMultipleFileNewsV1(ctx *gin.Context) {
	const publicURL = "http://localhost:8080/images/"

	var params PostNewsV1Param

	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid multipart form"})
		return
	}

	images := form.File["images"]
	if len(images) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	var successFiles []string
	var filedFile []map[string]string
	for _, image := range images {
		filename, err := utils.ValidateAndSaveFile(image, "./uploads")
		if err != nil {
			filedFile = append(filedFile, map[string]string{
				"filename": image.Filename,
				"errors":   err.Error(),
			})
			continue
		}

		publicImageUrl := publicURL + filename
		successFiles = append(successFiles, publicImageUrl)
	}

	resp := gin.H{
		"message":      "Post news (v1)",
		"title":        params.Title,
		"status":       params.Status,
		"success_file": successFiles,
	}

	if len(filedFile) > 0 {
		resp["message"] = "Upload completed with partial errors"
		resp["error_files"] = filedFile
	}

	ctx.JSON(http.StatusOK, resp)
}
