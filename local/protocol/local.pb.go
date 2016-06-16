// Code generated by protoc-gen-go.
// source: local.proto
// DO NOT EDIT!

/*
Package protocol is a generated protocol buffer package.

It is generated from these files:
	local.proto

It has these top-level messages:
	EchoRequest
	EchoReply
	Command
	ProcessStatus
	StatusQuery
*/
package protocol

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ProcessState int32

const (
	ProcessState_UNKNOWN   ProcessState = 0
	ProcessState_PENDING   ProcessState = 1
	ProcessState_RUNNING   ProcessState = 2
	ProcessState_COMPLETED ProcessState = 3
	ProcessState_FAILED    ProcessState = 4
)

var ProcessState_name = map[int32]string{
	0: "UNKNOWN",
	1: "PENDING",
	2: "RUNNING",
	3: "COMPLETED",
	4: "FAILED",
}
var ProcessState_value = map[string]int32{
	"UNKNOWN":   0,
	"PENDING":   1,
	"RUNNING":   2,
	"COMPLETED": 3,
	"FAILED":    4,
}

func (x ProcessState) String() string {
	return proto.EnumName(ProcessState_name, int32(x))
}
func (ProcessState) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type EchoRequest struct {
	Ping string `protobuf:"bytes,1,opt,name=ping" json:"ping,omitempty"`
}

func (m *EchoRequest) Reset()                    { *m = EchoRequest{} }
func (m *EchoRequest) String() string            { return proto.CompactTextString(m) }
func (*EchoRequest) ProtoMessage()               {}
func (*EchoRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type EchoReply struct {
	Pong string `protobuf:"bytes,1,opt,name=pong" json:"pong,omitempty"`
}

func (m *EchoReply) Reset()                    { *m = EchoReply{} }
func (m *EchoReply) String() string            { return proto.CompactTextString(m) }
func (*EchoReply) ProtoMessage()               {}
func (*EchoReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type Command struct {
	Argv []string          `protobuf:"bytes,1,rep,name=argv" json:"argv,omitempty"`
	Env  map[string]string `protobuf:"bytes,2,rep,name=env" json:"env,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	// Timeout in nanoseconds
	Timeout int64 `protobuf:"varint,3,opt,name=timeout" json:"timeout,omitempty"`
}

func (m *Command) Reset()                    { *m = Command{} }
func (m *Command) String() string            { return proto.CompactTextString(m) }
func (*Command) ProtoMessage()               {}
func (*Command) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Command) GetEnv() map[string]string {
	if m != nil {
		return m.Env
	}
	return nil
}

type ProcessStatus struct {
	RunId     string       `protobuf:"bytes,1,opt,name=run_id,json=runId" json:"run_id,omitempty"`
	State     ProcessState `protobuf:"varint,2,opt,name=state,enum=protocol.ProcessState" json:"state,omitempty"`
	StdoutRef string       `protobuf:"bytes,3,opt,name=stdout_ref,json=stdoutRef" json:"stdout_ref,omitempty"`
	StderrRef string       `protobuf:"bytes,4,opt,name=stderr_ref,json=stderrRef" json:"stderr_ref,omitempty"`
	ExitCode  int32        `protobuf:"varint,5,opt,name=exit_code,json=exitCode" json:"exit_code,omitempty"`
	Error     string       `protobuf:"bytes,6,opt,name=error" json:"error,omitempty"`
}

func (m *ProcessStatus) Reset()                    { *m = ProcessStatus{} }
func (m *ProcessStatus) String() string            { return proto.CompactTextString(m) }
func (*ProcessStatus) ProtoMessage()               {}
func (*ProcessStatus) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type StatusQuery struct {
	RunId string `protobuf:"bytes,1,opt,name=run_id,json=runId" json:"run_id,omitempty"`
}

func (m *StatusQuery) Reset()                    { *m = StatusQuery{} }
func (m *StatusQuery) String() string            { return proto.CompactTextString(m) }
func (*StatusQuery) ProtoMessage()               {}
func (*StatusQuery) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func init() {
	proto.RegisterType((*EchoRequest)(nil), "protocol.EchoRequest")
	proto.RegisterType((*EchoReply)(nil), "protocol.EchoReply")
	proto.RegisterType((*Command)(nil), "protocol.Command")
	proto.RegisterType((*ProcessStatus)(nil), "protocol.ProcessStatus")
	proto.RegisterType((*StatusQuery)(nil), "protocol.StatusQuery")
	proto.RegisterEnum("protocol.ProcessState", ProcessState_name, ProcessState_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for LocalScoot service

type LocalScootClient interface {
	Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoReply, error)
	Run(ctx context.Context, in *Command, opts ...grpc.CallOption) (*ProcessStatus, error)
	Status(ctx context.Context, in *StatusQuery, opts ...grpc.CallOption) (*ProcessStatus, error)
}

type localScootClient struct {
	cc *grpc.ClientConn
}

func NewLocalScootClient(cc *grpc.ClientConn) LocalScootClient {
	return &localScootClient{cc}
}

func (c *localScootClient) Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoReply, error) {
	out := new(EchoReply)
	err := grpc.Invoke(ctx, "/protocol.LocalScoot/Echo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *localScootClient) Run(ctx context.Context, in *Command, opts ...grpc.CallOption) (*ProcessStatus, error) {
	out := new(ProcessStatus)
	err := grpc.Invoke(ctx, "/protocol.LocalScoot/Run", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *localScootClient) Status(ctx context.Context, in *StatusQuery, opts ...grpc.CallOption) (*ProcessStatus, error) {
	out := new(ProcessStatus)
	err := grpc.Invoke(ctx, "/protocol.LocalScoot/Status", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for LocalScoot service

type LocalScootServer interface {
	Echo(context.Context, *EchoRequest) (*EchoReply, error)
	Run(context.Context, *Command) (*ProcessStatus, error)
	Status(context.Context, *StatusQuery) (*ProcessStatus, error)
}

func RegisterLocalScootServer(s *grpc.Server, srv LocalScootServer) {
	s.RegisterService(&_LocalScoot_serviceDesc, srv)
}

func _LocalScoot_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EchoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocalScootServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.LocalScoot/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocalScootServer).Echo(ctx, req.(*EchoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocalScoot_Run_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Command)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocalScootServer).Run(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.LocalScoot/Run",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocalScootServer).Run(ctx, req.(*Command))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocalScoot_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocalScootServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.LocalScoot/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocalScootServer).Status(ctx, req.(*StatusQuery))
	}
	return interceptor(ctx, in, info, handler)
}

var _LocalScoot_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protocol.LocalScoot",
	HandlerType: (*LocalScootServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _LocalScoot_Echo_Handler,
		},
		{
			MethodName: "Run",
			Handler:    _LocalScoot_Run_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _LocalScoot_Status_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("local.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 452 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x52, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0xae, 0xe3, 0xd8, 0x89, 0xc7, 0x14, 0x85, 0x81, 0x82, 0x15, 0x84, 0x00, 0x8b, 0x03, 0x42,
	0x28, 0x87, 0x14, 0x21, 0xd4, 0x1b, 0x4a, 0x0d, 0x8a, 0x08, 0x6e, 0xd8, 0x52, 0x71, 0x8c, 0x8c,
	0xbd, 0x2d, 0x11, 0x8e, 0x37, 0xac, 0x77, 0x23, 0xf2, 0x30, 0x3c, 0x05, 0xcf, 0xc0, 0x7b, 0x31,
	0x6b, 0xc7, 0xc4, 0x82, 0x8a, 0x93, 0xe7, 0xfb, 0x19, 0xef, 0xb7, 0xb3, 0x03, 0x7e, 0x2e, 0xd2,
	0x24, 0x1f, 0xad, 0xa5, 0x50, 0x02, 0xfb, 0xd5, 0x27, 0x15, 0x79, 0xf8, 0x18, 0xfc, 0x28, 0xfd,
	0x22, 0x18, 0xff, 0xa6, 0x79, 0xa9, 0x10, 0xa1, 0xbb, 0x5e, 0x16, 0x57, 0x81, 0xf5, 0xc8, 0x7a,
	0xea, 0xb1, 0xaa, 0x0e, 0x1f, 0x82, 0x57, 0x5b, 0xd6, 0xf9, 0xb6, 0x32, 0x88, 0x96, 0x81, 0xea,
	0xf0, 0x87, 0x05, 0xbd, 0x89, 0x58, 0xad, 0x92, 0x22, 0x33, 0x7a, 0x22, 0xaf, 0x36, 0xa4, 0xdb,
	0x46, 0x37, 0x35, 0x3e, 0x07, 0x9b, 0x17, 0x9b, 0xa0, 0x43, 0x94, 0x3f, 0x1e, 0x8e, 0x9a, 0xb3,
	0x47, 0xbb, 0x9e, 0x51, 0x54, 0x6c, 0xa2, 0x42, 0xc9, 0x2d, 0x33, 0x36, 0x0c, 0xa0, 0xa7, 0x96,
	0x2b, 0x2e, 0xb4, 0x0a, 0x6c, 0x3a, 0xc4, 0x66, 0x0d, 0x1c, 0xbe, 0x84, 0x7e, 0x63, 0xc5, 0x01,
	0xd8, 0x5f, 0xf9, 0x76, 0x17, 0xc3, 0x94, 0x78, 0x07, 0x9c, 0x4d, 0x92, 0x6b, 0x4e, 0xe7, 0x18,
	0xae, 0x06, 0x27, 0x9d, 0x57, 0x56, 0xf8, 0xcb, 0x82, 0xc3, 0xb9, 0x14, 0x29, 0x2f, 0xcb, 0x73,
	0x95, 0x28, 0x5d, 0xe2, 0x11, 0xb8, 0x52, 0x17, 0x8b, 0x65, 0xb6, 0xfb, 0x81, 0x43, 0x68, 0x9a,
	0x51, 0x50, 0xa7, 0x24, 0x43, 0xfd, 0x8b, 0x9b, 0xe3, 0xbb, 0xfb, 0xa8, 0xad, 0x76, 0xce, 0x6a,
	0x13, 0x3e, 0x00, 0x28, 0x55, 0x46, 0xc1, 0x16, 0x92, 0x5f, 0x56, 0x59, 0x3d, 0xe6, 0xd5, 0x0c,
	0xe3, 0x97, 0x3b, 0x99, 0x4b, 0x59, 0xc9, 0xdd, 0x3f, 0x32, 0x31, 0x46, 0xbe, 0x0f, 0x1e, 0xff,
	0xbe, 0x54, 0x8b, 0x54, 0x64, 0x3c, 0x70, 0x48, 0x75, 0x58, 0xdf, 0x10, 0x13, 0xc2, 0xe6, 0x2e,
	0x64, 0x13, 0x32, 0x70, 0xeb, 0x78, 0x15, 0x08, 0x9f, 0x80, 0x5f, 0xe7, 0xff, 0xa0, 0x39, 0x8d,
	0xe0, 0xfa, 0x4b, 0x3c, 0x9b, 0xc3, 0x8d, 0x76, 0x5a, 0xf4, 0xa1, 0x77, 0x11, 0xbf, 0x8b, 0xcf,
	0x3e, 0xc5, 0x83, 0x03, 0x03, 0xe6, 0x51, 0x7c, 0x3a, 0x8d, 0xdf, 0x0e, 0x2c, 0x03, 0xd8, 0x45,
	0x1c, 0x1b, 0xd0, 0xc1, 0x43, 0xf0, 0x26, 0x67, 0xef, 0xe7, 0xb3, 0xe8, 0x63, 0x74, 0x3a, 0xb0,
	0x11, 0xc0, 0x7d, 0xf3, 0x7a, 0x3a, 0xa3, 0xba, 0x3b, 0xfe, 0x69, 0x01, 0xcc, 0xcc, 0xf6, 0x9c,
	0xa7, 0x42, 0x28, 0x7c, 0x01, 0x5d, 0xb3, 0x0f, 0x78, 0xb4, 0x1f, 0x4f, 0x6b, 0x85, 0x86, 0xb7,
	0xff, 0xa6, 0x69, 0x6d, 0xc2, 0x03, 0x3c, 0x06, 0x9b, 0xe9, 0x02, 0x6f, 0xfd, 0xf3, 0xfc, 0xc3,
	0x7b, 0xd7, 0x8e, 0x59, 0x97, 0xd4, 0x74, 0x02, 0x6e, 0xf3, 0x62, 0x7b, 0x53, 0x6b, 0x06, 0xff,
	0xe9, 0xfd, 0xec, 0x56, 0xca, 0xf1, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0x2a, 0xc2, 0xbd, 0x06,
	0xf9, 0x02, 0x00, 0x00,
}
