package input

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/core"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/perrors"
)

func ConvertMessageFromProto(file *protogen.File, _ core.FlagKindSet) (*Service, error) {
	if len(file.Services) != 1 {
		return nil, perrors.New("this file should define only one service").SetValues(map[string]any{
			"file": file.Desc.FullName(),
		})
	}
	if len(file.Services) == 0 {
		return nil, nil
	}

	service := file.Services[0]
	s := &Service{
		PkgName:   core.ToPkgName(file.Proto.GetOptions().GetGoPackage()),
		GoPkgName: file.Proto.GetOptions().GetGoPackage(),
		SnakeName: core.ToSnakeCase(string(service.Desc.Name())),
		CamelName: core.ToGolangCamelCase(string(service.Desc.Name())),
		GoName:    core.ToGolangPascalCase(string(service.Desc.Name())),
		Methods:   make([]*Method, 0, len(service.Methods)),
	}

	for _, method := range service.Methods {
		messageOption, ok := method.Desc.Options().(*descriptorpb.MethodOptions)
		if !ok {
			return nil, perrors.New("fail to assert message option")
		}
		var level string
		switch messageOption.GetIdempotencyLevel() {
		case descriptorpb.MethodOptions_IDEMPOTENCY_UNKNOWN:
			level = "Unknown"
		case descriptorpb.MethodOptions_NO_SIDE_EFFECTS:
			level = "NoSideEffects"
		case descriptorpb.MethodOptions_IDEMPOTENT:
			level = "Idempotent"
		}
		s.Methods = append(s.Methods, &Method{
			GoName:           service.GoName + method.GoName,
			Comment:          core.CommentReplacer.Replace(method.Comments.Leading.String()),
			IdempotencyLevel: level,
		})
	}

	return s, nil
}
