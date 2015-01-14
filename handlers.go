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
