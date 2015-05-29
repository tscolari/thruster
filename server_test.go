package thruster_test

import (
	"crypto/tls"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/tscolari/thruster"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func startServer(server *thruster.Server) {
	go func() {
		err := server.Run()
		Expect(err).ToNot(HaveOccurred())
	}()
}

var _ = Describe("Server", func() {
	var subject *thruster.Server

	handleFunc := func(c *gin.Context) {
		c.String(200, "OK")
	}

	Describe("#AddHandler", func() {
		var testServer *httptest.Server

		BeforeEach(func() {
			engine := gin.Default()
			subject = thruster.NewServerWithEngine(thruster.Config{}, engine)
			testServer = httptest.NewUnstartedServer(engine)
		})

		AfterEach(func() {
			testServer.Close()
		})

		requestTypes := []string{
			thruster.GET,
			thruster.POST,
			thruster.DELETE,
			thruster.PUT,
		}

		for _, requestType := range requestTypes {
			Context(requestType, func() {
				It("registers and serve a handler function", func() {
					subject.AddHandler(requestType, "/path", handleFunc)
					testServer.Start()

					client := &http.Client{}
					request, err := http.NewRequest(requestType, testServer.URL+"/path", nil)
					resp, err := client.Do(request)

					Expect(err).ToNot(HaveOccurred())
					Expect(resp.StatusCode).To(Equal(http.StatusOK))

					body, err := ioutil.ReadAll(resp.Body)
					Expect(err).ToNot(HaveOccurred())
					Expect(string(body)).To(Equal("OK"))
				})
			})
		}
	})

	Context("Server configuration", func() {
		var config thruster.Config
		var port int

		BeforeEach(func() {
			config = thruster.Config{Hostname: "localhost"}
		})

		JustBeforeEach(func() {
			port = rand.Intn(8000) + 3000
			config.Port = port
			subject = thruster.NewServer(config)
			subject.AddHandler(thruster.GET, "/test", handleFunc)
		})

		Context("server url and port", func() {
			It("listens to the hostname and port in the configuration", func() {
				startServer(subject)
				resp, err := http.Get("http://localhost:" + strconv.Itoa(port) + "/test")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})
		})

		Context("http auth", func() {
			var httpAuth []thruster.HTTPAuth

			BeforeEach(func() {
				httpAuth = []thruster.HTTPAuth{
					thruster.HTTPAuth{Username: "admin", Password: "passwd"},
					thruster.HTTPAuth{Username: "root", Password: "passwd2"},
				}
				config.HTTPAuth = httpAuth
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

		Context("TLS", func() {

			BeforeEach(func() {
				config.TLS = true
				config.PublicKey = "fixtures/server.key"
				config.Certificate = "fixtures/server.crt"
			})

			It("listens only to https when in TLS mode", func() {
				go func() {
					err := subject.Run()
					Expect(err).ToNot(HaveOccurred())
				}()

				tr := &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				}
				client := &http.Client{Transport: tr}

				url := "https://localhost:" + strconv.Itoa(port) + "/test"

				Eventually(func() int {
					resp, err := client.Get(url)
					if err != nil {
						return 0
					}
					return resp.StatusCode
				}).Should(Equal(http.StatusOK))

				url = "http://localhost:" + strconv.Itoa(port) + "/test"
				_, err := http.Get(url)
				Expect(err).ToNot(MatchError("malformed HTTP response"))
			})

			Context("missing certificates", func() {
				BeforeEach(func() {
					config.Certificate = ""
					config.PublicKey = ""
				})

				It("fails to start", func() {
					err := subject.Run()
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})
