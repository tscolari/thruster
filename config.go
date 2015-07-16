package thruster

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Hostname string     `yaml:"hostname"`
	Port     int        `yaml:"port"`
	HTTPAuth []HTTPAuth `yaml:"http_auth"`
	TLS      bool       `yaml:"tls"`

	Certificate string `yaml:"certificate"`
	PublicKey   string `yaml:"public_key"`
}

type HTTPAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func NewHTTPAuth(username, password string) HTTPAuth {
	return HTTPAuth{
		Username: username,
		Password: password,
	}
}

func NewConfig(path string) (Config, error) {
	config := Config{}
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	return config, err
}
