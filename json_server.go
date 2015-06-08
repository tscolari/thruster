package thruster

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSONServer struct {
	Server
}

func NewJSONServerWithEngine(config Config, engine *gin.Engine) *JSONServer {
	server := Server{
		config: config,
		engine: engine,
	}
	return &JSONServer{server}
}

type Controller interface {
	Index(context *gin.Context) (interface{}, error)
	Show(context *gin.Context) (interface{}, error)
	Create(context *gin.Context) (interface{}, error)
	Update(context *gin.Context) (interface{}, error)
	Destroy(context *gin.Context) (interface{}, error)
}
type ControllerHandler func(context *gin.Context) (interface{}, error)

func (s *JSONServer) AddResource(path string, controller Controller) {
	jsonWrapper := func(okStatus int, handler ControllerHandler) func(c *gin.Context) {
		return func(c *gin.Context) {
			data, err := handler(c)
			if err != nil {
				errorStatus := http.StatusInternalServerError
				if err == ErrNotFound {
					errorStatus = http.StatusNotFound
				}

				c.JSON(errorStatus, map[string]string{"error": err.Error()})
				return
			}
			c.JSON(okStatus, data)
		}
	}

	s.AddHandler(GET, path, jsonWrapper(http.StatusOK, controller.Index))
	s.AddHandler(GET, path+"/:id", jsonWrapper(http.StatusOK, controller.Show))
	s.AddHandler(POST, path, jsonWrapper(http.StatusCreated, controller.Create))
	s.AddHandler(PUT, path+"/:id", jsonWrapper(http.StatusOK, controller.Update))
	s.AddHandler(DELETE, path+"/:id", jsonWrapper(http.StatusOK, controller.Destroy))
}
