// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: api/m3uetcpb/queue.proto

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

type QueueAction int32

const (
	QueueAction_Q_NONE    QueueAction = 0
	QueueAction_Q_APPEND  QueueAction = 1
	QueueAction_Q_INSERT  QueueAction = 2
	QueueAction_Q_PREPEND QueueAction = 3
	QueueAction_Q_DELETE  QueueAction = 4
	QueueAction_Q_CLEAR   QueueAction = 5
	QueueAction_Q_MOVE    QueueAction = 6
)

// Enum value maps for QueueAction.
var (
	QueueAction_name = map[int32]string{
		0: "Q_NONE",
		1: "Q_APPEND",
		2: "Q_INSERT",
		3: "Q_PREPEND",
		4: "Q_DELETE",
		5: "Q_CLEAR",
		6: "Q_MOVE",
	}
	QueueAction_value = map[string]int32{
		"Q_NONE":    0,
		"Q_APPEND":  1,
		"Q_INSERT":  2,
		"Q_PREPEND": 3,
		"Q_DELETE":  4,
		"Q_CLEAR":   5,
		"Q_MOVE":    6,
	}
)

func (x QueueAction) Enum() *QueueAction {
	p := new(QueueAction)
	*p = x
	return p
}

func (x QueueAction) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (QueueAction) Descriptor() protoreflect.EnumDescriptor {
	return file_api_m3uetcpb_queue_proto_enumTypes[0].Descriptor()
}

func (QueueAction) Type() protoreflect.EnumType {
	return &file_api_m3uetcpb_queue_proto_enumTypes[0]
}

func (x QueueAction) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use QueueAction.Descriptor instead.
func (QueueAction) EnumDescriptor() ([]byte, []int) {
	return file_api_m3uetcpb_queue_proto_rawDescGZIP(), []int{0}
}

type GetQueueRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Perspective Perspective `protobuf:"varint,1,opt,name=perspective,proto3,enum=m3uetcpb.Perspective" json:"perspective,omitempty"`
	Limit       int32       `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *GetQueueRequest) Reset() {
	*x = GetQueueRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_queue_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetQueueRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetQueueRequest) ProtoMessage() {}

func (x *GetQueueRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_m3uetcpb_queue_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetQueueRequest.ProtoReflect.Descriptor instead.
func (*GetQueueRequest) Descriptor() ([]byte, []int) {
	return file_api_m3uetcpb_queue_proto_rawDescGZIP(), []int{0}
}

func (x *GetQueueRequest) GetPerspective() Perspective {
	if x != nil {
		return x.Perspective
	}
	return Perspective_MUSIC
}

func (x *GetQueueRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type GetQueueResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueueTracks []*QueueTrack `protobuf:"bytes,1,rep,name=queue_tracks,json=queueTracks,proto3" json:"queue_tracks,omitempty"`
	Tracks      []*Track      `protobuf:"bytes,2,rep,name=tracks,proto3" json:"tracks,omitempty"`
	Duration    int64         `protobuf:"varint,3,opt,name=duration,proto3" json:"duration,omitempty"`
}

func (x *GetQueueResponse) Reset() {
	*x = GetQueueResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_queue_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetQueueResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetQueueResponse) ProtoMessage() {}

func (x *GetQueueResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_m3uetcpb_queue_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetQueueResponse.ProtoReflect.Descriptor instead.
func (*GetQueueResponse) Descriptor() ([]byte, []int) {
	return file_api_m3uetcpb_queue_proto_rawDescGZIP(), []int{1}
}

func (x *GetQueueResponse) GetQueueTracks() []*QueueTrack {
	if x != nil {
		return x.QueueTracks
	}
	return nil
}

func (x *GetQueueResponse) GetTracks() []*Track {
	if x != nil {
		return x.Tracks
	}
	return nil
}

func (x *GetQueueResponse) GetDuration() int64 {
	if x != nil {
		return x.Duration
	}
	return 0
}

type ExecuteQueueActionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Action       QueueAction `protobuf:"varint,1,opt,name=action,proto3,enum=m3uetcpb.QueueAction" json:"action,omitempty"`
	Position     int32       `protobuf:"varint,2,opt,name=position,proto3" json:"position,omitempty"`
	FromPosition int32       `protobuf:"varint,3,opt,name=from_position,json=fromPosition,proto3" json:"from_position,omitempty"`
	Perspective  Perspective `protobuf:"varint,4,opt,name=perspective,proto3,enum=m3uetcpb.Perspective" json:"perspective,omitempty"`
	Ids          []int64     `protobuf:"varint,5,rep,packed,name=ids,proto3" json:"ids,omitempty"`
	Locations    []string    `protobuf:"bytes,6,rep,name=locations,proto3" json:"locations,omitempty"`
}

func (x *ExecuteQueueActionRequest) Reset() {
	*x = ExecuteQueueActionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_queue_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecuteQueueActionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecuteQueueActionRequest) ProtoMessage() {}

func (x *ExecuteQueueActionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_m3uetcpb_queue_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecuteQueueActionRequest.ProtoReflect.Descriptor instead.
func (*ExecuteQueueActionRequest) Descriptor() ([]byte, []int) {
	return file_api_m3uetcpb_queue_proto_rawDescGZIP(), []int{2}
}

func (x *ExecuteQueueActionRequest) GetAction() QueueAction {
	if x != nil {
		return x.Action
	}
	return QueueAction_Q_NONE
}

func (x *ExecuteQueueActionRequest) GetPosition() int32 {
	if x != nil {
		return x.Position
	}
	return 0
}

func (x *ExecuteQueueActionRequest) GetFromPosition() int32 {
	if x != nil {
		return x.FromPosition
	}
	return 0
}

func (x *ExecuteQueueActionRequest) GetPerspective() Perspective {
	if x != nil {
		return x.Perspective
	}
	return Perspective_MUSIC
}

func (x *ExecuteQueueActionRequest) GetIds() []int64 {
	if x != nil {
		return x.Ids
	}
	return nil
}

func (x *ExecuteQueueActionRequest) GetLocations() []string {
	if x != nil {
		return x.Locations
	}
	return nil
}

type SubscribeToQueueStoreResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SubscriptionId string               `protobuf:"bytes,1,opt,name=subscription_id,json=subscriptionId,proto3" json:"subscription_id,omitempty"`
	QueueTracks    []*QueueTrack        `protobuf:"bytes,2,rep,name=queue_tracks,json=queueTracks,proto3" json:"queue_tracks,omitempty"`
	Tracks         []*Track             `protobuf:"bytes,3,rep,name=tracks,proto3" json:"tracks,omitempty"`
	Digest         []*PerspectiveDigest `protobuf:"bytes,4,rep,name=digest,proto3" json:"digest,omitempty"`
}

func (x *SubscribeToQueueStoreResponse) Reset() {
	*x = SubscribeToQueueStoreResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_queue_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscribeToQueueStoreResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeToQueueStoreResponse) ProtoMessage() {}

func (x *SubscribeToQueueStoreResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_m3uetcpb_queue_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeToQueueStoreResponse.ProtoReflect.Descriptor instead.
func (*SubscribeToQueueStoreResponse) Descriptor() ([]byte, []int) {
	return file_api_m3uetcpb_queue_proto_rawDescGZIP(), []int{3}
}

func (x *SubscribeToQueueStoreResponse) GetSubscriptionId() string {
	if x != nil {
		return x.SubscriptionId
	}
	return ""
}

func (x *SubscribeToQueueStoreResponse) GetQueueTracks() []*QueueTrack {
	if x != nil {
		return x.QueueTracks
	}
	return nil
}

func (x *SubscribeToQueueStoreResponse) GetTracks() []*Track {
	if x != nil {
		return x.Tracks
	}
	return nil
}

func (x *SubscribeToQueueStoreResponse) GetDigest() []*PerspectiveDigest {
	if x != nil {
		return x.Digest
	}
	return nil
}

type UnsubscribeFromQueueStoreRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SubscriptionId string `protobuf:"bytes,1,opt,name=subscription_id,json=subscriptionId,proto3" json:"subscription_id,omitempty"`
}

func (x *UnsubscribeFromQueueStoreRequest) Reset() {
	*x = UnsubscribeFromQueueStoreRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_queue_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnsubscribeFromQueueStoreRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnsubscribeFromQueueStoreRequest) ProtoMessage() {}

func (x *UnsubscribeFromQueueStoreRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_m3uetcpb_queue_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnsubscribeFromQueueStoreRequest.ProtoReflect.Descriptor instead.
func (*UnsubscribeFromQueueStoreRequest) Descriptor() ([]byte, []int) {
	return file_api_m3uetcpb_queue_proto_rawDescGZIP(), []int{4}
}

func (x *UnsubscribeFromQueueStoreRequest) GetSubscriptionId() string {
	if x != nil {
		return x.SubscriptionId
	}
	return ""
}

type QueueTrack struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          int64       `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Position    int32       `protobuf:"varint,2,opt,name=position,proto3" json:"position,omitempty"`
	Played      bool        `protobuf:"varint,3,opt,name=played,proto3" json:"played,omitempty"`
	Location    string      `protobuf:"bytes,4,opt,name=location,proto3" json:"location,omitempty"`
	Perspective Perspective `protobuf:"varint,5,opt,name=perspective,proto3,enum=m3uetcpb.Perspective" json:"perspective,omitempty"`
	TrackId     int64       `protobuf:"varint,6,opt,name=track_id,json=trackId,proto3" json:"track_id,omitempty"`
	CreatedAt   int64       `protobuf:"varint,101,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt   int64       `protobuf:"varint,102,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *QueueTrack) Reset() {
	*x = QueueTrack{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_queue_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueTrack) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueTrack) ProtoMessage() {}

func (x *QueueTrack) ProtoReflect() protoreflect.Message {
	mi := &file_api_m3uetcpb_queue_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueTrack.ProtoReflect.Descriptor instead.
func (*QueueTrack) Descriptor() ([]byte, []int) {
	return file_api_m3uetcpb_queue_proto_rawDescGZIP(), []int{5}
}

func (x *QueueTrack) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *QueueTrack) GetPosition() int32 {
	if x != nil {
		return x.Position
	}
	return 0
}

func (x *QueueTrack) GetPlayed() bool {
	if x != nil {
		return x.Played
	}
	return false
}

func (x *QueueTrack) GetLocation() string {
	if x != nil {
		return x.Location
	}
	return ""
}

func (x *QueueTrack) GetPerspective() Perspective {
	if x != nil {
		return x.Perspective
	}
	return Perspective_MUSIC
}

func (x *QueueTrack) GetTrackId() int64 {
	if x != nil {
		return x.TrackId
	}
	return 0
}

func (x *QueueTrack) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *QueueTrack) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

var File_api_m3uetcpb_queue_proto protoreflect.FileDescriptor

var file_api_m3uetcpb_queue_proto_rawDesc = []byte{
	0x0a, 0x18, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2f, 0x71,
	0x75, 0x65, 0x75, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x6d, 0x33, 0x75, 0x65,
	0x74, 0x63, 0x70, 0x62, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1e, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2f,
	0x70, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x18, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2f,
	0x74, 0x72, 0x61, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x60, 0x0a, 0x0f, 0x47,
	0x65, 0x74, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x37,
	0x0a, 0x0b, 0x70, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x50,
	0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x52, 0x0b, 0x70, 0x65, 0x72, 0x73,
	0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22, 0x90, 0x01,
	0x0a, 0x10, 0x47, 0x65, 0x74, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x37, 0x0a, 0x0c, 0x71, 0x75, 0x65, 0x75, 0x65, 0x5f, 0x74, 0x72, 0x61, 0x63,
	0x6b, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74,
	0x63, 0x70, 0x62, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x52, 0x0b,
	0x71, 0x75, 0x65, 0x75, 0x65, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x73, 0x12, 0x27, 0x0a, 0x06, 0x74,
	0x72, 0x61, 0x63, 0x6b, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d, 0x33,
	0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x52, 0x06, 0x74, 0x72,
	0x61, 0x63, 0x6b, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x22, 0xf4, 0x01, 0x0a, 0x19, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x51, 0x75, 0x65, 0x75,
	0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2d,
	0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x15,
	0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x41,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a,
	0x08, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x08, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x0d, 0x66, 0x72, 0x6f,
	0x6d, 0x5f, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0c, 0x66, 0x72, 0x6f, 0x6d, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x37,
	0x0a, 0x0b, 0x70, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x50,
	0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x52, 0x0b, 0x70, 0x65, 0x72, 0x73,
	0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64, 0x73, 0x18, 0x05,
	0x20, 0x03, 0x28, 0x03, 0x52, 0x03, 0x69, 0x64, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x6f, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x6c, 0x6f,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0xdf, 0x01, 0x0a, 0x1d, 0x53, 0x75, 0x62, 0x73,
	0x63, 0x72, 0x69, 0x62, 0x65, 0x54, 0x6f, 0x51, 0x75, 0x65, 0x75, 0x65, 0x53, 0x74, 0x6f, 0x72,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x73, 0x75, 0x62,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x49, 0x64, 0x12, 0x37, 0x0a, 0x0c, 0x71, 0x75, 0x65, 0x75, 0x65, 0x5f, 0x74, 0x72, 0x61, 0x63,
	0x6b, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74,
	0x63, 0x70, 0x62, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x52, 0x0b,
	0x71, 0x75, 0x65, 0x75, 0x65, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x73, 0x12, 0x27, 0x0a, 0x06, 0x74,
	0x72, 0x61, 0x63, 0x6b, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d, 0x33,
	0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x52, 0x06, 0x74, 0x72,
	0x61, 0x63, 0x6b, 0x73, 0x12, 0x33, 0x0a, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e,
	0x50, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x44, 0x69, 0x67, 0x65, 0x73,
	0x74, 0x52, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x22, 0x4b, 0x0a, 0x20, 0x55, 0x6e, 0x73,
	0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x51, 0x75, 0x65, 0x75,
	0x65, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a,
	0x0f, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0xfe, 0x01, 0x0a, 0x0a, 0x51, 0x75, 0x65, 0x75, 0x65,
	0x54, 0x72, 0x61, 0x63, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x06, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x6f, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x6f, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x37, 0x0a, 0x0b, 0x70, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63,
	0x74, 0x69, 0x76, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e, 0x6d, 0x33, 0x75,
	0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x50, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76,
	0x65, 0x52, 0x0b, 0x70, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x12, 0x19,
	0x0a, 0x08, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x07, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x65, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x66, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x2a, 0x6b, 0x0a, 0x0b, 0x51, 0x75, 0x65, 0x75, 0x65,
	0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0a, 0x0a, 0x06, 0x51, 0x5f, 0x4e, 0x4f, 0x4e, 0x45,
	0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x51, 0x5f, 0x41, 0x50, 0x50, 0x45, 0x4e, 0x44, 0x10, 0x01,
	0x12, 0x0c, 0x0a, 0x08, 0x51, 0x5f, 0x49, 0x4e, 0x53, 0x45, 0x52, 0x54, 0x10, 0x02, 0x12, 0x0d,
	0x0a, 0x09, 0x51, 0x5f, 0x50, 0x52, 0x45, 0x50, 0x45, 0x4e, 0x44, 0x10, 0x03, 0x12, 0x0c, 0x0a,
	0x08, 0x51, 0x5f, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x10, 0x04, 0x12, 0x0b, 0x0a, 0x07, 0x51,
	0x5f, 0x43, 0x4c, 0x45, 0x41, 0x52, 0x10, 0x05, 0x12, 0x0a, 0x0a, 0x06, 0x51, 0x5f, 0x4d, 0x4f,
	0x56, 0x45, 0x10, 0x06, 0x32, 0xdd, 0x02, 0x0a, 0x08, 0x51, 0x75, 0x65, 0x75, 0x65, 0x53, 0x76,
	0x63, 0x12, 0x41, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x51, 0x75, 0x65, 0x75, 0x65, 0x12, 0x19, 0x2e,
	0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x51, 0x75, 0x65, 0x75,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74,
	0x63, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x51, 0x0a, 0x12, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x51,
	0x75, 0x65, 0x75, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x2e, 0x6d, 0x33, 0x75,
	0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x51, 0x75, 0x65,
	0x75, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x5a, 0x0a, 0x15, 0x53, 0x75, 0x62, 0x73, 0x63,
	0x72, 0x69, 0x62, 0x65, 0x54, 0x6f, 0x51, 0x75, 0x65, 0x75, 0x65, 0x53, 0x74, 0x6f, 0x72, 0x65,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x27, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74,
	0x63, 0x70, 0x62, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x54, 0x6f, 0x51,
	0x75, 0x65, 0x75, 0x65, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x30, 0x01, 0x12, 0x5f, 0x0a, 0x19, 0x55, 0x6e, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69,
	0x62, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x51, 0x75, 0x65, 0x75, 0x65, 0x53, 0x74, 0x6f, 0x72, 0x65,
	0x12, 0x2a, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x55, 0x6e, 0x73, 0x75,
	0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x51, 0x75, 0x65, 0x75, 0x65,
	0x53, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_m3uetcpb_queue_proto_rawDescOnce sync.Once
	file_api_m3uetcpb_queue_proto_rawDescData = file_api_m3uetcpb_queue_proto_rawDesc
)

func file_api_m3uetcpb_queue_proto_rawDescGZIP() []byte {
	file_api_m3uetcpb_queue_proto_rawDescOnce.Do(func() {
		file_api_m3uetcpb_queue_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_m3uetcpb_queue_proto_rawDescData)
	})
	return file_api_m3uetcpb_queue_proto_rawDescData
}

var file_api_m3uetcpb_queue_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_m3uetcpb_queue_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_api_m3uetcpb_queue_proto_goTypes = []interface{}{
	(QueueAction)(0),                         // 0: m3uetcpb.QueueAction
	(*GetQueueRequest)(nil),                  // 1: m3uetcpb.GetQueueRequest
	(*GetQueueResponse)(nil),                 // 2: m3uetcpb.GetQueueResponse
	(*ExecuteQueueActionRequest)(nil),        // 3: m3uetcpb.ExecuteQueueActionRequest
	(*SubscribeToQueueStoreResponse)(nil),    // 4: m3uetcpb.SubscribeToQueueStoreResponse
	(*UnsubscribeFromQueueStoreRequest)(nil), // 5: m3uetcpb.UnsubscribeFromQueueStoreRequest
	(*QueueTrack)(nil),                       // 6: m3uetcpb.QueueTrack
	(Perspective)(0),                         // 7: m3uetcpb.Perspective
	(*Track)(nil),                            // 8: m3uetcpb.Track
	(*PerspectiveDigest)(nil),                // 9: m3uetcpb.PerspectiveDigest
	(*empty.Empty)(nil),                      // 10: google.protobuf.Empty
}
var file_api_m3uetcpb_queue_proto_depIdxs = []int32{
	7,  // 0: m3uetcpb.GetQueueRequest.perspective:type_name -> m3uetcpb.Perspective
	6,  // 1: m3uetcpb.GetQueueResponse.queue_tracks:type_name -> m3uetcpb.QueueTrack
	8,  // 2: m3uetcpb.GetQueueResponse.tracks:type_name -> m3uetcpb.Track
	0,  // 3: m3uetcpb.ExecuteQueueActionRequest.action:type_name -> m3uetcpb.QueueAction
	7,  // 4: m3uetcpb.ExecuteQueueActionRequest.perspective:type_name -> m3uetcpb.Perspective
	6,  // 5: m3uetcpb.SubscribeToQueueStoreResponse.queue_tracks:type_name -> m3uetcpb.QueueTrack
	8,  // 6: m3uetcpb.SubscribeToQueueStoreResponse.tracks:type_name -> m3uetcpb.Track
	9,  // 7: m3uetcpb.SubscribeToQueueStoreResponse.digest:type_name -> m3uetcpb.PerspectiveDigest
	7,  // 8: m3uetcpb.QueueTrack.perspective:type_name -> m3uetcpb.Perspective
	1,  // 9: m3uetcpb.QueueSvc.GetQueue:input_type -> m3uetcpb.GetQueueRequest
	3,  // 10: m3uetcpb.QueueSvc.ExecuteQueueAction:input_type -> m3uetcpb.ExecuteQueueActionRequest
	10, // 11: m3uetcpb.QueueSvc.SubscribeToQueueStore:input_type -> google.protobuf.Empty
	5,  // 12: m3uetcpb.QueueSvc.UnsubscribeFromQueueStore:input_type -> m3uetcpb.UnsubscribeFromQueueStoreRequest
	2,  // 13: m3uetcpb.QueueSvc.GetQueue:output_type -> m3uetcpb.GetQueueResponse
	10, // 14: m3uetcpb.QueueSvc.ExecuteQueueAction:output_type -> google.protobuf.Empty
	4,  // 15: m3uetcpb.QueueSvc.SubscribeToQueueStore:output_type -> m3uetcpb.SubscribeToQueueStoreResponse
	10, // 16: m3uetcpb.QueueSvc.UnsubscribeFromQueueStore:output_type -> google.protobuf.Empty
	13, // [13:17] is the sub-list for method output_type
	9,  // [9:13] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_api_m3uetcpb_queue_proto_init() }
func file_api_m3uetcpb_queue_proto_init() {
	if File_api_m3uetcpb_queue_proto != nil {
		return
	}
	file_api_m3uetcpb_perspective_proto_init()
	file_api_m3uetcpb_track_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_m3uetcpb_queue_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetQueueRequest); i {
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
		file_api_m3uetcpb_queue_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetQueueResponse); i {
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
		file_api_m3uetcpb_queue_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExecuteQueueActionRequest); i {
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
		file_api_m3uetcpb_queue_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubscribeToQueueStoreResponse); i {
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
		file_api_m3uetcpb_queue_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnsubscribeFromQueueStoreRequest); i {
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
		file_api_m3uetcpb_queue_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueTrack); i {
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
			RawDescriptor: file_api_m3uetcpb_queue_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_m3uetcpb_queue_proto_goTypes,
		DependencyIndexes: file_api_m3uetcpb_queue_proto_depIdxs,
		EnumInfos:         file_api_m3uetcpb_queue_proto_enumTypes,
		MessageInfos:      file_api_m3uetcpb_queue_proto_msgTypes,
	}.Build()
	File_api_m3uetcpb_queue_proto = out.File
	file_api_m3uetcpb_queue_proto_rawDesc = nil
	file_api_m3uetcpb_queue_proto_goTypes = nil
	file_api_m3uetcpb_queue_proto_depIdxs = nil
}
