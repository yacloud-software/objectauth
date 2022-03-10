// Code generated by protoc-gen-go.
// source: golang.conradwood.net/apis/publicproxy/publicproxy.proto
// DO NOT EDIT!

/*
Package publicproxy is a generated protocol buffer package.

It is generated from these files:
	golang.conradwood.net/apis/publicproxy/publicproxy.proto

It has these top-level messages:
	PingResponse
	Mapper
	ProxyTarget
*/
package publicproxy

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import common "golang.conradwood.net/apis/common"
import h2gproxy "golang.conradwood.net/apis/h2gproxy"

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

// comment: message pingresponse
type PingResponse struct {
	// comment: field pingresponse.response
	Response string `protobuf:"bytes,1,opt,name=Response" json:"Response,omitempty"`
}

func (m *PingResponse) Reset()                    { *m = PingResponse{} }
func (m *PingResponse) String() string            { return proto.CompactTextString(m) }
func (*PingResponse) ProtoMessage()               {}
func (*PingResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *PingResponse) GetResponse() string {
	if m != nil {
		return m.Response
	}
	return ""
}

// activate HOSTNAME to proxy for
type Mapper struct {
	ID       uint64 `protobuf:"varint,1,opt,name=ID" json:"ID,omitempty"`
	HostName string `protobuf:"bytes,2,opt,name=HostName" json:"HostName,omitempty"`
}

func (m *Mapper) Reset()                    { *m = Mapper{} }
func (m *Mapper) String() string            { return proto.CompactTextString(m) }
func (*Mapper) ProtoMessage()               {}
func (*Mapper) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Mapper) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *Mapper) GetHostName() string {
	if m != nil {
		return m.HostName
	}
	return ""
}

// this is the definition of where and how to create a request to a target
// it consists of "matches" - which determine wether or not this ProxyTarget is applicable to a specific request
// and "out" parameters - which are needed to generate the request to the target
type ProxyTarget struct {
	ID            uint64  `protobuf:"varint,1,opt,name=ID" json:"ID,omitempty"`
	Mapper        *Mapper `protobuf:"bytes,2,opt,name=Mapper" json:"Mapper,omitempty"`
	MatchPort     uint32  `protobuf:"varint,3,opt,name=MatchPort" json:"MatchPort,omitempty"`
	OutTargetHost string  `protobuf:"bytes,4,opt,name=OutTargetHost" json:"OutTargetHost,omitempty"`
	OutTargetPort uint32  `protobuf:"varint,5,opt,name=OutTargetPort" json:"OutTargetPort,omitempty"`
	OutTLS        bool    `protobuf:"varint,6,opt,name=OutTLS" json:"OutTLS,omitempty"`
}

func (m *ProxyTarget) Reset()                    { *m = ProxyTarget{} }
func (m *ProxyTarget) String() string            { return proto.CompactTextString(m) }
func (*ProxyTarget) ProtoMessage()               {}
func (*ProxyTarget) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ProxyTarget) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *ProxyTarget) GetMapper() *Mapper {
	if m != nil {
		return m.Mapper
	}
	return nil
}

func (m *ProxyTarget) GetMatchPort() uint32 {
	if m != nil {
		return m.MatchPort
	}
	return 0
}

func (m *ProxyTarget) GetOutTargetHost() string {
	if m != nil {
		return m.OutTargetHost
	}
	return ""
}

func (m *ProxyTarget) GetOutTargetPort() uint32 {
	if m != nil {
		return m.OutTargetPort
	}
	return 0
}

func (m *ProxyTarget) GetOutTLS() bool {
	if m != nil {
		return m.OutTLS
	}
	return false
}

func init() {
	proto.RegisterType((*PingResponse)(nil), "publicproxy.PingResponse")
	proto.RegisterType((*Mapper)(nil), "publicproxy.Mapper")
	proto.RegisterType((*ProxyTarget)(nil), "publicproxy.ProxyTarget")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for PublicProxy service

type PublicProxyClient interface {
	// comment: rpc ping
	Ping(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*PingResponse, error)
	StreamHTTP(ctx context.Context, in *h2gproxy.StreamRequest, opts ...grpc.CallOption) (PublicProxy_StreamHTTPClient, error)
	BiStreamHTTP(ctx context.Context, opts ...grpc.CallOption) (PublicProxy_BiStreamHTTPClient, error)
}

type publicProxyClient struct {
	cc *grpc.ClientConn
}

func NewPublicProxyClient(cc *grpc.ClientConn) PublicProxyClient {
	return &publicProxyClient{cc}
}

func (c *publicProxyClient) Ping(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := grpc.Invoke(ctx, "/publicproxy.PublicProxy/Ping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publicProxyClient) StreamHTTP(ctx context.Context, in *h2gproxy.StreamRequest, opts ...grpc.CallOption) (PublicProxy_StreamHTTPClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_PublicProxy_serviceDesc.Streams[0], c.cc, "/publicproxy.PublicProxy/StreamHTTP", opts...)
	if err != nil {
		return nil, err
	}
	x := &publicProxyStreamHTTPClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PublicProxy_StreamHTTPClient interface {
	Recv() (*h2gproxy.StreamDataResponse, error)
	grpc.ClientStream
}

type publicProxyStreamHTTPClient struct {
	grpc.ClientStream
}

func (x *publicProxyStreamHTTPClient) Recv() (*h2gproxy.StreamDataResponse, error) {
	m := new(h2gproxy.StreamDataResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *publicProxyClient) BiStreamHTTP(ctx context.Context, opts ...grpc.CallOption) (PublicProxy_BiStreamHTTPClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_PublicProxy_serviceDesc.Streams[1], c.cc, "/publicproxy.PublicProxy/BiStreamHTTP", opts...)
	if err != nil {
		return nil, err
	}
	x := &publicProxyBiStreamHTTPClient{stream}
	return x, nil
}

type PublicProxy_BiStreamHTTPClient interface {
	Send(*h2gproxy.StreamDataRequest) error
	Recv() (*h2gproxy.StreamDataResponse, error)
	grpc.ClientStream
}

type publicProxyBiStreamHTTPClient struct {
	grpc.ClientStream
}

func (x *publicProxyBiStreamHTTPClient) Send(m *h2gproxy.StreamDataRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *publicProxyBiStreamHTTPClient) Recv() (*h2gproxy.StreamDataResponse, error) {
	m := new(h2gproxy.StreamDataResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for PublicProxy service

type PublicProxyServer interface {
	// comment: rpc ping
	Ping(context.Context, *common.Void) (*PingResponse, error)
	StreamHTTP(*h2gproxy.StreamRequest, PublicProxy_StreamHTTPServer) error
	BiStreamHTTP(PublicProxy_BiStreamHTTPServer) error
}

func RegisterPublicProxyServer(s *grpc.Server, srv PublicProxyServer) {
	s.RegisterService(&_PublicProxy_serviceDesc, srv)
}

func _PublicProxy_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PublicProxyServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/publicproxy.PublicProxy/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PublicProxyServer).Ping(ctx, req.(*common.Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _PublicProxy_StreamHTTP_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(h2gproxy.StreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PublicProxyServer).StreamHTTP(m, &publicProxyStreamHTTPServer{stream})
}

type PublicProxy_StreamHTTPServer interface {
	Send(*h2gproxy.StreamDataResponse) error
	grpc.ServerStream
}

type publicProxyStreamHTTPServer struct {
	grpc.ServerStream
}

func (x *publicProxyStreamHTTPServer) Send(m *h2gproxy.StreamDataResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _PublicProxy_BiStreamHTTP_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PublicProxyServer).BiStreamHTTP(&publicProxyBiStreamHTTPServer{stream})
}

type PublicProxy_BiStreamHTTPServer interface {
	Send(*h2gproxy.StreamDataResponse) error
	Recv() (*h2gproxy.StreamDataRequest, error)
	grpc.ServerStream
}

type publicProxyBiStreamHTTPServer struct {
	grpc.ServerStream
}

func (x *publicProxyBiStreamHTTPServer) Send(m *h2gproxy.StreamDataResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *publicProxyBiStreamHTTPServer) Recv() (*h2gproxy.StreamDataRequest, error) {
	m := new(h2gproxy.StreamDataRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _PublicProxy_serviceDesc = grpc.ServiceDesc{
	ServiceName: "publicproxy.PublicProxy",
	HandlerType: (*PublicProxyServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _PublicProxy_Ping_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamHTTP",
			Handler:       _PublicProxy_StreamHTTP_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "BiStreamHTTP",
			Handler:       _PublicProxy_BiStreamHTTP_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "golang.conradwood.net/apis/publicproxy/publicproxy.proto",
}

func init() {
	proto.RegisterFile("golang.conradwood.net/apis/publicproxy/publicproxy.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 403 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x92, 0xb1, 0xce, 0xd3, 0x30,
	0x10, 0xc7, 0xe5, 0x12, 0xa2, 0xd4, 0xed, 0xc7, 0x60, 0x24, 0x9a, 0x84, 0x4a, 0x94, 0xaa, 0x43,
	0xc4, 0xe0, 0x56, 0x01, 0x21, 0x26, 0x86, 0x28, 0x48, 0xad, 0x44, 0x21, 0x72, 0x2b, 0x76, 0x37,
	0xb1, 0xd2, 0x48, 0x8d, 0x1d, 0x1c, 0x47, 0xc0, 0xda, 0xc7, 0x60, 0x64, 0xe4, 0x79, 0x78, 0x05,
	0xde, 0x03, 0xc5, 0x69, 0x4a, 0x02, 0xa8, 0x62, 0xf2, 0xf9, 0x7f, 0x77, 0x3f, 0xdf, 0xf9, 0x0e,
	0xbe, 0x4a, 0xc5, 0x89, 0xf2, 0x14, 0xc7, 0x82, 0x4b, 0x9a, 0x7c, 0x12, 0x22, 0xc1, 0x9c, 0xa9,
	0x25, 0x2d, 0xb2, 0x72, 0x59, 0x54, 0x87, 0x53, 0x16, 0x17, 0x52, 0x7c, 0xfe, 0xd2, 0xb5, 0x71,
	0x21, 0x85, 0x12, 0x68, 0xd4, 0x91, 0x5c, 0x7c, 0x03, 0x13, 0x8b, 0x3c, 0x17, 0xfc, 0x72, 0x34,
	0xc9, 0xae, 0x7f, 0x23, 0xfe, 0xe8, 0xa7, 0xcd, 0x9b, 0xad, 0xd1, 0xe4, 0xcc, 0x9f, 0xc1, 0x71,
	0x94, 0xf1, 0x94, 0xb0, 0xb2, 0x10, 0xbc, 0x64, 0xc8, 0x85, 0x56, 0x6b, 0xdb, 0x60, 0x06, 0xbc,
	0x21, 0xb9, 0xde, 0xe7, 0xaf, 0xa1, 0xb9, 0xa5, 0x45, 0xc1, 0x24, 0x7a, 0x00, 0x07, 0x9b, 0x50,
	0xfb, 0x0d, 0x32, 0xd8, 0x84, 0x68, 0x01, 0xad, 0xb5, 0x28, 0xd5, 0x3b, 0x9a, 0x33, 0x7b, 0x50,
	0x67, 0x05, 0xd6, 0xf7, 0xb3, 0x63, 0x28, 0x59, 0x31, 0x72, 0xf5, 0xcc, 0x7f, 0x02, 0x38, 0x8a,
	0xea, 0xb7, 0xf7, 0x54, 0xa6, 0x4c, 0xfd, 0x45, 0x09, 0x5b, 0xbe, 0x66, 0x8c, 0xfc, 0x87, 0xb8,
	0xfb, 0x41, 0x8d, 0x2b, 0x98, 0x7c, 0x3d, 0x3b, 0x66, 0x95, 0x71, 0xf5, 0xf2, 0xc5, 0xb7, 0xb3,
	0x33, 0xcc, 0xb5, 0x8a, 0xb3, 0x84, 0xb4, 0xb5, 0x4d, 0xe1, 0x70, 0x4b, 0x55, 0x7c, 0x8c, 0x84,
	0x54, 0xf6, 0xbd, 0x19, 0xf0, 0xee, 0xc8, 0x6f, 0x01, 0x2d, 0xe0, 0xdd, 0xfb, 0x4a, 0x35, 0x05,
	0xd4, 0x85, 0xd9, 0x86, 0x6e, 0xb2, 0x2f, 0xf6, 0xa2, 0x34, 0xe7, 0xbe, 0xe6, 0xf4, 0x45, 0xf4,
	0x08, 0x9a, 0xb5, 0xf0, 0x76, 0x67, 0x9b, 0x33, 0xe0, 0x59, 0xe4, 0x72, 0xf3, 0x7f, 0xd4, 0x7d,
	0xea, 0xca, 0x75, 0xb7, 0x68, 0x09, 0x8d, 0xfa, 0x8f, 0xd1, 0x18, 0x5f, 0xc6, 0xf5, 0x41, 0x64,
	0x89, 0xeb, 0xf4, 0xba, 0xeb, 0x0d, 0xe1, 0x0d, 0x84, 0x3b, 0x25, 0x19, 0xcd, 0xd7, 0xfb, 0x7d,
	0x84, 0x26, 0xf8, 0x3a, 0xb3, 0x46, 0x25, 0xec, 0x63, 0xc5, 0x4a, 0xe5, 0x4e, 0xff, 0x74, 0x84,
	0x54, 0xd1, 0x16, 0xb2, 0x02, 0x68, 0x0b, 0xc7, 0x41, 0xd6, 0x01, 0x3d, 0xfe, 0x77, 0xfc, 0x7f,
	0xc0, 0x3c, 0xb0, 0x02, 0xc1, 0x53, 0xf8, 0x84, 0x33, 0xd5, 0xdd, 0xae, 0x7a, 0xb3, 0xba, 0x5d,
	0x1c, 0x4c, 0xbd, 0x54, 0xcf, 0x7f, 0x05, 0x00, 0x00, 0xff, 0xff, 0xad, 0xa0, 0x45, 0xc9, 0x01,
	0x03, 0x00, 0x00,
}
