package thruster

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type JSONHandler func(*gin.Context) (interface{}, error)

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

func (s *Server) Run() error {
	address := s.config.Hostname + ":" + strconv.Itoa(s.config.Port)

	var err error
	if s.config.TLS {
		certificate, err := s.certificate()
		if err != nil {
			return err
		}
		publicKey, err := s.publicKey()
		if err != nil {
			return err
		}

		err = http.ListenAndServeTLS(address, certificate, publicKey, s.engine)
	} else {
		err = http.ListenAndServe(address, s.engine)
	}

	return err
}

func (s *Server) AddHandler(method, path string, handler gin.HandlerFunc) {
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
	}
}

func (s *Server) AddJSONHandler(method, path string, handler JSONHandler) {
	ginHandler := func(c *gin.Context) {
		data, err := handler(c)
		if err != nil {
			c.JSON(s.statusError(err), map[string]string{"error": err.Error()})
			return
		}
		c.JSON(s.statusOK(method), data)
	}
	s.AddHandler(method, path, ginHandler)
}

func (s *Server) statusError(err error) int {
	if err == ErrNotFound {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}

func (s *Server) statusOK(method string) int {
	if method == POST {
		return http.StatusCreated
	}
	return http.StatusOK
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

func (s *Server) certificate() (string, error) {
	if len(s.config.Certificate) < 500 {
		return s.config.Certificate, nil
	}

	certFile, err := ioutil.TempFile("", "cert")
	if err != nil {
		return "", err
	}

	_, err = certFile.Write([]byte(s.config.Certificate))
	if err != nil {
		return "", err
	}

	certFile.Close()

	return certFile.Name(), nil
}

func (s *Server) publicKey() (string, error) {
	if len(s.config.PublicKey) < 500 {
		return s.config.PublicKey, nil
	}

	keyFile, err := ioutil.TempFile("", "key")
	if err != nil {
		return "", err
	}

	_, err = keyFile.Write([]byte(s.config.PublicKey))
	if err != nil {
		return "", err
	}

	keyFile.Close()

	return keyFile.Name(), nil
}
