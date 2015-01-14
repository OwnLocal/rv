package rv

type Source int

func (s Source) String() string {
	return sources[s]
}

const (
	UNDEFINED Source = iota
	PATH
	QUERY
	JSON
	FORM
)

var sources = []string{
	"UNDEFINED",
	"PATH",
	"QUERY",
	"JSON",
	"FORM",
}

var sourceMap = map[string]Source{
	"path":  PATH,
	"query": QUERY,
	"json":  JSON,
	"form":  FORM,
}
