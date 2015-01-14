package rv

type FieldHandlerCreator func(args []string) (FieldHandler, error)

type FieldHandler interface {
	Run(Request, *Field)
}

type PrecidenceFieldHandler interface {
	FieldHandler
	Precidence() int
}

type Field struct {
	Value  interface{}
	Errors []error
}

type DefaultHandler struct {
	Default string
}

func NewDefaultHandler(args []string) (FieldHandler, error) {
	return DefaultHandler{args[0]}, nil
}

func (h DefaultHandler) Run(req Request, field *Field) {
	if field.Value == nil {
		field.Value = h.Default
	}
}

func (h DefaultHandler) Precidence() int { return 900 }
