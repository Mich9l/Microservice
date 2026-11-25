package api

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const _ = grpc.SupportPackageIsVersion9

const (
	RecipeService_GetRecipe_FullMethodName = "/grpc.RecipeService/GetRecipe"
)

type RecipeServiceClient interface {
	GetRecipe(ctx context.Context, in *GetRecipeRequest, opts ...grpc.CallOption) (*GetRecipeResponse, error)
}

type recipeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRecipeServiceClient(cc grpc.ClientConnInterface) RecipeServiceClient {
	return &recipeServiceClient{cc}
}

func (c *recipeServiceClient) GetRecipe(ctx context.Context, in *GetRecipeRequest, opts ...grpc.CallOption) (*GetRecipeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRecipeResponse)
	err := c.cc.Invoke(ctx, RecipeService_GetRecipe_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type RecipeServiceServer interface {
	GetRecipe(context.Context, *GetRecipeRequest) (*GetRecipeResponse, error)
	mustEmbedUnimplementedRecipeServiceServer()
}

type UnimplementedRecipeServiceServer struct{}

func (UnimplementedRecipeServiceServer) GetRecipe(context.Context, *GetRecipeRequest) (*GetRecipeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetRecipe not implemented")
}
func (UnimplementedRecipeServiceServer) mustEmbedUnimplementedRecipeServiceServer() {}
func (UnimplementedRecipeServiceServer) testEmbeddedByValue()                       {}

type UnsafeRecipeServiceServer interface {
	mustEmbedUnimplementedRecipeServiceServer()
}

func RegisterRecipeServiceServer(s grpc.ServiceRegistrar, srv RecipeServiceServer) {

	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RecipeService_ServiceDesc, srv)
}

func _RecipeService_GetRecipe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRecipeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecipeServiceServer).GetRecipe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecipeService_GetRecipe_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecipeServiceServer).GetRecipe(ctx, req.(*GetRecipeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var RecipeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.RecipeService",
	HandlerType: (*RecipeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRecipe",
			Handler:    _RecipeService_GetRecipe_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}
