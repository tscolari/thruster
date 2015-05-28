package thruster

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func New(config Config) *Server {
	return &Server{
		config: config,
		router: gin.Default(),
	}
}

type Server struct {
	config      Config
	router      *gin.Engine
	routerGroup *gin.RouterGroup
}

func (s *Server) Router() *gin.Engine {
	return s.router
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
		s.routerGroup = s.router.Group("/")
		return s.routerGroup
	}

	accounts := gin.Accounts{}
	for _, account := range s.config.HTTPAuth {
		accounts[account.Username] = account.Password
	}

	s.routerGroup = s.router.Group("/", gin.BasicAuth(accounts))
	return s.routerGroup
}

func (s *Server) Run() error {
	address := s.config.Hostname + ":" + strconv.Itoa(s.config.Port)

	var err error
	if s.config.TLS {
		err = http.ListenAndServeTLS(address, s.config.Certificate, s.config.PublicKey, s.router)
	} else {
		err = http.ListenAndServe(address, s.router)
	}

	return err
}
