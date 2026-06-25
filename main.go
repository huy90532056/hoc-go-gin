package main

import (
	v1handler "hoc-thuat-toan/internal/api/v1/handler"
	v2handler "hoc-thuat-toan/internal/api/v2/handler"
	"hoc-thuat-toan/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	if err := utils.RegisterValidators(); err != nil {
		panic(err)
	}

	v1 := r.Group("/api/v1")
	{
		user := v1.Group("/users")
		{
			userHandlerV1 := v1handler.NewUserHandler()

			user.GET("/", userHandlerV1.GetUsersV1)
			user.GET("/:id", userHandlerV1.GetUsersByIdV1)
			user.GET("/admin/:uuid", userHandlerV1.GetUsersByUuidV1)
			user.POST("/", userHandlerV1.PostUsersV1)
			user.PUT("/:id", userHandlerV1.PutUsersV1)
			user.DELETE("/:id", userHandlerV1.DeleteUsersV1)
		}

		product := v1.Group("/products")
		{
			productHandlerV1 := v1handler.NewProductHandler()

			product.GET("/", productHandlerV1.GetProductsV1)
			product.GET("/:slug", productHandlerV1.GetProductsBySlugV1)
			product.POST("/", productHandlerV1.PostProductsV1)
			product.PUT("/:id", productHandlerV1.PutProductsV1)
			product.DELETE("/:id", productHandlerV1.DeleteProductsV1)
		}

		category := v1.Group("/categories")
		{
			categoryHandlerV1 := v1handler.NewCategoryHandler()

			category.GET("/:category", categoryHandlerV1.GetCategoryByCategoryV1)
			category.POST("/", categoryHandlerV1.PostCategoriesV1)
		}

		news := v1.Group("/news")
		{
			categoryHandlerV1 := v1handler.NewNewsHandler()

			news.GET("/", categoryHandlerV1.GetNewsV1)
			news.GET("/:slug", categoryHandlerV1.GetNewsV1)
			news.POST("/", categoryHandlerV1.PostNewsV1)
			news.POST("/upload-file", categoryHandlerV1.PostUploadFileNewsV1)
			news.POST("/upload-multiple-file", categoryHandlerV1.PostUploadMultipleFileNewsV1)
		}
	}

	v2 := r.Group("/api/v2")
	{
		userHandlerV2 := v2handler.NewUserHandler()

		v2.GET("/users", userHandlerV2.GetUsersV2)
		v2.GET("/users/:id", userHandlerV2.GetUsersByIdV2)
		v2.POST("/users", userHandlerV2.PostUsersV2)
		v2.PUT("/users/:id", userHandlerV2.PutUsersV2)
		v2.DELETE("/users/:id", userHandlerV2.DeleteUsersV2)
	}

	r.Run(":8080")
}
