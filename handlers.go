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

type TypeHandler struct {
	Type string
}

func NewTypeHandler(args []string) (FieldHandler, error) {
	// TODO: Error for types that can't be processed
	// TODO: add support for array types of basic types (right now the type of array elements is not passed along)
	return TypeHandler{args[0]}, nil
}

func (h TypeHandler) Precidence() int { return 999 }
func (h TypeHandler) Run(r Request, f *Field) {
	// TODO: Validate type
}
