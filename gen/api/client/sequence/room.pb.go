// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: sequence/room.proto

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

// Room module operation number
type RoomSeq int32

const (
	RoomSeq_RoomUnspecified RoomSeq = 0
	// Room list
	RoomSeq_RoomList RoomSeq = 1
	// Room detail
	RoomSeq_RoomDetail RoomSeq = 2
	// Create room
	RoomSeq_CreateRoom RoomSeq = 3
	// Invite to join room
	RoomSeq_InviteToJoinRoom RoomSeq = 4
	// Agree to invite to join room
	RoomSeq_AgreeToInviteJoinRoom RoomSeq = 5
	// Request to join room
	RoomSeq_RequestToJoinRoom RoomSeq = 6
	// Approve request to join room
	RoomSeq_ApproveRequestToJoinRoom RoomSeq = 7
	// @push joined room
	RoomSeq_PushJoinedRoom RoomSeq = 8
	// Kick user from room
	RoomSeq_KickUserFromRoom RoomSeq = 9
	// Close room
	RoomSeq_CloseRoom RoomSeq = 10
	// Leave room
	RoomSeq_LeaveRoom RoomSeq = 11
	// @push Removed from room. Leave, be kicked or room closed
	RoomSeq_PushRemovedFromRoom RoomSeq = 12
)

// Enum value maps for RoomSeq.
var (
	RoomSeq_name = map[int32]string{
		0:  "RoomUnspecified",
		1:  "RoomList",
		2:  "RoomDetail",
		3:  "CreateRoom",
		4:  "InviteToJoinRoom",
		5:  "AgreeToInviteJoinRoom",
		6:  "RequestToJoinRoom",
		7:  "ApproveRequestToJoinRoom",
		8:  "PushJoinedRoom",
		9:  "KickUserFromRoom",
		10: "CloseRoom",
		11: "LeaveRoom",
		12: "PushRemovedFromRoom",
	}
	RoomSeq_value = map[string]int32{
		"RoomUnspecified":          0,
		"RoomList":                 1,
		"RoomDetail":               2,
		"CreateRoom":               3,
		"InviteToJoinRoom":         4,
		"AgreeToInviteJoinRoom":    5,
		"RequestToJoinRoom":        6,
		"ApproveRequestToJoinRoom": 7,
		"PushJoinedRoom":           8,
		"KickUserFromRoom":         9,
		"CloseRoom":                10,
		"LeaveRoom":                11,
		"PushRemovedFromRoom":      12,
	}
)

func (x RoomSeq) Enum() *RoomSeq {
	p := new(RoomSeq)
	*p = x
	return p
}

func (x RoomSeq) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RoomSeq) Descriptor() protoreflect.EnumDescriptor {
	return file_sequence_room_proto_enumTypes[0].Descriptor()
}

func (RoomSeq) Type() protoreflect.EnumType {
	return &file_sequence_room_proto_enumTypes[0]
}

func (x RoomSeq) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RoomSeq.Descriptor instead.
func (RoomSeq) EnumDescriptor() ([]byte, []int) {
	return file_sequence_room_proto_rawDescGZIP(), []int{0}
}

var File_sequence_room_proto protoreflect.FileDescriptor

var file_sequence_room_proto_rawDesc = string([]byte{
	0x0a, 0x13, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x2f, 0x72, 0x6f, 0x6f, 0x6d, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x2a,
	0x93, 0x02, 0x0a, 0x07, 0x52, 0x6f, 0x6f, 0x6d, 0x53, 0x65, 0x71, 0x12, 0x13, 0x0a, 0x0f, 0x52,
	0x6f, 0x6f, 0x6d, 0x55, 0x6e, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x65, 0x64, 0x10, 0x00,
	0x12, 0x0c, 0x0a, 0x08, 0x52, 0x6f, 0x6f, 0x6d, 0x4c, 0x69, 0x73, 0x74, 0x10, 0x01, 0x12, 0x0e,
	0x0a, 0x0a, 0x52, 0x6f, 0x6f, 0x6d, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x10, 0x02, 0x12, 0x0e,
	0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x6f, 0x6f, 0x6d, 0x10, 0x03, 0x12, 0x14,
	0x0a, 0x10, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x54, 0x6f, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x6f,
	0x6f, 0x6d, 0x10, 0x04, 0x12, 0x19, 0x0a, 0x15, 0x41, 0x67, 0x72, 0x65, 0x65, 0x54, 0x6f, 0x49,
	0x6e, 0x76, 0x69, 0x74, 0x65, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x6f, 0x6f, 0x6d, 0x10, 0x05, 0x12,
	0x15, 0x0a, 0x11, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x6f, 0x4a, 0x6f, 0x69, 0x6e,
	0x52, 0x6f, 0x6f, 0x6d, 0x10, 0x06, 0x12, 0x1c, 0x0a, 0x18, 0x41, 0x70, 0x70, 0x72, 0x6f, 0x76,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x6f, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x6f,
	0x6f, 0x6d, 0x10, 0x07, 0x12, 0x12, 0x0a, 0x0e, 0x50, 0x75, 0x73, 0x68, 0x4a, 0x6f, 0x69, 0x6e,
	0x65, 0x64, 0x52, 0x6f, 0x6f, 0x6d, 0x10, 0x08, 0x12, 0x14, 0x0a, 0x10, 0x4b, 0x69, 0x63, 0x6b,
	0x55, 0x73, 0x65, 0x72, 0x46, 0x72, 0x6f, 0x6d, 0x52, 0x6f, 0x6f, 0x6d, 0x10, 0x09, 0x12, 0x0d,
	0x0a, 0x09, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x6f, 0x6f, 0x6d, 0x10, 0x0a, 0x12, 0x0d, 0x0a,
	0x09, 0x4c, 0x65, 0x61, 0x76, 0x65, 0x52, 0x6f, 0x6f, 0x6d, 0x10, 0x0b, 0x12, 0x17, 0x0a, 0x13,
	0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x64, 0x46, 0x72, 0x6f, 0x6d, 0x52,
	0x6f, 0x6f, 0x6d, 0x10, 0x0c, 0x42, 0x1c, 0x5a, 0x1a, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x2f, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x3b, 0x63, 0x6c, 0x69,
	0x73, 0x65, 0x71, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_sequence_room_proto_rawDescOnce sync.Once
	file_sequence_room_proto_rawDescData []byte
)

func file_sequence_room_proto_rawDescGZIP() []byte {
	file_sequence_room_proto_rawDescOnce.Do(func() {
		file_sequence_room_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_sequence_room_proto_rawDesc), len(file_sequence_room_proto_rawDesc)))
	})
	return file_sequence_room_proto_rawDescData
}

var file_sequence_room_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_sequence_room_proto_goTypes = []any{
	(RoomSeq)(0), // 0: sequence.RoomSeq
}
var file_sequence_room_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_sequence_room_proto_init() }
func file_sequence_room_proto_init() {
	if File_sequence_room_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_sequence_room_proto_rawDesc), len(file_sequence_room_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_sequence_room_proto_goTypes,
		DependencyIndexes: file_sequence_room_proto_depIdxs,
		EnumInfos:         file_sequence_room_proto_enumTypes,
	}.Build()
	File_sequence_room_proto = out.File
	file_sequence_room_proto_goTypes = nil
	file_sequence_room_proto_depIdxs = nil
}
