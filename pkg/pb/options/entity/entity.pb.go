// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        (unknown)
// source: options/entity/entity.proto

package entity

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MessageOption struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Schema *MessageOption_Schema `protobuf:"bytes,1,opt,name=schema,proto3" json:"schema,omitempty"`
}

func (x *MessageOption) Reset() {
	*x = MessageOption{}
	if protoimpl.UnsafeEnabled {
		mi := &file_options_entity_entity_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageOption) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageOption) ProtoMessage() {}

func (x *MessageOption) ProtoReflect() protoreflect.Message {
	mi := &file_options_entity_entity_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageOption.ProtoReflect.Descriptor instead.
func (*MessageOption) Descriptor() ([]byte, []int) {
	return file_options_entity_entity_proto_rawDescGZIP(), []int{0}
}

func (x *MessageOption) GetSchema() *MessageOption_Schema {
	if x != nil {
		return x.Schema
	}
	return nil
}

type FieldOption struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Schema *FieldOption_Schema `protobuf:"bytes,1,opt,name=schema,proto3" json:"schema,omitempty"`
}

func (x *FieldOption) Reset() {
	*x = FieldOption{}
	if protoimpl.UnsafeEnabled {
		mi := &file_options_entity_entity_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FieldOption) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FieldOption) ProtoMessage() {}

func (x *FieldOption) ProtoReflect() protoreflect.Message {
	mi := &file_options_entity_entity_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FieldOption.ProtoReflect.Descriptor instead.
func (*FieldOption) Descriptor() ([]byte, []int) {
	return file_options_entity_entity_proto_rawDescGZIP(), []int{1}
}

func (x *FieldOption) GetSchema() *FieldOption_Schema {
	if x != nil {
		return x.Schema
	}
	return nil
}

type MessageOption_Schema struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Indexes    []*MessageOption_Schema_Index    `protobuf:"bytes,1,rep,name=indexes,proto3" json:"indexes,omitempty"`
	Interleave *MessageOption_Schema_Interleave `protobuf:"bytes,2,opt,name=interleave,proto3" json:"interleave,omitempty"`
	Ttl        *MessageOption_Schema_TTL        `protobuf:"bytes,3,opt,name=ttl,proto3" json:"ttl,omitempty"`
}

func (x *MessageOption_Schema) Reset() {
	*x = MessageOption_Schema{}
	if protoimpl.UnsafeEnabled {
		mi := &file_options_entity_entity_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageOption_Schema) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageOption_Schema) ProtoMessage() {}

func (x *MessageOption_Schema) ProtoReflect() protoreflect.Message {
	mi := &file_options_entity_entity_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageOption_Schema.ProtoReflect.Descriptor instead.
func (*MessageOption_Schema) Descriptor() ([]byte, []int) {
	return file_options_entity_entity_proto_rawDescGZIP(), []int{0, 0}
}

func (x *MessageOption_Schema) GetIndexes() []*MessageOption_Schema_Index {
	if x != nil {
		return x.Indexes
	}
	return nil
}

func (x *MessageOption_Schema) GetInterleave() *MessageOption_Schema_Interleave {
	if x != nil {
		return x.Interleave
	}
	return nil
}

func (x *MessageOption_Schema) GetTtl() *MessageOption_Schema_TTL {
	if x != nil {
		return x.Ttl
	}
	return nil
}

type MessageOption_Schema_Index struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys         []*MessageOption_Schema_Index_Key `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
	Unique       bool                              `protobuf:"varint,2,opt,name=unique,proto3" json:"unique,omitempty"`
	NullFiltered bool                              `protobuf:"varint,3,opt,name=null_filtered,json=nullFiltered,proto3" json:"null_filtered,omitempty"`
	Storing      []string                          `protobuf:"bytes,4,rep,name=storing,proto3" json:"storing,omitempty"`
}

func (x *MessageOption_Schema_Index) Reset() {
	*x = MessageOption_Schema_Index{}
	if protoimpl.UnsafeEnabled {
		mi := &file_options_entity_entity_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageOption_Schema_Index) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageOption_Schema_Index) ProtoMessage() {}

func (x *MessageOption_Schema_Index) ProtoReflect() protoreflect.Message {
	mi := &file_options_entity_entity_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageOption_Schema_Index.ProtoReflect.Descriptor instead.
func (*MessageOption_Schema_Index) Descriptor() ([]byte, []int) {
	return file_options_entity_entity_proto_rawDescGZIP(), []int{0, 0, 0}
}

func (x *MessageOption_Schema_Index) GetKeys() []*MessageOption_Schema_Index_Key {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *MessageOption_Schema_Index) GetUnique() bool {
	if x != nil {
		return x.Unique
	}
	return false
}

func (x *MessageOption_Schema_Index) GetNullFiltered() bool {
	if x != nil {
		return x.NullFiltered
	}
	return false
}

func (x *MessageOption_Schema_Index) GetStoring() []string {
	if x != nil {
		return x.Storing
	}
	return nil
}

type MessageOption_Schema_Interleave struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Parent string `protobuf:"bytes,1,opt,name=parent,proto3" json:"parent,omitempty"`
}

func (x *MessageOption_Schema_Interleave) Reset() {
	*x = MessageOption_Schema_Interleave{}
	if protoimpl.UnsafeEnabled {
		mi := &file_options_entity_entity_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageOption_Schema_Interleave) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageOption_Schema_Interleave) ProtoMessage() {}

func (x *MessageOption_Schema_Interleave) ProtoReflect() protoreflect.Message {
	mi := &file_options_entity_entity_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageOption_Schema_Interleave.ProtoReflect.Descriptor instead.
func (*MessageOption_Schema_Interleave) Descriptor() ([]byte, []int) {
	return file_options_entity_entity_proto_rawDescGZIP(), []int{0, 0, 1}
}

func (x *MessageOption_Schema_Interleave) GetParent() string {
	if x != nil {
		return x.Parent
	}
	return ""
}

type MessageOption_Schema_TTL struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TimestampColumn string `protobuf:"bytes,1,opt,name=timestamp_column,json=timestampColumn,proto3" json:"timestamp_column,omitempty"`
	Days            int32  `protobuf:"varint,2,opt,name=days,proto3" json:"days,omitempty"`
}

func (x *MessageOption_Schema_TTL) Reset() {
	*x = MessageOption_Schema_TTL{}
	if protoimpl.UnsafeEnabled {
		mi := &file_options_entity_entity_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageOption_Schema_TTL) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageOption_Schema_TTL) ProtoMessage() {}

func (x *MessageOption_Schema_TTL) ProtoReflect() protoreflect.Message {
	mi := &file_options_entity_entity_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageOption_Schema_TTL.ProtoReflect.Descriptor instead.
func (*MessageOption_Schema_TTL) Descriptor() ([]byte, []int) {
	return file_options_entity_entity_proto_rawDescGZIP(), []int{0, 0, 2}
}

func (x *MessageOption_Schema_TTL) GetTimestampColumn() string {
	if x != nil {
		return x.TimestampColumn
	}
	return ""
}

func (x *MessageOption_Schema_TTL) GetDays() int32 {
	if x != nil {
		return x.Days
	}
	return 0
}

type MessageOption_Schema_Index_Key struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Column string `protobuf:"bytes,1,opt,name=column,proto3" json:"column,omitempty"`
	Desc   bool   `protobuf:"varint,2,opt,name=desc,proto3" json:"desc,omitempty"`
}

func (x *MessageOption_Schema_Index_Key) Reset() {
	*x = MessageOption_Schema_Index_Key{}
	if protoimpl.UnsafeEnabled {
		mi := &file_options_entity_entity_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageOption_Schema_Index_Key) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageOption_Schema_Index_Key) ProtoMessage() {}

func (x *MessageOption_Schema_Index_Key) ProtoReflect() protoreflect.Message {
	mi := &file_options_entity_entity_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageOption_Schema_Index_Key.ProtoReflect.Descriptor instead.
func (*MessageOption_Schema_Index_Key) Descriptor() ([]byte, []int) {
	return file_options_entity_entity_proto_rawDescGZIP(), []int{0, 0, 0, 0}
}

func (x *MessageOption_Schema_Index_Key) GetColumn() string {
	if x != nil {
		return x.Column
	}
	return ""
}

func (x *MessageOption_Schema_Index_Key) GetDesc() bool {
	if x != nil {
		return x.Desc
	}
	return false
}

type FieldOption_Schema struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pk   bool `protobuf:"varint,1,opt,name=pk,proto3" json:"pk,omitempty"`
	Desc bool `protobuf:"varint,2,opt,name=desc,proto3" json:"desc,omitempty"`
}

func (x *FieldOption_Schema) Reset() {
	*x = FieldOption_Schema{}
	if protoimpl.UnsafeEnabled {
		mi := &file_options_entity_entity_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FieldOption_Schema) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FieldOption_Schema) ProtoMessage() {}

func (x *FieldOption_Schema) ProtoReflect() protoreflect.Message {
	mi := &file_options_entity_entity_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FieldOption_Schema.ProtoReflect.Descriptor instead.
func (*FieldOption_Schema) Descriptor() ([]byte, []int) {
	return file_options_entity_entity_proto_rawDescGZIP(), []int{1, 0}
}

func (x *FieldOption_Schema) GetPk() bool {
	if x != nil {
		return x.Pk
	}
	return false
}

func (x *FieldOption_Schema) GetDesc() bool {
	if x != nil {
		return x.Desc
	}
	return false
}

var file_options_entity_entity_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*MessageOption)(nil),
		Field:         50040,
		Name:          "options.entity.message",
		Tag:           "bytes,50040,opt,name=message",
		Filename:      "options/entity/entity.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*FieldOption)(nil),
		Field:         500440,
		Name:          "options.entity.field",
		Tag:           "bytes,500440,opt,name=field",
		Filename:      "options/entity/entity.proto",
	},
}

// Extension fields to descriptorpb.MessageOptions.
var (
	// optional options.entity.MessageOption message = 50040;
	E_Message = &file_options_entity_entity_proto_extTypes[0]
)

// Extension fields to descriptorpb.FieldOptions.
var (
	// optional options.entity.FieldOption field = 500440;
	E_Field = &file_options_entity_entity_proto_extTypes[1]
)

var File_options_entity_entity_proto protoreflect.FileDescriptor

var file_options_entity_entity_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x6f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x1a, 0x20, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xef, 0x04, 0x0a, 0x0d, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x3c, 0x0a, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x24, 0x2e, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x52, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x1a,
	0x9f, 0x04, 0x0a, 0x06, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x44, 0x0a, 0x07, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x6f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x6d,
	0x61, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x52, 0x07, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x73,
	0x12, 0x4f, 0x0a, 0x0a, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6c, 0x65, 0x61, 0x76, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72,
	0x6c, 0x65, 0x61, 0x76, 0x65, 0x52, 0x0a, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6c, 0x65, 0x61, 0x76,
	0x65, 0x12, 0x3a, 0x0a, 0x03, 0x74, 0x74, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28,
	0x2e, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x63,
	0x68, 0x65, 0x6d, 0x61, 0x2e, 0x54, 0x54, 0x4c, 0x52, 0x03, 0x74, 0x74, 0x6c, 0x1a, 0xd5, 0x01,
	0x0a, 0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x42, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e, 0x49, 0x6e, 0x64, 0x65,
	0x78, 0x2e, 0x4b, 0x65, 0x79, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x75,
	0x6e, 0x69, 0x71, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x75, 0x6e, 0x69,
	0x71, 0x75, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x6e, 0x75, 0x6c, 0x6c, 0x5f, 0x66, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x6e, 0x75, 0x6c, 0x6c,
	0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x65, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x74, 0x6f, 0x72,
	0x69, 0x6e, 0x67, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x73, 0x74, 0x6f, 0x72, 0x69,
	0x6e, 0x67, 0x1a, 0x31, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x6f, 0x6c,
	0x75, 0x6d, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x6f, 0x6c, 0x75, 0x6d,
	0x6e, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x65, 0x73, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x04, 0x64, 0x65, 0x73, 0x63, 0x1a, 0x24, 0x0a, 0x0a, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6c, 0x65,
	0x61, 0x76, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x1a, 0x44, 0x0a, 0x03, 0x54,
	0x54, 0x4c, 0x12, 0x29, 0x0a, 0x10, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x5f,
	0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x12, 0x12, 0x0a,
	0x04, 0x64, 0x61, 0x79, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x64, 0x61, 0x79,
	0x73, 0x22, 0x77, 0x0a, 0x0b, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x3a, 0x0a, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x22, 0x2e, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x63,
	0x68, 0x65, 0x6d, 0x61, 0x52, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x1a, 0x2c, 0x0a, 0x06,
	0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x0e, 0x0a, 0x02, 0x70, 0x6b, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x02, 0x70, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x65, 0x73, 0x63, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x64, 0x65, 0x73, 0x63, 0x3a, 0x5a, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xf8, 0x86, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d,
	0x2e, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x3a, 0x52, 0x0a, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x12,
	0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd8,
	0xc5, 0x1e, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x42, 0x36, 0x5a, 0x34, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x61, 0x72, 0x61, 0x6d, 0x61, 0x72,
	0x75, 0x2d, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2f, 0x64, 0x61, 0x79, 0x73, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x70, 0x62, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_options_entity_entity_proto_rawDescOnce sync.Once
	file_options_entity_entity_proto_rawDescData = file_options_entity_entity_proto_rawDesc
)

func file_options_entity_entity_proto_rawDescGZIP() []byte {
	file_options_entity_entity_proto_rawDescOnce.Do(func() {
		file_options_entity_entity_proto_rawDescData = protoimpl.X.CompressGZIP(file_options_entity_entity_proto_rawDescData)
	})
	return file_options_entity_entity_proto_rawDescData
}

var file_options_entity_entity_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_options_entity_entity_proto_goTypes = []interface{}{
	(*MessageOption)(nil),                   // 0: options.entity.MessageOption
	(*FieldOption)(nil),                     // 1: options.entity.FieldOption
	(*MessageOption_Schema)(nil),            // 2: options.entity.MessageOption.Schema
	(*MessageOption_Schema_Index)(nil),      // 3: options.entity.MessageOption.Schema.Index
	(*MessageOption_Schema_Interleave)(nil), // 4: options.entity.MessageOption.Schema.Interleave
	(*MessageOption_Schema_TTL)(nil),        // 5: options.entity.MessageOption.Schema.TTL
	(*MessageOption_Schema_Index_Key)(nil),  // 6: options.entity.MessageOption.Schema.Index.Key
	(*FieldOption_Schema)(nil),              // 7: options.entity.FieldOption.Schema
	(*descriptorpb.MessageOptions)(nil),     // 8: google.protobuf.MessageOptions
	(*descriptorpb.FieldOptions)(nil),       // 9: google.protobuf.FieldOptions
}
var file_options_entity_entity_proto_depIdxs = []int32{
	2,  // 0: options.entity.MessageOption.schema:type_name -> options.entity.MessageOption.Schema
	7,  // 1: options.entity.FieldOption.schema:type_name -> options.entity.FieldOption.Schema
	3,  // 2: options.entity.MessageOption.Schema.indexes:type_name -> options.entity.MessageOption.Schema.Index
	4,  // 3: options.entity.MessageOption.Schema.interleave:type_name -> options.entity.MessageOption.Schema.Interleave
	5,  // 4: options.entity.MessageOption.Schema.ttl:type_name -> options.entity.MessageOption.Schema.TTL
	6,  // 5: options.entity.MessageOption.Schema.Index.keys:type_name -> options.entity.MessageOption.Schema.Index.Key
	8,  // 6: options.entity.message:extendee -> google.protobuf.MessageOptions
	9,  // 7: options.entity.field:extendee -> google.protobuf.FieldOptions
	0,  // 8: options.entity.message:type_name -> options.entity.MessageOption
	1,  // 9: options.entity.field:type_name -> options.entity.FieldOption
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	8,  // [8:10] is the sub-list for extension type_name
	6,  // [6:8] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_options_entity_entity_proto_init() }
func file_options_entity_entity_proto_init() {
	if File_options_entity_entity_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_options_entity_entity_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageOption); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_options_entity_entity_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FieldOption); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_options_entity_entity_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageOption_Schema); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_options_entity_entity_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageOption_Schema_Index); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_options_entity_entity_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageOption_Schema_Interleave); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_options_entity_entity_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageOption_Schema_TTL); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_options_entity_entity_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageOption_Schema_Index_Key); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_options_entity_entity_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FieldOption_Schema); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_options_entity_entity_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 2,
			NumServices:   0,
		},
		GoTypes:           file_options_entity_entity_proto_goTypes,
		DependencyIndexes: file_options_entity_entity_proto_depIdxs,
		MessageInfos:      file_options_entity_entity_proto_msgTypes,
		ExtensionInfos:    file_options_entity_entity_proto_extTypes,
	}.Build()
	File_options_entity_entity_proto = out.File
	file_options_entity_entity_proto_rawDesc = nil
	file_options_entity_entity_proto_goTypes = nil
	file_options_entity_entity_proto_depIdxs = nil
}
