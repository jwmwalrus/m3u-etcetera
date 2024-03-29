// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: api/m3uetcpb/root.proto

package m3uetcpb

import (
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

type StatusResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Healthy bool `protobuf:"varint,1,opt,name=healthy,proto3" json:"healthy,omitempty"`
}

func (x *StatusResponse) Reset() {
	*x = StatusResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_root_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusResponse) ProtoMessage() {}

func (x *StatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_m3uetcpb_root_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusResponse.ProtoReflect.Descriptor instead.
func (*StatusResponse) Descriptor() ([]byte, []int) {
	return file_api_m3uetcpb_root_proto_rawDescGZIP(), []int{0}
}

func (x *StatusResponse) GetHealthy() bool {
	if x != nil {
		return x.Healthy
	}
	return false
}

type OffRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Force bool `protobuf:"varint,1,opt,name=force,proto3" json:"force,omitempty"`
}

func (x *OffRequest) Reset() {
	*x = OffRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_root_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OffRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OffRequest) ProtoMessage() {}

func (x *OffRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_m3uetcpb_root_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OffRequest.ProtoReflect.Descriptor instead.
func (*OffRequest) Descriptor() ([]byte, []int) {
	return file_api_m3uetcpb_root_proto_rawDescGZIP(), []int{1}
}

func (x *OffRequest) GetForce() bool {
	if x != nil {
		return x.Force
	}
	return false
}

type OffResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GoingOff bool   `protobuf:"varint,1,opt,name=going_off,json=goingOff,proto3" json:"going_off,omitempty"`
	Reason   string `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
}

func (x *OffResponse) Reset() {
	*x = OffResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_root_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OffResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OffResponse) ProtoMessage() {}

func (x *OffResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_m3uetcpb_root_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OffResponse.ProtoReflect.Descriptor instead.
func (*OffResponse) Descriptor() ([]byte, []int) {
	return file_api_m3uetcpb_root_proto_rawDescGZIP(), []int{2}
}

func (x *OffResponse) GetGoingOff() bool {
	if x != nil {
		return x.GoingOff
	}
	return false
}

func (x *OffResponse) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

var File_api_m3uetcpb_root_proto protoreflect.FileDescriptor

var file_api_m3uetcpb_root_proto_rawDesc = []byte{
	0x0a, 0x17, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2f, 0x72,
	0x6f, 0x6f, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x6d, 0x33, 0x75, 0x65, 0x74,
	0x63, 0x70, 0x62, 0x1a, 0x18, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70,
	0x62, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2a, 0x0a,
	0x0e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x07, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x79, 0x22, 0x22, 0x0a, 0x0a, 0x4f, 0x66, 0x66,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x6f, 0x72, 0x63, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x22, 0x42, 0x0a,
	0x0b, 0x4f, 0x66, 0x66, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x67, 0x6f, 0x69, 0x6e, 0x67, 0x5f, 0x6f, 0x66, 0x66, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x08, 0x67, 0x6f, 0x69, 0x6e, 0x67, 0x4f, 0x66, 0x66, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x61,
	0x73, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f,
	0x6e, 0x32, 0x72, 0x0a, 0x07, 0x52, 0x6f, 0x6f, 0x74, 0x53, 0x76, 0x63, 0x12, 0x33, 0x0a, 0x06,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0f, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70,
	0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x18, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63,
	0x70, 0x62, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x32, 0x0a, 0x03, 0x4f, 0x66, 0x66, 0x12, 0x14, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74,
	0x63, 0x70, 0x62, 0x2e, 0x4f, 0x66, 0x66, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15,
	0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x4f, 0x66, 0x66, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x6d, 0x33, 0x75, 0x65, 0x74,
	0x63, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_m3uetcpb_root_proto_rawDescOnce sync.Once
	file_api_m3uetcpb_root_proto_rawDescData = file_api_m3uetcpb_root_proto_rawDesc
)

func file_api_m3uetcpb_root_proto_rawDescGZIP() []byte {
	file_api_m3uetcpb_root_proto_rawDescOnce.Do(func() {
		file_api_m3uetcpb_root_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_m3uetcpb_root_proto_rawDescData)
	})
	return file_api_m3uetcpb_root_proto_rawDescData
}

var file_api_m3uetcpb_root_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_api_m3uetcpb_root_proto_goTypes = []interface{}{
	(*StatusResponse)(nil), // 0: m3uetcpb.StatusResponse
	(*OffRequest)(nil),     // 1: m3uetcpb.OffRequest
	(*OffResponse)(nil),    // 2: m3uetcpb.OffResponse
	(*Empty)(nil),          // 3: m3uetcpb.Empty
}
var file_api_m3uetcpb_root_proto_depIdxs = []int32{
	3, // 0: m3uetcpb.RootSvc.Status:input_type -> m3uetcpb.Empty
	1, // 1: m3uetcpb.RootSvc.Off:input_type -> m3uetcpb.OffRequest
	0, // 2: m3uetcpb.RootSvc.Status:output_type -> m3uetcpb.StatusResponse
	2, // 3: m3uetcpb.RootSvc.Off:output_type -> m3uetcpb.OffResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_m3uetcpb_root_proto_init() }
func file_api_m3uetcpb_root_proto_init() {
	if File_api_m3uetcpb_root_proto != nil {
		return
	}
	file_api_m3uetcpb_empty_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_m3uetcpb_root_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatusResponse); i {
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
		file_api_m3uetcpb_root_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OffRequest); i {
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
		file_api_m3uetcpb_root_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OffResponse); i {
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
			RawDescriptor: file_api_m3uetcpb_root_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_m3uetcpb_root_proto_goTypes,
		DependencyIndexes: file_api_m3uetcpb_root_proto_depIdxs,
		MessageInfos:      file_api_m3uetcpb_root_proto_msgTypes,
	}.Build()
	File_api_m3uetcpb_root_proto = out.File
	file_api_m3uetcpb_root_proto_rawDesc = nil
	file_api_m3uetcpb_root_proto_goTypes = nil
	file_api_m3uetcpb_root_proto_depIdxs = nil
}
