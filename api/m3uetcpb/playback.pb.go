// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.12.4
// source: api/m3uetcpb/playback.proto

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

type PlaybackAction int32

const (
	PlaybackAction_NONE     PlaybackAction = 0
	PlaybackAction_PLAY     PlaybackAction = 1
	PlaybackAction_NEXT     PlaybackAction = 2
	PlaybackAction_PREVIOUS PlaybackAction = 3
	PlaybackAction_SEEK     PlaybackAction = 4
	PlaybackAction_PAUSE    PlaybackAction = 5
	PlaybackAction_STOP     PlaybackAction = 6
)

// Enum value maps for PlaybackAction.
var (
	PlaybackAction_name = map[int32]string{
		0: "NONE",
		1: "PLAY",
		2: "NEXT",
		3: "PREVIOUS",
		4: "SEEK",
		5: "PAUSE",
		6: "STOP",
	}
	PlaybackAction_value = map[string]int32{
		"NONE":     0,
		"PLAY":     1,
		"NEXT":     2,
		"PREVIOUS": 3,
		"SEEK":     4,
		"PAUSE":    5,
		"STOP":     6,
	}
)

func (x PlaybackAction) Enum() *PlaybackAction {
	p := new(PlaybackAction)
	*p = x
	return p
}

func (x PlaybackAction) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PlaybackAction) Descriptor() protoreflect.EnumDescriptor {
	return file_api_m3uetcpb_playback_proto_enumTypes[0].Descriptor()
}

func (PlaybackAction) Type() protoreflect.EnumType {
	return &file_api_m3uetcpb_playback_proto_enumTypes[0]
}

func (x PlaybackAction) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PlaybackAction.Descriptor instead.
func (PlaybackAction) EnumDescriptor() ([]byte, []int) {
	return file_api_m3uetcpb_playback_proto_rawDescGZIP(), []int{0}
}

type Playback struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Location  string `protobuf:"bytes,2,opt,name=location,proto3" json:"location,omitempty"`
	Played    bool   `protobuf:"varint,3,opt,name=played,proto3" json:"played,omitempty"`
	Skip      int64  `protobuf:"varint,4,opt,name=skip,proto3" json:"skip,omitempty"`
	TrackId   int64  `protobuf:"varint,5,opt,name=track_id,json=trackId,proto3" json:"track_id,omitempty"`
	CreatedAt int64  `protobuf:"varint,101,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt int64  `protobuf:"varint,102,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *Playback) Reset() {
	*x = Playback{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_playback_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Playback) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Playback) ProtoMessage() {}

func (x *Playback) ProtoReflect() protoreflect.Message {
	mi := &file_api_m3uetcpb_playback_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Playback.ProtoReflect.Descriptor instead.
func (*Playback) Descriptor() ([]byte, []int) {
	return file_api_m3uetcpb_playback_proto_rawDescGZIP(), []int{0}
}

func (x *Playback) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Playback) GetLocation() string {
	if x != nil {
		return x.Location
	}
	return ""
}

func (x *Playback) GetPlayed() bool {
	if x != nil {
		return x.Played
	}
	return false
}

func (x *Playback) GetSkip() int64 {
	if x != nil {
		return x.Skip
	}
	return 0
}

func (x *Playback) GetTrackId() int64 {
	if x != nil {
		return x.TrackId
	}
	return 0
}

func (x *Playback) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Playback) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

type GetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Playing:
	//	*GetResponse_Empty
	//	*GetResponse_Playback
	//	*GetResponse_Track
	Playing isGetResponse_Playing `protobuf_oneof:"playing"`
}

func (x *GetResponse) Reset() {
	*x = GetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_playback_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetResponse) ProtoMessage() {}

func (x *GetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_m3uetcpb_playback_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetResponse.ProtoReflect.Descriptor instead.
func (*GetResponse) Descriptor() ([]byte, []int) {
	return file_api_m3uetcpb_playback_proto_rawDescGZIP(), []int{1}
}

func (m *GetResponse) GetPlaying() isGetResponse_Playing {
	if m != nil {
		return m.Playing
	}
	return nil
}

func (x *GetResponse) GetEmpty() *Empty {
	if x, ok := x.GetPlaying().(*GetResponse_Empty); ok {
		return x.Empty
	}
	return nil
}

func (x *GetResponse) GetPlayback() *Playback {
	if x, ok := x.GetPlaying().(*GetResponse_Playback); ok {
		return x.Playback
	}
	return nil
}

func (x *GetResponse) GetTrack() *Track {
	if x, ok := x.GetPlaying().(*GetResponse_Track); ok {
		return x.Track
	}
	return nil
}

type isGetResponse_Playing interface {
	isGetResponse_Playing()
}

type GetResponse_Empty struct {
	Empty *Empty `protobuf:"bytes,1,opt,name=empty,proto3,oneof"`
}

type GetResponse_Playback struct {
	Playback *Playback `protobuf:"bytes,2,opt,name=playback,proto3,oneof"`
}

type GetResponse_Track struct {
	Track *Track `protobuf:"bytes,3,opt,name=track,proto3,oneof"`
}

func (*GetResponse_Empty) isGetResponse_Playing() {}

func (*GetResponse_Playback) isGetResponse_Playing() {}

func (*GetResponse_Track) isGetResponse_Playing() {}

type ExecuteActionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Action      PlaybackAction `protobuf:"varint,1,opt,name=action,proto3,enum=m3uetcpb.PlaybackAction" json:"action,omitempty"`
	Force       bool           `protobuf:"varint,2,opt,name=force,proto3" json:"force,omitempty"`
	Seek        int64          `protobuf:"varint,3,opt,name=seek,proto3" json:"seek,omitempty"`
	Perspective int32          `protobuf:"varint,4,opt,name=perspective,proto3" json:"perspective,omitempty"`
	Ids         []int64        `protobuf:"varint,5,rep,packed,name=ids,proto3" json:"ids,omitempty"`
	Locations   []string       `protobuf:"bytes,6,rep,name=locations,proto3" json:"locations,omitempty"`
}

func (x *ExecuteActionRequest) Reset() {
	*x = ExecuteActionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_playback_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecuteActionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecuteActionRequest) ProtoMessage() {}

func (x *ExecuteActionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_m3uetcpb_playback_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecuteActionRequest.ProtoReflect.Descriptor instead.
func (*ExecuteActionRequest) Descriptor() ([]byte, []int) {
	return file_api_m3uetcpb_playback_proto_rawDescGZIP(), []int{2}
}

func (x *ExecuteActionRequest) GetAction() PlaybackAction {
	if x != nil {
		return x.Action
	}
	return PlaybackAction_NONE
}

func (x *ExecuteActionRequest) GetForce() bool {
	if x != nil {
		return x.Force
	}
	return false
}

func (x *ExecuteActionRequest) GetSeek() int64 {
	if x != nil {
		return x.Seek
	}
	return 0
}

func (x *ExecuteActionRequest) GetPerspective() int32 {
	if x != nil {
		return x.Perspective
	}
	return 0
}

func (x *ExecuteActionRequest) GetIds() []int64 {
	if x != nil {
		return x.Ids
	}
	return nil
}

func (x *ExecuteActionRequest) GetLocations() []string {
	if x != nil {
		return x.Locations
	}
	return nil
}

var File_api_m3uetcpb_playback_proto protoreflect.FileDescriptor

var file_api_m3uetcpb_playback_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2f, 0x70,
	0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x6d,
	0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x1a, 0x19, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x33, 0x75,
	0x65, 0x74, 0x63, 0x70, 0x62, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x18, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62,
	0x2f, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xbb, 0x01, 0x0a,
	0x08, 0x50, 0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x6f, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x6f, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x73, 0x6b, 0x69, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x73, 0x6b, 0x69,
	0x70, 0x12, 0x19, 0x0a, 0x08, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x07, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x65, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x66, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x9c, 0x01, 0x0a, 0x0b, 0x47,
	0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x27, 0x0a, 0x05, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d, 0x33, 0x75, 0x65,
	0x74, 0x63, 0x70, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x48, 0x00, 0x52, 0x05, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x12, 0x30, 0x0a, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62,
	0x2e, 0x50, 0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x48, 0x00, 0x52, 0x08, 0x70, 0x6c, 0x61,
	0x79, 0x62, 0x61, 0x63, 0x6b, 0x12, 0x27, 0x0a, 0x05, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e,
	0x54, 0x72, 0x61, 0x63, 0x6b, 0x48, 0x00, 0x52, 0x05, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x42, 0x09,
	0x0a, 0x07, 0x70, 0x6c, 0x61, 0x79, 0x69, 0x6e, 0x67, 0x22, 0xc4, 0x01, 0x0a, 0x14, 0x45, 0x78,
	0x65, 0x63, 0x75, 0x74, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x30, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x18, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x50, 0x6c,
	0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x05, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x65,
	0x65, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x73, 0x65, 0x65, 0x6b, 0x12, 0x20,
	0x0a, 0x0b, 0x70, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0b, 0x70, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x69, 0x64, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x03, 0x52, 0x03, 0x69,
	0x64, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2a, 0x5b, 0x0a, 0x0e, 0x50, 0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x41, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04,
	0x50, 0x4c, 0x41, 0x59, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x45, 0x58, 0x54, 0x10, 0x02,
	0x12, 0x0c, 0x0a, 0x08, 0x50, 0x52, 0x45, 0x56, 0x49, 0x4f, 0x55, 0x53, 0x10, 0x03, 0x12, 0x08,
	0x0a, 0x04, 0x53, 0x45, 0x45, 0x4b, 0x10, 0x04, 0x12, 0x09, 0x0a, 0x05, 0x50, 0x41, 0x55, 0x53,
	0x45, 0x10, 0x05, 0x12, 0x08, 0x0a, 0x04, 0x53, 0x54, 0x4f, 0x50, 0x10, 0x06, 0x32, 0x7e, 0x0a,
	0x0b, 0x50, 0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x53, 0x76, 0x63, 0x12, 0x2d, 0x0a, 0x03,
	0x47, 0x65, 0x74, 0x12, 0x0f, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x1a, 0x15, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e,
	0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x40, 0x0a, 0x0d, 0x45,
	0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x2e, 0x6d,
	0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x41,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x6d,
	0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x0c, 0x5a,
	0x0a, 0x2e, 0x2f, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_api_m3uetcpb_playback_proto_rawDescOnce sync.Once
	file_api_m3uetcpb_playback_proto_rawDescData = file_api_m3uetcpb_playback_proto_rawDesc
)

func file_api_m3uetcpb_playback_proto_rawDescGZIP() []byte {
	file_api_m3uetcpb_playback_proto_rawDescOnce.Do(func() {
		file_api_m3uetcpb_playback_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_m3uetcpb_playback_proto_rawDescData)
	})
	return file_api_m3uetcpb_playback_proto_rawDescData
}

var file_api_m3uetcpb_playback_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_m3uetcpb_playback_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_api_m3uetcpb_playback_proto_goTypes = []interface{}{
	(PlaybackAction)(0),          // 0: m3uetcpb.PlaybackAction
	(*Playback)(nil),             // 1: m3uetcpb.Playback
	(*GetResponse)(nil),          // 2: m3uetcpb.GetResponse
	(*ExecuteActionRequest)(nil), // 3: m3uetcpb.ExecuteActionRequest
	(*Empty)(nil),                // 4: m3uetcpb.Empty
	(*Track)(nil),                // 5: m3uetcpb.Track
}
var file_api_m3uetcpb_playback_proto_depIdxs = []int32{
	4, // 0: m3uetcpb.GetResponse.empty:type_name -> m3uetcpb.Empty
	1, // 1: m3uetcpb.GetResponse.playback:type_name -> m3uetcpb.Playback
	5, // 2: m3uetcpb.GetResponse.track:type_name -> m3uetcpb.Track
	0, // 3: m3uetcpb.ExecuteActionRequest.action:type_name -> m3uetcpb.PlaybackAction
	4, // 4: m3uetcpb.PlaybackSvc.Get:input_type -> m3uetcpb.Empty
	3, // 5: m3uetcpb.PlaybackSvc.ExecuteAction:input_type -> m3uetcpb.ExecuteActionRequest
	2, // 6: m3uetcpb.PlaybackSvc.Get:output_type -> m3uetcpb.GetResponse
	4, // 7: m3uetcpb.PlaybackSvc.ExecuteAction:output_type -> m3uetcpb.Empty
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_api_m3uetcpb_playback_proto_init() }
func file_api_m3uetcpb_playback_proto_init() {
	if File_api_m3uetcpb_playback_proto != nil {
		return
	}
	file_api_m3uetcpb_common_proto_init()
	file_api_m3uetcpb_track_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_m3uetcpb_playback_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Playback); i {
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
		file_api_m3uetcpb_playback_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetResponse); i {
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
		file_api_m3uetcpb_playback_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExecuteActionRequest); i {
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
	file_api_m3uetcpb_playback_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*GetResponse_Empty)(nil),
		(*GetResponse_Playback)(nil),
		(*GetResponse_Track)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_m3uetcpb_playback_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_m3uetcpb_playback_proto_goTypes,
		DependencyIndexes: file_api_m3uetcpb_playback_proto_depIdxs,
		EnumInfos:         file_api_m3uetcpb_playback_proto_enumTypes,
		MessageInfos:      file_api_m3uetcpb_playback_proto_msgTypes,
	}.Build()
	File_api_m3uetcpb_playback_proto = out.File
	file_api_m3uetcpb_playback_proto_rawDesc = nil
	file_api_m3uetcpb_playback_proto_goTypes = nil
	file_api_m3uetcpb_playback_proto_depIdxs = nil
}
