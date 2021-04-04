// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protos

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

// SdkMgrTelemetryServiceClient is the client API for SdkMgrTelemetryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SdkMgrTelemetryServiceClient interface {
	/// Add or update telemetry data
	TelemetryAddOrUpdate(ctx context.Context, in *TelemetryUpdateRequest, opts ...grpc.CallOption) (*TelemetryUpdateResponse, error)
	/// Delete telemetry data
	TelemetryDelete(ctx context.Context, in *TelemetryDeleteRequest, opts ...grpc.CallOption) (*TelemetryDeleteResponse, error)
}

type sdkMgrTelemetryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSdkMgrTelemetryServiceClient(cc grpc.ClientConnInterface) SdkMgrTelemetryServiceClient {
	return &sdkMgrTelemetryServiceClient{cc}
}

func (c *sdkMgrTelemetryServiceClient) TelemetryAddOrUpdate(ctx context.Context, in *TelemetryUpdateRequest, opts ...grpc.CallOption) (*TelemetryUpdateResponse, error) {
	out := new(TelemetryUpdateResponse)
	err := c.cc.Invoke(ctx, "/srlinux.sdk.SdkMgrTelemetryService/TelemetryAddOrUpdate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sdkMgrTelemetryServiceClient) TelemetryDelete(ctx context.Context, in *TelemetryDeleteRequest, opts ...grpc.CallOption) (*TelemetryDeleteResponse, error) {
	out := new(TelemetryDeleteResponse)
	err := c.cc.Invoke(ctx, "/srlinux.sdk.SdkMgrTelemetryService/TelemetryDelete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SdkMgrTelemetryServiceServer is the server API for SdkMgrTelemetryService service.
// All implementations must embed UnimplementedSdkMgrTelemetryServiceServer
// for forward compatibility
type SdkMgrTelemetryServiceServer interface {
	/// Add or update telemetry data
	TelemetryAddOrUpdate(context.Context, *TelemetryUpdateRequest) (*TelemetryUpdateResponse, error)
	/// Delete telemetry data
	TelemetryDelete(context.Context, *TelemetryDeleteRequest) (*TelemetryDeleteResponse, error)
	mustEmbedUnimplementedSdkMgrTelemetryServiceServer()
}

// UnimplementedSdkMgrTelemetryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSdkMgrTelemetryServiceServer struct {
}

func (UnimplementedSdkMgrTelemetryServiceServer) TelemetryAddOrUpdate(context.Context, *TelemetryUpdateRequest) (*TelemetryUpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TelemetryAddOrUpdate not implemented")
}
func (UnimplementedSdkMgrTelemetryServiceServer) TelemetryDelete(context.Context, *TelemetryDeleteRequest) (*TelemetryDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TelemetryDelete not implemented")
}
func (UnimplementedSdkMgrTelemetryServiceServer) mustEmbedUnimplementedSdkMgrTelemetryServiceServer() {
}

// UnsafeSdkMgrTelemetryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SdkMgrTelemetryServiceServer will
// result in compilation errors.
type UnsafeSdkMgrTelemetryServiceServer interface {
	mustEmbedUnimplementedSdkMgrTelemetryServiceServer()
}

func RegisterSdkMgrTelemetryServiceServer(s grpc.ServiceRegistrar, srv SdkMgrTelemetryServiceServer) {
	s.RegisterService(&SdkMgrTelemetryService_ServiceDesc, srv)
}

func _SdkMgrTelemetryService_TelemetryAddOrUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TelemetryUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SdkMgrTelemetryServiceServer).TelemetryAddOrUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/srlinux.sdk.SdkMgrTelemetryService/TelemetryAddOrUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SdkMgrTelemetryServiceServer).TelemetryAddOrUpdate(ctx, req.(*TelemetryUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SdkMgrTelemetryService_TelemetryDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TelemetryDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SdkMgrTelemetryServiceServer).TelemetryDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/srlinux.sdk.SdkMgrTelemetryService/TelemetryDelete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SdkMgrTelemetryServiceServer).TelemetryDelete(ctx, req.(*TelemetryDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SdkMgrTelemetryService_ServiceDesc is the grpc.ServiceDesc for SdkMgrTelemetryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SdkMgrTelemetryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "srlinux.sdk.SdkMgrTelemetryService",
	HandlerType: (*SdkMgrTelemetryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TelemetryAddOrUpdate",
			Handler:    _SdkMgrTelemetryService_TelemetryAddOrUpdate_Handler,
		},
		{
			MethodName: "TelemetryDelete",
			Handler:    _SdkMgrTelemetryService_TelemetryDelete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "telemetry_service.proto",
}
