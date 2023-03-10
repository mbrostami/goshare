// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: api/grpc/pb/goshare.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GoShareClient is the sharing API for GoShare service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GoShareClient interface {
	Ping(ctx context.Context, in *PingMsg, opts ...grpc.CallOption) (*PongMsg, error)
	ShareInit(ctx context.Context, in *ShareInitRequest, opts ...grpc.CallOption) (*ShareInitResponse, error)
	Share(ctx context.Context, opts ...grpc.CallOption) (GoShare_ShareClient, error)
	ReceiveInit(ctx context.Context, in *ReceiveRequest, opts ...grpc.CallOption) (*ReceiveInitResponse, error)
	Receive(ctx context.Context, in *ReceiveRequest, opts ...grpc.CallOption) (GoShare_ReceiveClient, error)
}

type goShareClient struct {
	cc grpc.ClientConnInterface
}

func NewGoShareClient(cc grpc.ClientConnInterface) GoShareClient {
	return &goShareClient{cc}
}

func (c *goShareClient) Ping(ctx context.Context, in *PingMsg, opts ...grpc.CallOption) (*PongMsg, error) {
	out := new(PongMsg)
	err := c.cc.Invoke(ctx, "/GoShare.GoShare/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *goShareClient) ShareInit(ctx context.Context, in *ShareInitRequest, opts ...grpc.CallOption) (*ShareInitResponse, error) {
	out := new(ShareInitResponse)
	err := c.cc.Invoke(ctx, "/GoShare.GoShare/ShareInit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *goShareClient) Share(ctx context.Context, opts ...grpc.CallOption) (GoShare_ShareClient, error) {
	stream, err := c.cc.NewStream(ctx, &GoShare_ServiceDesc.Streams[0], "/GoShare.GoShare/Share", opts...)
	if err != nil {
		return nil, err
	}
	x := &goShareShareClient{stream}
	return x, nil
}

type GoShare_ShareClient interface {
	Send(*ShareRequest) error
	Recv() (*ShareResponse, error)
	grpc.ClientStream
}

type goShareShareClient struct {
	grpc.ClientStream
}

func (x *goShareShareClient) Send(m *ShareRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *goShareShareClient) Recv() (*ShareResponse, error) {
	m := new(ShareResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *goShareClient) ReceiveInit(ctx context.Context, in *ReceiveRequest, opts ...grpc.CallOption) (*ReceiveInitResponse, error) {
	out := new(ReceiveInitResponse)
	err := c.cc.Invoke(ctx, "/GoShare.GoShare/ReceiveInit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *goShareClient) Receive(ctx context.Context, in *ReceiveRequest, opts ...grpc.CallOption) (GoShare_ReceiveClient, error) {
	stream, err := c.cc.NewStream(ctx, &GoShare_ServiceDesc.Streams[1], "/GoShare.GoShare/Receive", opts...)
	if err != nil {
		return nil, err
	}
	x := &goShareReceiveClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GoShare_ReceiveClient interface {
	Recv() (*ReceiveResponse, error)
	grpc.ClientStream
}

type goShareReceiveClient struct {
	grpc.ClientStream
}

func (x *goShareReceiveClient) Recv() (*ReceiveResponse, error) {
	m := new(ReceiveResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GoShareServer is the server API for GoShare service.
// All implementations must embed UnimplementedGoShareServer
// for forward compatibility
type GoShareServer interface {
	Ping(context.Context, *PingMsg) (*PongMsg, error)
	ShareInit(context.Context, *ShareInitRequest) (*ShareInitResponse, error)
	Share(GoShare_ShareServer) error
	ReceiveInit(context.Context, *ReceiveRequest) (*ReceiveInitResponse, error)
	Receive(*ReceiveRequest, GoShare_ReceiveServer) error
	mustEmbedUnimplementedGoShareServer()
}

// UnimplementedGoShareServer must be embedded to have forward compatible implementations.
type UnimplementedGoShareServer struct {
}

func (UnimplementedGoShareServer) Ping(context.Context, *PingMsg) (*PongMsg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedGoShareServer) ShareInit(context.Context, *ShareInitRequest) (*ShareInitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShareInit not implemented")
}
func (UnimplementedGoShareServer) Share(GoShare_ShareServer) error {
	return status.Errorf(codes.Unimplemented, "method Share not implemented")
}
func (UnimplementedGoShareServer) ReceiveInit(context.Context, *ReceiveRequest) (*ReceiveInitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReceiveInit not implemented")
}
func (UnimplementedGoShareServer) Receive(*ReceiveRequest, GoShare_ReceiveServer) error {
	return status.Errorf(codes.Unimplemented, "method Receive not implemented")
}
func (UnimplementedGoShareServer) mustEmbedUnimplementedGoShareServer() {}

// UnsafeGoShareServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GoShareServer will
// result in compilation errors.
type UnsafeGoShareServer interface {
	mustEmbedUnimplementedGoShareServer()
}

func RegisterGoShareServer(s grpc.ServiceRegistrar, srv GoShareServer) {
	s.RegisterService(&GoShare_ServiceDesc, srv)
}

func _GoShare_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoShareServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/GoShare.GoShare/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoShareServer).Ping(ctx, req.(*PingMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _GoShare_ShareInit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShareInitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoShareServer).ShareInit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/GoShare.GoShare/ShareInit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoShareServer).ShareInit(ctx, req.(*ShareInitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GoShare_Share_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GoShareServer).Share(&goShareShareServer{stream})
}

type GoShare_ShareServer interface {
	Send(*ShareResponse) error
	Recv() (*ShareRequest, error)
	grpc.ServerStream
}

type goShareShareServer struct {
	grpc.ServerStream
}

func (x *goShareShareServer) Send(m *ShareResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *goShareShareServer) Recv() (*ShareRequest, error) {
	m := new(ShareRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _GoShare_ReceiveInit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReceiveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoShareServer).ReceiveInit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/GoShare.GoShare/ReceiveInit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoShareServer).ReceiveInit(ctx, req.(*ReceiveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GoShare_Receive_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ReceiveRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GoShareServer).Receive(m, &goShareReceiveServer{stream})
}

type GoShare_ReceiveServer interface {
	Send(*ReceiveResponse) error
	grpc.ServerStream
}

type goShareReceiveServer struct {
	grpc.ServerStream
}

func (x *goShareReceiveServer) Send(m *ReceiveResponse) error {
	return x.ServerStream.SendMsg(m)
}

// GoShare_ServiceDesc is the grpc.ServiceDesc for GoShare service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GoShare_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "GoShare.GoShare",
	HandlerType: (*GoShareServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _GoShare_Ping_Handler,
		},
		{
			MethodName: "ShareInit",
			Handler:    _GoShare_ShareInit_Handler,
		},
		{
			MethodName: "ReceiveInit",
			Handler:    _GoShare_ReceiveInit_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Share",
			Handler:       _GoShare_Share_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Receive",
			Handler:       _GoShare_Receive_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/grpc/pb/goshare.proto",
}
