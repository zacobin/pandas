// Code generated by protoc-gen-go. DO NOT EDIT.
// source: mixer.proto

package grpc_mixer_v1

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type AdaptorOptions struct {
	Name                 string   `protobuf:"bytes,1,opt,name=Name" json:"Name,omitempty", bson:"Name,omitempty"`
	Protocol             string   `protobuf:"bytes,2,opt,name=Protocol" json:"Protocol,omitempty", bson:"Protocol,omitempty"`
	IsProvider           bool     `protobuf:"varint,3,opt,name=IsProvider" json:"IsProvider,omitempty", bson:"IsProvider,omitempty"`
	ServicePort          string   `protobuf:"bytes,4,opt,name=ServicePort" json:"ServicePort,omitempty", bson:"ServicePort,omitempty"`
	ConnectURL           string   `protobuf:"bytes,5,opt,name=ConnectURL" json:"ConnectURL,omitempty", bson:"ConnectURL,omitempty"`
	IsTLSEnabled         bool     `protobuf:"varint,6,opt,name=IsTLSEnabled" json:"IsTLSEnabled,omitempty", bson:"IsTLSEnabled,omitempty"`
	KeyFile              []byte   `protobuf:"bytes,7,opt,name=KeyFile,proto3" json:"KeyFile,omitempty", bson:"KeyFile,omitempty"`
	CertFile             []byte   `protobuf:"bytes,8,opt,name=CertFile,proto3" json:"CertFile,omitempty", bson:"CertFile,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AdaptorOptions) Reset()         { *m = AdaptorOptions{} }
func (m *AdaptorOptions) String() string { return proto.CompactTextString(m) }
func (*AdaptorOptions) ProtoMessage()    {}
func (*AdaptorOptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_mixer_2a4fdb5a6a956ce4, []int{0}
}
func (m *AdaptorOptions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AdaptorOptions.Unmarshal(m, b)
}
func (m *AdaptorOptions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AdaptorOptions.Marshal(b, m, deterministic)
}
func (dst *AdaptorOptions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AdaptorOptions.Merge(dst, src)
}
func (m *AdaptorOptions) XXX_Size() int {
	return xxx_messageInfo_AdaptorOptions.Size(m)
}
func (m *AdaptorOptions) XXX_DiscardUnknown() {
	xxx_messageInfo_AdaptorOptions.DiscardUnknown(m)
}

var xxx_messageInfo_AdaptorOptions proto.InternalMessageInfo

func (m *AdaptorOptions) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *AdaptorOptions) GetProtocol() string {
	if m != nil {
		return m.Protocol
	}
	return ""
}

func (m *AdaptorOptions) GetIsProvider() bool {
	if m != nil {
		return m.IsProvider
	}
	return false
}

func (m *AdaptorOptions) GetServicePort() string {
	if m != nil {
		return m.ServicePort
	}
	return ""
}

func (m *AdaptorOptions) GetConnectURL() string {
	if m != nil {
		return m.ConnectURL
	}
	return ""
}

func (m *AdaptorOptions) GetIsTLSEnabled() bool {
	if m != nil {
		return m.IsTLSEnabled
	}
	return false
}

func (m *AdaptorOptions) GetKeyFile() []byte {
	if m != nil {
		return m.KeyFile
	}
	return nil
}

func (m *AdaptorOptions) GetCertFile() []byte {
	if m != nil {
		return m.CertFile
	}
	return nil
}

type CreateAdaptorRequest struct {
	UserID               string          `protobuf:"bytes,1,opt,name=UserID" json:"UserID,omitempty", bson:"UserID,omitempty"`
	AdaptorOptions       *AdaptorOptions `protobuf:"bytes,2,opt,name=AdaptorOptions" json:"AdaptorOptions,omitempty", bson:"AdaptorOptions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *CreateAdaptorRequest) Reset()         { *m = CreateAdaptorRequest{} }
func (m *CreateAdaptorRequest) String() string { return proto.CompactTextString(m) }
func (*CreateAdaptorRequest) ProtoMessage()    {}
func (*CreateAdaptorRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_mixer_2a4fdb5a6a956ce4, []int{1}
}
func (m *CreateAdaptorRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateAdaptorRequest.Unmarshal(m, b)
}
func (m *CreateAdaptorRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateAdaptorRequest.Marshal(b, m, deterministic)
}
func (dst *CreateAdaptorRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateAdaptorRequest.Merge(dst, src)
}
func (m *CreateAdaptorRequest) XXX_Size() int {
	return xxx_messageInfo_CreateAdaptorRequest.Size(m)
}
func (m *CreateAdaptorRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateAdaptorRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateAdaptorRequest proto.InternalMessageInfo

func (m *CreateAdaptorRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *CreateAdaptorRequest) GetAdaptorOptions() *AdaptorOptions {
	if m != nil {
		return m.AdaptorOptions
	}
	return nil
}

type CreateAdaptorResponse struct {
	AdaptorID            string   `protobuf:"bytes,1,opt,name=AdaptorID" json:"AdaptorID,omitempty", bson:"AdaptorID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateAdaptorResponse) Reset()         { *m = CreateAdaptorResponse{} }
func (m *CreateAdaptorResponse) String() string { return proto.CompactTextString(m) }
func (*CreateAdaptorResponse) ProtoMessage()    {}
func (*CreateAdaptorResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_mixer_2a4fdb5a6a956ce4, []int{2}
}
func (m *CreateAdaptorResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateAdaptorResponse.Unmarshal(m, b)
}
func (m *CreateAdaptorResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateAdaptorResponse.Marshal(b, m, deterministic)
}
func (dst *CreateAdaptorResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateAdaptorResponse.Merge(dst, src)
}
func (m *CreateAdaptorResponse) XXX_Size() int {
	return xxx_messageInfo_CreateAdaptorResponse.Size(m)
}
func (m *CreateAdaptorResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateAdaptorResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateAdaptorResponse proto.InternalMessageInfo

func (m *CreateAdaptorResponse) GetAdaptorID() string {
	if m != nil {
		return m.AdaptorID
	}
	return ""
}

type DeleteAdaptorRequest struct {
	UserID               string   `protobuf:"bytes,1,opt,name=UserID" json:"UserID,omitempty", bson:"UserID,omitempty"`
	AdaptorID            string   `protobuf:"bytes,2,opt,name=AdaptorID" json:"AdaptorID,omitempty", bson:"AdaptorID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteAdaptorRequest) Reset()         { *m = DeleteAdaptorRequest{} }
func (m *DeleteAdaptorRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteAdaptorRequest) ProtoMessage()    {}
func (*DeleteAdaptorRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_mixer_2a4fdb5a6a956ce4, []int{3}
}
func (m *DeleteAdaptorRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteAdaptorRequest.Unmarshal(m, b)
}
func (m *DeleteAdaptorRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteAdaptorRequest.Marshal(b, m, deterministic)
}
func (dst *DeleteAdaptorRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteAdaptorRequest.Merge(dst, src)
}
func (m *DeleteAdaptorRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteAdaptorRequest.Size(m)
}
func (m *DeleteAdaptorRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteAdaptorRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteAdaptorRequest proto.InternalMessageInfo

func (m *DeleteAdaptorRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *DeleteAdaptorRequest) GetAdaptorID() string {
	if m != nil {
		return m.AdaptorID
	}
	return ""
}

type DeleteAdaptorResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteAdaptorResponse) Reset()         { *m = DeleteAdaptorResponse{} }
func (m *DeleteAdaptorResponse) String() string { return proto.CompactTextString(m) }
func (*DeleteAdaptorResponse) ProtoMessage()    {}
func (*DeleteAdaptorResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_mixer_2a4fdb5a6a956ce4, []int{4}
}
func (m *DeleteAdaptorResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteAdaptorResponse.Unmarshal(m, b)
}
func (m *DeleteAdaptorResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteAdaptorResponse.Marshal(b, m, deterministic)
}
func (dst *DeleteAdaptorResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteAdaptorResponse.Merge(dst, src)
}
func (m *DeleteAdaptorResponse) XXX_Size() int {
	return xxx_messageInfo_DeleteAdaptorResponse.Size(m)
}
func (m *DeleteAdaptorResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteAdaptorResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteAdaptorResponse proto.InternalMessageInfo

type GetAdaptorFactoriesRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAdaptorFactoriesRequest) Reset()         { *m = GetAdaptorFactoriesRequest{} }
func (m *GetAdaptorFactoriesRequest) String() string { return proto.CompactTextString(m) }
func (*GetAdaptorFactoriesRequest) ProtoMessage()    {}
func (*GetAdaptorFactoriesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_mixer_2a4fdb5a6a956ce4, []int{5}
}
func (m *GetAdaptorFactoriesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAdaptorFactoriesRequest.Unmarshal(m, b)
}
func (m *GetAdaptorFactoriesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAdaptorFactoriesRequest.Marshal(b, m, deterministic)
}
func (dst *GetAdaptorFactoriesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAdaptorFactoriesRequest.Merge(dst, src)
}
func (m *GetAdaptorFactoriesRequest) XXX_Size() int {
	return xxx_messageInfo_GetAdaptorFactoriesRequest.Size(m)
}
func (m *GetAdaptorFactoriesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAdaptorFactoriesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetAdaptorFactoriesRequest proto.InternalMessageInfo

type GetAdaptorFactoriesResponse struct {
	AdaptorFactoryNames  []string `protobuf:"bytes,1,rep,name=AdaptorFactoryNames" json:"AdaptorFactoryNames,omitempty", bson:"AdaptorFactoryNames,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAdaptorFactoriesResponse) Reset()         { *m = GetAdaptorFactoriesResponse{} }
func (m *GetAdaptorFactoriesResponse) String() string { return proto.CompactTextString(m) }
func (*GetAdaptorFactoriesResponse) ProtoMessage()    {}
func (*GetAdaptorFactoriesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_mixer_2a4fdb5a6a956ce4, []int{6}
}
func (m *GetAdaptorFactoriesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAdaptorFactoriesResponse.Unmarshal(m, b)
}
func (m *GetAdaptorFactoriesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAdaptorFactoriesResponse.Marshal(b, m, deterministic)
}
func (dst *GetAdaptorFactoriesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAdaptorFactoriesResponse.Merge(dst, src)
}
func (m *GetAdaptorFactoriesResponse) XXX_Size() int {
	return xxx_messageInfo_GetAdaptorFactoriesResponse.Size(m)
}
func (m *GetAdaptorFactoriesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAdaptorFactoriesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetAdaptorFactoriesResponse proto.InternalMessageInfo

func (m *GetAdaptorFactoriesResponse) GetAdaptorFactoryNames() []string {
	if m != nil {
		return m.AdaptorFactoryNames
	}
	return nil
}

type GetAdaptorsRequest struct {
	UserID               string   `protobuf:"bytes,1,opt,name=UserID" json:"UserID,omitempty", bson:"UserID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAdaptorsRequest) Reset()         { *m = GetAdaptorsRequest{} }
func (m *GetAdaptorsRequest) String() string { return proto.CompactTextString(m) }
func (*GetAdaptorsRequest) ProtoMessage()    {}
func (*GetAdaptorsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_mixer_2a4fdb5a6a956ce4, []int{7}
}
func (m *GetAdaptorsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAdaptorsRequest.Unmarshal(m, b)
}
func (m *GetAdaptorsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAdaptorsRequest.Marshal(b, m, deterministic)
}
func (dst *GetAdaptorsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAdaptorsRequest.Merge(dst, src)
}
func (m *GetAdaptorsRequest) XXX_Size() int {
	return xxx_messageInfo_GetAdaptorsRequest.Size(m)
}
func (m *GetAdaptorsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAdaptorsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetAdaptorsRequest proto.InternalMessageInfo

func (m *GetAdaptorsRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

type GetAdaptorsResponse struct {
	AdaptorOptions       []*AdaptorOptions `protobuf:"bytes,1,rep,name=AdaptorOptions" json:"AdaptorOptions,omitempty", bson:"AdaptorOptions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *GetAdaptorsResponse) Reset()         { *m = GetAdaptorsResponse{} }
func (m *GetAdaptorsResponse) String() string { return proto.CompactTextString(m) }
func (*GetAdaptorsResponse) ProtoMessage()    {}
func (*GetAdaptorsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_mixer_2a4fdb5a6a956ce4, []int{8}
}
func (m *GetAdaptorsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAdaptorsResponse.Unmarshal(m, b)
}
func (m *GetAdaptorsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAdaptorsResponse.Marshal(b, m, deterministic)
}
func (dst *GetAdaptorsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAdaptorsResponse.Merge(dst, src)
}
func (m *GetAdaptorsResponse) XXX_Size() int {
	return xxx_messageInfo_GetAdaptorsResponse.Size(m)
}
func (m *GetAdaptorsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAdaptorsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetAdaptorsResponse proto.InternalMessageInfo

func (m *GetAdaptorsResponse) GetAdaptorOptions() []*AdaptorOptions {
	if m != nil {
		return m.AdaptorOptions
	}
	return nil
}

func init() {
	proto.RegisterType((*AdaptorOptions)(nil), "grpc.mixer.v1.AdaptorOptions")
	proto.RegisterType((*CreateAdaptorRequest)(nil), "grpc.mixer.v1.CreateAdaptorRequest")
	proto.RegisterType((*CreateAdaptorResponse)(nil), "grpc.mixer.v1.CreateAdaptorResponse")
	proto.RegisterType((*DeleteAdaptorRequest)(nil), "grpc.mixer.v1.DeleteAdaptorRequest")
	proto.RegisterType((*DeleteAdaptorResponse)(nil), "grpc.mixer.v1.DeleteAdaptorResponse")
	proto.RegisterType((*GetAdaptorFactoriesRequest)(nil), "grpc.mixer.v1.GetAdaptorFactoriesRequest")
	proto.RegisterType((*GetAdaptorFactoriesResponse)(nil), "grpc.mixer.v1.GetAdaptorFactoriesResponse")
	proto.RegisterType((*GetAdaptorsRequest)(nil), "grpc.mixer.v1.GetAdaptorsRequest")
	proto.RegisterType((*GetAdaptorsResponse)(nil), "grpc.mixer.v1.GetAdaptorsResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MixerClient is the client API for Mixer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MixerClient interface {
	CreateAdaptor(ctx context.Context, in *CreateAdaptorRequest, opts ...grpc.CallOption) (*CreateAdaptorResponse, error)
	DeleteAdaptor(ctx context.Context, in *DeleteAdaptorRequest, opts ...grpc.CallOption) (*DeleteAdaptorResponse, error)
	GetAdaptorFactories(ctx context.Context, in *GetAdaptorFactoriesRequest, opts ...grpc.CallOption) (*GetAdaptorFactoriesResponse, error)
	GetAdaptors(ctx context.Context, in *GetAdaptorsRequest, opts ...grpc.CallOption) (*GetAdaptorsResponse, error)
}

type mixerClient struct {
	cc *grpc.ClientConn
}

func NewMixerClient(cc *grpc.ClientConn) MixerClient {
	return &mixerClient{cc}
}

func (c *mixerClient) CreateAdaptor(ctx context.Context, in *CreateAdaptorRequest, opts ...grpc.CallOption) (*CreateAdaptorResponse, error) {
	out := new(CreateAdaptorResponse)
	err := c.cc.Invoke(ctx, "/grpc.mixer.v1.Mixer/CreateAdaptor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mixerClient) DeleteAdaptor(ctx context.Context, in *DeleteAdaptorRequest, opts ...grpc.CallOption) (*DeleteAdaptorResponse, error) {
	out := new(DeleteAdaptorResponse)
	err := c.cc.Invoke(ctx, "/grpc.mixer.v1.Mixer/DeleteAdaptor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mixerClient) GetAdaptorFactories(ctx context.Context, in *GetAdaptorFactoriesRequest, opts ...grpc.CallOption) (*GetAdaptorFactoriesResponse, error) {
	out := new(GetAdaptorFactoriesResponse)
	err := c.cc.Invoke(ctx, "/grpc.mixer.v1.Mixer/GetAdaptorFactories", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mixerClient) GetAdaptors(ctx context.Context, in *GetAdaptorsRequest, opts ...grpc.CallOption) (*GetAdaptorsResponse, error) {
	out := new(GetAdaptorsResponse)
	err := c.cc.Invoke(ctx, "/grpc.mixer.v1.Mixer/GetAdaptors", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MixerServer is the server API for Mixer service.
type MixerServer interface {
	CreateAdaptor(context.Context, *CreateAdaptorRequest) (*CreateAdaptorResponse, error)
	DeleteAdaptor(context.Context, *DeleteAdaptorRequest) (*DeleteAdaptorResponse, error)
	GetAdaptorFactories(context.Context, *GetAdaptorFactoriesRequest) (*GetAdaptorFactoriesResponse, error)
	GetAdaptors(context.Context, *GetAdaptorsRequest) (*GetAdaptorsResponse, error)
}

func RegisterMixerServer(s *grpc.Server, srv MixerServer) {
	s.RegisterService(&_Mixer_serviceDesc, srv)
}

func _Mixer_CreateAdaptor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAdaptorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MixerServer).CreateAdaptor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.mixer.v1.Mixer/CreateAdaptor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MixerServer).CreateAdaptor(ctx, req.(*CreateAdaptorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mixer_DeleteAdaptor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAdaptorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MixerServer).DeleteAdaptor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.mixer.v1.Mixer/DeleteAdaptor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MixerServer).DeleteAdaptor(ctx, req.(*DeleteAdaptorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mixer_GetAdaptorFactories_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAdaptorFactoriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MixerServer).GetAdaptorFactories(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.mixer.v1.Mixer/GetAdaptorFactories",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MixerServer).GetAdaptorFactories(ctx, req.(*GetAdaptorFactoriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mixer_GetAdaptors_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAdaptorsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MixerServer).GetAdaptors(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.mixer.v1.Mixer/GetAdaptors",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MixerServer).GetAdaptors(ctx, req.(*GetAdaptorsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Mixer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.mixer.v1.Mixer",
	HandlerType: (*MixerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAdaptor",
			Handler:    _Mixer_CreateAdaptor_Handler,
		},
		{
			MethodName: "DeleteAdaptor",
			Handler:    _Mixer_DeleteAdaptor_Handler,
		},
		{
			MethodName: "GetAdaptorFactories",
			Handler:    _Mixer_GetAdaptorFactories_Handler,
		},
		{
			MethodName: "GetAdaptors",
			Handler:    _Mixer_GetAdaptors_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mixer.proto",
}

func init() { proto.RegisterFile("mixer.proto", fileDescriptor_mixer_2a4fdb5a6a956ce4) }

var fileDescriptor_mixer_2a4fdb5a6a956ce4 = []byte{
	// 451 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xd1, 0x4e, 0xd4, 0x40,
	0x14, 0xa5, 0x2c, 0x2c, 0xbb, 0x77, 0xc1, 0x87, 0x0b, 0xe8, 0x64, 0x45, 0x53, 0x47, 0x1f, 0xaa,
	0x31, 0x8d, 0x62, 0xfc, 0x00, 0xb3, 0x80, 0xd9, 0xb8, 0xca, 0xa6, 0x88, 0x4f, 0xbc, 0x94, 0xee,
	0x8d, 0x69, 0x52, 0x3a, 0x75, 0x66, 0xd8, 0xc8, 0x37, 0xf8, 0xcb, 0x3e, 0x98, 0x0e, 0x43, 0xe9,
	0x74, 0x0b, 0x2e, 0x6f, 0xbd, 0xe7, 0xde, 0x7b, 0xce, 0xcc, 0x39, 0x93, 0xc2, 0xe0, 0x22, 0xfd,
	0x4d, 0x32, 0x2c, 0xa4, 0xd0, 0x02, 0xb7, 0x7e, 0xca, 0x22, 0x09, 0xaf, 0x91, 0xf9, 0x7b, 0xfe,
	0xd7, 0x83, 0x47, 0x9f, 0x66, 0x71, 0xa1, 0x85, 0x3c, 0x2e, 0x74, 0x2a, 0x72, 0x85, 0x08, 0x6b,
	0xdf, 0xe2, 0x0b, 0x62, 0x9e, 0xef, 0x05, 0xfd, 0xc8, 0x7c, 0xe3, 0x10, 0x7a, 0xd3, 0x72, 0x3d,
	0x11, 0x19, 0x5b, 0x35, 0x78, 0x55, 0xe3, 0x73, 0x80, 0xb1, 0x9a, 0x4a, 0x31, 0x4f, 0x67, 0x24,
	0x59, 0xc7, 0xf7, 0x82, 0x5e, 0x54, 0x43, 0xd0, 0x87, 0xc1, 0x09, 0xc9, 0x79, 0x9a, 0xd0, 0x54,
	0x48, 0xcd, 0xd6, 0xcc, 0x7a, 0x1d, 0x2a, 0x19, 0x46, 0x22, 0xcf, 0x29, 0xd1, 0xa7, 0xd1, 0x84,
	0xad, 0x9b, 0x81, 0x1a, 0x82, 0x1c, 0x36, 0xc7, 0xea, 0xfb, 0xe4, 0xe4, 0x30, 0x8f, 0xcf, 0x33,
	0x9a, 0xb1, 0xae, 0xd1, 0x70, 0x30, 0x64, 0xb0, 0xf1, 0x85, 0xae, 0x8e, 0xd2, 0x8c, 0xd8, 0x86,
	0xef, 0x05, 0x9b, 0xd1, 0x4d, 0x59, 0x9e, 0x7d, 0x44, 0x52, 0x9b, 0x56, 0xcf, 0xb4, 0xaa, 0x9a,
	0x5f, 0xc2, 0xce, 0x48, 0x52, 0xac, 0xc9, 0x7a, 0x10, 0xd1, 0xaf, 0x4b, 0x52, 0x1a, 0x1f, 0x43,
	0xf7, 0x54, 0x91, 0x1c, 0x1f, 0x58, 0x17, 0x6c, 0x85, 0x87, 0x4d, 0xb7, 0x8c, 0x1b, 0x83, 0xfd,
	0x67, 0xa1, 0x63, 0x6b, 0xe8, 0x0e, 0x45, 0x8d, 0x25, 0xfe, 0x11, 0x76, 0x1b, 0xb2, 0xaa, 0x10,
	0xb9, 0x22, 0xdc, 0x83, 0xbe, 0x85, 0x2a, 0xe9, 0x5b, 0x80, 0x4f, 0x60, 0xe7, 0x80, 0x32, 0x5a,
	0xfa, 0xb4, 0x0e, 0xdb, 0x6a, 0x93, 0xed, 0x09, 0xec, 0x36, 0xd8, 0xae, 0x0f, 0xc1, 0xf7, 0x60,
	0xf8, 0x99, 0xb4, 0x45, 0x8f, 0xe2, 0x44, 0x0b, 0x99, 0x92, 0xb2, 0x62, 0xfc, 0x18, 0x9e, 0xb6,
	0x76, 0xed, 0x0d, 0xde, 0xc1, 0xb6, 0xd3, 0xbb, 0x2a, 0xdf, 0x8f, 0x62, 0x9e, 0xdf, 0x09, 0xfa,
	0x51, 0x5b, 0x8b, 0xbf, 0x05, 0xbc, 0x25, 0x54, 0xff, 0xb9, 0x13, 0x3f, 0x83, 0x6d, 0x67, 0xda,
	0xca, 0x2e, 0x06, 0x53, 0x2a, 0x3e, 0x34, 0x98, 0xfd, 0x3f, 0x1d, 0x58, 0xff, 0x5a, 0xce, 0xe2,
	0x19, 0x6c, 0x39, 0x11, 0xe1, 0xcb, 0x06, 0x53, 0xdb, 0xbb, 0x19, 0xbe, 0xba, 0x7f, 0xc8, 0x1a,
	0xbc, 0x52, 0xb2, 0x3b, 0xde, 0x2f, 0xb0, 0xb7, 0xe5, 0xbc, 0xc0, 0xde, 0x1e, 0xdf, 0x0a, 0xe6,
	0x75, 0x8f, 0xaa, 0x88, 0xf0, 0x75, 0x63, 0xfd, 0xee, 0x90, 0x87, 0x6f, 0x96, 0x19, 0xad, 0xf4,
	0x7e, 0xc0, 0xa0, 0x96, 0x09, 0xbe, 0xb8, 0x73, 0xb9, 0xe2, 0xe7, 0xf7, 0x8d, 0xdc, 0xf0, 0x9e,
	0x77, 0xcd, 0x2f, 0xeb, 0xc3, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x6d, 0x4c, 0xfc, 0x67, 0xc1,
	0x04, 0x00, 0x00,
}
