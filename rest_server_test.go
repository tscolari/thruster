package thruster_test

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tscolari/thruster"
	"github.com/tscolari/thruster/fakes"
)

var _ = Describe("RestServer", func() {
	var (
		subject    *thruster.RestServer
		testServer *httptest.Server
	)

	BeforeEach(func() {
		engine := gin.Default()
		subject = thruster.NewRestServerWithEngine(thruster.Config{}, engine)
		testServer = httptest.NewUnstartedServer(engine)
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
})
