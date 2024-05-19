package input

type Method struct {
	GoName           string
	Comment          string
	IdempotencyLevel string
}

type Service struct {
	PkgName   string
	GoPkgName string
	CamelName string
	GoName    string
	SnakeName string
	Methods   []*Method
}

func (m *Service) GetPkgName() string {
	return m.PkgName
}

func (m *Service) GetSnakeName() string {
	return m.SnakeName
}

func (m *Service) Nil() bool {
	return m == nil
}
