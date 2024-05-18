package input

type Value struct {
	RawName   string
	GoName    string
	SnakeName string
	CamelName string
	Number    int32
	Comment   string
}

type Enum struct {
	GoName    string
	SnakeName string
	Comment   string
	Values    []*Value
}

func (e *Enum) GetPkgName() string {
	return "enum"
}

func (e *Enum) GetSnakeName() string {
	return e.SnakeName
}

func (e *Enum) Nil() bool {
	return e == nil
}
