package plugin

type Statuser interface {
	Status(token string, app string) []int
}
