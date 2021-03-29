// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: player/v1/hero.proto

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

// Hero List
type UserHeroListProto struct {
	state         protoimpl.MessageState   `protogen:"open.v1"`
	Heroes        map[int64]*UserHeroProto `protobuf:"bytes,1,rep,name=heroes,proto3" json:"heroes,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value" bson:"heroes"` // @gotags: bson:"heroes" All Heroes
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserHeroListProto) Reset() {
	*x = UserHeroListProto{}
	mi := &file_player_v1_hero_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserHeroListProto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserHeroListProto) ProtoMessage() {}

func (x *UserHeroListProto) ProtoReflect() protoreflect.Message {
	mi := &file_player_v1_hero_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserHeroListProto.ProtoReflect.Descriptor instead.
func (*UserHeroListProto) Descriptor() ([]byte, []int) {
	return file_player_v1_hero_proto_rawDescGZIP(), []int{0}
}

func (x *UserHeroListProto) GetHeroes() map[int64]*UserHeroProto {
	if x != nil {
		return x.Heroes
	}
	return nil
}

// Hero
type UserHeroProto struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" bson:"id"`                                                                                    // @gotags: bson:"id" Hero ID
	DataId        int64                  `protobuf:"varint,2,opt,name=data_id,json=dataId,proto3" json:"data_id,omitempty" bson:"data_id"`                                                              // @gotags: bson:"data_id" Hero DataID
	Level         int64                  `protobuf:"varint,3,opt,name=level,proto3" json:"level,omitempty" bson:"level"`                                                                              // @gotags: bson:"level" Level
	Skills        map[int64]int64        `protobuf:"bytes,4,rep,name=skills,proto3" json:"skills,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"varint,2,opt,name=value" bson:"skills"` // @gotags: bson:"skills" Skill ID -> Skill Level
	Equips        []int64                `protobuf:"varint,5,rep,packed,name=equips,proto3" json:"equips,omitempty" bson:"equips"`                                                                     // @gotags: bson:"equips" Wearing Equipments
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserHeroProto) Reset() {
	*x = UserHeroProto{}
	mi := &file_player_v1_hero_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserHeroProto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserHeroProto) ProtoMessage() {}

func (x *UserHeroProto) ProtoReflect() protoreflect.Message {
	mi := &file_player_v1_hero_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserHeroProto.ProtoReflect.Descriptor instead.
func (*UserHeroProto) Descriptor() ([]byte, []int) {
	return file_player_v1_hero_proto_rawDescGZIP(), []int{1}
}

func (x *UserHeroProto) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UserHeroProto) GetDataId() int64 {
	if x != nil {
		return x.DataId
	}
	return 0
}

func (x *UserHeroProto) GetLevel() int64 {
	if x != nil {
		return x.Level
	}
	return 0
}

func (x *UserHeroProto) GetSkills() map[int64]int64 {
	if x != nil {
		return x.Skills
	}
	return nil
}

func (x *UserHeroProto) GetEquips() []int64 {
	if x != nil {
		return x.Equips
	}
	return nil
}

var File_player_v1_hero_proto protoreflect.FileDescriptor

var file_player_v1_hero_proto_rawDesc = string([]byte{
	0x0a, 0x14, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x68, 0x65, 0x72, 0x6f,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2e, 0x76,
	0x31, 0x22, 0xaa, 0x01, 0x0a, 0x11, 0x55, 0x73, 0x65, 0x72, 0x48, 0x65, 0x72, 0x6f, 0x4c, 0x69,
	0x73, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x40, 0x0a, 0x06, 0x68, 0x65, 0x72, 0x6f, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x48, 0x65, 0x72, 0x6f, 0x4c, 0x69, 0x73, 0x74,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x48, 0x65, 0x72, 0x6f, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x06, 0x68, 0x65, 0x72, 0x6f, 0x65, 0x73, 0x1a, 0x53, 0x0a, 0x0b, 0x48, 0x65, 0x72,
	0x6f, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2e, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x70, 0x6c, 0x61, 0x79,
	0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x48, 0x65, 0x72, 0x6f, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xdf,
	0x01, 0x0a, 0x0d, 0x55, 0x73, 0x65, 0x72, 0x48, 0x65, 0x72, 0x6f, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x17, 0x0a, 0x07, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x64, 0x61, 0x74, 0x61, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x65, 0x76,
	0x65, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x12,
	0x3c, 0x0a, 0x06, 0x73, 0x6b, 0x69, 0x6c, 0x6c, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x24, 0x2e, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72,
	0x48, 0x65, 0x72, 0x6f, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x6b, 0x69, 0x6c, 0x6c, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x73, 0x6b, 0x69, 0x6c, 0x6c, 0x73, 0x12, 0x16, 0x0a,
	0x06, 0x65, 0x71, 0x75, 0x69, 0x70, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x03, 0x52, 0x06, 0x65,
	0x71, 0x75, 0x69, 0x70, 0x73, 0x1a, 0x39, 0x0a, 0x0b, 0x53, 0x6b, 0x69, 0x6c, 0x6c, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x42, 0x17, 0x5a, 0x15, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x62, 0x2f, 0x70, 0x6c, 0x61, 0x79, 0x65,
	0x72, 0x2f, 0x76, 0x31, 0x3b, 0x64, 0x62, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
})

var (
	file_player_v1_hero_proto_rawDescOnce sync.Once
	file_player_v1_hero_proto_rawDescData []byte
)

func file_player_v1_hero_proto_rawDescGZIP() []byte {
	file_player_v1_hero_proto_rawDescOnce.Do(func() {
		file_player_v1_hero_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_player_v1_hero_proto_rawDesc), len(file_player_v1_hero_proto_rawDesc)))
	})
	return file_player_v1_hero_proto_rawDescData
}

var file_player_v1_hero_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_player_v1_hero_proto_goTypes = []any{
	(*UserHeroListProto)(nil), // 0: player.v1.UserHeroListProto
	(*UserHeroProto)(nil),     // 1: player.v1.UserHeroProto
	nil,                       // 2: player.v1.UserHeroListProto.HeroesEntry
	nil,                       // 3: player.v1.UserHeroProto.SkillsEntry
}
var file_player_v1_hero_proto_depIdxs = []int32{
	2, // 0: player.v1.UserHeroListProto.heroes:type_name -> player.v1.UserHeroListProto.HeroesEntry
	3, // 1: player.v1.UserHeroProto.skills:type_name -> player.v1.UserHeroProto.SkillsEntry
	1, // 2: player.v1.UserHeroListProto.HeroesEntry.value:type_name -> player.v1.UserHeroProto
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_player_v1_hero_proto_init() }
func file_player_v1_hero_proto_init() {
	if File_player_v1_hero_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_player_v1_hero_proto_rawDesc), len(file_player_v1_hero_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_player_v1_hero_proto_goTypes,
		DependencyIndexes: file_player_v1_hero_proto_depIdxs,
		MessageInfos:      file_player_v1_hero_proto_msgTypes,
	}.Build()
	File_player_v1_hero_proto = out.File
	file_player_v1_hero_proto_goTypes = nil
	file_player_v1_hero_proto_depIdxs = nil
}
