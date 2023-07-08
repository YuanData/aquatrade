// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: service_aquatrade.proto

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

const (
	AquaTrade_CreateMember_FullMethodName = "/pb.AquaTrade/CreateMember"
	AquaTrade_LoginMember_FullMethodName  = "/pb.AquaTrade/LoginMember"
)

// AquaTradeClient is the client API for AquaTrade service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AquaTradeClient interface {
	CreateMember(ctx context.Context, in *CreateMemberRequest, opts ...grpc.CallOption) (*CreateMemberResponse, error)
	LoginMember(ctx context.Context, in *LoginMemberRequest, opts ...grpc.CallOption) (*LoginMemberResponse, error)
}

type aquaTradeClient struct {
	cc grpc.ClientConnInterface
}

func NewAquaTradeClient(cc grpc.ClientConnInterface) AquaTradeClient {
	return &aquaTradeClient{cc}
}

func (c *aquaTradeClient) CreateMember(ctx context.Context, in *CreateMemberRequest, opts ...grpc.CallOption) (*CreateMemberResponse, error) {
	out := new(CreateMemberResponse)
	err := c.cc.Invoke(ctx, AquaTrade_CreateMember_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aquaTradeClient) LoginMember(ctx context.Context, in *LoginMemberRequest, opts ...grpc.CallOption) (*LoginMemberResponse, error) {
	out := new(LoginMemberResponse)
	err := c.cc.Invoke(ctx, AquaTrade_LoginMember_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AquaTradeServer is the server API for AquaTrade service.
// All implementations must embed UnimplementedAquaTradeServer
// for forward compatibility
type AquaTradeServer interface {
	CreateMember(context.Context, *CreateMemberRequest) (*CreateMemberResponse, error)
	LoginMember(context.Context, *LoginMemberRequest) (*LoginMemberResponse, error)
	mustEmbedUnimplementedAquaTradeServer()
}

// UnimplementedAquaTradeServer must be embedded to have forward compatible implementations.
type UnimplementedAquaTradeServer struct {
}

func (UnimplementedAquaTradeServer) CreateMember(context.Context, *CreateMemberRequest) (*CreateMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMember not implemented")
}
func (UnimplementedAquaTradeServer) LoginMember(context.Context, *LoginMemberRequest) (*LoginMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginMember not implemented")
}
func (UnimplementedAquaTradeServer) mustEmbedUnimplementedAquaTradeServer() {}

// UnsafeAquaTradeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AquaTradeServer will
// result in compilation errors.
type UnsafeAquaTradeServer interface {
	mustEmbedUnimplementedAquaTradeServer()
}

func RegisterAquaTradeServer(s grpc.ServiceRegistrar, srv AquaTradeServer) {
	s.RegisterService(&AquaTrade_ServiceDesc, srv)
}

func _AquaTrade_CreateMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AquaTradeServer).CreateMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AquaTrade_CreateMember_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AquaTradeServer).CreateMember(ctx, req.(*CreateMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AquaTrade_LoginMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AquaTradeServer).LoginMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AquaTrade_LoginMember_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AquaTradeServer).LoginMember(ctx, req.(*LoginMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AquaTrade_ServiceDesc is the grpc.ServiceDesc for AquaTrade service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AquaTrade_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.AquaTrade",
	HandlerType: (*AquaTradeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateMember",
			Handler:    _AquaTrade_CreateMember_Handler,
		},
		{
			MethodName: "LoginMember",
			Handler:    _AquaTrade_LoginMember_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service_aquatrade.proto",
}
