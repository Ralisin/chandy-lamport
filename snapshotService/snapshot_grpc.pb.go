// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: snapshotService/snapshot.proto

package snapshotService

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

const (
	ServiceRegistry_RegisterPeer_FullMethodName = "/ServiceRegistry/RegisterPeer"
)

// ServiceRegistryClient is the client API for ServiceRegistry service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceRegistryClient interface {
	RegisterPeer(ctx context.Context, in *Peer, opts ...grpc.CallOption) (*RegisterPeerResponse, error)
}

type serviceRegistryClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceRegistryClient(cc grpc.ClientConnInterface) ServiceRegistryClient {
	return &serviceRegistryClient{cc}
}

func (c *serviceRegistryClient) RegisterPeer(ctx context.Context, in *Peer, opts ...grpc.CallOption) (*RegisterPeerResponse, error) {
	out := new(RegisterPeerResponse)
	err := c.cc.Invoke(ctx, ServiceRegistry_RegisterPeer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceRegistryServer is the server API for ServiceRegistry service.
// All implementations must embed UnimplementedServiceRegistryServer
// for forward compatibility
type ServiceRegistryServer interface {
	RegisterPeer(context.Context, *Peer) (*RegisterPeerResponse, error)
	mustEmbedUnimplementedServiceRegistryServer()
}

// UnimplementedServiceRegistryServer must be embedded to have forward compatible implementations.
type UnimplementedServiceRegistryServer struct {
}

func (UnimplementedServiceRegistryServer) RegisterPeer(context.Context, *Peer) (*RegisterPeerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterPeer not implemented")
}
func (UnimplementedServiceRegistryServer) mustEmbedUnimplementedServiceRegistryServer() {}

// UnsafeServiceRegistryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceRegistryServer will
// result in compilation errors.
type UnsafeServiceRegistryServer interface {
	mustEmbedUnimplementedServiceRegistryServer()
}

func RegisterServiceRegistryServer(s grpc.ServiceRegistrar, srv ServiceRegistryServer) {
	s.RegisterService(&ServiceRegistry_ServiceDesc, srv)
}

func _ServiceRegistry_RegisterPeer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Peer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceRegistryServer).RegisterPeer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ServiceRegistry_RegisterPeer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceRegistryServer).RegisterPeer(ctx, req.(*Peer))
	}
	return interceptor(ctx, in, info, handler)
}

// ServiceRegistry_ServiceDesc is the grpc.ServiceDesc for ServiceRegistry service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ServiceRegistry_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ServiceRegistry",
	HandlerType: (*ServiceRegistryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterPeer",
			Handler:    _ServiceRegistry_RegisterPeer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "snapshotService/snapshot.proto",
}

const (
	PeerFunction_NewPeerAdded_FullMethodName = "/PeerFunction/NewPeerAdded"
	PeerFunction_SendMessage_FullMethodName  = "/PeerFunction/SendMessage"
)

// PeerFunctionClient is the client API for PeerFunction service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PeerFunctionClient interface {
	NewPeerAdded(ctx context.Context, in *Peer, opts ...grpc.CallOption) (*Empty, error)
	SendMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Empty, error)
}

type peerFunctionClient struct {
	cc grpc.ClientConnInterface
}

func NewPeerFunctionClient(cc grpc.ClientConnInterface) PeerFunctionClient {
	return &peerFunctionClient{cc}
}

func (c *peerFunctionClient) NewPeerAdded(ctx context.Context, in *Peer, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, PeerFunction_NewPeerAdded_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *peerFunctionClient) SendMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, PeerFunction_SendMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PeerFunctionServer is the server API for PeerFunction service.
// All implementations must embed UnimplementedPeerFunctionServer
// for forward compatibility
type PeerFunctionServer interface {
	NewPeerAdded(context.Context, *Peer) (*Empty, error)
	SendMessage(context.Context, *Message) (*Empty, error)
	mustEmbedUnimplementedPeerFunctionServer()
}

// UnimplementedPeerFunctionServer must be embedded to have forward compatible implementations.
type UnimplementedPeerFunctionServer struct {
}

func (UnimplementedPeerFunctionServer) NewPeerAdded(context.Context, *Peer) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewPeerAdded not implemented")
}
func (UnimplementedPeerFunctionServer) SendMessage(context.Context, *Message) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedPeerFunctionServer) mustEmbedUnimplementedPeerFunctionServer() {}

// UnsafePeerFunctionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PeerFunctionServer will
// result in compilation errors.
type UnsafePeerFunctionServer interface {
	mustEmbedUnimplementedPeerFunctionServer()
}

func RegisterPeerFunctionServer(s grpc.ServiceRegistrar, srv PeerFunctionServer) {
	s.RegisterService(&PeerFunction_ServiceDesc, srv)
}

func _PeerFunction_NewPeerAdded_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Peer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PeerFunctionServer).NewPeerAdded(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PeerFunction_NewPeerAdded_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PeerFunctionServer).NewPeerAdded(ctx, req.(*Peer))
	}
	return interceptor(ctx, in, info, handler)
}

func _PeerFunction_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PeerFunctionServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PeerFunction_SendMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PeerFunctionServer).SendMessage(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

// PeerFunction_ServiceDesc is the grpc.ServiceDesc for PeerFunction service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PeerFunction_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "PeerFunction",
	HandlerType: (*PeerFunctionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewPeerAdded",
			Handler:    _PeerFunction_NewPeerAdded_Handler,
		},
		{
			MethodName: "SendMessage",
			Handler:    _PeerFunction_SendMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "snapshotService/snapshot.proto",
}
