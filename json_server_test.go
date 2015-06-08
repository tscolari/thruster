package thruster_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tscolari/thruster"
	"github.com/tscolari/thruster/fakes"
)

type JSONResponse struct {
	Action string
}

var _ = Describe("JSONServer", func() {
	var (
		controller *fakes.FakeController
		subject    *thruster.JSONServer
		testServer *httptest.Server
	)

	BeforeEach(func() {
		controller = &fakes.FakeController{}
		engine := gin.Default()
		subject = thruster.NewJSONServerWithEngine(thruster.Config{}, engine)
		testServer = httptest.NewUnstartedServer(engine)
	})

	Describe("AddResource", func() {
		BeforeEach(func() {
			subject.AddResource("/users", controller)
			testServer.Start()
		})

		Context("CRUD actions", func() {
			Context("response codes", func() {
				It("Index returns 200 on success", func() {
					resp := makeSimpleRequest("GET", testServer.URL+"/users")
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
				})

				It("Index returns 500 on error", func() {
					controller.IndexReturns(nil, errors.New("failed"))
					resp := makeSimpleRequest("GET", testServer.URL+"/users")
					Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
				})

				It("Show returns 200 on success", func() {
					resp := makeSimpleRequest("GET", testServer.URL+"/users/1")
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
				})

				It("Show returns 500 on error", func() {
					controller.ShowReturns(nil, errors.New("failed"))
					resp := makeSimpleRequest("GET", testServer.URL+"/users/1")
					Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
				})

				It("Show returns 404 on NotFound error", func() {
					controller.ShowReturns(nil, thruster.ErrNotFound)
					resp := makeSimpleRequest("GET", testServer.URL+"/users/1")
					Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
				})

				It("Create returns 200 on success", func() {
					resp := makeSimpleRequest("POST", testServer.URL+"/users")
					Expect(resp.StatusCode).To(Equal(http.StatusCreated))
				})

				It("Create returns 500 on error", func() {
					controller.CreateReturns(nil, errors.New("failed"))
					resp := makeSimpleRequest("POST", testServer.URL+"/users")
					Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
				})

				It("Update returns 200 on success", func() {
					resp := makeSimpleRequest("PUT", testServer.URL+"/users/1")
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
				})

				It("Update returns 500 on error", func() {
					controller.UpdateReturns(nil, errors.New("failed"))
					resp := makeSimpleRequest("PUT", testServer.URL+"/users/1")
					Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
				})

				It("Update returns 404 on NotFound error", func() {
					controller.UpdateReturns(nil, thruster.ErrNotFound)
					resp := makeSimpleRequest("PUT", testServer.URL+"/users/1")
					Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
				})

				It("Delete returns 200 on success", func() {
					resp := makeSimpleRequest("DELETE", testServer.URL+"/users/1")
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
				})

				It("Delete returns 500 on error", func() {
					controller.DestroyReturns(nil, errors.New("failed"))
					resp := makeSimpleRequest("DELETE", testServer.URL+"/users/1")
					Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
				})

				It("Delete returns 404 on NotFound error", func() {
					controller.DestroyReturns(nil, thruster.ErrNotFound)
					resp := makeSimpleRequest("DELETE", testServer.URL+"/users/1")
					Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
				})
			})

			Context("response bodies", func() {
				It("maps GET / to Index and returns a JSON", func() {
					controller.IndexReturns(map[string]string{"action": "index"}, nil)
					resp := makeSimpleRequest("GET", testServer.URL+"/users")

					var jsonResponse JSONResponse
					body, err := ioutil.ReadAll(resp.Body)
					Expect(err).ToNot(HaveOccurred())
					json.Unmarshal(body, &jsonResponse)

					Expect(jsonResponse.Action).To(Equal("index"))
				})

				It("maps GET /:id to Show and returns a JSON", func() {
					controller.ShowReturns(map[string]string{"action": "show"}, nil)
					resp := makeSimpleRequest("GET", testServer.URL+"/users/10")

					var jsonResponse JSONResponse
					body, err := ioutil.ReadAll(resp.Body)
					Expect(err).ToNot(HaveOccurred())
					json.Unmarshal(body, &jsonResponse)

					Expect(jsonResponse.Action).To(Equal("show"))
				})

				It("maps POST / to Create and returns a JSON", func() {
					controller.CreateReturns(map[string]string{"action": "create"}, nil)
					resp := makeSimpleRequest("POST", testServer.URL+"/users")

					var jsonResponse JSONResponse
					body, err := ioutil.ReadAll(resp.Body)
					Expect(err).ToNot(HaveOccurred())
					json.Unmarshal(body, &jsonResponse)

					Expect(jsonResponse.Action).To(Equal("create"))
				})

				It("maps PUT /:id to Update and returns a JSON", func() {
					controller.UpdateReturns(map[string]string{"action": "update"}, nil)
					resp := makeSimpleRequest("PUT", testServer.URL+"/users/5")

					var jsonResponse JSONResponse
					body, err := ioutil.ReadAll(resp.Body)
					Expect(err).ToNot(HaveOccurred())
					json.Unmarshal(body, &jsonResponse)

					Expect(jsonResponse.Action).To(Equal("update"))
				})

				It("maps DELETE /:id to Destroy and returns a JSON", func() {
					controller.DestroyReturns(map[string]string{"action": "destroy"}, nil)
					resp := makeSimpleRequest("DELETE", testServer.URL+"/users/5")

					var jsonResponse JSONResponse
					body, err := ioutil.ReadAll(resp.Body)
					Expect(err).ToNot(HaveOccurred())
					json.Unmarshal(body, &jsonResponse)

					Expect(jsonResponse.Action).To(Equal("destroy"))
				})
			})
		})
	})
})
