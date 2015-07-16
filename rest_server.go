package thruster

import "github.com/gin-gonic/gin"

type RestServer struct {
	Server
}

func NewRestServerWithEngine(config Config, engine *gin.Engine) *RestServer {
	server := Server{
		config: config,
		engine: engine,
	}
	return &RestServer{server}
}

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

type ControllerHandler func(context *gin.Context) (interface{}, error)

func (s *RestServer) AddJSONResource(path string, controller JSONController) {
	s.AddJSONHandler(GET, path, controller.Index)
	s.AddJSONHandler(GET, path+"/:id", controller.Show)
	s.AddJSONHandler(POST, path, controller.Create)
	s.AddJSONHandler(PUT, path+"/:id", controller.Update)
	s.AddJSONHandler(DELETE, path+"/:id", controller.Destroy)
}

func (s *RestServer) AddResource(path string, controller Controller) {
	s.AddHandler(GET, path, controller.Index)
	s.AddHandler(GET, path+"/:id", controller.Show)
	s.AddHandler(POST, path, controller.Create)
	s.AddHandler(PUT, path+"/:id", controller.Update)
	s.AddHandler(DELETE, path+"/:id", controller.Destroy)
}
