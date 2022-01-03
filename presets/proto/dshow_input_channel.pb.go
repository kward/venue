// Generate proto code by running the following from the parent directory:
// $ protoc --go_out=. proto/*.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.1
// source: proto/dshow_input_channel.proto

package proto

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

type Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Bytes []byte `protobuf:"bytes,2,opt,name=bytes,proto3" json:"bytes,omitempty"`
	Int32 int32  `protobuf:"varint,3,opt,name=int32,proto3" json:"int32,omitempty"`
	Str   string `protobuf:"bytes,4,opt,name=str,proto3" json:"str,omitempty"`
}

func (x *Data) Reset() {
	*x = Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_dshow_input_channel_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data) ProtoMessage() {}

func (x *Data) ProtoReflect() protoreflect.Message {
	mi := &file_proto_dshow_input_channel_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data.ProtoReflect.Descriptor instead.
func (*Data) Descriptor() ([]byte, []int) {
	return file_proto_dshow_input_channel_proto_rawDescGZIP(), []int{0}
}

func (x *Data) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *Data) GetBytes() []byte {
	if x != nil {
		return x.Bytes
	}
	return nil
}

func (x *Data) GetInt32() int32 {
	if x != nil {
		return x.Int32
	}
	return 0
}

func (x *Data) GetStr() string {
	if x != nil {
		return x.Str
	}
	return ""
}

type DShowInputChannel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Header            *Data `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Version           *Data `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	FileType          *Data `protobuf:"bytes,3,opt,name=fileType,proto3" json:"fileType,omitempty"`
	UserComment       *Data `protobuf:"bytes,4,opt,name=userComment,proto3" json:"userComment,omitempty"`
	AudioMasterStrip  *Data `protobuf:"bytes,5,opt,name=audioMasterStrip,proto3" json:"audioMasterStrip,omitempty"`
	AudioStrip        *Data `protobuf:"bytes,6,opt,name=audioStrip,proto3" json:"audioStrip,omitempty"`
	AuxBussesOptions  *Data `protobuf:"bytes,7,opt,name=auxBussesOptions,proto3" json:"auxBussesOptions,omitempty"`
	AuxBussesOptions2 *Data `protobuf:"bytes,8,opt,name=auxBussesOptions2,proto3" json:"auxBussesOptions2,omitempty"`
	BusConfigMode     *Data `protobuf:"bytes,9,opt,name=busConfigMode,proto3" json:"busConfigMode,omitempty"`
	InputStrip        *Data `protobuf:"bytes,10,opt,name=inputStrip,proto3" json:"inputStrip,omitempty"`
	MatrixMasterStrip *Data `protobuf:"bytes,11,opt,name=matrixMasterStrip,proto3" json:"matrixMasterStrip,omitempty"`
	MicLineStrips     *Data `protobuf:"bytes,12,opt,name=micLineStrips,proto3" json:"micLineStrips,omitempty"`
	Strip             *Data `protobuf:"bytes,13,opt,name=strip,proto3" json:"strip,omitempty"`
	StripType         *Data `protobuf:"bytes,14,opt,name=stripType,proto3" json:"stripType,omitempty"`
}

func (x *DShowInputChannel) Reset() {
	*x = DShowInputChannel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_dshow_input_channel_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DShowInputChannel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DShowInputChannel) ProtoMessage() {}

func (x *DShowInputChannel) ProtoReflect() protoreflect.Message {
	mi := &file_proto_dshow_input_channel_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DShowInputChannel.ProtoReflect.Descriptor instead.
func (*DShowInputChannel) Descriptor() ([]byte, []int) {
	return file_proto_dshow_input_channel_proto_rawDescGZIP(), []int{1}
}

func (x *DShowInputChannel) GetHeader() *Data {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *DShowInputChannel) GetVersion() *Data {
	if x != nil {
		return x.Version
	}
	return nil
}

func (x *DShowInputChannel) GetFileType() *Data {
	if x != nil {
		return x.FileType
	}
	return nil
}

func (x *DShowInputChannel) GetUserComment() *Data {
	if x != nil {
		return x.UserComment
	}
	return nil
}

func (x *DShowInputChannel) GetAudioMasterStrip() *Data {
	if x != nil {
		return x.AudioMasterStrip
	}
	return nil
}

func (x *DShowInputChannel) GetAudioStrip() *Data {
	if x != nil {
		return x.AudioStrip
	}
	return nil
}

func (x *DShowInputChannel) GetAuxBussesOptions() *Data {
	if x != nil {
		return x.AuxBussesOptions
	}
	return nil
}

func (x *DShowInputChannel) GetAuxBussesOptions2() *Data {
	if x != nil {
		return x.AuxBussesOptions2
	}
	return nil
}

func (x *DShowInputChannel) GetBusConfigMode() *Data {
	if x != nil {
		return x.BusConfigMode
	}
	return nil
}

func (x *DShowInputChannel) GetInputStrip() *Data {
	if x != nil {
		return x.InputStrip
	}
	return nil
}

func (x *DShowInputChannel) GetMatrixMasterStrip() *Data {
	if x != nil {
		return x.MatrixMasterStrip
	}
	return nil
}

func (x *DShowInputChannel) GetMicLineStrips() *Data {
	if x != nil {
		return x.MicLineStrips
	}
	return nil
}

func (x *DShowInputChannel) GetStrip() *Data {
	if x != nil {
		return x.Strip
	}
	return nil
}

func (x *DShowInputChannel) GetStripType() *Data {
	if x != nil {
		return x.StripType
	}
	return nil
}

var File_proto_dshow_input_channel_proto protoreflect.FileDescriptor

var file_proto_dshow_input_channel_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x73, 0x68, 0x6f, 0x77, 0x5f, 0x69, 0x6e,
	0x70, 0x75, 0x74, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5a, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x62, 0x79, 0x74, 0x65, 0x73, 0x12, 0x14, 0x0a, 0x05,
	0x69, 0x6e, 0x74, 0x33, 0x32, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x69, 0x6e, 0x74,
	0x33, 0x32, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x74, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x73, 0x74, 0x72, 0x22, 0xad, 0x05, 0x0a, 0x11, 0x44, 0x53, 0x68, 0x6f, 0x77, 0x49, 0x6e,
	0x70, 0x75, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x23, 0x0a, 0x06, 0x68, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12,
	0x25, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x27, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x2d, 0x0a, 0x0b, 0x75, 0x73, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x0b, 0x75, 0x73, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x37,
	0x0a, 0x10, 0x61, 0x75, 0x64, 0x69, 0x6f, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x53, 0x74, 0x72,
	0x69, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x10, 0x61, 0x75, 0x64, 0x69, 0x6f, 0x4d, 0x61, 0x73, 0x74,
	0x65, 0x72, 0x53, 0x74, 0x72, 0x69, 0x70, 0x12, 0x2b, 0x0a, 0x0a, 0x61, 0x75, 0x64, 0x69, 0x6f,
	0x53, 0x74, 0x72, 0x69, 0x70, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x0a, 0x61, 0x75, 0x64, 0x69, 0x6f, 0x53,
	0x74, 0x72, 0x69, 0x70, 0x12, 0x37, 0x0a, 0x10, 0x61, 0x75, 0x78, 0x42, 0x75, 0x73, 0x73, 0x65,
	0x73, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x10, 0x61, 0x75, 0x78,
	0x42, 0x75, 0x73, 0x73, 0x65, 0x73, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x39, 0x0a,
	0x11, 0x61, 0x75, 0x78, 0x42, 0x75, 0x73, 0x73, 0x65, 0x73, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x32, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x11, 0x61, 0x75, 0x78, 0x42, 0x75, 0x73, 0x73, 0x65, 0x73,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x32, 0x12, 0x31, 0x0a, 0x0d, 0x62, 0x75, 0x73, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x4d, 0x6f, 0x64, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x0d, 0x62, 0x75,
	0x73, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x4d, 0x6f, 0x64, 0x65, 0x12, 0x2b, 0x0a, 0x0a, 0x69,
	0x6e, 0x70, 0x75, 0x74, 0x53, 0x74, 0x72, 0x69, 0x70, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x0a, 0x69, 0x6e,
	0x70, 0x75, 0x74, 0x53, 0x74, 0x72, 0x69, 0x70, 0x12, 0x39, 0x0a, 0x11, 0x6d, 0x61, 0x74, 0x72,
	0x69, 0x78, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x53, 0x74, 0x72, 0x69, 0x70, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x61, 0x74, 0x61,
	0x52, 0x11, 0x6d, 0x61, 0x74, 0x72, 0x69, 0x78, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x53, 0x74,
	0x72, 0x69, 0x70, 0x12, 0x31, 0x0a, 0x0d, 0x6d, 0x69, 0x63, 0x4c, 0x69, 0x6e, 0x65, 0x53, 0x74,
	0x72, 0x69, 0x70, 0x73, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x0d, 0x6d, 0x69, 0x63, 0x4c, 0x69, 0x6e, 0x65,
	0x53, 0x74, 0x72, 0x69, 0x70, 0x73, 0x12, 0x21, 0x0a, 0x05, 0x73, 0x74, 0x72, 0x69, 0x70, 0x18,
	0x0d, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x61,
	0x74, 0x61, 0x52, 0x05, 0x73, 0x74, 0x72, 0x69, 0x70, 0x12, 0x29, 0x0a, 0x09, 0x73, 0x74, 0x72,
	0x69, 0x70, 0x54, 0x79, 0x70, 0x65, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x09, 0x73, 0x74, 0x72, 0x69, 0x70,
	0x54, 0x79, 0x70, 0x65, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_dshow_input_channel_proto_rawDescOnce sync.Once
	file_proto_dshow_input_channel_proto_rawDescData = file_proto_dshow_input_channel_proto_rawDesc
)

func file_proto_dshow_input_channel_proto_rawDescGZIP() []byte {
	file_proto_dshow_input_channel_proto_rawDescOnce.Do(func() {
		file_proto_dshow_input_channel_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_dshow_input_channel_proto_rawDescData)
	})
	return file_proto_dshow_input_channel_proto_rawDescData
}

var file_proto_dshow_input_channel_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_dshow_input_channel_proto_goTypes = []interface{}{
	(*Data)(nil),              // 0: proto.Data
	(*DShowInputChannel)(nil), // 1: proto.DShowInputChannel
}
var file_proto_dshow_input_channel_proto_depIdxs = []int32{
	0,  // 0: proto.DShowInputChannel.header:type_name -> proto.Data
	0,  // 1: proto.DShowInputChannel.version:type_name -> proto.Data
	0,  // 2: proto.DShowInputChannel.fileType:type_name -> proto.Data
	0,  // 3: proto.DShowInputChannel.userComment:type_name -> proto.Data
	0,  // 4: proto.DShowInputChannel.audioMasterStrip:type_name -> proto.Data
	0,  // 5: proto.DShowInputChannel.audioStrip:type_name -> proto.Data
	0,  // 6: proto.DShowInputChannel.auxBussesOptions:type_name -> proto.Data
	0,  // 7: proto.DShowInputChannel.auxBussesOptions2:type_name -> proto.Data
	0,  // 8: proto.DShowInputChannel.busConfigMode:type_name -> proto.Data
	0,  // 9: proto.DShowInputChannel.inputStrip:type_name -> proto.Data
	0,  // 10: proto.DShowInputChannel.matrixMasterStrip:type_name -> proto.Data
	0,  // 11: proto.DShowInputChannel.micLineStrips:type_name -> proto.Data
	0,  // 12: proto.DShowInputChannel.strip:type_name -> proto.Data
	0,  // 13: proto.DShowInputChannel.stripType:type_name -> proto.Data
	14, // [14:14] is the sub-list for method output_type
	14, // [14:14] is the sub-list for method input_type
	14, // [14:14] is the sub-list for extension type_name
	14, // [14:14] is the sub-list for extension extendee
	0,  // [0:14] is the sub-list for field type_name
}

func init() { file_proto_dshow_input_channel_proto_init() }
func file_proto_dshow_input_channel_proto_init() {
	if File_proto_dshow_input_channel_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_dshow_input_channel_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Data); i {
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
		file_proto_dshow_input_channel_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DShowInputChannel); i {
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
			RawDescriptor: file_proto_dshow_input_channel_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_dshow_input_channel_proto_goTypes,
		DependencyIndexes: file_proto_dshow_input_channel_proto_depIdxs,
		MessageInfos:      file_proto_dshow_input_channel_proto_msgTypes,
	}.Build()
	File_proto_dshow_input_channel_proto = out.File
	file_proto_dshow_input_channel_proto_rawDesc = nil
	file_proto_dshow_input_channel_proto_goTypes = nil
	file_proto_dshow_input_channel_proto_depIdxs = nil
}
