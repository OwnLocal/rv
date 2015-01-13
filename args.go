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

type Range struct {
	Start int
	End   int
}

type Args struct {
	Source  Source
	Field   string
	Range   Range
	Options []string
	Default string
}
