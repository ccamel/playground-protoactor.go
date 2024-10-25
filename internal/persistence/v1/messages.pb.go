// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: persistence/v1/messages.proto

package persistencev1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ConsiderSnapshot struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload *anypb.Any `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *ConsiderSnapshot) Reset() {
	*x = ConsiderSnapshot{}
	if protoimpl.UnsafeEnabled {
		mi := &file_persistence_v1_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConsiderSnapshot) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsiderSnapshot) ProtoMessage() {}

func (x *ConsiderSnapshot) ProtoReflect() protoreflect.Message {
	mi := &file_persistence_v1_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsiderSnapshot.ProtoReflect.Descriptor instead.
func (*ConsiderSnapshot) Descriptor() ([]byte, []int) {
	return file_persistence_v1_messages_proto_rawDescGZIP(), []int{0}
}

func (x *ConsiderSnapshot) GetPayload() *anypb.Any {
	if x != nil {
		return x.Payload
	}
	return nil
}

var File_persistence_v1_messages_proto protoreflect.FileDescriptor

var file_persistence_v1_messages_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x70, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x63, 0x65, 0x2f, 0x76, 0x31,
	0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0e, 0x70, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x1a,
	0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x42, 0x0a, 0x10, 0x43, 0x6f,
	0x6e, 0x73, 0x69, 0x64, 0x65, 0x72, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x12, 0x2e,
	0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0xce,
	0x01, 0x0a, 0x12, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x6e,
	0x63, 0x65, 0x2e, 0x76, 0x31, 0x42, 0x0d, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x50, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x63, 0x63, 0x61, 0x6d, 0x65, 0x6c, 0x2f, 0x70, 0x6c, 0x61, 0x79, 0x67, 0x72,
	0x6f, 0x75, 0x6e, 0x64, 0x2d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x2e,
	0x67, 0x6f, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x65, 0x72, 0x73,
	0x69, 0x73, 0x74, 0x65, 0x6e, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x70, 0x65, 0x72, 0x73, 0x69,
	0x73, 0x74, 0x65, 0x6e, 0x63, 0x65, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x50, 0x58, 0x58, 0xaa, 0x02,
	0x0e, 0x50, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x56, 0x31, 0xca,
	0x02, 0x0e, 0x50, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x63, 0x65, 0x5c, 0x56, 0x31,
	0xe2, 0x02, 0x1a, 0x50, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x63, 0x65, 0x5c, 0x56,
	0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0f,
	0x50, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x63, 0x65, 0x3a, 0x3a, 0x56, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_persistence_v1_messages_proto_rawDescOnce sync.Once
	file_persistence_v1_messages_proto_rawDescData = file_persistence_v1_messages_proto_rawDesc
)

func file_persistence_v1_messages_proto_rawDescGZIP() []byte {
	file_persistence_v1_messages_proto_rawDescOnce.Do(func() {
		file_persistence_v1_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_persistence_v1_messages_proto_rawDescData)
	})
	return file_persistence_v1_messages_proto_rawDescData
}

var file_persistence_v1_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_persistence_v1_messages_proto_goTypes = []interface{}{
	(*ConsiderSnapshot)(nil), // 0: persistence.v1.ConsiderSnapshot
	(*anypb.Any)(nil),        // 1: google.protobuf.Any
}
var file_persistence_v1_messages_proto_depIdxs = []int32{
	1, // 0: persistence.v1.ConsiderSnapshot.payload:type_name -> google.protobuf.Any
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_persistence_v1_messages_proto_init() }
func file_persistence_v1_messages_proto_init() {
	if File_persistence_v1_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_persistence_v1_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConsiderSnapshot); i {
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
			RawDescriptor: file_persistence_v1_messages_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_persistence_v1_messages_proto_goTypes,
		DependencyIndexes: file_persistence_v1_messages_proto_depIdxs,
		MessageInfos:      file_persistence_v1_messages_proto_msgTypes,
	}.Build()
	File_persistence_v1_messages_proto = out.File
	file_persistence_v1_messages_proto_rawDesc = nil
	file_persistence_v1_messages_proto_goTypes = nil
	file_persistence_v1_messages_proto_depIdxs = nil
}
