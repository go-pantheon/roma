// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: sequence/equip.proto

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

// Equip module operation number
type EquipSeq int32

const (
	EquipSeq_EquipUnspecified EquipSeq = 0
	// @push equipment updated
	EquipSeq_PushEquipUpdated EquipSeq = 1
	// Wear equipment
	EquipSeq_EquipWear EquipSeq = 2
	// Take off equipment
	EquipSeq_EquipTakeOff EquipSeq = 3
	// Upgrade equipment
	EquipSeq_EquipUpgrade EquipSeq = 4
)

// Enum value maps for EquipSeq.
var (
	EquipSeq_name = map[int32]string{
		0: "EquipUnspecified",
		1: "PushEquipUpdated",
		2: "EquipWear",
		3: "EquipTakeOff",
		4: "EquipUpgrade",
	}
	EquipSeq_value = map[string]int32{
		"EquipUnspecified": 0,
		"PushEquipUpdated": 1,
		"EquipWear":        2,
		"EquipTakeOff":     3,
		"EquipUpgrade":     4,
	}
)

func (x EquipSeq) Enum() *EquipSeq {
	p := new(EquipSeq)
	*p = x
	return p
}

func (x EquipSeq) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EquipSeq) Descriptor() protoreflect.EnumDescriptor {
	return file_sequence_equip_proto_enumTypes[0].Descriptor()
}

func (EquipSeq) Type() protoreflect.EnumType {
	return &file_sequence_equip_proto_enumTypes[0]
}

func (x EquipSeq) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EquipSeq.Descriptor instead.
func (EquipSeq) EnumDescriptor() ([]byte, []int) {
	return file_sequence_equip_proto_rawDescGZIP(), []int{0}
}

var File_sequence_equip_proto protoreflect.FileDescriptor

var file_sequence_equip_proto_rawDesc = string([]byte{
	0x0a, 0x14, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x2f, 0x65, 0x71, 0x75, 0x69, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65,
	0x2a, 0x69, 0x0a, 0x08, 0x45, 0x71, 0x75, 0x69, 0x70, 0x53, 0x65, 0x71, 0x12, 0x14, 0x0a, 0x10,
	0x45, 0x71, 0x75, 0x69, 0x70, 0x55, 0x6e, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x65, 0x64,
	0x10, 0x00, 0x12, 0x14, 0x0a, 0x10, 0x50, 0x75, 0x73, 0x68, 0x45, 0x71, 0x75, 0x69, 0x70, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x45, 0x71, 0x75, 0x69,
	0x70, 0x57, 0x65, 0x61, 0x72, 0x10, 0x02, 0x12, 0x10, 0x0a, 0x0c, 0x45, 0x71, 0x75, 0x69, 0x70,
	0x54, 0x61, 0x6b, 0x65, 0x4f, 0x66, 0x66, 0x10, 0x03, 0x12, 0x10, 0x0a, 0x0c, 0x45, 0x71, 0x75,
	0x69, 0x70, 0x55, 0x70, 0x67, 0x72, 0x61, 0x64, 0x65, 0x10, 0x04, 0x42, 0x1c, 0x5a, 0x1a, 0x61,
	0x70, 0x69, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e,
	0x63, 0x65, 0x3b, 0x63, 0x6c, 0x69, 0x73, 0x65, 0x71, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
})

var (
	file_sequence_equip_proto_rawDescOnce sync.Once
	file_sequence_equip_proto_rawDescData []byte
)

func file_sequence_equip_proto_rawDescGZIP() []byte {
	file_sequence_equip_proto_rawDescOnce.Do(func() {
		file_sequence_equip_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_sequence_equip_proto_rawDesc), len(file_sequence_equip_proto_rawDesc)))
	})
	return file_sequence_equip_proto_rawDescData
}

var file_sequence_equip_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_sequence_equip_proto_goTypes = []any{
	(EquipSeq)(0), // 0: sequence.EquipSeq
}
var file_sequence_equip_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_sequence_equip_proto_init() }
func file_sequence_equip_proto_init() {
	if File_sequence_equip_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_sequence_equip_proto_rawDesc), len(file_sequence_equip_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_sequence_equip_proto_goTypes,
		DependencyIndexes: file_sequence_equip_proto_depIdxs,
		EnumInfos:         file_sequence_equip_proto_enumTypes,
	}.Build()
	File_sequence_equip_proto = out.File
	file_sequence_equip_proto_goTypes = nil
	file_sequence_equip_proto_depIdxs = nil
}
