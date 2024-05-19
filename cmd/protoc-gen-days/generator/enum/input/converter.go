package input

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/core"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/perrors"
)

func ConvertMessageFromProto(file *protogen.File, _ core.FlagKindSet) (*Enum, error) {
	if len(file.Enums) != 1 {
		return nil, perrors.New("this file should define only one enum").SetValues(map[string]any{
			"file": file.Desc.FullName(),
		})
	}

	enum := file.Enums[0]
	ret := &Enum{
		GoName:    core.ToGolangPascalCase(string(enum.Desc.Name())),
		SnakeName: core.ToSnakeCase(string(enum.Desc.Name())),
		Comment:   core.CommentReplacer.Replace(enum.Comments.Leading.String()),
		Values:    make([]*Value, 0, len(enum.Values)),
	}

	for _, value := range enum.Values {
		ret.Values = append(ret.Values, &Value{
			RawName:   string(value.Desc.Name()),
			GoName:    core.ToGolangPascalCase(string(value.Desc.Name())),
			SnakeName: core.ToSnakeCase(string(value.Desc.Name())),
			CamelName: core.ToGolangCamelCase(string(value.Desc.Name())),
			Comment:   core.CommentReplacer.Replace(value.Comments.Leading.String()),
			Number:    int32(value.Desc.Number()),
		})
	}
	return ret, nil
}
