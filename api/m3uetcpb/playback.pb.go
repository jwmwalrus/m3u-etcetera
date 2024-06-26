// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: api/m3uetcpb/playback.proto

package m3uetcpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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
	PlaybackAction_PB_NONE     PlaybackAction = 0
	PlaybackAction_PB_PLAY     PlaybackAction = 1
	PlaybackAction_PB_NEXT     PlaybackAction = 2
	PlaybackAction_PB_PREVIOUS PlaybackAction = 3
	PlaybackAction_PB_SEEK     PlaybackAction = 4
	PlaybackAction_PB_PAUSE    PlaybackAction = 5
	PlaybackAction_PB_STOP     PlaybackAction = 6
)

// Enum value maps for PlaybackAction.
var (
	PlaybackAction_name = map[int32]string{
		0: "PB_NONE",
		1: "PB_PLAY",
		2: "PB_NEXT",
		3: "PB_PREVIOUS",
		4: "PB_SEEK",
		5: "PB_PAUSE",
		6: "PB_STOP",
	}
	PlaybackAction_value = map[string]int32{
		"PB_NONE":     0,
		"PB_PLAY":     1,
		"PB_NEXT":     2,
		"PB_PREVIOUS": 3,
		"PB_SEEK":     4,
		"PB_PAUSE":    5,
		"PB_STOP":     6,
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

type GetPlaybackResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsStreaming bool      `protobuf:"varint,1,opt,name=is_streaming,json=isStreaming,proto3" json:"is_streaming,omitempty"`
	IsPlaying   bool      `protobuf:"varint,2,opt,name=is_playing,json=isPlaying,proto3" json:"is_playing,omitempty"`
	IsPaused    bool      `protobuf:"varint,3,opt,name=is_paused,json=isPaused,proto3" json:"is_paused,omitempty"`
	IsStopped   bool      `protobuf:"varint,4,opt,name=is_stopped,json=isStopped,proto3" json:"is_stopped,omitempty"`
	IsReady     bool      `protobuf:"varint,5,opt,name=is_ready,json=isReady,proto3" json:"is_ready,omitempty"`
	Playback    *Playback `protobuf:"bytes,6,opt,name=playback,proto3" json:"playback,omitempty"`
	Track       *Track    `protobuf:"bytes,7,opt,name=track,proto3" json:"track,omitempty"`
}

func (x *GetPlaybackResponse) Reset() {
	*x = GetPlaybackResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_playback_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPlaybackResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPlaybackResponse) ProtoMessage() {}

func (x *GetPlaybackResponse) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use GetPlaybackResponse.ProtoReflect.Descriptor instead.
func (*GetPlaybackResponse) Descriptor() ([]byte, []int) {
	return file_api_m3uetcpb_playback_proto_rawDescGZIP(), []int{0}
}

func (x *GetPlaybackResponse) GetIsStreaming() bool {
	if x != nil {
		return x.IsStreaming
	}
	return false
}

func (x *GetPlaybackResponse) GetIsPlaying() bool {
	if x != nil {
		return x.IsPlaying
	}
	return false
}

func (x *GetPlaybackResponse) GetIsPaused() bool {
	if x != nil {
		return x.IsPaused
	}
	return false
}

func (x *GetPlaybackResponse) GetIsStopped() bool {
	if x != nil {
		return x.IsStopped
	}
	return false
}

func (x *GetPlaybackResponse) GetIsReady() bool {
	if x != nil {
		return x.IsReady
	}
	return false
}

func (x *GetPlaybackResponse) GetPlayback() *Playback {
	if x != nil {
		return x.Playback
	}
	return nil
}

func (x *GetPlaybackResponse) GetTrack() *Track {
	if x != nil {
		return x.Track
	}
	return nil
}

type GetPlaybackListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlaybackEntries []*Playback `protobuf:"bytes,1,rep,name=playback_entries,json=playbackEntries,proto3" json:"playback_entries,omitempty"`
}

func (x *GetPlaybackListResponse) Reset() {
	*x = GetPlaybackListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_playback_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPlaybackListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPlaybackListResponse) ProtoMessage() {}

func (x *GetPlaybackListResponse) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use GetPlaybackListResponse.ProtoReflect.Descriptor instead.
func (*GetPlaybackListResponse) Descriptor() ([]byte, []int) {
	return file_api_m3uetcpb_playback_proto_rawDescGZIP(), []int{1}
}

func (x *GetPlaybackListResponse) GetPlaybackEntries() []*Playback {
	if x != nil {
		return x.PlaybackEntries
	}
	return nil
}

type ExecutePlaybackActionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Action      PlaybackAction `protobuf:"varint,1,opt,name=action,proto3,enum=m3uetcpb.PlaybackAction" json:"action,omitempty"`
	Force       bool           `protobuf:"varint,2,opt,name=force,proto3" json:"force,omitempty"`
	Seek        int64          `protobuf:"varint,3,opt,name=seek,proto3" json:"seek,omitempty"`
	Perspective Perspective    `protobuf:"varint,4,opt,name=perspective,proto3,enum=m3uetcpb.Perspective" json:"perspective,omitempty"`
	Ids         []int64        `protobuf:"varint,5,rep,packed,name=ids,proto3" json:"ids,omitempty"`
	Locations   []string       `protobuf:"bytes,6,rep,name=locations,proto3" json:"locations,omitempty"`
}

func (x *ExecutePlaybackActionRequest) Reset() {
	*x = ExecutePlaybackActionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_playback_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecutePlaybackActionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecutePlaybackActionRequest) ProtoMessage() {}

func (x *ExecutePlaybackActionRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use ExecutePlaybackActionRequest.ProtoReflect.Descriptor instead.
func (*ExecutePlaybackActionRequest) Descriptor() ([]byte, []int) {
	return file_api_m3uetcpb_playback_proto_rawDescGZIP(), []int{2}
}

func (x *ExecutePlaybackActionRequest) GetAction() PlaybackAction {
	if x != nil {
		return x.Action
	}
	return PlaybackAction_PB_NONE
}

func (x *ExecutePlaybackActionRequest) GetForce() bool {
	if x != nil {
		return x.Force
	}
	return false
}

func (x *ExecutePlaybackActionRequest) GetSeek() int64 {
	if x != nil {
		return x.Seek
	}
	return 0
}

func (x *ExecutePlaybackActionRequest) GetPerspective() Perspective {
	if x != nil {
		return x.Perspective
	}
	return Perspective_MUSIC
}

func (x *ExecutePlaybackActionRequest) GetIds() []int64 {
	if x != nil {
		return x.Ids
	}
	return nil
}

func (x *ExecutePlaybackActionRequest) GetLocations() []string {
	if x != nil {
		return x.Locations
	}
	return nil
}

type SubscribeToPlaybackResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SubscriptionId string    `protobuf:"bytes,1,opt,name=subscription_id,json=subscriptionId,proto3" json:"subscription_id,omitempty"`
	IsStreaming    bool      `protobuf:"varint,2,opt,name=is_streaming,json=isStreaming,proto3" json:"is_streaming,omitempty"`
	IsPlaying      bool      `protobuf:"varint,3,opt,name=is_playing,json=isPlaying,proto3" json:"is_playing,omitempty"`
	IsPaused       bool      `protobuf:"varint,4,opt,name=is_paused,json=isPaused,proto3" json:"is_paused,omitempty"`
	IsStopped      bool      `protobuf:"varint,5,opt,name=is_stopped,json=isStopped,proto3" json:"is_stopped,omitempty"`
	IsReady        bool      `protobuf:"varint,6,opt,name=is_ready,json=isReady,proto3" json:"is_ready,omitempty"`
	Playback       *Playback `protobuf:"bytes,7,opt,name=playback,proto3" json:"playback,omitempty"`
	Track          *Track    `protobuf:"bytes,8,opt,name=track,proto3" json:"track,omitempty"`
}

func (x *SubscribeToPlaybackResponse) Reset() {
	*x = SubscribeToPlaybackResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_playback_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscribeToPlaybackResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeToPlaybackResponse) ProtoMessage() {}

func (x *SubscribeToPlaybackResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_m3uetcpb_playback_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeToPlaybackResponse.ProtoReflect.Descriptor instead.
func (*SubscribeToPlaybackResponse) Descriptor() ([]byte, []int) {
	return file_api_m3uetcpb_playback_proto_rawDescGZIP(), []int{3}
}

func (x *SubscribeToPlaybackResponse) GetSubscriptionId() string {
	if x != nil {
		return x.SubscriptionId
	}
	return ""
}

func (x *SubscribeToPlaybackResponse) GetIsStreaming() bool {
	if x != nil {
		return x.IsStreaming
	}
	return false
}

func (x *SubscribeToPlaybackResponse) GetIsPlaying() bool {
	if x != nil {
		return x.IsPlaying
	}
	return false
}

func (x *SubscribeToPlaybackResponse) GetIsPaused() bool {
	if x != nil {
		return x.IsPaused
	}
	return false
}

func (x *SubscribeToPlaybackResponse) GetIsStopped() bool {
	if x != nil {
		return x.IsStopped
	}
	return false
}

func (x *SubscribeToPlaybackResponse) GetIsReady() bool {
	if x != nil {
		return x.IsReady
	}
	return false
}

func (x *SubscribeToPlaybackResponse) GetPlayback() *Playback {
	if x != nil {
		return x.Playback
	}
	return nil
}

func (x *SubscribeToPlaybackResponse) GetTrack() *Track {
	if x != nil {
		return x.Track
	}
	return nil
}

type UnsubscribeFromPlaybackRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SubscriptionId string `protobuf:"bytes,1,opt,name=subscription_id,json=subscriptionId,proto3" json:"subscription_id,omitempty"`
}

func (x *UnsubscribeFromPlaybackRequest) Reset() {
	*x = UnsubscribeFromPlaybackRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_playback_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnsubscribeFromPlaybackRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnsubscribeFromPlaybackRequest) ProtoMessage() {}

func (x *UnsubscribeFromPlaybackRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_m3uetcpb_playback_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnsubscribeFromPlaybackRequest.ProtoReflect.Descriptor instead.
func (*UnsubscribeFromPlaybackRequest) Descriptor() ([]byte, []int) {
	return file_api_m3uetcpb_playback_proto_rawDescGZIP(), []int{4}
}

func (x *UnsubscribeFromPlaybackRequest) GetSubscriptionId() string {
	if x != nil {
		return x.SubscriptionId
	}
	return ""
}

type Playback struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Location  string                 `protobuf:"bytes,2,opt,name=location,proto3" json:"location,omitempty"`
	Played    bool                   `protobuf:"varint,3,opt,name=played,proto3" json:"played,omitempty"`
	Skip      int64                  `protobuf:"varint,4,opt,name=skip,proto3" json:"skip,omitempty"`
	TrackId   int64                  `protobuf:"varint,5,opt,name=track_id,json=trackId,proto3" json:"track_id,omitempty"`
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,101,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,102,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *Playback) Reset() {
	*x = Playback{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_playback_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Playback) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Playback) ProtoMessage() {}

func (x *Playback) ProtoReflect() protoreflect.Message {
	mi := &file_api_m3uetcpb_playback_proto_msgTypes[5]
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
	return file_api_m3uetcpb_playback_proto_rawDescGZIP(), []int{5}
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

func (x *Playback) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Playback) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

var File_api_m3uetcpb_playback_proto protoreflect.FileDescriptor

var file_api_m3uetcpb_playback_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2f, 0x70,
	0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x6d,
	0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x18, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x33,
	0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1e, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62,
	0x2f, 0x70, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x18, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62,
	0x2f, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x85, 0x02, 0x0a,
	0x13, 0x47, 0x65, 0x74, 0x50, 0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x69, 0x73, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x69, 0x73, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x69, 0x6e, 0x67, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x73, 0x5f, 0x70, 0x6c,
	0x61, 0x79, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x69, 0x73, 0x50,
	0x6c, 0x61, 0x79, 0x69, 0x6e, 0x67, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x70, 0x61, 0x75,
	0x73, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x50, 0x61, 0x75,
	0x73, 0x65, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x73, 0x5f, 0x73, 0x74, 0x6f, 0x70, 0x70, 0x65,
	0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x69, 0x73, 0x53, 0x74, 0x6f, 0x70, 0x70,
	0x65, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x69, 0x73, 0x5f, 0x72, 0x65, 0x61, 0x64, 0x79, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x69, 0x73, 0x52, 0x65, 0x61, 0x64, 0x79, 0x12, 0x2e, 0x0a,
	0x08, 0x70, 0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x12, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x62,
	0x61, 0x63, 0x6b, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x12, 0x25, 0x0a,
	0x05, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d,
	0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x52, 0x05, 0x74,
	0x72, 0x61, 0x63, 0x6b, 0x22, 0x58, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x50, 0x6c, 0x61, 0x79, 0x62,
	0x61, 0x63, 0x6b, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x3d, 0x0a, 0x10, 0x70, 0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x5f, 0x65, 0x6e, 0x74, 0x72,
	0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x33, 0x75, 0x65,
	0x74, 0x63, 0x70, 0x62, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x52, 0x0f, 0x70,
	0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x22, 0xe3,
	0x01, 0x0a, 0x1c, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x50, 0x6c, 0x61, 0x79, 0x62, 0x61,
	0x63, 0x6b, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x30, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x18, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x62,
	0x61, 0x63, 0x6b, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x05, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x65, 0x65, 0x6b, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x73, 0x65, 0x65, 0x6b, 0x12, 0x37, 0x0a, 0x0b, 0x70,
	0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x15, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x50, 0x65, 0x72, 0x73,
	0x70, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x52, 0x0b, 0x70, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63,
	0x74, 0x69, 0x76, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28,
	0x03, 0x52, 0x03, 0x69, 0x64, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x6c, 0x6f, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x22, 0xb6, 0x02, 0x0a, 0x1b, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69,
	0x62, 0x65, 0x54, 0x6f, 0x50, 0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x73,
	0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x21, 0x0a,
	0x0c, 0x69, 0x73, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x0b, 0x69, 0x73, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x69, 0x6e, 0x67,
	0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x73, 0x5f, 0x70, 0x6c, 0x61, 0x79, 0x69, 0x6e, 0x67, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x69, 0x73, 0x50, 0x6c, 0x61, 0x79, 0x69, 0x6e, 0x67, 0x12,
	0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x70, 0x61, 0x75, 0x73, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x50, 0x61, 0x75, 0x73, 0x65, 0x64, 0x12, 0x1d, 0x0a, 0x0a,
	0x69, 0x73, 0x5f, 0x73, 0x74, 0x6f, 0x70, 0x70, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x09, 0x69, 0x73, 0x53, 0x74, 0x6f, 0x70, 0x70, 0x65, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x69,
	0x73, 0x5f, 0x72, 0x65, 0x61, 0x64, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x69,
	0x73, 0x52, 0x65, 0x61, 0x64, 0x79, 0x12, 0x2e, 0x0a, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x62, 0x61,
	0x63, 0x6b, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74,
	0x63, 0x70, 0x62, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x52, 0x08, 0x70, 0x6c,
	0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x12, 0x25, 0x0a, 0x05, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62,
	0x2e, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x52, 0x05, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x22, 0x49, 0x0a,
	0x1e, 0x55, 0x6e, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x46, 0x72, 0x6f, 0x6d,
	0x50, 0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x27, 0x0a, 0x0f, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0xf3, 0x01, 0x0a, 0x08, 0x50, 0x6c, 0x61,
	0x79, 0x62, 0x61, 0x63, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x06, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6b, 0x69,
	0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x73, 0x6b, 0x69, 0x70, 0x12, 0x19, 0x0a,
	0x08, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x07, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x49, 0x64, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x65, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x18, 0x66, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x2a, 0x70,
	0x0a, 0x0e, 0x50, 0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x0b, 0x0a, 0x07, 0x50, 0x42, 0x5f, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x0b, 0x0a,
	0x07, 0x50, 0x42, 0x5f, 0x50, 0x4c, 0x41, 0x59, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x42,
	0x5f, 0x4e, 0x45, 0x58, 0x54, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x50, 0x42, 0x5f, 0x50, 0x52,
	0x45, 0x56, 0x49, 0x4f, 0x55, 0x53, 0x10, 0x03, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x42, 0x5f, 0x53,
	0x45, 0x45, 0x4b, 0x10, 0x04, 0x12, 0x0c, 0x0a, 0x08, 0x50, 0x42, 0x5f, 0x50, 0x41, 0x55, 0x53,
	0x45, 0x10, 0x05, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x42, 0x5f, 0x53, 0x54, 0x4f, 0x50, 0x10, 0x06,
	0x32, 0x8c, 0x03, 0x0a, 0x0b, 0x50, 0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x53, 0x76, 0x63,
	0x12, 0x3d, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x50, 0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x12,
	0x0f, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x1a, 0x1d, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x50,
	0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x45, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x50, 0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x4c, 0x69,
	0x73, 0x74, 0x12, 0x0f, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x1a, 0x21, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x47,
	0x65, 0x74, 0x50, 0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x50, 0x0a, 0x15, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74,
	0x65, 0x50, 0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x26, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x45, 0x78, 0x65, 0x63, 0x75,
	0x74, 0x65, 0x50, 0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63,
	0x70, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x4f, 0x0a, 0x13, 0x53, 0x75, 0x62, 0x73,
	0x63, 0x72, 0x69, 0x62, 0x65, 0x54, 0x6f, 0x50, 0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x12,
	0x0f, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x1a, 0x25, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x53, 0x75, 0x62, 0x73,
	0x63, 0x72, 0x69, 0x62, 0x65, 0x54, 0x6f, 0x50, 0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x54, 0x0a, 0x17, 0x55, 0x6e, 0x73,
	0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x50, 0x6c, 0x61, 0x79,
	0x62, 0x61, 0x63, 0x6b, 0x12, 0x28, 0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e,
	0x55, 0x6e, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x50,
	0x6c, 0x61, 0x79, 0x62, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f,
	0x2e, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42,
	0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
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
var file_api_m3uetcpb_playback_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_api_m3uetcpb_playback_proto_goTypes = []interface{}{
	(PlaybackAction)(0),                    // 0: m3uetcpb.PlaybackAction
	(*GetPlaybackResponse)(nil),            // 1: m3uetcpb.GetPlaybackResponse
	(*GetPlaybackListResponse)(nil),        // 2: m3uetcpb.GetPlaybackListResponse
	(*ExecutePlaybackActionRequest)(nil),   // 3: m3uetcpb.ExecutePlaybackActionRequest
	(*SubscribeToPlaybackResponse)(nil),    // 4: m3uetcpb.SubscribeToPlaybackResponse
	(*UnsubscribeFromPlaybackRequest)(nil), // 5: m3uetcpb.UnsubscribeFromPlaybackRequest
	(*Playback)(nil),                       // 6: m3uetcpb.Playback
	(*Track)(nil),                          // 7: m3uetcpb.Track
	(Perspective)(0),                       // 8: m3uetcpb.Perspective
	(*timestamppb.Timestamp)(nil),          // 9: google.protobuf.Timestamp
	(*Empty)(nil),                          // 10: m3uetcpb.Empty
}
var file_api_m3uetcpb_playback_proto_depIdxs = []int32{
	6,  // 0: m3uetcpb.GetPlaybackResponse.playback:type_name -> m3uetcpb.Playback
	7,  // 1: m3uetcpb.GetPlaybackResponse.track:type_name -> m3uetcpb.Track
	6,  // 2: m3uetcpb.GetPlaybackListResponse.playback_entries:type_name -> m3uetcpb.Playback
	0,  // 3: m3uetcpb.ExecutePlaybackActionRequest.action:type_name -> m3uetcpb.PlaybackAction
	8,  // 4: m3uetcpb.ExecutePlaybackActionRequest.perspective:type_name -> m3uetcpb.Perspective
	6,  // 5: m3uetcpb.SubscribeToPlaybackResponse.playback:type_name -> m3uetcpb.Playback
	7,  // 6: m3uetcpb.SubscribeToPlaybackResponse.track:type_name -> m3uetcpb.Track
	9,  // 7: m3uetcpb.Playback.created_at:type_name -> google.protobuf.Timestamp
	9,  // 8: m3uetcpb.Playback.updated_at:type_name -> google.protobuf.Timestamp
	10, // 9: m3uetcpb.PlaybackSvc.GetPlayback:input_type -> m3uetcpb.Empty
	10, // 10: m3uetcpb.PlaybackSvc.GetPlaybackList:input_type -> m3uetcpb.Empty
	3,  // 11: m3uetcpb.PlaybackSvc.ExecutePlaybackAction:input_type -> m3uetcpb.ExecutePlaybackActionRequest
	10, // 12: m3uetcpb.PlaybackSvc.SubscribeToPlayback:input_type -> m3uetcpb.Empty
	5,  // 13: m3uetcpb.PlaybackSvc.UnsubscribeFromPlayback:input_type -> m3uetcpb.UnsubscribeFromPlaybackRequest
	1,  // 14: m3uetcpb.PlaybackSvc.GetPlayback:output_type -> m3uetcpb.GetPlaybackResponse
	2,  // 15: m3uetcpb.PlaybackSvc.GetPlaybackList:output_type -> m3uetcpb.GetPlaybackListResponse
	10, // 16: m3uetcpb.PlaybackSvc.ExecutePlaybackAction:output_type -> m3uetcpb.Empty
	4,  // 17: m3uetcpb.PlaybackSvc.SubscribeToPlayback:output_type -> m3uetcpb.SubscribeToPlaybackResponse
	10, // 18: m3uetcpb.PlaybackSvc.UnsubscribeFromPlayback:output_type -> m3uetcpb.Empty
	14, // [14:19] is the sub-list for method output_type
	9,  // [9:14] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_api_m3uetcpb_playback_proto_init() }
func file_api_m3uetcpb_playback_proto_init() {
	if File_api_m3uetcpb_playback_proto != nil {
		return
	}
	file_api_m3uetcpb_empty_proto_init()
	file_api_m3uetcpb_perspective_proto_init()
	file_api_m3uetcpb_track_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_m3uetcpb_playback_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPlaybackResponse); i {
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
			switch v := v.(*GetPlaybackListResponse); i {
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
			switch v := v.(*ExecutePlaybackActionRequest); i {
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
		file_api_m3uetcpb_playback_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubscribeToPlaybackResponse); i {
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
		file_api_m3uetcpb_playback_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnsubscribeFromPlaybackRequest); i {
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
		file_api_m3uetcpb_playback_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_m3uetcpb_playback_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
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
