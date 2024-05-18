package input

const (
	GoTypeFloat32 = "float32"
	GoTypeInt64   = "int64"
	GoTypeInt32   = "int32"
	GoTypeBool    = "bool"
	GoTypeString  = "string"
	GoTypeBytes   = "[]byte"
)

type ValidateOption struct {
	Key   string
	Value string
}

type Field struct {
	GoName    string
	CamelName string
	Comment   string
	Type      string
	IsList    bool
	IsEnum    bool
	PK        bool
}

type IndexKey struct {
	GoName string
}

type Index struct {
	Keys          []*IndexKey
	PascalStoring []string
}

type Message struct {
	FileDirName string
	PkgName     string
	GoName      string
	SnakeName   string
	Comment     string
	Fields      []*Field
}

func (m *Message) GetPkgName() string {
	return m.PkgName
}

func (m *Message) GetSnakeName() string {
	return m.SnakeName
}

func (m *Message) Nil() bool {
	return m == nil
}
