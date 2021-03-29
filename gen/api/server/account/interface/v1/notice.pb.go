// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: account/interface/v1/notice.proto

package interfacev1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type NoticeListResponse_Code int32

const (
	NoticeListResponse_CODE_ERR_UNSPECIFIED NoticeListResponse_Code = 0 // Please try again later
	NoticeListResponse_CODE_SUCCEEDED       NoticeListResponse_Code = 1 // Succeeded
)

// Enum value maps for NoticeListResponse_Code.
var (
	NoticeListResponse_Code_name = map[int32]string{
		0: "CODE_ERR_UNSPECIFIED",
		1: "CODE_SUCCEEDED",
	}
	NoticeListResponse_Code_value = map[string]int32{
		"CODE_ERR_UNSPECIFIED": 0,
		"CODE_SUCCEEDED":       1,
	}
)

func (x NoticeListResponse_Code) Enum() *NoticeListResponse_Code {
	p := new(NoticeListResponse_Code)
	*p = x
	return p
}

func (x NoticeListResponse_Code) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (NoticeListResponse_Code) Descriptor() protoreflect.EnumDescriptor {
	return file_account_interface_v1_notice_proto_enumTypes[0].Descriptor()
}

func (NoticeListResponse_Code) Type() protoreflect.EnumType {
	return &file_account_interface_v1_notice_proto_enumTypes[0]
}

func (x NoticeListResponse_Code) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use NoticeListResponse_Code.Descriptor instead.
func (NoticeListResponse_Code) EnumDescriptor() ([]byte, []int) {
	return file_account_interface_v1_notice_proto_rawDescGZIP(), []int{1, 0}
}

type NoticeListRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NoticeListRequest) Reset() {
	*x = NoticeListRequest{}
	mi := &file_account_interface_v1_notice_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NoticeListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoticeListRequest) ProtoMessage() {}

func (x *NoticeListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_account_interface_v1_notice_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NoticeListRequest.ProtoReflect.Descriptor instead.
func (*NoticeListRequest) Descriptor() ([]byte, []int) {
	return file_account_interface_v1_notice_proto_rawDescGZIP(), []int{0}
}

type NoticeListResponse struct {
	state         protoimpl.MessageState  `protogen:"open.v1"`
	Code          NoticeListResponse_Code `protobuf:"varint,1,opt,name=code,proto3,enum=account.interface.v1.NoticeListResponse_Code" json:"code,omitempty"` // Response code
	List          []*Notice               `protobuf:"bytes,2,rep,name=list,proto3" json:"list,omitempty"`                                                    // Notice list, up to 10 items
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NoticeListResponse) Reset() {
	*x = NoticeListResponse{}
	mi := &file_account_interface_v1_notice_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NoticeListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoticeListResponse) ProtoMessage() {}

func (x *NoticeListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_account_interface_v1_notice_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NoticeListResponse.ProtoReflect.Descriptor instead.
func (*NoticeListResponse) Descriptor() ([]byte, []int) {
	return file_account_interface_v1_notice_proto_rawDescGZIP(), []int{1}
}

func (x *NoticeListResponse) GetCode() NoticeListResponse_Code {
	if x != nil {
		return x.Code
	}
	return NoticeListResponse_CODE_ERR_UNSPECIFIED
}

func (x *NoticeListResponse) GetList() []*Notice {
	if x != nil {
		return x.List
	}
	return nil
}

// Notice
type Notice struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Title         string                 `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`     // Title
	Content       string                 `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"` // Content
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Notice) Reset() {
	*x = Notice{}
	mi := &file_account_interface_v1_notice_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Notice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Notice) ProtoMessage() {}

func (x *Notice) ProtoReflect() protoreflect.Message {
	mi := &file_account_interface_v1_notice_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Notice.ProtoReflect.Descriptor instead.
func (*Notice) Descriptor() ([]byte, []int) {
	return file_account_interface_v1_notice_proto_rawDescGZIP(), []int{2}
}

func (x *Notice) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Notice) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

var File_account_interface_v1_notice_proto protoreflect.FileDescriptor

var file_account_interface_v1_notice_proto_rawDesc = string([]byte{
	0x0a, 0x21, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66,
	0x61, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x14, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x13, 0x0a, 0x11, 0x4e, 0x6f, 0x74, 0x69, 0x63,
	0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0xbf, 0x01, 0x0a,
	0x12, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x2d, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x66, 0x61, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x4c,
	0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x43, 0x6f, 0x64, 0x65,
	0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x30, 0x0a, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x6f, 0x74, 0x69,
	0x63, 0x65, 0x52, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x22, 0x34, 0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x18, 0x0a, 0x14, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x45, 0x52, 0x52, 0x5f, 0x55, 0x4e, 0x53,
	0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e, 0x43, 0x4f,
	0x44, 0x45, 0x5f, 0x53, 0x55, 0x43, 0x43, 0x45, 0x45, 0x44, 0x45, 0x44, 0x10, 0x01, 0x22, 0x38,
	0x0a, 0x06, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x32, 0x8b, 0x01, 0x0a, 0x0f, 0x4e, 0x6f, 0x74,
	0x69, 0x63, 0x65, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x12, 0x78, 0x0a, 0x0a,
	0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x27, 0x2e, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x63, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x63,
	0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x17, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x11, 0x12, 0x0f, 0x2f, 0x76, 0x31, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x63,
	0x65, 0x2f, 0x6c, 0x69, 0x73, 0x74, 0x42, 0x2d, 0x5a, 0x2b, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2f, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66,
	0x61, 0x63, 0x65, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_account_interface_v1_notice_proto_rawDescOnce sync.Once
	file_account_interface_v1_notice_proto_rawDescData []byte
)

func file_account_interface_v1_notice_proto_rawDescGZIP() []byte {
	file_account_interface_v1_notice_proto_rawDescOnce.Do(func() {
		file_account_interface_v1_notice_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_account_interface_v1_notice_proto_rawDesc), len(file_account_interface_v1_notice_proto_rawDesc)))
	})
	return file_account_interface_v1_notice_proto_rawDescData
}

var file_account_interface_v1_notice_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_account_interface_v1_notice_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_account_interface_v1_notice_proto_goTypes = []any{
	(NoticeListResponse_Code)(0), // 0: account.interface.v1.NoticeListResponse.Code
	(*NoticeListRequest)(nil),    // 1: account.interface.v1.NoticeListRequest
	(*NoticeListResponse)(nil),   // 2: account.interface.v1.NoticeListResponse
	(*Notice)(nil),               // 3: account.interface.v1.Notice
}
var file_account_interface_v1_notice_proto_depIdxs = []int32{
	0, // 0: account.interface.v1.NoticeListResponse.code:type_name -> account.interface.v1.NoticeListResponse.Code
	3, // 1: account.interface.v1.NoticeListResponse.list:type_name -> account.interface.v1.Notice
	1, // 2: account.interface.v1.NoticeInterface.NoticeList:input_type -> account.interface.v1.NoticeListRequest
	2, // 3: account.interface.v1.NoticeInterface.NoticeList:output_type -> account.interface.v1.NoticeListResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_account_interface_v1_notice_proto_init() }
func file_account_interface_v1_notice_proto_init() {
	if File_account_interface_v1_notice_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_account_interface_v1_notice_proto_rawDesc), len(file_account_interface_v1_notice_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_account_interface_v1_notice_proto_goTypes,
		DependencyIndexes: file_account_interface_v1_notice_proto_depIdxs,
		EnumInfos:         file_account_interface_v1_notice_proto_enumTypes,
		MessageInfos:      file_account_interface_v1_notice_proto_msgTypes,
	}.Build()
	File_account_interface_v1_notice_proto = out.File
	file_account_interface_v1_notice_proto_goTypes = nil
	file_account_interface_v1_notice_proto_depIdxs = nil
}
