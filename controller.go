package thruster

import "github.com/gin-gonic/gin"

type JSONController interface {
	Index(context *gin.Context) (interface{}, error)
	Show(context *gin.Context) (interface{}, error)
	Create(context *gin.Context) (interface{}, error)
	Update(context *gin.Context) (interface{}, error)
	Destroy(context *gin.Context) (interface{}, error)
}

type Controller interface {
	Index(context *gin.Context)
	Show(context *gin.Context)
	Create(context *gin.Context)
	Update(context *gin.Context)
	Destroy(context *gin.Context)
}
