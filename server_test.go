package thruster_test

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/tscolari/thruster"
	"github.com/tscolari/thruster/fakes"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	var subject *thruster.Server
	var testServer *httptest.Server

	handleFunc := func(c *gin.Context) {
		c.String(200, "OK")
	}

	BeforeEach(func() {
		engine := gin.Default()
		subject = thruster.NewServerWithEngine(thruster.Config{}, engine)
		testServer = httptest.NewUnstartedServer(engine)
	})

	AfterEach(func() {
		testServer.Close()
	})

	Describe("#AddHandler", func() {

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

					resp := makeSimpleRequest(requestType, testServer.URL+"/path")
					Expect(resp.StatusCode).To(Equal(http.StatusOK))

					body, err := ioutil.ReadAll(resp.Body)
					Expect(err).ToNot(HaveOccurred())
					Expect(string(body)).To(Equal("OK"))
				})
			})
		}
	})

	Describe("#AddJSONHandler", func() {
		var jsonHandler thruster.JSONHandler

		Context("GET", func() {
			BeforeEach(func() {
				jsonHandler = func(c *gin.Context) (interface{}, error) {
					return map[string]string{"key": "value"}, nil
				}
			})

			It("returns 200 on success", func() {
				subject.AddJSONHandler(thruster.GET, "/path", jsonHandler)
				testServer.Start()

				resp := makeSimpleRequest(thruster.GET, testServer.URL+"/path")
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			It("responds with a JSON", func() {
				subject.AddJSONHandler(thruster.GET, "/path", jsonHandler)
				testServer.Start()

				resp := makeSimpleRequest(thruster.GET, testServer.URL+"/path")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(strings.Trim(string(body), "\n")).To(Equal(`{"key":"value"}`))
			})
		})

		Context("POST", func() {
			It("returns 201 on success", func() {
				subject.AddJSONHandler(thruster.POST, "/path", jsonHandler)
				testServer.Start()

				resp := makeSimpleRequest(thruster.POST, testServer.URL+"/path")
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))
			})

			It("responds with a JSON", func() {
				subject.AddJSONHandler(thruster.POST, "/path", jsonHandler)
				testServer.Start()

				resp := makeSimpleRequest(thruster.POST, testServer.URL+"/path")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(strings.Trim(string(body), "\n")).To(Equal(`{"key":"value"}`))
			})
		})

		Context("when the handler returns an error", func() {
			Context("an unknown error", func() {
				It("returns 500 on the registered route", func() {
					jsonHandler = func(c *gin.Context) (interface{}, error) {
						return nil, errors.New("failed")
					}

					subject.AddJSONHandler("GET", "/path", jsonHandler)
					testServer.Start()

					resp := makeSimpleRequest("GET", testServer.URL+"/path")
					Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
				})
			})

			Context("not found error", func() {
				It("returns 404 on the registered route", func() {
					jsonHandler = func(c *gin.Context) (interface{}, error) {
						return nil, thruster.ErrNotFound
					}

					subject.AddJSONHandler("GET", "/path", jsonHandler)
					testServer.Start()

					resp := makeSimpleRequest("GET", testServer.URL+"/path")
					Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
				})
			})
		})
	})

	Describe("AddJSONResource", func() {
		var controller *fakes.FakeJSONController

		BeforeEach(func() {
			controller = &fakes.FakeJSONController{}

			subject.AddJSONResource("/users", controller)
			testServer.Start()
		})

		Context("CRUD actions", func() {
			Context("GET /resources", func() {
				It("calls the controller Index method", func() {
					makeSimpleRequest(thruster.GET, testServer.URL+"/users")
					Expect(controller.IndexCallCount()).To(Equal(1))
				})
			})

			Context("GET /resources/ID", func() {
				It("calls the controller Show method", func() {
					makeSimpleRequest(thruster.GET, testServer.URL+"/users/1")
					Expect(controller.ShowCallCount()).To(Equal(1))
				})
			})

			Context("POST /resources", func() {
				It("calls the controller Create method", func() {
					makeSimpleRequest(thruster.POST, testServer.URL+"/users")
					Expect(controller.CreateCallCount()).To(Equal(1))
				})
			})

			Context("PUT /resources/ID", func() {
				It("calls the controller Update method", func() {
					makeSimpleRequest(thruster.PUT, testServer.URL+"/users/1")
					Expect(controller.UpdateCallCount()).To(Equal(1))
				})
			})

			Context("DELETE /resources/1", func() {
				It("calls the controller Destroy method", func() {
					makeSimpleRequest(thruster.DELETE, testServer.URL+"/users/1")
					Expect(controller.DestroyCallCount()).To(Equal(1))
				})
			})
		})
	})

	Describe("AddResource", func() {
		var controller *fakes.FakeController

		BeforeEach(func() {
			controller = &fakes.FakeController{}

			subject.AddResource("/users", controller)
			testServer.Start()
		})

		Context("CRUD actions", func() {
			Context("GET /resources", func() {
				It("calls the controller Index method", func() {
					makeSimpleRequest(thruster.GET, testServer.URL+"/users")
					Expect(controller.IndexCallCount()).To(Equal(1))
				})
			})

			Context("GET /resources/ID", func() {
				It("calls the controller Show method", func() {
					makeSimpleRequest(thruster.GET, testServer.URL+"/users/1")
					Expect(controller.ShowCallCount()).To(Equal(1))
				})
			})

			Context("POST /resources", func() {
				It("calls the controller Create method", func() {
					makeSimpleRequest(thruster.POST, testServer.URL+"/users")
					Expect(controller.CreateCallCount()).To(Equal(1))
				})
			})

			Context("PUT /resources/ID", func() {
				It("calls the controller Update method", func() {
					makeSimpleRequest(thruster.PUT, testServer.URL+"/users/1")
					Expect(controller.UpdateCallCount()).To(Equal(1))
				})
			})

			Context("DELETE /resources/1", func() {
				It("calls the controller Destroy method", func() {
					makeSimpleRequest(thruster.DELETE, testServer.URL+"/users/1")
					Expect(controller.DestroyCallCount()).To(Equal(1))
				})
			})
		})
	})

	Context("Server configuration", func() {
		var config thruster.Config
		var port int

		BeforeEach(func() {
			config = thruster.Config{Hostname: "localhost"}
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
	})
})
