package thruster_test

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tscolari/thruster"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HTTP Auth", func() {
	var subject *thruster.Server
	var config thruster.Config
	var port int
	var httpAuth []thruster.HTTPAuth

	handlerFunc := func(c *gin.Context) {
		c.String(200, "OK")
	}

	BeforeEach(func() {
		config = thruster.Config{Hostname: "localhost"}
		httpAuth = []thruster.HTTPAuth{
			thruster.NewHTTPAuth("admin", "passwd"),
			thruster.NewHTTPAuth("root", "passwd2"),
		}
		config.HTTPAuth = httpAuth
	})

	JustBeforeEach(func() {
		port = rand.Intn(8000) + 3000
		config.Port = port
		subject = thruster.NewServer(config)
		subject.AddHandler(thruster.GET, "/test", handlerFunc)
	})

	It("returns 401 when wrong credentials are given", func() {
		startServer(subject)
		resp, err := http.Get("http://user:passwd@localhost:" + strconv.Itoa(port) + "/test")
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
	})

	It("returns 401 when no credential is given", func() {
		startServer(subject)
		resp, err := http.Get("http://localhost:" + strconv.Itoa(port) + "/test")
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
	})

	It("returns 200 when the correct credentials are given", func() {
		startServer(subject)
		for _, credential := range httpAuth {
			url := "http://" + credential.Username + ":" + credential.Password + "@localhost:" + strconv.Itoa(port) + "/test"
			resp, err := http.Get(url)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		}
	})
})
