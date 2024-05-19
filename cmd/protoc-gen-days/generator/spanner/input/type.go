package input

const (
	GoTypeInt64  = "int64"
	GoTypeInt32  = "int32"
	GoTypeBool   = "bool"
	GoTypeString = "string"
	GoTypeBytes  = "byte"
	GoTypeTime   = "time.Time"

	DBTypeInt    = "INT64"
	DBTypeBool   = "BOOL"
	DBTypeString = "STRING(MAX)"
	DBTypeBytes  = "BYTES(MAX)"
	DBTypeTime   = "TIMESTAMP"

	SetTypeInt32  = "i32set"
	SetTypeInt64  = "i64set"
	SetTypeString = "strset"
)

type Field struct {
	GoName    string
	SnakeName string
	CamelName string
	Comment   string
	GoType    string
	DBType    string
	SetType   string
	IsEnum    bool
	IsList    bool
	PK        bool
	Desc      bool
}

type IndexKey struct {
	GoName string
	Desc   bool
}

type Index struct {
	Keys          []*IndexKey
	Unique        bool
	NullFiltered  bool
	PascalStoring []string
}

type Interleave struct {
	GoName string
}

type TTL struct {
	TimestampColumnGoName string
	Days                  int32
}

type Message struct {
	PkgName            string
	GoName             string
	SnakeName          string
	CamelName          string
	Comment            string
	Fields             []*Field
	Indexes            []*Index
	Interleave         *Interleave
	TTL                *TTL
	NeedCommonResponse bool
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
