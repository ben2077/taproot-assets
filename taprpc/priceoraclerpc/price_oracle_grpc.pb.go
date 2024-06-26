// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package priceoraclerpc

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

// PriceOracleClient is the client API for PriceOracle service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PriceOracleClient interface {
	// QueryRateTick queries the rate tick for a given transaction type, subject
	// asset, and payment asset. The rate tick is the exchange rate between the
	// subject asset and the payment asset.
	QueryRateTick(ctx context.Context, in *QueryRateTickRequest, opts ...grpc.CallOption) (*QueryRateTickResponse, error)
}

type priceOracleClient struct {
	cc grpc.ClientConnInterface
}

func NewPriceOracleClient(cc grpc.ClientConnInterface) PriceOracleClient {
	return &priceOracleClient{cc}
}

func (c *priceOracleClient) QueryRateTick(ctx context.Context, in *QueryRateTickRequest, opts ...grpc.CallOption) (*QueryRateTickResponse, error) {
	out := new(QueryRateTickResponse)
	err := c.cc.Invoke(ctx, "/priceoraclerpc.PriceOracle/QueryRateTick", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PriceOracleServer is the server API for PriceOracle service.
// All implementations must embed UnimplementedPriceOracleServer
// for forward compatibility
type PriceOracleServer interface {
	// QueryRateTick queries the rate tick for a given transaction type, subject
	// asset, and payment asset. The rate tick is the exchange rate between the
	// subject asset and the payment asset.
	QueryRateTick(context.Context, *QueryRateTickRequest) (*QueryRateTickResponse, error)
	mustEmbedUnimplementedPriceOracleServer()
}

// UnimplementedPriceOracleServer must be embedded to have forward compatible implementations.
type UnimplementedPriceOracleServer struct {
}

func (UnimplementedPriceOracleServer) QueryRateTick(context.Context, *QueryRateTickRequest) (*QueryRateTickResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryRateTick not implemented")
}
func (UnimplementedPriceOracleServer) mustEmbedUnimplementedPriceOracleServer() {}

// UnsafePriceOracleServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PriceOracleServer will
// result in compilation errors.
type UnsafePriceOracleServer interface {
	mustEmbedUnimplementedPriceOracleServer()
}

func RegisterPriceOracleServer(s grpc.ServiceRegistrar, srv PriceOracleServer) {
	s.RegisterService(&PriceOracle_ServiceDesc, srv)
}

func _PriceOracle_QueryRateTick_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRateTickRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PriceOracleServer).QueryRateTick(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/priceoraclerpc.PriceOracle/QueryRateTick",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PriceOracleServer).QueryRateTick(ctx, req.(*QueryRateTickRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PriceOracle_ServiceDesc is the grpc.ServiceDesc for PriceOracle service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PriceOracle_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "priceoraclerpc.PriceOracle",
	HandlerType: (*PriceOracleServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "QueryRateTick",
			Handler:    _PriceOracle_QueryRateTick_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "priceoraclerpc/price_oracle.proto",
}
