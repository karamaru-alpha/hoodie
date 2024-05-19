package input

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/core"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/perrors"
	"github.com/karamaru-alpha/days/pkg/pb/options/entity"
)

const (
	createdTimeSnakeName = "created_time"
	updatedTimeSnakeName = "updated_time"
)

func ConvertMessageFromProto(file *protogen.File, _ core.FlagKindSet) (*Message, error) {
	if len(file.Messages) != 1 {
		return nil, perrors.New("this file should define only one message").SetValues(map[string]any{
			"file": file.Desc.FullName(),
		})
	}

	message := file.Messages[0]
	messageOption, ok := proto.GetExtension(message.Desc.Options(), entity.E_Message).(*entity.MessageOption)
	if !ok {
		return nil, perrors.New("fail to assert message option")
	}

	result := &Message{
		PkgName:   core.ToPkgName(file.Proto.GetOptions().GetGoPackage()),
		SnakeName: core.ToSnakeCase(string(message.Desc.Name())),
		CamelName: core.ToGolangCamelCase(string(message.Desc.Name())),
		GoName:    core.ToGolangPascalCase(string(message.Desc.Name())),
		Comment:   core.CommentReplacer.Replace(message.Comments.Leading.String()),
		Fields:    make([]*Field, 0, len(message.Fields)+2),
		Indexes:   make([]*Index, 0, len(messageOption.GetSchema().GetIndexes())),
	}
	for _, field := range message.Fields {
		fieldOption, ok := proto.GetExtension(field.Desc.Options(), entity.E_Field).(*entity.FieldOption)
		if !ok {
			return nil, perrors.New("fail to assert field option")
		}

		var goTypeName string
		var dbTypeName string
		var setType string
		var isEnum bool
		isList := field.Desc.IsList()
		switch field.Desc.Kind() {
		case protoreflect.Int64Kind:
			goTypeName = GoTypeInt64
			dbTypeName = DBTypeInt
			setType = SetTypeInt64
		case protoreflect.Int32Kind:
			goTypeName = GoTypeInt64
			dbTypeName = DBTypeInt
			setType = SetTypeInt64
		case protoreflect.BoolKind:
			goTypeName = GoTypeBool
			dbTypeName = DBTypeBool
		case protoreflect.StringKind:
			goTypeName = GoTypeString
			dbTypeName = DBTypeString
			setType = SetTypeString
		case protoreflect.BytesKind:
			goTypeName = GoTypeBytes
			isList = true
			dbTypeName = DBTypeBytes
		case protoreflect.EnumKind:
			goTypeName = "enum." + string(field.Desc.Enum().Name())
			dbTypeName = DBTypeInt
			setType = SetTypeInt32
			isEnum = true
		case protoreflect.MessageKind, protoreflect.DoubleKind, protoreflect.FloatKind, protoreflect.Fixed32Kind, protoreflect.Fixed64Kind,
			protoreflect.GroupKind, protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind,
			protoreflect.Sint32Kind, protoreflect.Sint64Kind, protoreflect.Uint32Kind, protoreflect.Uint64Kind:
			return nil, perrors.New("this kind is not supported").SetValues(map[string]any{
				"kind":    field.Desc.Kind().String(),
				"rawName": message.Desc.FullName(),
			})
		default:
			return nil, perrors.New("this kind is not supported").SetValues(map[string]any{
				"kind":    field.Desc.Kind().String(),
				"rawName": message.Desc.FullName(),
			})
		}
		if core.IsTimeField(core.ToSnakeCase(field.Desc.TextName())) {
			goTypeName = GoTypeTime
			dbTypeName = DBTypeTime
		}

		inputField := &Field{
			GoName:    core.ToGolangPascalCase(field.Desc.TextName()),
			SnakeName: core.ToSnakeCase(field.Desc.TextName()),
			CamelName: core.ToGolangCamelCase(field.Desc.TextName()),
			Comment:   core.CommentReplacer.Replace(field.Comments.Leading.String()),
			GoType:    goTypeName,
			DBType:    dbTypeName,
			SetType:   setType,
			IsEnum:    isEnum,
			IsList:    isList,
			PK:        fieldOption.GetSchema().GetPk(),
			Desc:      fieldOption.GetSchema().GetDesc(),
		}
		result.Fields = append(result.Fields, inputField)
	}
	result.Fields = append(result.Fields,
		&Field{
			GoName:    core.ToGolangPascalCase(createdTimeSnakeName),
			SnakeName: createdTimeSnakeName,
			CamelName: core.ToGolangCamelCase(createdTimeSnakeName),
			Comment:   "CreatedTime",
			GoType:    GoTypeTime,
			DBType:    DBTypeTime,
		},
		&Field{
			GoName:    core.ToGolangPascalCase(updatedTimeSnakeName),
			SnakeName: updatedTimeSnakeName,
			CamelName: core.ToGolangCamelCase(updatedTimeSnakeName),
			Comment:   "UpdatedTime",
			GoType:    GoTypeTime,
			DBType:    DBTypeTime,
		},
	)

	for _, index := range messageOption.GetSchema().GetIndexes() {
		keys := make([]*IndexKey, 0, len(index.GetKeys()))
		for _, key := range index.GetKeys() {
			keys = append(keys, &IndexKey{
				GoName: core.ToGolangPascalCase(key.GetColumn()),
				Desc:   key.GetDesc(),
			})
		}

		storing := make([]string, 0, len(index.GetStoring()))
		for _, s := range index.GetStoring() {
			storing = append(storing, core.ToGolangPascalCase(s))
		}

		result.Indexes = append(result.Indexes, &Index{
			Keys:          keys,
			Unique:        index.GetUnique(),
			NullFiltered:  index.GetNullFiltered(),
			PascalStoring: storing,
		})
	}

	if messageOption.GetSchema().GetInterleave() != nil {
		result.Interleave = &Interleave{
			GoName: core.ToGolangPascalCase(messageOption.GetSchema().GetInterleave().GetParent()),
		}
	}
	if messageOption.GetSchema().GetTtl() != nil {
		result.TTL = &TTL{
			TimestampColumnGoName: core.ToGolangPascalCase(messageOption.GetSchema().GetTtl().GetTimestampColumn()),
			Days:                  messageOption.GetSchema().GetTtl().GetDays(),
		}
	}

	return result, nil
}
