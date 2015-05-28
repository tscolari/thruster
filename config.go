package thruster

type Config struct {
	Hostname string
	Port     int
	HTTPAuth []HTTPAuth
	TLS      bool

	Certificate string
	PublicKey   string
}

type HTTPAuth struct {
	Username string
	Password string
}
