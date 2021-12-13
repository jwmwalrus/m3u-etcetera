// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.12.4
// source: api/m3uetcpb/track.proto

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

type Track struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Location    string `protobuf:"bytes,2,opt,name=location,proto3" json:"location,omitempty"`
	Format      string `protobuf:"bytes,3,opt,name=format,proto3" json:"format,omitempty"`
	Type        string `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	Title       string `protobuf:"bytes,5,opt,name=title,proto3" json:"title,omitempty"`
	Album       string `protobuf:"bytes,6,opt,name=album,proto3" json:"album,omitempty"`
	Artist      string `protobuf:"bytes,7,opt,name=artist,proto3" json:"artist,omitempty"`
	Albumartist string `protobuf:"bytes,8,opt,name=albumartist,proto3" json:"albumartist,omitempty"`
	Composer    string `protobuf:"bytes,9,opt,name=composer,proto3" json:"composer,omitempty"`
	Genre       string `protobuf:"bytes,10,opt,name=genre,proto3" json:"genre,omitempty"`
	Year        int32  `protobuf:"varint,11,opt,name=year,proto3" json:"year,omitempty"`
	Tracknumber int32  `protobuf:"varint,12,opt,name=tracknumber,proto3" json:"tracknumber,omitempty"`
	Tracktotal  int32  `protobuf:"varint,13,opt,name=tracktotal,proto3" json:"tracktotal,omitempty"`
	Discnumber  int32  `protobuf:"varint,14,opt,name=discnumber,proto3" json:"discnumber,omitempty"`
	Disctotal   int32  `protobuf:"varint,15,opt,name=disctotal,proto3" json:"disctotal,omitempty"`
	Lyrics      string `protobuf:"bytes,16,opt,name=lyrics,proto3" json:"lyrics,omitempty"`
	Comment     string `protobuf:"bytes,17,opt,name=comment,proto3" json:"comment,omitempty"`
	Tags        string `protobuf:"bytes,18,opt,name=tags,proto3" json:"tags,omitempty"`
	Sum         string `protobuf:"bytes,19,opt,name=sum,proto3" json:"sum,omitempty"`
	Playcount   int32  `protobuf:"varint,20,opt,name=playcount,proto3" json:"playcount,omitempty"`
	Rating      int32  `protobuf:"varint,21,opt,name=rating,proto3" json:"rating,omitempty"`
	Duration    int64  `protobuf:"varint,22,opt,name=duration,proto3" json:"duration,omitempty"`
	Remote      bool   `protobuf:"varint,23,opt,name=remote,proto3" json:"remote,omitempty"`
	Lastplayed  int64  `protobuf:"varint,24,opt,name=lastplayed,proto3" json:"lastplayed,omitempty"`
	CreatedAt   int64  `protobuf:"varint,101,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt   int64  `protobuf:"varint,102,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *Track) Reset() {
	*x = Track{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_m3uetcpb_track_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Track) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Track) ProtoMessage() {}

func (x *Track) ProtoReflect() protoreflect.Message {
	mi := &file_api_m3uetcpb_track_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Track.ProtoReflect.Descriptor instead.
func (*Track) Descriptor() ([]byte, []int) {
	return file_api_m3uetcpb_track_proto_rawDescGZIP(), []int{0}
}

func (x *Track) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Track) GetLocation() string {
	if x != nil {
		return x.Location
	}
	return ""
}

func (x *Track) GetFormat() string {
	if x != nil {
		return x.Format
	}
	return ""
}

func (x *Track) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Track) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Track) GetAlbum() string {
	if x != nil {
		return x.Album
	}
	return ""
}

func (x *Track) GetArtist() string {
	if x != nil {
		return x.Artist
	}
	return ""
}

func (x *Track) GetAlbumartist() string {
	if x != nil {
		return x.Albumartist
	}
	return ""
}

func (x *Track) GetComposer() string {
	if x != nil {
		return x.Composer
	}
	return ""
}

func (x *Track) GetGenre() string {
	if x != nil {
		return x.Genre
	}
	return ""
}

func (x *Track) GetYear() int32 {
	if x != nil {
		return x.Year
	}
	return 0
}

func (x *Track) GetTracknumber() int32 {
	if x != nil {
		return x.Tracknumber
	}
	return 0
}

func (x *Track) GetTracktotal() int32 {
	if x != nil {
		return x.Tracktotal
	}
	return 0
}

func (x *Track) GetDiscnumber() int32 {
	if x != nil {
		return x.Discnumber
	}
	return 0
}

func (x *Track) GetDisctotal() int32 {
	if x != nil {
		return x.Disctotal
	}
	return 0
}

func (x *Track) GetLyrics() string {
	if x != nil {
		return x.Lyrics
	}
	return ""
}

func (x *Track) GetComment() string {
	if x != nil {
		return x.Comment
	}
	return ""
}

func (x *Track) GetTags() string {
	if x != nil {
		return x.Tags
	}
	return ""
}

func (x *Track) GetSum() string {
	if x != nil {
		return x.Sum
	}
	return ""
}

func (x *Track) GetPlaycount() int32 {
	if x != nil {
		return x.Playcount
	}
	return 0
}

func (x *Track) GetRating() int32 {
	if x != nil {
		return x.Rating
	}
	return 0
}

func (x *Track) GetDuration() int64 {
	if x != nil {
		return x.Duration
	}
	return 0
}

func (x *Track) GetRemote() bool {
	if x != nil {
		return x.Remote
	}
	return false
}

func (x *Track) GetLastplayed() int64 {
	if x != nil {
		return x.Lastplayed
	}
	return 0
}

func (x *Track) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Track) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

var File_api_m3uetcpb_track_proto protoreflect.FileDescriptor

var file_api_m3uetcpb_track_proto_rawDesc = []byte{
	0x0a, 0x18, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62, 0x2f, 0x74,
	0x72, 0x61, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x6d, 0x33, 0x75, 0x65,
	0x74, 0x63, 0x70, 0x62, 0x22, 0xab, 0x05, 0x0a, 0x05, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a,
	0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x6f,
	0x72, 0x6d, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x6f, 0x72, 0x6d,
	0x61, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x61, 0x6c, 0x62, 0x75, 0x6d, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x61, 0x6c, 0x62,
	0x75, 0x6d, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x72, 0x74, 0x69, 0x73, 0x74, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x61, 0x72, 0x74, 0x69, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x61, 0x6c,
	0x62, 0x75, 0x6d, 0x61, 0x72, 0x74, 0x69, 0x73, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x61, 0x6c, 0x62, 0x75, 0x6d, 0x61, 0x72, 0x74, 0x69, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08,
	0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x73, 0x65, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x73, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x67, 0x65, 0x6e, 0x72,
	0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x67, 0x65, 0x6e, 0x72, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x79, 0x65, 0x61, 0x72, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x79, 0x65,
	0x61, 0x72, 0x12, 0x20, 0x0a, 0x0b, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x6e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x6e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x74,
	0x6f, 0x74, 0x61, 0x6c, 0x12, 0x1e, 0x0a, 0x0a, 0x64, 0x69, 0x73, 0x63, 0x6e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x64, 0x69, 0x73, 0x63, 0x6e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x64, 0x69, 0x73, 0x63, 0x74, 0x6f, 0x74, 0x61,
	0x6c, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x64, 0x69, 0x73, 0x63, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x79, 0x72, 0x69, 0x63, 0x73, 0x18, 0x10, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x6c, 0x79, 0x72, 0x69, 0x63, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f,
	0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x11, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6d,
	0x6d, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x12, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x75, 0x6d, 0x18,
	0x13, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73, 0x75, 0x6d, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x6c,
	0x61, 0x79, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x14, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x70,
	0x6c, 0x61, 0x79, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x61, 0x74, 0x69,
	0x6e, 0x67, 0x18, 0x15, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67,
	0x12, 0x1a, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x16, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06,
	0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x18, 0x17, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x72, 0x65,
	0x6d, 0x6f, 0x74, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x6c, 0x61, 0x73, 0x74, 0x70, 0x6c, 0x61, 0x79,
	0x65, 0x64, 0x18, 0x18, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x6c, 0x61, 0x73, 0x74, 0x70, 0x6c,
	0x61, 0x79, 0x65, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x65, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x18, 0x66, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x6d, 0x33, 0x75, 0x65, 0x74, 0x63, 0x70, 0x62,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_m3uetcpb_track_proto_rawDescOnce sync.Once
	file_api_m3uetcpb_track_proto_rawDescData = file_api_m3uetcpb_track_proto_rawDesc
)

func file_api_m3uetcpb_track_proto_rawDescGZIP() []byte {
	file_api_m3uetcpb_track_proto_rawDescOnce.Do(func() {
		file_api_m3uetcpb_track_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_m3uetcpb_track_proto_rawDescData)
	})
	return file_api_m3uetcpb_track_proto_rawDescData
}

var file_api_m3uetcpb_track_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_api_m3uetcpb_track_proto_goTypes = []interface{}{
	(*Track)(nil), // 0: m3uetcpb.Track
}
var file_api_m3uetcpb_track_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_m3uetcpb_track_proto_init() }
func file_api_m3uetcpb_track_proto_init() {
	if File_api_m3uetcpb_track_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_m3uetcpb_track_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Track); i {
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
			RawDescriptor: file_api_m3uetcpb_track_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_m3uetcpb_track_proto_goTypes,
		DependencyIndexes: file_api_m3uetcpb_track_proto_depIdxs,
		MessageInfos:      file_api_m3uetcpb_track_proto_msgTypes,
	}.Build()
	File_api_m3uetcpb_track_proto = out.File
	file_api_m3uetcpb_track_proto_rawDesc = nil
	file_api_m3uetcpb_track_proto_goTypes = nil
	file_api_m3uetcpb_track_proto_depIdxs = nil
}
