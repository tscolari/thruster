package thruster_test

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

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
		var testServer *httptest.Server
		var jsonHandler thruster.JSONHandler

		BeforeEach(func() {
			engine := gin.Default()
			subject = thruster.NewServerWithEngine(thruster.Config{}, engine)
			testServer = httptest.NewUnstartedServer(engine)
		})

		AfterEach(func() {
			testServer.Close()
		})

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

			Context("certificates are inlined in the config", func() {
				BeforeEach(func() {
					config.Certificate = `-----BEGIN CERTIFICATE-----
MIIDBjCCAe4CCQCRvfTE6erY5DANBgkqhkiG9w0BAQUFADBFMQswCQYDVQQGEwJB
VTETMBEGA1UECBMKU29tZS1TdGF0ZTEhMB8GA1UEChMYSW50ZXJuZXQgV2lkZ2l0
cyBQdHkgTHRkMB4XDTE1MDUyNzIyMTIyOVoXDTE2MDUyNjIyMTIyOVowRTELMAkG
A1UEBhMCQVUxEzARBgNVBAgTClNvbWUtU3RhdGUxITAfBgNVBAoTGEludGVybmV0
IFdpZGdpdHMgUHR5IEx0ZDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
AK8RlyC/+7Q1l5aZI87nlx4bmTK14UkCIM8gNGqbdkRNptZZ+AwbAsR3ZlfO2fG+
CtQGln+/tVTzZMWHelEJDXoAInCmecBD01n1jCbUWc+z9hXGSqIn722k1lWK3nMH
txb+5RjTuekAqgc95Rg3iXZ6HADdU7NwJshATzS7F6eZwwJf58/8XNWYuqL8c8Z+
Gg1cTMVx5XcBCXFQszpZnLGj98ANVDvQqo4hB0FNnrMi1NjPmueBRROV9p5Cju9s
uR3v5Hh3M1Uecel8DIyn9lDg8VCvXLQggTXe26JZgjgWfbLQDWU7YnXGRz9SGWnv
yFp7W8FTMrA2xG5OgDmjlx8CAwEAATANBgkqhkiG9w0BAQUFAAOCAQEAT6DRSI1D
CrQIPupf/AYCE6jiF4qKDKtzl7nQ7AQfzNygTYPkgmYCLjSefQNiosrfc+CxxPGX
5AHeIdEpBno797JJtPzahP3XvOt0UEhg0MyI3q4z4xzah8vqOSBVc1GzS3z+QWs0
+uxVGOJ+6/WLC0AFCq/FJqchHNWnUjgpIfWMd6IjScu+v6+D/3/m6a08Y2gSZrpN
rb2y9kORgnM7qK2fyutaippjkir91nK1ggjrHE3to2WMGHEE/J/l6VIVl90qwzMc
dBbJmQuRs/gY9qzEY7Ga+V6b4R3GN4LXI7w86v26cR+Dk02yGs9FNYMpCyAz9ET6
ZUEbxc1l9V0aQQ==
-----END CERTIFICATE-----`
					config.PublicKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEArxGXIL/7tDWXlpkjzueXHhuZMrXhSQIgzyA0apt2RE2m1ln4
DBsCxHdmV87Z8b4K1AaWf7+1VPNkxYd6UQkNegAicKZ5wEPTWfWMJtRZz7P2FcZK
oifvbaTWVYrecwe3Fv7lGNO56QCqBz3lGDeJdnocAN1Ts3AmyEBPNLsXp5nDAl/n
z/xc1Zi6ovxzxn4aDVxMxXHldwEJcVCzOlmcsaP3wA1UO9CqjiEHQU2esyLU2M+a
54FFE5X2nkKO72y5He/keHczVR5x6XwMjKf2UODxUK9ctCCBNd7bolmCOBZ9stAN
ZTtidcZHP1IZae/IWntbwVMysDbEbk6AOaOXHwIDAQABAoIBAGT6R2pLcfoywznJ
IN9Rs1dZYdbfE4+R26y8jZ9EBkZFZ8rRYAJTfhgmKnDRTeJi1EoRdrM+t2/FZ8WL
bCDbkNtiwnqpeyZLuNd1ix5Gc3sa+QD8O8YmNLLQVhRHIiHFPHTWFvxn+x6LFIdS
yxZZyj79FbPl9UZVlPkCJu1qUK2JTMkNuGgabUzh8q9cdiRg/zX6CNc26RZjbsrS
HxHU24fzaBnyiPM5SJYJlzBymNWC/lwxOEyq5GGVyyfodU6/SJCxaYoG53xm05qC
3xkcFJW67zGvx0MPYQ25Ms5xUu9fn+lMb1EVUuN7oG6mPlEtG/2Zom/YvBThIfl+
Fv9cxsECgYEA5ynqQKGvO3GSqL8alHOjf6Mr2DeZJgp6/njyVjynzgcRYEQbdZrl
F8iPLRBY/b5YDnzGNtwMy2Fpob55e5sF0dPRw+VrV4dSSRy+GAB5+QPAqOCLLzvd
ybVvN0YBPI+CRGAt4pG9P18ywha1SY5GuVWEjzWj+N4Lr8i/myNK9r8CgYEAweDM
tTwOld9P1/Nza1zNLulPg6sOrHBqBkurJ7okKzYP8KVd4AZ7oNzdgtCNoHYK8X23
3Mg+zDhYTk6adBdOEvbb7SY4F4NlEj7mSzTGaLmTBtSb0wbVavJkWxR7pCzFSZA0
dxzYlNgNAlfG75jVy5+UELZh+YHJS9mvJ6K816ECgYEApegTVCe22Hb+x1XBAeKs
6aJ2iUv+Aqtq8tBjPTlzRg8UjX7UJmfxHEy5VaJx/Etsb5lluWHdXOqhIZDPJ8Nv
PdVEq9AwZjWc/RQ/6oINCIeE8q+VtWTGHUq2c3ku0gQ9fk15IS9wH9d3Wo1pt00B
vWp/JTleYfMbeCIgQnvmBYkCgYA4SmClLi697Pxtos2cGnGocS0Y+Y1lG65s7YNg
IXdm5Gd0Y08CQF+csQPPe2XjdOJwgyPjAnDZMnLRKZlGo42TjAEGtdYLXab2yTRs
GYKR3W+GyCwF9TH5vy7MEwJjBGyzkx7ohoOLk78TMxEbd7B7UnXW9F016CzdzPJB
+8oAgQKBgQCRSZ+ZyqZSozW9icsT028wKMjI65ev/DAxHezlqXuyQUOHZ+aMrf6x
RA4b1DdDpcN6+IcFsAwPhBTs0HP33VfPxSdVDo5tYyc7RVdmTFVKa/VNf2mNtKJR
uuWtL8DK0AMuP2tpmAxTfB6BGitzCFDXiK9aM2Igc7a/OhBIJNqwNQ==
-----END RSA PRIVATE KEY-----`
				})

				It("works the same", func() {
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
				})
			})
		})
	})
})
