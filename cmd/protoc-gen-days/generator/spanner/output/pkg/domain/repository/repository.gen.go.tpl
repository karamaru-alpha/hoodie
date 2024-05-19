{{ template "autogen_comment" }}
{{ $pkgName := .PkgName }}
package {{ .PkgName }}

import (
	"context"

	"github.com/karamaru-alpha/days/pkg/domain/database"
	"github.com/karamaru-alpha/days/pkg/domain/entity/{{ .PkgName }}"
	"github.com/karamaru-alpha/days/pkg/domain/enum"
)

type {{ .GoName }}Repository interface {
	LoadByPK(ctx context.Context, tx database.ROTx, pk *{{ .PkgName }}.{{ .GoName }}PK) (*{{ .PkgName }}.{{ .GoName }}, error)
	LoadByPKs(ctx context.Context, tx database.ROTx, pks {{ .PkgName }}.{{ .GoName }}PKs) ({{ .PkgName }}.{{ .GoName }}Slice, error)
	SelectByPK(ctx context.Context, tx database.ROTx, pk *{{ .PkgName }}.{{ .GoName }}PK) (*{{ .PkgName }}.{{ .GoName }}, error)
	SelectByPKs(ctx context.Context, tx database.ROTx, pks {{ .PkgName }}.{{ .GoName }}PKs) ({{ .PkgName }}.{{ .GoName }}Slice, error)
	SelectAll(ctx context.Context, tx database.ROTx, limit int, offset int) ({{ .PkgName }}.{{ .GoName }}Slice, error)
	{{ range .Methods -}}
	SelectBy{{ .Name }}(ctx context.Context, tx database.ROTx, {{ .Args }}) ({{ $pkgName }}.{{ .ReturnName }}Slice, error)
	{{ end -}}
	Insert(ctx context.Context, tx database.RWTx, row *{{ .PkgName }}.{{ .GoName }}) error
	BulkInsert(ctx context.Context, tx database.RWTx, rows {{ .PkgName }}.{{ .GoName }}Slice) error
	Update(ctx context.Context, tx database.RWTx, row *{{ .PkgName }}.{{ .GoName }}) error
	Save(ctx context.Context, tx database.RWTx, row *{{ .PkgName }}.{{ .GoName }}) error
	Delete(ctx context.Context, tx database.RWTx, pk *{{ .PkgName }}.{{ .GoName }}PK) error
	BulkDelete(ctx context.Context, tx database.RWTx, pks {{ .PkgName }}.{{ .GoName }}PKs) error
}
