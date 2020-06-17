// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0
// 	protoc        v3.12.3
// source: book_entities.proto

package booklend

import (
	proto "github.com/golang/protobuf/proto"
	duration "github.com/golang/protobuf/ptypes/duration"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type BookEntity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id               string               `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Title            string               `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Isbn             string               `protobuf:"bytes,3,opt,name=isbn,proto3" json:"isbn,omitempty"`
	Borrower         string               `protobuf:"bytes,4,opt,name=borrower,proto3" json:"borrower,omitempty"`
	Date             *timestamp.Timestamp `protobuf:"bytes,5,opt,name=date,proto3" json:"date,omitempty"`
	ExpectedDuration *duration.Duration   `protobuf:"bytes,6,opt,name=expectedDuration,proto3" json:"expectedDuration,omitempty"`
}

func (x *BookEntity) Reset() {
	*x = BookEntity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_book_entities_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BookEntity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BookEntity) ProtoMessage() {}

func (x *BookEntity) ProtoReflect() protoreflect.Message {
	mi := &file_book_entities_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BookEntity.ProtoReflect.Descriptor instead.
func (*BookEntity) Descriptor() ([]byte, []int) {
	return file_book_entities_proto_rawDescGZIP(), []int{0}
}

func (x *BookEntity) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *BookEntity) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *BookEntity) GetIsbn() string {
	if x != nil {
		return x.Isbn
	}
	return ""
}

func (x *BookEntity) GetBorrower() string {
	if x != nil {
		return x.Borrower
	}
	return ""
}

func (x *BookEntity) GetDate() *timestamp.Timestamp {
	if x != nil {
		return x.Date
	}
	return nil
}

func (x *BookEntity) GetExpectedDuration() *duration.Duration {
	if x != nil {
		return x.ExpectedDuration
	}
	return nil
}

var File_book_entities_proto protoreflect.FileDescriptor

var file_book_entities_proto_rawDesc = []byte{
	0x0a, 0x13, 0x62, 0x6f, 0x6f, 0x6b, 0x5f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x62, 0x6f, 0x6f, 0x6b, 0x6c, 0x65, 0x6e, 0x64, 0x1a,
	0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xd9, 0x01, 0x0a, 0x0a, 0x42, 0x6f, 0x6f, 0x6b, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x69, 0x73, 0x62, 0x6e, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x69, 0x73, 0x62, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x6f, 0x72,
	0x72, 0x6f, 0x77, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x62, 0x6f, 0x72,
	0x72, 0x6f, 0x77, 0x65, 0x72, 0x12, 0x2e, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x65, 0x12, 0x45, 0x0a, 0x10, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65,
	0x64, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x10, 0x65, 0x78, 0x70, 0x65,
	0x63, 0x74, 0x65, 0x64, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_book_entities_proto_rawDescOnce sync.Once
	file_book_entities_proto_rawDescData = file_book_entities_proto_rawDesc
)

func file_book_entities_proto_rawDescGZIP() []byte {
	file_book_entities_proto_rawDescOnce.Do(func() {
		file_book_entities_proto_rawDescData = protoimpl.X.CompressGZIP(file_book_entities_proto_rawDescData)
	})
	return file_book_entities_proto_rawDescData
}

var file_book_entities_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_book_entities_proto_goTypes = []interface{}{
	(*BookEntity)(nil),          // 0: booklend.BookEntity
	(*timestamp.Timestamp)(nil), // 1: google.protobuf.Timestamp
	(*duration.Duration)(nil),   // 2: google.protobuf.Duration
}
var file_book_entities_proto_depIdxs = []int32{
	1, // 0: booklend.BookEntity.date:type_name -> google.protobuf.Timestamp
	2, // 1: booklend.BookEntity.expectedDuration:type_name -> google.protobuf.Duration
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_book_entities_proto_init() }
func file_book_entities_proto_init() {
	if File_book_entities_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_book_entities_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BookEntity); i {
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
			RawDescriptor: file_book_entities_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_book_entities_proto_goTypes,
		DependencyIndexes: file_book_entities_proto_depIdxs,
		MessageInfos:      file_book_entities_proto_msgTypes,
	}.Build()
	File_book_entities_proto = out.File
	file_book_entities_proto_rawDesc = nil
	file_book_entities_proto_goTypes = nil
	file_book_entities_proto_depIdxs = nil
}
