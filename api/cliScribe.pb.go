// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cliScribe.proto

package api

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

type VersionRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VersionRequest) Reset()         { *m = VersionRequest{} }
func (m *VersionRequest) String() string { return proto.CompactTextString(m) }
func (*VersionRequest) ProtoMessage()    {}
func (*VersionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cliScribe_41e2a648a039d414, []int{0}
}
func (m *VersionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VersionRequest.Unmarshal(m, b)
}
func (m *VersionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VersionRequest.Marshal(b, m, deterministic)
}
func (dst *VersionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VersionRequest.Merge(dst, src)
}
func (m *VersionRequest) XXX_Size() int {
	return xxx_messageInfo_VersionRequest.Size(m)
}
func (m *VersionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_VersionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_VersionRequest proto.InternalMessageInfo

type VersionResponse struct {
	Version              string   `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VersionResponse) Reset()         { *m = VersionResponse{} }
func (m *VersionResponse) String() string { return proto.CompactTextString(m) }
func (*VersionResponse) ProtoMessage()    {}
func (*VersionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cliScribe_41e2a648a039d414, []int{1}
}
func (m *VersionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VersionResponse.Unmarshal(m, b)
}
func (m *VersionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VersionResponse.Marshal(b, m, deterministic)
}
func (dst *VersionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VersionResponse.Merge(dst, src)
}
func (m *VersionResponse) XXX_Size() int {
	return xxx_messageInfo_VersionResponse.Size(m)
}
func (m *VersionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_VersionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_VersionResponse proto.InternalMessageInfo

func (m *VersionResponse) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func init() {
	proto.RegisterType((*VersionRequest)(nil), "com.romanostrechlis.scribe.api.VersionRequest")
	proto.RegisterType((*VersionResponse)(nil), "com.romanostrechlis.scribe.api.VersionResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// CLIScribeClient is the client API for CLIScribe service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CLIScribeClient interface {
	GetVersion(ctx context.Context, in *VersionRequest, opts ...grpc.CallOption) (*VersionResponse, error)
}

type cLIScribeClient struct {
	cc *grpc.ClientConn
}

func NewCLIScribeClient(cc *grpc.ClientConn) CLIScribeClient {
	return &cLIScribeClient{cc}
}

func (c *cLIScribeClient) GetVersion(ctx context.Context, in *VersionRequest, opts ...grpc.CallOption) (*VersionResponse, error) {
	out := new(VersionResponse)
	err := c.cc.Invoke(ctx, "/com.romanostrechlis.scribe.api.CLIScribe/GetVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CLIScribeServer is the server API for CLIScribe service.
type CLIScribeServer interface {
	GetVersion(context.Context, *VersionRequest) (*VersionResponse, error)
}

func RegisterCLIScribeServer(s *grpc.Server, srv CLIScribeServer) {
	s.RegisterService(&_CLIScribe_serviceDesc, srv)
}

func _CLIScribe_GetVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CLIScribeServer).GetVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.romanostrechlis.scribe.api.CLIScribe/GetVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CLIScribeServer).GetVersion(ctx, req.(*VersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _CLIScribe_serviceDesc = grpc.ServiceDesc{
	ServiceName: "com.romanostrechlis.scribe.api.CLIScribe",
	HandlerType: (*CLIScribeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVersion",
			Handler:    _CLIScribe_GetVersion_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cliScribe.proto",
}

func init() { proto.RegisterFile("cliScribe.proto", fileDescriptor_cliScribe_41e2a648a039d414) }

var fileDescriptor_cliScribe_41e2a648a039d414 = []byte{
	// 162 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4f, 0xce, 0xc9, 0x0c,
	0x4e, 0x2e, 0xca, 0x4c, 0x4a, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x92, 0x4b, 0xce, 0xcf,
	0xd5, 0x2b, 0xca, 0xcf, 0x4d, 0xcc, 0xcb, 0x2f, 0x2e, 0x29, 0x4a, 0x4d, 0xce, 0xc8, 0xc9, 0x2c,
	0xd6, 0x2b, 0x86, 0xa8, 0x48, 0x2c, 0xc8, 0x54, 0x12, 0xe0, 0xe2, 0x0b, 0x4b, 0x2d, 0x2a, 0xce,
	0xcc, 0xcf, 0x0b, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x51, 0xd2, 0xe6, 0xe2, 0x87, 0x8b, 0x14,
	0x17, 0xe4, 0xe7, 0x15, 0xa7, 0x0a, 0x49, 0x70, 0xb1, 0x97, 0x41, 0x84, 0x24, 0x18, 0x15, 0x18,
	0x35, 0x38, 0x83, 0x60, 0x5c, 0xa3, 0x1a, 0x2e, 0x4e, 0x67, 0x1f, 0x4f, 0x88, 0x8d, 0x42, 0xf9,
	0x5c, 0x5c, 0xee, 0xa9, 0x25, 0x50, 0xcd, 0x42, 0x7a, 0x7a, 0xf8, 0xad, 0xd6, 0x43, 0xb5, 0x57,
	0x4a, 0x9f, 0x68, 0xf5, 0x10, 0x57, 0x29, 0x31, 0x38, 0xb1, 0x46, 0x31, 0x27, 0x16, 0x64, 0x26,
	0xb1, 0x81, 0xbd, 0x6a, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x30, 0x5e, 0x38, 0xb9, 0xfd, 0x00,
	0x00, 0x00,
}
