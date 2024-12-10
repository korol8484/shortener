// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v3.12.4
// source: proto/service/contract/request.proto

package contract

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RequestShort struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data string `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *RequestShort) Reset() {
	*x = RequestShort{}
	mi := &file_proto_service_contract_request_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RequestShort) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestShort) ProtoMessage() {}

func (x *RequestShort) ProtoReflect() protoreflect.Message {
	mi := &file_proto_service_contract_request_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestShort.ProtoReflect.Descriptor instead.
func (*RequestShort) Descriptor() ([]byte, []int) {
	return file_proto_service_contract_request_proto_rawDescGZIP(), []int{0}
}

func (x *RequestShort) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

type RequestFindByAlias struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Alias string `protobuf:"bytes,1,opt,name=alias,proto3" json:"alias,omitempty"`
}

func (x *RequestFindByAlias) Reset() {
	*x = RequestFindByAlias{}
	mi := &file_proto_service_contract_request_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RequestFindByAlias) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestFindByAlias) ProtoMessage() {}

func (x *RequestFindByAlias) ProtoReflect() protoreflect.Message {
	mi := &file_proto_service_contract_request_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestFindByAlias.ProtoReflect.Descriptor instead.
func (*RequestFindByAlias) Descriptor() ([]byte, []int) {
	return file_proto_service_contract_request_proto_rawDescGZIP(), []int{1}
}

func (x *RequestFindByAlias) GetAlias() string {
	if x != nil {
		return x.Alias
	}
	return ""
}

type RequestBatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Batch []*Batch `protobuf:"bytes,1,rep,name=batch,proto3" json:"batch,omitempty"`
}

func (x *RequestBatch) Reset() {
	*x = RequestBatch{}
	mi := &file_proto_service_contract_request_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RequestBatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestBatch) ProtoMessage() {}

func (x *RequestBatch) ProtoReflect() protoreflect.Message {
	mi := &file_proto_service_contract_request_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestBatch.ProtoReflect.Descriptor instead.
func (*RequestBatch) Descriptor() ([]byte, []int) {
	return file_proto_service_contract_request_proto_rawDescGZIP(), []int{2}
}

func (x *RequestBatch) GetBatch() []*Batch {
	if x != nil {
		return x.Batch
	}
	return nil
}

type RequestBatchDelete struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Alias []string `protobuf:"bytes,1,rep,name=alias,proto3" json:"alias,omitempty"`
}

func (x *RequestBatchDelete) Reset() {
	*x = RequestBatchDelete{}
	mi := &file_proto_service_contract_request_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RequestBatchDelete) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestBatchDelete) ProtoMessage() {}

func (x *RequestBatchDelete) ProtoReflect() protoreflect.Message {
	mi := &file_proto_service_contract_request_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestBatchDelete.ProtoReflect.Descriptor instead.
func (*RequestBatchDelete) Descriptor() ([]byte, []int) {
	return file_proto_service_contract_request_proto_rawDescGZIP(), []int{3}
}

func (x *RequestBatchDelete) GetAlias() []string {
	if x != nil {
		return x.Alias
	}
	return nil
}

var File_proto_service_contract_request_proto protoreflect.FileDescriptor

var file_proto_service_contract_request_proto_rawDesc = []byte{
	0x0a, 0x24, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x2f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x1a, 0x23, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74,
	0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x22, 0x0a,
	0x0c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x22, 0x2a, 0x0a, 0x12, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x46, 0x69, 0x6e, 0x64,
	0x42, 0x79, 0x41, 0x6c, 0x69, 0x61, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x22, 0x3d, 0x0a,
	0x0c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12, 0x2d, 0x0a,
	0x05, 0x62, 0x61, 0x74, 0x63, 0x68, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x2e,
	0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x05, 0x62, 0x61, 0x74, 0x63, 0x68, 0x22, 0x2a, 0x0a, 0x12,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x61, 0x74, 0x63, 0x68, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x42, 0x43, 0x5a, 0x41, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x6f, 0x72, 0x6f, 0x6c, 0x38, 0x34, 0x38, 0x34,
	0x2f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_service_contract_request_proto_rawDescOnce sync.Once
	file_proto_service_contract_request_proto_rawDescData = file_proto_service_contract_request_proto_rawDesc
)

func file_proto_service_contract_request_proto_rawDescGZIP() []byte {
	file_proto_service_contract_request_proto_rawDescOnce.Do(func() {
		file_proto_service_contract_request_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_service_contract_request_proto_rawDescData)
	})
	return file_proto_service_contract_request_proto_rawDescData
}

var file_proto_service_contract_request_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_service_contract_request_proto_goTypes = []any{
	(*RequestShort)(nil),       // 0: service.contract.RequestShort
	(*RequestFindByAlias)(nil), // 1: service.contract.RequestFindByAlias
	(*RequestBatch)(nil),       // 2: service.contract.RequestBatch
	(*RequestBatchDelete)(nil), // 3: service.contract.RequestBatchDelete
	(*Batch)(nil),              // 4: service.contract.Batch
}
var file_proto_service_contract_request_proto_depIdxs = []int32{
	4, // 0: service.contract.RequestBatch.batch:type_name -> service.contract.Batch
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_service_contract_request_proto_init() }
func file_proto_service_contract_request_proto_init() {
	if File_proto_service_contract_request_proto != nil {
		return
	}
	file_proto_service_contract_entity_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_service_contract_request_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_service_contract_request_proto_goTypes,
		DependencyIndexes: file_proto_service_contract_request_proto_depIdxs,
		MessageInfos:      file_proto_service_contract_request_proto_msgTypes,
	}.Build()
	File_proto_service_contract_request_proto = out.File
	file_proto_service_contract_request_proto_rawDesc = nil
	file_proto_service_contract_request_proto_goTypes = nil
	file_proto_service_contract_request_proto_depIdxs = nil
}