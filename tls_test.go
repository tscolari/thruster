package thruster_test

import (
	"crypto/tls"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/tscolari/thruster"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TLS", func() {
	var subject *thruster.Server
	var config thruster.Config
	var port int

	handlerFunc := func(c *gin.Context) {
		c.String(200, "OK")
	}

	BeforeEach(func() {
		config = thruster.Config{Hostname: "localhost"}
		config.TLS = true
		config.PublicKey = "fixtures/server.key"
		config.Certificate = "fixtures/server.crt"
	})

	JustBeforeEach(func() {
		port = rand.Intn(8000) + 3000
		config.Port = port
		subject = thruster.NewServer(config)
		subject.AddHandler(thruster.GET, "/test", handlerFunc)
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
