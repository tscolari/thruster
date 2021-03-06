package thruster_test

import (
	"net/http"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tscolari/thruster"

	"testing"
)

func TestThruster(t *testing.T) {
	RegisterFailHandler(Fail)
	gin.SetMode(gin.TestMode)
	RunSpecs(t, "Thuster Suite")
}

func makeSimpleRequest(method, url string) *http.Response {
	client := &http.Client{}
	request, err := http.NewRequest(method, url, nil)
	Expect(err).ToNot(HaveOccurred())
	resp, err := client.Do(request)
	Expect(err).ToNot(HaveOccurred())

	return resp
}

func startServer(server *thruster.Server) {
	go func() {
		err := server.Run()
		Expect(err).ToNot(HaveOccurred())
	}()
}
