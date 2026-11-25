package api

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetRecipeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetRecipeRequest) Reset() {
	*x = GetRecipeRequest{}
	mi := &file_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRecipeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRecipeRequest) ProtoMessage() {}

func (x *GetRecipeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetRecipeRequest) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetRecipeRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type GetRecipeResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Recipe        *Recipe                `protobuf:"bytes,1,opt,name=recipe,proto3" json:"recipe,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetRecipeResponse) Reset() {
	*x = GetRecipeResponse{}
	mi := &file_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRecipeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRecipeResponse) ProtoMessage() {}

func (x *GetRecipeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetRecipeResponse) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{1}
}

func (x *GetRecipeResponse) GetRecipe() *Recipe {
	if x != nil {
		return x.Recipe
	}
	return nil
}

type Recipe struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Ingredients   []*Ingredient          `protobuf:"bytes,3,rep,name=ingredients,proto3" json:"ingredients,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Recipe) Reset() {
	*x = Recipe{}
	mi := &file_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Recipe) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Recipe) ProtoMessage() {}

func (x *Recipe) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*Recipe) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{2}
}

func (x *Recipe) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Recipe) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Recipe) GetIngredients() []*Ingredient {
	if x != nil {
		return x.Ingredients
	}
	return nil
}

type Ingredient struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Unit          string                 `protobuf:"bytes,3,opt,name=unit,proto3" json:"unit,omitempty"`
	Quantity      int32                  `protobuf:"varint,4,opt,name=quantity,proto3" json:"quantity,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Ingredient) Reset() {
	*x = Ingredient{}
	mi := &file_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Ingredient) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ingredient) ProtoMessage() {}

func (x *Ingredient) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*Ingredient) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{3}
}

func (x *Ingredient) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Ingredient) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Ingredient) GetUnit() string {
	if x != nil {
		return x.Unit
	}
	return ""
}

func (x *Ingredient) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

var File_service_proto protoreflect.FileDescriptor

const file_service_proto_rawDesc = "" +
	"\n" +
	"\rservice.proto\x12\x04grpc\"\"\n" +
	"\x10GetRecipeRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\"9\n" +
	"\x11GetRecipeResponse\x12$\n" +
	"\x06recipe\x18\x01 \x01(\v2\f.grpc.RecipeR\x06recipe\"`\n" +
	"\x06Recipe\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x122\n" +
	"\vingredients\x18\x03 \x03(\v2\x10.grpc.IngredientR\vingredients\"`\n" +
	"\n" +
	"Ingredient\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12\x12\n" +
	"\x04unit\x18\x03 \x01(\tR\x04unit\x12\x1a\n" +
	"\bquantity\x18\x04 \x01(\x05R\bquantity2M\n" +
	"\rRecipeService\x12<\n" +
	"\tGetRecipe\x12\x16.grpc.GetRecipeRequest\x1a\x17.grpc.GetRecipeResponseB\x0fZ\rmyservice/apib\x06proto3"

var (
	file_service_proto_rawDescOnce sync.Once
	file_service_proto_rawDescData []byte
)

func file_service_proto_rawDescGZIP() []byte {
	file_service_proto_rawDescOnce.Do(func() {
		file_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_service_proto_rawDesc), len(file_service_proto_rawDesc)))
	})
	return file_service_proto_rawDescData
}

var file_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_service_proto_goTypes = []any{
	(*GetRecipeRequest)(nil),
	(*GetRecipeResponse)(nil),
	(*Recipe)(nil),
	(*Ingredient)(nil),
}
var file_service_proto_depIdxs = []int32{
	2,
	3,
	0,
	1,
	3,
	2,
	2,
	2,
	0,
}

func init() { file_service_proto_init() }
func file_service_proto_init() {
	if File_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_service_proto_rawDesc), len(file_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_proto_goTypes,
		DependencyIndexes: file_service_proto_depIdxs,
		MessageInfos:      file_service_proto_msgTypes,
	}.Build()
	File_service_proto = out.File
	file_service_proto_goTypes = nil
	file_service_proto_depIdxs = nil
}
