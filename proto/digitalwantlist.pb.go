// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0-devel
// 	protoc        (unknown)
// source: digitalwantlist.proto

package proto

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type TrackMaster_State int32

const (
	TrackMaster_UNKNOWN     TrackMaster_State = 0
	TrackMaster_TRACKING    TrackMaster_State = 1
	TrackMaster_COVERED     TrackMaster_State = 2
	TrackMaster_INELLIGIBLE TrackMaster_State = 3
)

// Enum value maps for TrackMaster_State.
var (
	TrackMaster_State_name = map[int32]string{
		0: "UNKNOWN",
		1: "TRACKING",
		2: "COVERED",
		3: "INELLIGIBLE",
	}
	TrackMaster_State_value = map[string]int32{
		"UNKNOWN":     0,
		"TRACKING":    1,
		"COVERED":     2,
		"INELLIGIBLE": 3,
	}
)

func (x TrackMaster_State) Enum() *TrackMaster_State {
	p := new(TrackMaster_State)
	*p = x
	return p
}

func (x TrackMaster_State) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TrackMaster_State) Descriptor() protoreflect.EnumDescriptor {
	return file_digitalwantlist_proto_enumTypes[0].Descriptor()
}

func (TrackMaster_State) Type() protoreflect.EnumType {
	return &file_digitalwantlist_proto_enumTypes[0]
}

func (x TrackMaster_State) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TrackMaster_State.Descriptor instead.
func (TrackMaster_State) EnumDescriptor() ([]byte, []int) {
	return file_digitalwantlist_proto_rawDescGZIP(), []int{1, 0}
}

type Config struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Masters   []*TrackMaster `protobuf:"bytes,1,rep,name=masters,proto3" json:"masters,omitempty"`
	Purchased []int32        `protobuf:"varint,2,rep,packed,name=purchased,proto3" json:"purchased,omitempty"`
}

func (x *Config) Reset() {
	*x = Config{}
	if protoimpl.UnsafeEnabled {
		mi := &file_digitalwantlist_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_digitalwantlist_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Config.ProtoReflect.Descriptor instead.
func (*Config) Descriptor() ([]byte, []int) {
	return file_digitalwantlist_proto_rawDescGZIP(), []int{0}
}

func (x *Config) GetMasters() []*TrackMaster {
	if x != nil {
		return x.Masters
	}
	return nil
}

func (x *Config) GetPurchased() []int32 {
	if x != nil {
		return x.Purchased
	}
	return nil
}

type TrackMaster struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MasterId   int32             `protobuf:"varint,1,opt,name=master_id,json=masterId,proto3" json:"master_id,omitempty"`
	DigitalIds []int32           `protobuf:"varint,2,rep,packed,name=digital_ids,json=digitalIds,proto3" json:"digital_ids,omitempty"`
	State      TrackMaster_State `protobuf:"varint,3,opt,name=state,proto3,enum=digitalwantlist.TrackMaster_State" json:"state,omitempty"`
}

func (x *TrackMaster) Reset() {
	*x = TrackMaster{}
	if protoimpl.UnsafeEnabled {
		mi := &file_digitalwantlist_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TrackMaster) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TrackMaster) ProtoMessage() {}

func (x *TrackMaster) ProtoReflect() protoreflect.Message {
	mi := &file_digitalwantlist_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TrackMaster.ProtoReflect.Descriptor instead.
func (*TrackMaster) Descriptor() ([]byte, []int) {
	return file_digitalwantlist_proto_rawDescGZIP(), []int{1}
}

func (x *TrackMaster) GetMasterId() int32 {
	if x != nil {
		return x.MasterId
	}
	return 0
}

func (x *TrackMaster) GetDigitalIds() []int32 {
	if x != nil {
		return x.DigitalIds
	}
	return nil
}

func (x *TrackMaster) GetState() TrackMaster_State {
	if x != nil {
		return x.State
	}
	return TrackMaster_UNKNOWN
}

var File_digitalwantlist_proto protoreflect.FileDescriptor

var file_digitalwantlist_proto_rawDesc = []byte{
	0x0a, 0x15, 0x64, 0x69, 0x67, 0x69, 0x74, 0x61, 0x6c, 0x77, 0x61, 0x6e, 0x74, 0x6c, 0x69, 0x73,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x64, 0x69, 0x67, 0x69, 0x74, 0x61, 0x6c,
	0x77, 0x61, 0x6e, 0x74, 0x6c, 0x69, 0x73, 0x74, 0x22, 0x5e, 0x0a, 0x06, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x12, 0x36, 0x0a, 0x07, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x64, 0x69, 0x67, 0x69, 0x74, 0x61, 0x6c, 0x77, 0x61, 0x6e,
	0x74, 0x6c, 0x69, 0x73, 0x74, 0x2e, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x4d, 0x61, 0x73, 0x74, 0x65,
	0x72, 0x52, 0x07, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x75,
	0x72, 0x63, 0x68, 0x61, 0x73, 0x65, 0x64, 0x18, 0x02, 0x20, 0x03, 0x28, 0x05, 0x52, 0x09, 0x70,
	0x75, 0x72, 0x63, 0x68, 0x61, 0x73, 0x65, 0x64, 0x22, 0xc7, 0x01, 0x0a, 0x0b, 0x54, 0x72, 0x61,
	0x63, 0x6b, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x61, 0x73, 0x74,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x6d, 0x61, 0x73,
	0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x69, 0x67, 0x69, 0x74, 0x61, 0x6c,
	0x5f, 0x69, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x05, 0x52, 0x0a, 0x64, 0x69, 0x67, 0x69,
	0x74, 0x61, 0x6c, 0x49, 0x64, 0x73, 0x12, 0x38, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x22, 0x2e, 0x64, 0x69, 0x67, 0x69, 0x74, 0x61, 0x6c, 0x77,
	0x61, 0x6e, 0x74, 0x6c, 0x69, 0x73, 0x74, 0x2e, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x4d, 0x61, 0x73,
	0x74, 0x65, 0x72, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x22, 0x40, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b,
	0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x54, 0x52, 0x41, 0x43, 0x4b, 0x49,
	0x4e, 0x47, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x4f, 0x56, 0x45, 0x52, 0x45, 0x44, 0x10,
	0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x49, 0x4e, 0x45, 0x4c, 0x4c, 0x49, 0x47, 0x49, 0x42, 0x4c, 0x45,
	0x10, 0x03, 0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x62, 0x72, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x2f, 0x64, 0x69,
	0x67, 0x69, 0x74, 0x61, 0x6c, 0x77, 0x61, 0x6e, 0x74, 0x6c, 0x69, 0x73, 0x74, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_digitalwantlist_proto_rawDescOnce sync.Once
	file_digitalwantlist_proto_rawDescData = file_digitalwantlist_proto_rawDesc
)

func file_digitalwantlist_proto_rawDescGZIP() []byte {
	file_digitalwantlist_proto_rawDescOnce.Do(func() {
		file_digitalwantlist_proto_rawDescData = protoimpl.X.CompressGZIP(file_digitalwantlist_proto_rawDescData)
	})
	return file_digitalwantlist_proto_rawDescData
}

var file_digitalwantlist_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_digitalwantlist_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_digitalwantlist_proto_goTypes = []interface{}{
	(TrackMaster_State)(0), // 0: digitalwantlist.TrackMaster.State
	(*Config)(nil),         // 1: digitalwantlist.Config
	(*TrackMaster)(nil),    // 2: digitalwantlist.TrackMaster
}
var file_digitalwantlist_proto_depIdxs = []int32{
	2, // 0: digitalwantlist.Config.masters:type_name -> digitalwantlist.TrackMaster
	0, // 1: digitalwantlist.TrackMaster.state:type_name -> digitalwantlist.TrackMaster.State
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_digitalwantlist_proto_init() }
func file_digitalwantlist_proto_init() {
	if File_digitalwantlist_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_digitalwantlist_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Config); i {
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
		file_digitalwantlist_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TrackMaster); i {
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
			RawDescriptor: file_digitalwantlist_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_digitalwantlist_proto_goTypes,
		DependencyIndexes: file_digitalwantlist_proto_depIdxs,
		EnumInfos:         file_digitalwantlist_proto_enumTypes,
		MessageInfos:      file_digitalwantlist_proto_msgTypes,
	}.Build()
	File_digitalwantlist_proto = out.File
	file_digitalwantlist_proto_rawDesc = nil
	file_digitalwantlist_proto_goTypes = nil
	file_digitalwantlist_proto_depIdxs = nil
}