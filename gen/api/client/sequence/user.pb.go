// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: sequence/user.proto

package cliseq

import (
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

// User module operation number
type UserSeq int32

const (
	UserSeq_UserUnspecified UserSeq = 0
	// Login
	UserSeq_Login UserSeq = 1
	// @push latest user data. The client receives the data and updates its own data to avoid data inconsistency between the client and the server when GM modifies data or the server restarts
	UserSeq_PushSyncUser UserSeq = 2
	// Update name
	UserSeq_UpdateName UserSeq = 3
	// Set gender
	UserSeq_SetGender UserSeq = 4
)

// Enum value maps for UserSeq.
var (
	UserSeq_name = map[int32]string{
		0: "UserUnspecified",
		1: "Login",
		2: "PushSyncUser",
		3: "UpdateName",
		4: "SetGender",
	}
	UserSeq_value = map[string]int32{
		"UserUnspecified": 0,
		"Login":           1,
		"PushSyncUser":    2,
		"UpdateName":      3,
		"SetGender":       4,
	}
)

func (x UserSeq) Enum() *UserSeq {
	p := new(UserSeq)
	*p = x
	return p
}

func (x UserSeq) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UserSeq) Descriptor() protoreflect.EnumDescriptor {
	return file_sequence_user_proto_enumTypes[0].Descriptor()
}

func (UserSeq) Type() protoreflect.EnumType {
	return &file_sequence_user_proto_enumTypes[0]
}

func (x UserSeq) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UserSeq.Descriptor instead.
func (UserSeq) EnumDescriptor() ([]byte, []int) {
	return file_sequence_user_proto_rawDescGZIP(), []int{0}
}

var File_sequence_user_proto protoreflect.FileDescriptor

var file_sequence_user_proto_rawDesc = string([]byte{
	0x0a, 0x13, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x2a,
	0x5a, 0x0a, 0x07, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x71, 0x12, 0x13, 0x0a, 0x0f, 0x55, 0x73,
	0x65, 0x72, 0x55, 0x6e, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x65, 0x64, 0x10, 0x00, 0x12,
	0x09, 0x0a, 0x05, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x10, 0x01, 0x12, 0x10, 0x0a, 0x0c, 0x50, 0x75,
	0x73, 0x68, 0x53, 0x79, 0x6e, 0x63, 0x55, 0x73, 0x65, 0x72, 0x10, 0x02, 0x12, 0x0e, 0x0a, 0x0a,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x10, 0x03, 0x12, 0x0d, 0x0a, 0x09,
	0x53, 0x65, 0x74, 0x47, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x10, 0x04, 0x42, 0x1c, 0x5a, 0x1a, 0x61,
	0x70, 0x69, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e,
	0x63, 0x65, 0x3b, 0x63, 0x6c, 0x69, 0x73, 0x65, 0x71, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
})

var (
	file_sequence_user_proto_rawDescOnce sync.Once
	file_sequence_user_proto_rawDescData []byte
)

func file_sequence_user_proto_rawDescGZIP() []byte {
	file_sequence_user_proto_rawDescOnce.Do(func() {
		file_sequence_user_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_sequence_user_proto_rawDesc), len(file_sequence_user_proto_rawDesc)))
	})
	return file_sequence_user_proto_rawDescData
}

var file_sequence_user_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_sequence_user_proto_goTypes = []any{
	(UserSeq)(0), // 0: sequence.UserSeq
}
var file_sequence_user_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_sequence_user_proto_init() }
func file_sequence_user_proto_init() {
	if File_sequence_user_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_sequence_user_proto_rawDesc), len(file_sequence_user_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_sequence_user_proto_goTypes,
		DependencyIndexes: file_sequence_user_proto_depIdxs,
		EnumInfos:         file_sequence_user_proto_enumTypes,
	}.Build()
	File_sequence_user_proto = out.File
	file_sequence_user_proto_goTypes = nil
	file_sequence_user_proto_depIdxs = nil
}
