package thruster

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewServer(config Config) *Server {
	return &Server{
		config: config,
		engine: gin.Default(),
	}
}

func NewServerWithEngine(config Config, engine *gin.Engine) *Server {
	return &Server{
		config: config,
		engine: engine,
	}
}

type Server struct {
	config      Config
	engine      *gin.Engine
	routerGroup *gin.RouterGroup
}

const (
	GET    string = "GET"
	POST   string = "POST"
	PUT    string = "PUT"
	DELETE string = "DELETE"
)

func (s *Server) AddHandler(method, path string, handler func(*gin.Context)) {
	method = strings.ToUpper(method)
	switch method {
	case GET:
		s.group().GET(path, handler)
	case POST:
		s.group().POST(path, handler)
	case PUT:
		s.group().PUT(path, handler)
	case DELETE:
		s.group().DELETE(path, handler)
		//default:
		//log.Error("Couldn't register route: '%s:%s'", method, path)
	}
}

func (s *Server) group() *gin.RouterGroup {
	if s.routerGroup != nil {
		return s.routerGroup
	}

	if len(s.config.HTTPAuth) == 0 {
		s.routerGroup = s.engine.Group("/")
		return s.routerGroup
	}

	accounts := gin.Accounts{}
	for _, account := range s.config.HTTPAuth {
		accounts[account.Username] = account.Password
	}

	s.routerGroup = s.engine.Group("/", gin.BasicAuth(accounts))
	return s.routerGroup
}

func (s *Server) Run() error {
	address := s.config.Hostname + ":" + strconv.Itoa(s.config.Port)

	var err error
	if s.config.TLS {
		err = http.ListenAndServeTLS(address, s.config.Certificate, s.config.PublicKey, s.engine)
	} else {
		err = http.ListenAndServe(address, s.engine)
	}

	return err
}
