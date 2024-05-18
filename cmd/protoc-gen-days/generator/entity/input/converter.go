package input

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/core"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/perrors"
	options "github.com/karamaru-alpha/days/pkg/pb/options/entity"
)

func ConvertMessageFromProto(file *protogen.File, flagKindSet core.FlagKindSet) (*Message, error) {
	if len(file.Messages) != 1 {
		return nil, perrors.New("The proto file should define only one message").SetValues(map[string]any{
			"file": file.Desc.FullName(),
		})
	}
	message := file.Messages[0]
	result := &Message{
		FileDirName: core.GoPackageNameToFileDirName(file.Proto.GetOptions().GetGoPackage()),
		PkgName:     core.ToPkgName(file.Proto.GetOptions().GetGoPackage()),
		SnakeName:   core.ToSnakeCase(string(message.Desc.Name())),
		GoName:      core.ToGolangPascalCase(string(message.Desc.Name())),
		Comment:     core.CommentReplacer.Replace(message.Comments.Leading.String()),
		Fields:      make([]*Field, 0, len(message.Fields)+1),
	}

	for _, field := range message.Fields {
		var typeName string
		var isEnum bool
		var isList bool
		switch field.Desc.Kind() {
		case protoreflect.FloatKind:
			typeName = GoTypeFloat32
		case protoreflect.BoolKind:
			typeName = GoTypeBool
		case protoreflect.Int32Kind:
			typeName = GoTypeInt32
		case protoreflect.Int64Kind:
			typeName = GoTypeInt64
		case protoreflect.StringKind:
			typeName = GoTypeString
		case protoreflect.BytesKind:
			typeName = GoTypeBytes
		case protoreflect.EnumKind:
			typeName = "enum." + string(field.Desc.Enum().Name())
			isEnum = true
		case protoreflect.MessageKind, protoreflect.DoubleKind, protoreflect.Fixed32Kind, protoreflect.Fixed64Kind,
			protoreflect.GroupKind, protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind,
			protoreflect.Sint32Kind, protoreflect.Sint64Kind, protoreflect.Uint32Kind, protoreflect.Uint64Kind:
			return nil, perrors.New("The kind is not supported").SetValues(map[string]any{
				"kind":    field.Desc.Kind().String(),
				"rawName": message.Desc.FullName(),
			})
		default:
			return nil, perrors.New("The kind is not supported").SetValues(map[string]any{
				"kind":    field.Desc.Kind().String(),
				"rawName": message.Desc.FullName(),
			})
		}

		fieldOption, ok := proto.GetExtension(field.Desc.Options(), options.E_Field).(*options.FieldOption)
		if !ok {
			return nil, perrors.New("fail to assert field option")
		}
		inputField := &Field{
			GoName:    core.ToGolangPascalCase(field.Desc.TextName()),
			CamelName: core.ToGolangCamelCase(field.Desc.TextName()),
			Comment:   core.CommentReplacer.Replace(field.Comments.Leading.String()),
			Type:      typeName,
			IsList:    isList,
			IsEnum:    isEnum,
			PK:        fieldOption.GetSchema().GetPk(),
		}
		result.Fields = append(result.Fields, inputField)
	}

	return result, nil
}
