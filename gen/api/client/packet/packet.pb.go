// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: packet/packet.proto

package clipkt

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

// TCP packet structure definition
// For public network access
// The complete message format: 4(len, bigEndian) + encrypt(byte[](Marshal(Packet))). The client and server send messages in this format.
// After the handshake protocol, all protocols use AES encryption and decryption
// The message index number is incremented by 1 each time, and the message index number is unique within the same module
// mod + seq + obj forms the unique ID of data
type Packet struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          []byte                 `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`                                   // Serialized bytes of the cs/sc protocol in the message. If the corresponding protocol specifies that the data needs to be cached on the client, use data_version to compare the data version number
	DataVersion   uint64                 `protobuf:"varint,2,opt,name=data_version,json=dataVersion,proto3" json:"data_version,omitempty"` // Data version number, default is 0. When greater than 0 and greater than the version of the same message (unique ID is the same), the data field is valid
	Obj           int64                  `protobuf:"varint,3,opt,name=obj,proto3" json:"obj,omitempty"`                                    // Module object ID, according to the business agreement to pass the corresponding object ID
	Mod           int32                  `protobuf:"varint,4,opt,name=mod,proto3" json:"mod,omitempty"`                                    // Module ID, globally unique
	Seq           int32                  `protobuf:"varint,5,opt,name=seq,proto3" json:"seq,omitempty"`                                    // Message ID within the module, unique within the module
	Ver           int32                  `protobuf:"varint,6,opt,name=ver,proto3" json:"ver,omitempty"`                                    // Version
	Index         int32                  `protobuf:"varint,7,opt,name=index,proto3" json:"index,omitempty"`                                // Message index number, increment
	Compress      bool                   `protobuf:"varint,8,opt,name=compress,proto3" json:"compress,omitempty"`                          // Whether the data in the body is compressed. The default compression method is zlib
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Packet) Reset() {
	*x = Packet{}
	mi := &file_packet_packet_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Packet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Packet) ProtoMessage() {}

func (x *Packet) ProtoReflect() protoreflect.Message {
	mi := &file_packet_packet_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Packet.ProtoReflect.Descriptor instead.
func (*Packet) Descriptor() ([]byte, []int) {
	return file_packet_packet_proto_rawDescGZIP(), []int{0}
}

func (x *Packet) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Packet) GetDataVersion() uint64 {
	if x != nil {
		return x.DataVersion
	}
	return 0
}

func (x *Packet) GetObj() int64 {
	if x != nil {
		return x.Obj
	}
	return 0
}

func (x *Packet) GetMod() int32 {
	if x != nil {
		return x.Mod
	}
	return 0
}

func (x *Packet) GetSeq() int32 {
	if x != nil {
		return x.Seq
	}
	return 0
}

func (x *Packet) GetVer() int32 {
	if x != nil {
		return x.Ver
	}
	return 0
}

func (x *Packet) GetIndex() int32 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *Packet) GetCompress() bool {
	if x != nil {
		return x.Compress
	}
	return false
}

var File_packet_packet_proto protoreflect.FileDescriptor

var file_packet_packet_proto_rawDesc = string([]byte{
	0x0a, 0x13, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x2f, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x22, 0xb9, 0x01,
	0x0a, 0x06, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x21, 0x0a, 0x0c,
	0x64, 0x61, 0x74, 0x61, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x0b, 0x64, 0x61, 0x74, 0x61, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x10, 0x0a, 0x03, 0x6f, 0x62, 0x6a, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x6f, 0x62,
	0x6a, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x6f, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03,
	0x6d, 0x6f, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x65, 0x71, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x03, 0x73, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x76, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x03, 0x76, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1a, 0x0a,
	0x08, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x08, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x42, 0x1a, 0x5a, 0x18, 0x61, 0x70, 0x69,
	0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x3b, 0x63,
	0x6c, 0x69, 0x70, 0x6b, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_packet_packet_proto_rawDescOnce sync.Once
	file_packet_packet_proto_rawDescData []byte
)

func file_packet_packet_proto_rawDescGZIP() []byte {
	file_packet_packet_proto_rawDescOnce.Do(func() {
		file_packet_packet_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_packet_packet_proto_rawDesc), len(file_packet_packet_proto_rawDesc)))
	})
	return file_packet_packet_proto_rawDescData
}

var file_packet_packet_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_packet_packet_proto_goTypes = []any{
	(*Packet)(nil), // 0: packet.Packet
}
var file_packet_packet_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_packet_packet_proto_init() }
func file_packet_packet_proto_init() {
	if File_packet_packet_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_packet_packet_proto_rawDesc), len(file_packet_packet_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_packet_packet_proto_goTypes,
		DependencyIndexes: file_packet_packet_proto_depIdxs,
		MessageInfos:      file_packet_packet_proto_msgTypes,
	}.Build()
	File_packet_packet_proto = out.File
	file_packet_packet_proto_goTypes = nil
	file_packet_packet_proto_depIdxs = nil
}
