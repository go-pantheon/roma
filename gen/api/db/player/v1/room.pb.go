// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: player/v1/room.proto

package dbv1

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

type UserRoomProto struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" bson:"id"`                                // @gotags: bson:"id" Room ID
	IsCreator     bool                   `protobuf:"varint,2,opt,name=is_creator,json=isCreator,proto3" json:"is_creator,omitempty" bson:"is_creator"` // @gotags: bson:"is_creator" Is Creator of the Room
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserRoomProto) Reset() {
	*x = UserRoomProto{}
	mi := &file_player_v1_room_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserRoomProto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserRoomProto) ProtoMessage() {}

func (x *UserRoomProto) ProtoReflect() protoreflect.Message {
	mi := &file_player_v1_room_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserRoomProto.ProtoReflect.Descriptor instead.
func (*UserRoomProto) Descriptor() ([]byte, []int) {
	return file_player_v1_room_proto_rawDescGZIP(), []int{0}
}

func (x *UserRoomProto) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UserRoomProto) GetIsCreator() bool {
	if x != nil {
		return x.IsCreator
	}
	return false
}

var File_player_v1_room_proto protoreflect.FileDescriptor

var file_player_v1_room_proto_rawDesc = string([]byte{
	0x0a, 0x14, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x6f, 0x6f, 0x6d,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2e, 0x76,
	0x31, 0x22, 0x3e, 0x0a, 0x0d, 0x55, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6f, 0x6d, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x73, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x69, 0x73, 0x43, 0x72, 0x65, 0x61, 0x74, 0x6f,
	0x72, 0x42, 0x17, 0x5a, 0x15, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x62, 0x2f, 0x70, 0x6c, 0x61, 0x79,
	0x65, 0x72, 0x2f, 0x76, 0x31, 0x3b, 0x64, 0x62, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
})

var (
	file_player_v1_room_proto_rawDescOnce sync.Once
	file_player_v1_room_proto_rawDescData []byte
)

func file_player_v1_room_proto_rawDescGZIP() []byte {
	file_player_v1_room_proto_rawDescOnce.Do(func() {
		file_player_v1_room_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_player_v1_room_proto_rawDesc), len(file_player_v1_room_proto_rawDesc)))
	})
	return file_player_v1_room_proto_rawDescData
}

var file_player_v1_room_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_player_v1_room_proto_goTypes = []any{
	(*UserRoomProto)(nil), // 0: player.v1.UserRoomProto
}
var file_player_v1_room_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_player_v1_room_proto_init() }
func file_player_v1_room_proto_init() {
	if File_player_v1_room_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_player_v1_room_proto_rawDesc), len(file_player_v1_room_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_player_v1_room_proto_goTypes,
		DependencyIndexes: file_player_v1_room_proto_depIdxs,
		MessageInfos:      file_player_v1_room_proto_msgTypes,
	}.Build()
	File_player_v1_room_proto = out.File
	file_player_v1_room_proto_goTypes = nil
	file_player_v1_room_proto_depIdxs = nil
}
