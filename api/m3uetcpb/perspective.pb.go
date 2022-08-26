// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: api/m3uetcpb/perspective.proto

package m3uetcpb

import (
	empty "github.com/golang/protobuf/ptypes/empty"
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

type Perspective int32

const (
	Perspective_MUSIC      Perspective = 0
	Perspective_RADIO      Perspective = 1
	Perspective_PODCASTS   Perspective = 2
	Perspective_AUDIOBOOKS Perspective = 3
)

// Enum value maps for Perspective.
var (
	Perspective_name = map[int32]string{
		0: "MUSIC",
		1: "RADIO",
		2: "PODCASTS",
		3: "AUDIOBOOKS",
	}
	Perspective_value = map[string]int32{
		"MUSIC":      0,
		"RADIO":      1,
		"PODCASTS":   2,
		"AUDIOBOOKS": 3,
	}
)

func (x Perspective) Enum() *Perspective {
	p := new(Perspective)
	*p = x
	return p
}

func (x Perspective) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Perspective) Descriptor() protoreflect.EnumDescriptor {
	return file_api_m3uetcpb_perspective_proto_enumTypes[0].Descriptor()
}

func (Perspective) Type() protoreflect.EnumType {
	return &file_api_m3uetcpb_perspective_proto_enumTypes[0]
}

func (x Perspective) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Perspective.Descriptor instead.
func (Perspective) EnumDescriptor() ([]byte, []int) {
	return file_api_m3uetcpb_perspective_proto_rawDescGZIP(), []int{0}
}

type GetActivePerspectiveResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Perspective Perspective `protobuf:"varint,1,opt,name=perspective,proto3,enum=m3uetcpb.Perspective" json:"perspective,omitempty"`
}

func (x *GetActivePerspectiveResponse) Reset() {
	*x = GetActivePerspectiveResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_perspective_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetActivePerspectiveResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetActivePerspectiveResponse) ProtoMessage() {}

func (x *GetActivePerspectiveResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_m3uetcpb_perspective_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetActivePerspectiveResponse.ProtoReflect.Descriptor instead.
func (*GetActivePerspectiveResponse) Descriptor() ([]byte, []int) {
	return file_api_m3uetcpb_perspective_proto_rawDescGZIP(), []int{0}
}

func (x *GetActivePerspectiveResponse) GetPerspective() Perspective {
	if x != nil {
		return x.Perspective
	}
	return Perspective_MUSIC
}

type SetActivePerspectiveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Perspective Perspective `protobuf:"varint,1,opt,name=perspective,proto3,enum=m3uetcpb.Perspective" json:"perspective,omitempty"`
}

func (x *SetActivePerspectiveRequest) Reset() {
	*x = SetActivePerspectiveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_perspective_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetActivePerspectiveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetActivePerspectiveRequest) ProtoMessage() {}

func (x *SetActivePerspectiveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_m3uetcpb_perspective_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetActivePerspectiveRequest.ProtoReflect.Descriptor instead.
func (*SetActivePerspectiveRequest) Descriptor() ([]byte, []int) {
	return file_api_m3uetcpb_perspective_proto_rawDescGZIP(), []int{1}
}

func (x *SetActivePerspectiveRequest) GetPerspective() Perspective {
	if x != nil {
		return x.Perspective
	}
	return Perspective_MUSIC
}

type SubscribeToPerspectiveResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SubscriptionId    string      `protobuf:"bytes,1,opt,name=subscription_id,json=subscriptionId,proto3" json:"subscription_id,omitempty"`
	ActivePerspective Perspective `protobuf:"varint,2,opt,name=active_perspective,json=activePerspective,proto3,enum=m3uetcpb.Perspective" json:"active_perspective,omitempty"`
}

func (x *SubscribeToPerspectiveResponse) Reset() {
	*x = SubscribeToPerspectiveResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_perspective_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscribeToPerspectiveResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeToPerspectiveResponse) ProtoMessage() {}

func (x *SubscribeToPerspectiveResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_m3uetcpb_perspective_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeToPerspectiveResponse.ProtoReflect.Descriptor instead.
func (*SubscribeToPerspectiveResponse) Descriptor() ([]byte, []int) {
	return file_api_m3uetcpb_perspective_proto_rawDescGZIP(), []int{2}
}

func (x *SubscribeToPerspectiveResponse) GetSubscriptionId() string {
	if x != nil {
		return x.SubscriptionId
	}
	return ""
}

func (x *SubscribeToPerspectiveResponse) GetActivePerspective() Perspective {
	if x != nil {
		return x.ActivePerspective
	}
	return Perspective_MUSIC
}

type UnsubscribeFromPerspectiveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SubscriptionId string `protobuf:"bytes,1,opt,name=subscription_id,json=subscriptionId,proto3" json:"subscription_id,omitempty"`
}

func (x *UnsubscribeFromPerspectiveRequest) Reset() {
	*x = UnsubscribeFromPerspectiveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_perspective_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnsubscribeFromPerspectiveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnsubscribeFromPerspectiveRequest) ProtoMessage() {}

func (x *UnsubscribeFromPerspectiveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_m3uetcpb_perspective_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnsubscribeFromPerspectiveRequest.ProtoReflect.Descriptor instead.
func (*UnsubscribeFromPerspectiveRequest) Descriptor() ([]byte, []int) {
	return file_api_m3uetcpb_perspective_proto_rawDescGZIP(), []int{3}
}

func (x *UnsubscribeFromPerspectiveRequest) GetSubscriptionId() string {
	if x != nil {
		return x.SubscriptionId
	}
	return ""
}

type PerspectiveDigest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Perspective Perspective `protobuf:"varint,1,opt,name=perspective,proto3,enum=m3uetcpb.Perspective" json:"perspective,omitempty"`
	Duration    int64       `protobuf:"varint,2,opt,name=duration,proto3" json:"duration,omitempty"`
}

func (x *PerspectiveDigest) Reset() {
	*x = PerspectiveDigest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_perspective_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PerspectiveDigest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PerspectiveDigest) ProtoMessage() {}

func (x *PerspectiveDigest) ProtoReflect() protoreflect.Message {
	mi := &file_api_m3uetcpb_perspective_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PerspectiveDigest.ProtoReflect.Descriptor instead.
func (*PerspectiveDigest) Descriptor() ([]byte, []int) {
	return file_api_m3uetcpb_perspective_proto_rawDescGZIP(), []int{4}
}

func (x *PerspectiveDigest) GetPerspective() Perspective {
	if x != nil {
		return x.Perspective
	}
	return Perspective_MUSIC
}

func (x *PerspectiveDigest) GetDuration() int64 {
	if x != nil {
		return x.Duration
	}
	return 0
}

var File_api_m3uetcpb_perspective_proto protoreflect.FileDescriptor

var file_api_m3uetcpb_perspective_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2f, 0x70,
	0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x57, 0x0a, 0x1c, 0x47, 0x65, 0x74, 0x41, 0x63,
	0x74, 0x69, 0x76, 0x65, 0x50, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x0b, 0x70, 0x65, 0x72, 0x73, 0x70,
	0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e, 0x6d,
	0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x50, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74,
	0x69, 0x76, 0x65, 0x52, 0x0b, 0x70, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65,
	0x22, 0x56, 0x0a, 0x1b, 0x53, 0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x50, 0x65, 0x72,
	0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x37, 0x0a, 0x0b, 0x70, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e,
	0x50, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x52, 0x0b, 0x70, 0x65, 0x72,
	0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x22, 0x8f, 0x01, 0x0a, 0x1e, 0x53, 0x75, 0x62,
	0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x54, 0x6f, 0x50, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74,
	0x69, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x73,
	0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x49, 0x64, 0x12, 0x44, 0x0a, 0x12, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x70,
	0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x15, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x50, 0x65, 0x72, 0x73,
	0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x52, 0x11, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x50,
	0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x22, 0x4c, 0x0a, 0x21, 0x55, 0x6e,
	0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x50, 0x65, 0x72,
	0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x27, 0x0a, 0x0f, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x68, 0x0a, 0x11, 0x50, 0x65, 0x72, 0x73,
	0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x12, 0x37, 0x0a,
	0x0b, 0x70, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x15, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x50, 0x65,
	0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x52, 0x0b, 0x70, 0x65, 0x72, 0x73, 0x70,
	0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2a, 0x41, 0x0a, 0x0b, 0x50, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76,
	0x65, 0x12, 0x09, 0x0a, 0x05, 0x4d, 0x55, 0x53, 0x49, 0x43, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05,
	0x52, 0x41, 0x44, 0x49, 0x4f, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x50, 0x4f, 0x44, 0x43, 0x41,
	0x53, 0x54, 0x53, 0x10, 0x02, 0x12, 0x0e, 0x0a, 0x0a, 0x41, 0x55, 0x44, 0x49, 0x4f, 0x42, 0x4f,
	0x4f, 0x4b, 0x53, 0x10, 0x03, 0x32, 0x80, 0x03, 0x0a, 0x0e, 0x50, 0x65, 0x72, 0x73, 0x70, 0x65,
	0x63, 0x74, 0x69, 0x76, 0x65, 0x53, 0x76, 0x63, 0x12, 0x56, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x41,
	0x63, 0x74, 0x69, 0x76, 0x65, 0x50, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x26, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74,
	0x63, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x50, 0x65, 0x72,
	0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x55, 0x0a, 0x14, 0x53, 0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x50, 0x65, 0x72,
	0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x12, 0x25, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74,
	0x63, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x50, 0x65, 0x72,
	0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x5c, 0x0a, 0x16, 0x53, 0x75, 0x62, 0x73, 0x63,
	0x72, 0x69, 0x62, 0x65, 0x54, 0x6f, 0x50, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76,
	0x65, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x28, 0x2e, 0x6d, 0x33, 0x75, 0x65,
	0x74, 0x63, 0x70, 0x62, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x54, 0x6f,
	0x50, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x61, 0x0a, 0x1a, 0x55, 0x6e, 0x73, 0x75, 0x62, 0x73, 0x63,
	0x72, 0x69, 0x62, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x50, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74,
	0x69, 0x76, 0x65, 0x12, 0x2b, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x55,
	0x6e, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x50, 0x65,
	0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x6d, 0x33,
	0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_m3uetcpb_perspective_proto_rawDescOnce sync.Once
	file_api_m3uetcpb_perspective_proto_rawDescData = file_api_m3uetcpb_perspective_proto_rawDesc
)

func file_api_m3uetcpb_perspective_proto_rawDescGZIP() []byte {
	file_api_m3uetcpb_perspective_proto_rawDescOnce.Do(func() {
		file_api_m3uetcpb_perspective_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_m3uetcpb_perspective_proto_rawDescData)
	})
	return file_api_m3uetcpb_perspective_proto_rawDescData
}

var file_api_m3uetcpb_perspective_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_m3uetcpb_perspective_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_api_m3uetcpb_perspective_proto_goTypes = []interface{}{
	(Perspective)(0),                          // 0: m3uetcpb.Perspective
	(*GetActivePerspectiveResponse)(nil),      // 1: m3uetcpb.GetActivePerspectiveResponse
	(*SetActivePerspectiveRequest)(nil),       // 2: m3uetcpb.SetActivePerspectiveRequest
	(*SubscribeToPerspectiveResponse)(nil),    // 3: m3uetcpb.SubscribeToPerspectiveResponse
	(*UnsubscribeFromPerspectiveRequest)(nil), // 4: m3uetcpb.UnsubscribeFromPerspectiveRequest
	(*PerspectiveDigest)(nil),                 // 5: m3uetcpb.PerspectiveDigest
	(*empty.Empty)(nil),                       // 6: google.protobuf.Empty
}
var file_api_m3uetcpb_perspective_proto_depIdxs = []int32{
	0, // 0: m3uetcpb.GetActivePerspectiveResponse.perspective:type_name -> m3uetcpb.Perspective
	0, // 1: m3uetcpb.SetActivePerspectiveRequest.perspective:type_name -> m3uetcpb.Perspective
	0, // 2: m3uetcpb.SubscribeToPerspectiveResponse.active_perspective:type_name -> m3uetcpb.Perspective
	0, // 3: m3uetcpb.PerspectiveDigest.perspective:type_name -> m3uetcpb.Perspective
	6, // 4: m3uetcpb.PerspectiveSvc.GetActivePerspective:input_type -> google.protobuf.Empty
	2, // 5: m3uetcpb.PerspectiveSvc.SetActivePerspective:input_type -> m3uetcpb.SetActivePerspectiveRequest
	6, // 6: m3uetcpb.PerspectiveSvc.SubscribeToPerspective:input_type -> google.protobuf.Empty
	4, // 7: m3uetcpb.PerspectiveSvc.UnsubscribeFromPerspective:input_type -> m3uetcpb.UnsubscribeFromPerspectiveRequest
	1, // 8: m3uetcpb.PerspectiveSvc.GetActivePerspective:output_type -> m3uetcpb.GetActivePerspectiveResponse
	6, // 9: m3uetcpb.PerspectiveSvc.SetActivePerspective:output_type -> google.protobuf.Empty
	3, // 10: m3uetcpb.PerspectiveSvc.SubscribeToPerspective:output_type -> m3uetcpb.SubscribeToPerspectiveResponse
	6, // 11: m3uetcpb.PerspectiveSvc.UnsubscribeFromPerspective:output_type -> google.protobuf.Empty
	8, // [8:12] is the sub-list for method output_type
	4, // [4:8] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_api_m3uetcpb_perspective_proto_init() }
func file_api_m3uetcpb_perspective_proto_init() {
	if File_api_m3uetcpb_perspective_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_m3uetcpb_perspective_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetActivePerspectiveResponse); i {
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
		file_api_m3uetcpb_perspective_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetActivePerspectiveRequest); i {
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
		file_api_m3uetcpb_perspective_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubscribeToPerspectiveResponse); i {
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
		file_api_m3uetcpb_perspective_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnsubscribeFromPerspectiveRequest); i {
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
		file_api_m3uetcpb_perspective_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PerspectiveDigest); i {
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
			RawDescriptor: file_api_m3uetcpb_perspective_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_m3uetcpb_perspective_proto_goTypes,
		DependencyIndexes: file_api_m3uetcpb_perspective_proto_depIdxs,
		EnumInfos:         file_api_m3uetcpb_perspective_proto_enumTypes,
		MessageInfos:      file_api_m3uetcpb_perspective_proto_msgTypes,
	}.Build()
	File_api_m3uetcpb_perspective_proto = out.File
	file_api_m3uetcpb_perspective_proto_rawDesc = nil
	file_api_m3uetcpb_perspective_proto_goTypes = nil
	file_api_m3uetcpb_perspective_proto_depIdxs = nil
}
