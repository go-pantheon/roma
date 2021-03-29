// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: player/admin/gamedata/v1/gamedata_error.proto

package adminv1

import (
	_ "github.com/go-kratos/kratos/errors"
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

type GamedataAdminErrorReason int32

const (
	GamedataAdminErrorReason_GAMEDATA_ADMIN_ERROR_REASON_UNSPECIFIED GamedataAdminErrorReason = 0
	GamedataAdminErrorReason_GAMEDATA_ADMIN_ERROR_REASON_SERVER      GamedataAdminErrorReason = 1
	GamedataAdminErrorReason_GAMEDATA_ADMIN_ERROR_REASON_ID          GamedataAdminErrorReason = 2
)

// Enum value maps for GamedataAdminErrorReason.
var (
	GamedataAdminErrorReason_name = map[int32]string{
		0: "GAMEDATA_ADMIN_ERROR_REASON_UNSPECIFIED",
		1: "GAMEDATA_ADMIN_ERROR_REASON_SERVER",
		2: "GAMEDATA_ADMIN_ERROR_REASON_ID",
	}
	GamedataAdminErrorReason_value = map[string]int32{
		"GAMEDATA_ADMIN_ERROR_REASON_UNSPECIFIED": 0,
		"GAMEDATA_ADMIN_ERROR_REASON_SERVER":      1,
		"GAMEDATA_ADMIN_ERROR_REASON_ID":          2,
	}
)

func (x GamedataAdminErrorReason) Enum() *GamedataAdminErrorReason {
	p := new(GamedataAdminErrorReason)
	*p = x
	return p
}

func (x GamedataAdminErrorReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GamedataAdminErrorReason) Descriptor() protoreflect.EnumDescriptor {
	return file_player_admin_gamedata_v1_gamedata_error_proto_enumTypes[0].Descriptor()
}

func (GamedataAdminErrorReason) Type() protoreflect.EnumType {
	return &file_player_admin_gamedata_v1_gamedata_error_proto_enumTypes[0]
}

func (x GamedataAdminErrorReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GamedataAdminErrorReason.Descriptor instead.
func (GamedataAdminErrorReason) EnumDescriptor() ([]byte, []int) {
	return file_player_admin_gamedata_v1_gamedata_error_proto_rawDescGZIP(), []int{0}
}

var File_player_admin_gamedata_v1_gamedata_error_proto protoreflect.FileDescriptor

var file_player_admin_gamedata_v1_gamedata_error_proto_rawDesc = string([]byte{
	0x0a, 0x2d, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x67,
	0x61, 0x6d, 0x65, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x76, 0x31, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x64,
	0x61, 0x74, 0x61, 0x5f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x18, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x67, 0x61,
	0x6d, 0x65, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x76, 0x31, 0x1a, 0x13, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x73, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0xab,
	0x01, 0x0a, 0x18, 0x47, 0x61, 0x6d, 0x65, 0x64, 0x61, 0x74, 0x61, 0x41, 0x64, 0x6d, 0x69, 0x6e,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x31, 0x0a, 0x27, 0x47,
	0x41, 0x4d, 0x45, 0x44, 0x41, 0x54, 0x41, 0x5f, 0x41, 0x44, 0x4d, 0x49, 0x4e, 0x5f, 0x45, 0x52,
	0x52, 0x4f, 0x52, 0x5f, 0x52, 0x45, 0x41, 0x53, 0x4f, 0x4e, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45,
	0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x1a, 0x04, 0xa8, 0x45, 0xf4, 0x03, 0x12, 0x2c,
	0x0a, 0x22, 0x47, 0x41, 0x4d, 0x45, 0x44, 0x41, 0x54, 0x41, 0x5f, 0x41, 0x44, 0x4d, 0x49, 0x4e,
	0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x52, 0x45, 0x41, 0x53, 0x4f, 0x4e, 0x5f, 0x53, 0x45,
	0x52, 0x56, 0x45, 0x52, 0x10, 0x01, 0x1a, 0x04, 0xa8, 0x45, 0xf4, 0x03, 0x12, 0x28, 0x0a, 0x1e,
	0x47, 0x41, 0x4d, 0x45, 0x44, 0x41, 0x54, 0x41, 0x5f, 0x41, 0x44, 0x4d, 0x49, 0x4e, 0x5f, 0x45,
	0x52, 0x52, 0x4f, 0x52, 0x5f, 0x52, 0x45, 0x41, 0x53, 0x4f, 0x4e, 0x5f, 0x49, 0x44, 0x10, 0x02,
	0x1a, 0x04, 0xa8, 0x45, 0x91, 0x03, 0x1a, 0x04, 0xa0, 0x45, 0xf4, 0x03, 0x42, 0x2d, 0x5a, 0x2b,
	0x61, 0x70, 0x69, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x6c, 0x61, 0x79, 0x65,
	0x72, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x64, 0x61, 0x74, 0x61,
	0x2f, 0x76, 0x31, 0x3b, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
})

var (
	file_player_admin_gamedata_v1_gamedata_error_proto_rawDescOnce sync.Once
	file_player_admin_gamedata_v1_gamedata_error_proto_rawDescData []byte
)

func file_player_admin_gamedata_v1_gamedata_error_proto_rawDescGZIP() []byte {
	file_player_admin_gamedata_v1_gamedata_error_proto_rawDescOnce.Do(func() {
		file_player_admin_gamedata_v1_gamedata_error_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_player_admin_gamedata_v1_gamedata_error_proto_rawDesc), len(file_player_admin_gamedata_v1_gamedata_error_proto_rawDesc)))
	})
	return file_player_admin_gamedata_v1_gamedata_error_proto_rawDescData
}

var file_player_admin_gamedata_v1_gamedata_error_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_player_admin_gamedata_v1_gamedata_error_proto_goTypes = []any{
	(GamedataAdminErrorReason)(0), // 0: player.admin.gamedata.v1.GamedataAdminErrorReason
}
var file_player_admin_gamedata_v1_gamedata_error_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_player_admin_gamedata_v1_gamedata_error_proto_init() }
func file_player_admin_gamedata_v1_gamedata_error_proto_init() {
	if File_player_admin_gamedata_v1_gamedata_error_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_player_admin_gamedata_v1_gamedata_error_proto_rawDesc), len(file_player_admin_gamedata_v1_gamedata_error_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_player_admin_gamedata_v1_gamedata_error_proto_goTypes,
		DependencyIndexes: file_player_admin_gamedata_v1_gamedata_error_proto_depIdxs,
		EnumInfos:         file_player_admin_gamedata_v1_gamedata_error_proto_enumTypes,
	}.Build()
	File_player_admin_gamedata_v1_gamedata_error_proto = out.File
	file_player_admin_gamedata_v1_gamedata_error_proto_goTypes = nil
	file_player_admin_gamedata_v1_gamedata_error_proto_depIdxs = nil
}
