package plugin

type Authenticator interface {
	Authenticate() string
}
