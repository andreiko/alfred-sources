package sources

type Item interface {
	Attributes() map[string]interface{}
	GetRank(string) int
	Autocomplete() string
	LessThan(Item) bool
}

type Source interface {
	Query(string) []Item
	Update() error
	Id() string
}
