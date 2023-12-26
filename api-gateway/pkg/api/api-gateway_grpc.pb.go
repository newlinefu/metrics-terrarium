// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: api-gateway.proto

package api

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

// MetricsGetterClient is the client API for MetricsGetter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetricsGetterClient interface {
	GetRawMetrics(ctx context.Context, in *RawMetricsRequestMessage, opts ...grpc.CallOption) (*MetricsResponse, error)
	GetPreparedMetrics(ctx context.Context, in *PreparedMetricsRequestMessage, opts ...grpc.CallOption) (*MetricsResponse, error)
}

type metricsGetterClient struct {
	cc grpc.ClientConnInterface
}

func NewMetricsGetterClient(cc grpc.ClientConnInterface) MetricsGetterClient {
	return &metricsGetterClient{cc}
}

func (c *metricsGetterClient) GetRawMetrics(ctx context.Context, in *RawMetricsRequestMessage, opts ...grpc.CallOption) (*MetricsResponse, error) {
	out := new(MetricsResponse)
	err := c.cc.Invoke(ctx, "/MetricsGetter/GetRawMetrics", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsGetterClient) GetPreparedMetrics(ctx context.Context, in *PreparedMetricsRequestMessage, opts ...grpc.CallOption) (*MetricsResponse, error) {
	out := new(MetricsResponse)
	err := c.cc.Invoke(ctx, "/MetricsGetter/GetPreparedMetrics", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetricsGetterServer is the server API for MetricsGetter service.
// All implementations must embed UnimplementedMetricsGetterServer
// for forward compatibility
type MetricsGetterServer interface {
	GetRawMetrics(context.Context, *RawMetricsRequestMessage) (*MetricsResponse, error)
	GetPreparedMetrics(context.Context, *PreparedMetricsRequestMessage) (*MetricsResponse, error)
	mustEmbedUnimplementedMetricsGetterServer()
}

// UnimplementedMetricsGetterServer must be embedded to have forward compatible implementations.
type UnimplementedMetricsGetterServer struct {
}

func (UnimplementedMetricsGetterServer) GetRawMetrics(context.Context, *RawMetricsRequestMessage) (*MetricsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRawMetrics not implemented")
}
func (UnimplementedMetricsGetterServer) GetPreparedMetrics(context.Context, *PreparedMetricsRequestMessage) (*MetricsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPreparedMetrics not implemented")
}
func (UnimplementedMetricsGetterServer) mustEmbedUnimplementedMetricsGetterServer() {}

// UnsafeMetricsGetterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricsGetterServer will
// result in compilation errors.
type UnsafeMetricsGetterServer interface {
	mustEmbedUnimplementedMetricsGetterServer()
}

func RegisterMetricsGetterServer(s grpc.ServiceRegistrar, srv MetricsGetterServer) {
	s.RegisterService(&MetricsGetter_ServiceDesc, srv)
}

func _MetricsGetter_GetRawMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RawMetricsRequestMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsGetterServer).GetRawMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MetricsGetter/GetRawMetrics",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsGetterServer).GetRawMetrics(ctx, req.(*RawMetricsRequestMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricsGetter_GetPreparedMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PreparedMetricsRequestMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsGetterServer).GetPreparedMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MetricsGetter/GetPreparedMetrics",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsGetterServer).GetPreparedMetrics(ctx, req.(*PreparedMetricsRequestMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// MetricsGetter_ServiceDesc is the grpc.ServiceDesc for MetricsGetter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MetricsGetter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "MetricsGetter",
	HandlerType: (*MetricsGetterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRawMetrics",
			Handler:    _MetricsGetter_GetRawMetrics_Handler,
		},
		{
			MethodName: "GetPreparedMetrics",
			Handler:    _MetricsGetter_GetPreparedMetrics_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api-gateway.proto",
}
