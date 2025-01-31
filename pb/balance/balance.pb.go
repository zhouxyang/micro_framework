// Code generated by protoc-gen-go. DO NOT EDIT.
// source: balance.proto

//protoc ./balance.proto --go_out=plugins=grpc:../pb

package balance

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// DeductrRequest
type DeductRequest struct {
	UserID               string   `protobuf:"bytes,1,opt,name=userID,proto3" json:"userID,omitempty"`
	ProductID            string   `protobuf:"bytes,2,opt,name=productID,proto3" json:"productID,omitempty"`
	Price                string   `protobuf:"bytes,3,opt,name=price,proto3" json:"price,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeductRequest) Reset()         { *m = DeductRequest{} }
func (m *DeductRequest) String() string { return proto.CompactTextString(m) }
func (*DeductRequest) ProtoMessage()    {}
func (*DeductRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee25a00b628521b1, []int{0}
}

func (m *DeductRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeductRequest.Unmarshal(m, b)
}
func (m *DeductRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeductRequest.Marshal(b, m, deterministic)
}
func (m *DeductRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeductRequest.Merge(m, src)
}
func (m *DeductRequest) XXX_Size() int {
	return xxx_messageInfo_DeductRequest.Size(m)
}
func (m *DeductRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeductRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeductRequest proto.InternalMessageInfo

func (m *DeductRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *DeductRequest) GetProductID() string {
	if m != nil {
		return m.ProductID
	}
	return ""
}

func (m *DeductRequest) GetPrice() string {
	if m != nil {
		return m.Price
	}
	return ""
}

type DeductResponse struct {
	Balance              string   `protobuf:"bytes,1,opt,name=balance,proto3" json:"balance,omitempty"`
	Result               *Result  `protobuf:"bytes,2,opt,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeductResponse) Reset()         { *m = DeductResponse{} }
func (m *DeductResponse) String() string { return proto.CompactTextString(m) }
func (*DeductResponse) ProtoMessage()    {}
func (*DeductResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee25a00b628521b1, []int{1}
}

func (m *DeductResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeductResponse.Unmarshal(m, b)
}
func (m *DeductResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeductResponse.Marshal(b, m, deterministic)
}
func (m *DeductResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeductResponse.Merge(m, src)
}
func (m *DeductResponse) XXX_Size() int {
	return xxx_messageInfo_DeductResponse.Size(m)
}
func (m *DeductResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeductResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeductResponse proto.InternalMessageInfo

func (m *DeductResponse) GetBalance() string {
	if m != nil {
		return m.Balance
	}
	return ""
}

func (m *DeductResponse) GetResult() *Result {
	if m != nil {
		return m.Result
	}
	return nil
}

// 返回码与返回状态
type Result struct {
	Code                 int64    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Result) Reset()         { *m = Result{} }
func (m *Result) String() string { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()    {}
func (*Result) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee25a00b628521b1, []int{2}
}

func (m *Result) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Result.Unmarshal(m, b)
}
func (m *Result) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Result.Marshal(b, m, deterministic)
}
func (m *Result) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Result.Merge(m, src)
}
func (m *Result) XXX_Size() int {
	return xxx_messageInfo_Result.Size(m)
}
func (m *Result) XXX_DiscardUnknown() {
	xxx_messageInfo_Result.DiscardUnknown(m)
}

var xxx_messageInfo_Result proto.InternalMessageInfo

func (m *Result) GetCode() int64 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *Result) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func init() {
	proto.RegisterType((*DeductRequest)(nil), "balance.DeductRequest")
	proto.RegisterType((*DeductResponse)(nil), "balance.DeductResponse")
	proto.RegisterType((*Result)(nil), "balance.Result")
}

func init() { proto.RegisterFile("balance.proto", fileDescriptor_ee25a00b628521b1) }

var fileDescriptor_ee25a00b628521b1 = []byte{
	// 216 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x90, 0x3b, 0x4f, 0xc4, 0x30,
	0x10, 0x84, 0x09, 0x01, 0x9f, 0x6e, 0xd1, 0x01, 0x5a, 0xa1, 0xc3, 0x42, 0x14, 0xc8, 0x0d, 0x57,
	0xa5, 0x08, 0x35, 0x0d, 0x4a, 0x41, 0x5a, 0x53, 0x52, 0x25, 0xce, 0x0a, 0x21, 0x85, 0xd8, 0xf8,
	0xf1, 0xff, 0x51, 0xfc, 0x00, 0xa1, 0xeb, 0x76, 0x66, 0xac, 0x6f, 0xc7, 0x0b, 0xbb, 0x71, 0x98,
	0x87, 0x45, 0x51, 0x63, 0xac, 0xf6, 0x1a, 0x37, 0x59, 0x8a, 0x77, 0xd8, 0x75, 0x34, 0x05, 0xe5,
	0x25, 0x7d, 0x07, 0x72, 0x1e, 0xf7, 0xc0, 0x82, 0x23, 0xdb, 0x77, 0xbc, 0x7a, 0xa8, 0x0e, 0x5b,
	0x99, 0x15, 0xde, 0xc3, 0xd6, 0x58, 0xbd, 0xbe, 0xec, 0x3b, 0x7e, 0x1a, 0xa3, 0x3f, 0x03, 0x6f,
	0xe0, 0xdc, 0xd8, 0x4f, 0x45, 0xbc, 0x8e, 0x49, 0x12, 0xe2, 0x0d, 0x2e, 0x0b, 0xdc, 0x19, 0xbd,
	0x38, 0x42, 0x0e, 0x65, 0x73, 0xc6, 0x17, 0x89, 0x8f, 0xc0, 0x2c, 0xb9, 0x30, 0xfb, 0x08, 0xbf,
	0x68, 0xaf, 0x9a, 0xd2, 0x58, 0x46, 0x5b, 0xe6, 0x58, 0x34, 0xc0, 0x92, 0x83, 0x08, 0x67, 0x4a,
	0x4f, 0x89, 0x54, 0xcb, 0x38, 0xe3, 0x35, 0xd4, 0x5f, 0xee, 0x23, 0x17, 0x5c, 0xc7, 0xf6, 0x15,
	0x36, 0x2f, 0x79, 0xc7, 0x33, 0xb0, 0xd4, 0x07, 0xf7, 0xbf, 0xf4, 0x7f, 0xbf, 0xbf, 0xbb, 0x3d,
	0xf2, 0x53, 0x71, 0x71, 0x72, 0xa8, 0x46, 0x16, 0x6f, 0xf7, 0xf4, 0x13, 0x00, 0x00, 0xff, 0xff,
	0xc6, 0xaa, 0x4d, 0xd9, 0x4c, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// BalanceClient is the client API for Balance service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BalanceClient interface {
	// Deduct
	Deduct(ctx context.Context, opts ...grpc.CallOption) (Balance_DeductClient, error)
}

type balanceClient struct {
	cc *grpc.ClientConn
}

func NewBalanceClient(cc *grpc.ClientConn) BalanceClient {
	return &balanceClient{cc}
}

func (c *balanceClient) Deduct(ctx context.Context, opts ...grpc.CallOption) (Balance_DeductClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Balance_serviceDesc.Streams[0], "/balance.Balance/Deduct", opts...)
	if err != nil {
		return nil, err
	}
	x := &balanceDeductClient{stream}
	return x, nil
}

type Balance_DeductClient interface {
	Send(*DeductRequest) error
	CloseAndRecv() (*DeductResponse, error)
	grpc.ClientStream
}

type balanceDeductClient struct {
	grpc.ClientStream
}

func (x *balanceDeductClient) Send(m *DeductRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *balanceDeductClient) CloseAndRecv() (*DeductResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(DeductResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BalanceServer is the server API for Balance service.
type BalanceServer interface {
	// Deduct
	Deduct(Balance_DeductServer) error
}

// UnimplementedBalanceServer can be embedded to have forward compatible implementations.
type UnimplementedBalanceServer struct {
}

func (*UnimplementedBalanceServer) Deduct(srv Balance_DeductServer) error {
	return status.Errorf(codes.Unimplemented, "method Deduct not implemented")
}

func RegisterBalanceServer(s *grpc.Server, srv BalanceServer) {
	s.RegisterService(&_Balance_serviceDesc, srv)
}

func _Balance_Deduct_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BalanceServer).Deduct(&balanceDeductServer{stream})
}

type Balance_DeductServer interface {
	SendAndClose(*DeductResponse) error
	Recv() (*DeductRequest, error)
	grpc.ServerStream
}

type balanceDeductServer struct {
	grpc.ServerStream
}

func (x *balanceDeductServer) SendAndClose(m *DeductResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *balanceDeductServer) Recv() (*DeductRequest, error) {
	m := new(DeductRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Balance_serviceDesc = grpc.ServiceDesc{
	ServiceName: "balance.Balance",
	HandlerType: (*BalanceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Deduct",
			Handler:       _Balance_Deduct_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "balance.proto",
}
