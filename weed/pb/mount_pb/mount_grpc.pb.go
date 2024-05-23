// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: mount.proto

package mount_pb

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
	SeaweedMount_Configure_FullMethodName = "/messaging_pb.SeaweedMount/Configure"
)

// SeaweedMountClient is the client API for SeaweedMount service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SeaweedMountClient interface {
	Configure(ctx context.Context, in *ConfigureRequest, opts ...grpc.CallOption) (*ConfigureResponse, error)
}

type seaweedMountClient struct {
	cc grpc.ClientConnInterface
}

func NewSeaweedMountClient(cc grpc.ClientConnInterface) SeaweedMountClient {
	return &seaweedMountClient{cc}
}

func (c *seaweedMountClient) Configure(ctx context.Context, in *ConfigureRequest, opts ...grpc.CallOption) (*ConfigureResponse, error) {
	out := new(ConfigureResponse)
	err := c.cc.Invoke(ctx, SeaweedMount_Configure_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SeaweedMountServer is the server API for SeaweedMount service.
// All implementations must embed UnimplementedSeaweedMountServer
// for forward compatibility
type SeaweedMountServer interface {
	Configure(context.Context, *ConfigureRequest) (*ConfigureResponse, error)
	mustEmbedUnimplementedSeaweedMountServer()
}

// UnimplementedSeaweedMountServer must be embedded to have forward compatible implementations.
type UnimplementedSeaweedMountServer struct {
}

func (UnimplementedSeaweedMountServer) Configure(context.Context, *ConfigureRequest) (*ConfigureResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Configure not implemented")
}
func (UnimplementedSeaweedMountServer) mustEmbedUnimplementedSeaweedMountServer() {}

// UnsafeSeaweedMountServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SeaweedMountServer will
// result in compilation errors.
type UnsafeSeaweedMountServer interface {
	mustEmbedUnimplementedSeaweedMountServer()
}

func RegisterSeaweedMountServer(s grpc.ServiceRegistrar, srv SeaweedMountServer) {
	s.RegisterService(&SeaweedMount_ServiceDesc, srv)
}

func _SeaweedMount_Configure_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedMountServer).Configure(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SeaweedMount_Configure_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedMountServer).Configure(ctx, req.(*ConfigureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SeaweedMount_ServiceDesc is the grpc.ServiceDesc for SeaweedMount service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SeaweedMount_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "messaging_pb.SeaweedMount",
	HandlerType: (*SeaweedMountServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Configure",
			Handler:    _SeaweedMount_Configure_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mount.proto",
}
