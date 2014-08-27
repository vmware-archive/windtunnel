package plugin

type Plugin interface {
	Authenticate() string
	Status(token string, app string) []int
}
