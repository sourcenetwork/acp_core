// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sourcenetwork/acp_core/playground/sandbox.proto

package playground

import (
	fmt "fmt"
	proto "github.com/cosmos/gogoproto/proto"
	errors "github.com/sourcenetwork/acp_core/pkg/errors"
	types "github.com/sourcenetwork/acp_core/pkg/types"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// SandboxRecord represents an instance of a sandbox
type SandboxRecord struct {
	// Handle is an opaque identifier to a sandbox
	Handle uint64 `protobuf:"varint,1,opt,name=handle,proto3" json:"handle,omitempty"`
	// name is a user given designation to a sandbox
	Name string       `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Data *SandboxData `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	// scratchpad acts as temporary storage for modifications in the sandbox state
	Scratchpad  *SandboxData `protobuf:"bytes,4,opt,name=scratchpad,proto3" json:"scratchpad,omitempty"`
	Ctx         *SandboxCtx  `protobuf:"bytes,5,opt,name=ctx,proto3" json:"ctx,omitempty"`
	Initialized bool         `protobuf:"varint,6,opt,name=initialized,proto3" json:"initialized,omitempty"`
}

func (m *SandboxRecord) Reset()         { *m = SandboxRecord{} }
func (m *SandboxRecord) String() string { return proto.CompactTextString(m) }
func (*SandboxRecord) ProtoMessage()    {}
func (*SandboxRecord) Descriptor() ([]byte, []int) {
	return fileDescriptor_96ae4e460be0ca31, []int{0}
}
func (m *SandboxRecord) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SandboxRecord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SandboxRecord.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SandboxRecord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SandboxRecord.Merge(m, src)
}
func (m *SandboxRecord) XXX_Size() int {
	return m.Size()
}
func (m *SandboxRecord) XXX_DiscardUnknown() {
	xxx_messageInfo_SandboxRecord.DiscardUnknown(m)
}

var xxx_messageInfo_SandboxRecord proto.InternalMessageInfo

func (m *SandboxRecord) GetHandle() uint64 {
	if m != nil {
		return m.Handle
	}
	return 0
}

func (m *SandboxRecord) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *SandboxRecord) GetData() *SandboxData {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *SandboxRecord) GetScratchpad() *SandboxData {
	if m != nil {
		return m.Scratchpad
	}
	return nil
}

func (m *SandboxRecord) GetCtx() *SandboxCtx {
	if m != nil {
		return m.Ctx
	}
	return nil
}

func (m *SandboxRecord) GetInitialized() bool {
	if m != nil {
		return m.Initialized
	}
	return false
}

// SandboxData encapsulates all the data necessary to create a Sandbox
type SandboxData struct {
	PolicyDefinition string `protobuf:"bytes,1,opt,name=policy_definition,json=policyDefinition,proto3" json:"policy_definition,omitempty"`
	Relationships    string `protobuf:"bytes,2,opt,name=relationships,proto3" json:"relationships,omitempty"`
	PolicyTheorem    string `protobuf:"bytes,3,opt,name=policy_theorem,json=policyTheorem,proto3" json:"policy_theorem,omitempty"`
}

func (m *SandboxData) Reset()         { *m = SandboxData{} }
func (m *SandboxData) String() string { return proto.CompactTextString(m) }
func (*SandboxData) ProtoMessage()    {}
func (*SandboxData) Descriptor() ([]byte, []int) {
	return fileDescriptor_96ae4e460be0ca31, []int{1}
}
func (m *SandboxData) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SandboxData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SandboxData.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SandboxData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SandboxData.Merge(m, src)
}
func (m *SandboxData) XXX_Size() int {
	return m.Size()
}
func (m *SandboxData) XXX_DiscardUnknown() {
	xxx_messageInfo_SandboxData.DiscardUnknown(m)
}

var xxx_messageInfo_SandboxData proto.InternalMessageInfo

func (m *SandboxData) GetPolicyDefinition() string {
	if m != nil {
		return m.PolicyDefinition
	}
	return ""
}

func (m *SandboxData) GetRelationships() string {
	if m != nil {
		return m.Relationships
	}
	return ""
}

func (m *SandboxData) GetPolicyTheorem() string {
	if m != nil {
		return m.PolicyTheorem
	}
	return ""
}

// SandboxCtx encapsulated all context data
// to execute an isolated theorem execution simulation
type SandboxCtx struct {
	Policy        *types.Policy         `protobuf:"bytes,1,opt,name=policy,proto3" json:"policy,omitempty"`
	Relationships []*types.Relationship `protobuf:"bytes,2,rep,name=relationships,proto3" json:"relationships,omitempty"`
	PolicyTheorem *types.PolicyTheorem  `protobuf:"bytes,3,opt,name=policy_theorem,json=policyTheorem,proto3" json:"policy_theorem,omitempty"`
}

func (m *SandboxCtx) Reset()         { *m = SandboxCtx{} }
func (m *SandboxCtx) String() string { return proto.CompactTextString(m) }
func (*SandboxCtx) ProtoMessage()    {}
func (*SandboxCtx) Descriptor() ([]byte, []int) {
	return fileDescriptor_96ae4e460be0ca31, []int{2}
}
func (m *SandboxCtx) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SandboxCtx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SandboxCtx.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SandboxCtx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SandboxCtx.Merge(m, src)
}
func (m *SandboxCtx) XXX_Size() int {
	return m.Size()
}
func (m *SandboxCtx) XXX_DiscardUnknown() {
	xxx_messageInfo_SandboxCtx.DiscardUnknown(m)
}

var xxx_messageInfo_SandboxCtx proto.InternalMessageInfo

func (m *SandboxCtx) GetPolicy() *types.Policy {
	if m != nil {
		return m.Policy
	}
	return nil
}

func (m *SandboxCtx) GetRelationships() []*types.Relationship {
	if m != nil {
		return m.Relationships
	}
	return nil
}

func (m *SandboxCtx) GetPolicyTheorem() *types.PolicyTheorem {
	if m != nil {
		return m.PolicyTheorem
	}
	return nil
}

type SandboxDataErrors struct {
	// policy_errors contains all errors encountered while
	// processing the given policy
	PolicyErrors []*errors.ParserMessage `protobuf:"bytes,1,rep,name=policy_errors,json=policyErrors,proto3" json:"policy_errors,omitempty"`
	// policy_errors contains all errors encountered while
	// processing the relationship set
	RelationshipsErrors []*errors.ParserMessage `protobuf:"bytes,2,rep,name=relationships_errors,json=relationshipsErrors,proto3" json:"relationships_errors,omitempty"`
	// policy_errors contains all errors encountered while
	// parsing the theorems
	TheoremsErrrors []*errors.ParserMessage `protobuf:"bytes,3,rep,name=theorems_errrors,json=theoremsErrrors,proto3" json:"theorems_errrors,omitempty"`
}

func (m *SandboxDataErrors) Reset()         { *m = SandboxDataErrors{} }
func (m *SandboxDataErrors) String() string { return proto.CompactTextString(m) }
func (*SandboxDataErrors) ProtoMessage()    {}
func (*SandboxDataErrors) Descriptor() ([]byte, []int) {
	return fileDescriptor_96ae4e460be0ca31, []int{3}
}
func (m *SandboxDataErrors) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SandboxDataErrors) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SandboxDataErrors.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SandboxDataErrors) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SandboxDataErrors.Merge(m, src)
}
func (m *SandboxDataErrors) XXX_Size() int {
	return m.Size()
}
func (m *SandboxDataErrors) XXX_DiscardUnknown() {
	xxx_messageInfo_SandboxDataErrors.DiscardUnknown(m)
}

var xxx_messageInfo_SandboxDataErrors proto.InternalMessageInfo

func (m *SandboxDataErrors) GetPolicyErrors() []*errors.ParserMessage {
	if m != nil {
		return m.PolicyErrors
	}
	return nil
}

func (m *SandboxDataErrors) GetRelationshipsErrors() []*errors.ParserMessage {
	if m != nil {
		return m.RelationshipsErrors
	}
	return nil
}

func (m *SandboxDataErrors) GetTheoremsErrrors() []*errors.ParserMessage {
	if m != nil {
		return m.TheoremsErrrors
	}
	return nil
}

func init() {
	proto.RegisterType((*SandboxRecord)(nil), "sourcenetwork.acp_core.playground.SandboxRecord")
	proto.RegisterType((*SandboxData)(nil), "sourcenetwork.acp_core.playground.SandboxData")
	proto.RegisterType((*SandboxCtx)(nil), "sourcenetwork.acp_core.playground.SandboxCtx")
	proto.RegisterType((*SandboxDataErrors)(nil), "sourcenetwork.acp_core.playground.SandboxDataErrors")
}

func init() {
	proto.RegisterFile("sourcenetwork/acp_core/playground/sandbox.proto", fileDescriptor_96ae4e460be0ca31)
}

var fileDescriptor_96ae4e460be0ca31 = []byte{
	// 520 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x94, 0xc1, 0x8e, 0xd3, 0x30,
	0x10, 0x86, 0xeb, 0xb6, 0x54, 0x74, 0x42, 0x61, 0xd7, 0x20, 0x14, 0xed, 0x21, 0x0a, 0xa5, 0x2b,
	0x05, 0x01, 0x09, 0x0a, 0x12, 0x57, 0xa4, 0x65, 0xf7, 0xb2, 0x02, 0xb4, 0x18, 0x24, 0x24, 0x2e,
	0x91, 0x9b, 0x98, 0x26, 0xda, 0x34, 0x8e, 0x1c, 0x57, 0x74, 0x79, 0x01, 0xae, 0x3c, 0x02, 0xcf,
	0xc0, 0x53, 0x70, 0xdc, 0x23, 0x27, 0x84, 0xda, 0x17, 0x41, 0xd8, 0x2e, 0x9b, 0x16, 0x02, 0xa8,
	0xb7, 0x64, 0xfc, 0xcf, 0x37, 0x33, 0xbf, 0x35, 0x86, 0xa0, 0xe2, 0x33, 0x11, 0xb3, 0x82, 0xc9,
	0x77, 0x5c, 0x9c, 0x06, 0x34, 0x2e, 0xa3, 0x98, 0x0b, 0x16, 0x94, 0x39, 0x3d, 0x9b, 0x08, 0x3e,
	0x2b, 0x92, 0xa0, 0xa2, 0x45, 0x32, 0xe6, 0x73, 0xbf, 0x14, 0x5c, 0x72, 0x7c, 0x6b, 0x2d, 0xc1,
	0x5f, 0x25, 0xf8, 0x17, 0x09, 0x7b, 0x61, 0x03, 0x93, 0x09, 0xc1, 0x45, 0x15, 0x94, 0x54, 0x54,
	0x4c, 0x44, 0x53, 0x56, 0x55, 0x74, 0xc2, 0x34, 0x76, 0xef, 0x4e, 0x43, 0x8e, 0x60, 0x39, 0x95,
	0x19, 0x2f, 0xaa, 0x34, 0x2b, 0x8d, 0xf4, 0x76, 0x53, 0xcb, 0x3c, 0xcf, 0xe2, 0x33, 0x23, 0x1a,
	0x35, 0x88, 0x64, 0xca, 0xb8, 0x60, 0x53, 0xad, 0x1a, 0x7e, 0x6e, 0xc3, 0xe0, 0xa5, 0x1e, 0x8f,
	0xb0, 0x98, 0x8b, 0x04, 0xdf, 0x84, 0x5e, 0x4a, 0x8b, 0x24, 0x67, 0x36, 0x72, 0x91, 0xd7, 0x25,
	0xe6, 0x0f, 0x63, 0xe8, 0x16, 0x74, 0xca, 0xec, 0xb6, 0x8b, 0xbc, 0x3e, 0x51, 0xdf, 0xf8, 0x00,
	0xba, 0x09, 0x95, 0xd4, 0xee, 0xb8, 0xc8, 0xb3, 0x42, 0xdf, 0xff, 0xa7, 0x33, 0xbe, 0xa9, 0x75,
	0x48, 0x25, 0x25, 0x2a, 0x17, 0x3f, 0x07, 0xa8, 0x62, 0x41, 0x65, 0x9c, 0x96, 0x34, 0xb1, 0xbb,
	0x5b, 0x91, 0x6a, 0x04, 0xfc, 0x18, 0x3a, 0xb1, 0x9c, 0xdb, 0x97, 0x14, 0xe8, 0xfe, 0xff, 0x83,
	0x9e, 0xc8, 0x39, 0xf9, 0x99, 0x89, 0x5d, 0xb0, 0xb2, 0x22, 0x93, 0x19, 0xcd, 0xb3, 0xf7, 0x2c,
	0xb1, 0x7b, 0x2e, 0xf2, 0x2e, 0x93, 0x7a, 0x68, 0xf8, 0x01, 0x81, 0x55, 0x2b, 0x8f, 0xef, 0xc2,
	0xae, 0xb6, 0x3e, 0x4a, 0xd8, 0x5b, 0x25, 0xe4, 0x85, 0x72, 0xaf, 0x4f, 0x76, 0xf4, 0xc1, 0xe1,
	0xaf, 0x38, 0x1e, 0xc1, 0xa0, 0x7e, 0xa5, 0x95, 0x31, 0x74, 0x3d, 0x88, 0xf7, 0xe1, 0xaa, 0x41,
	0x9a, 0xfb, 0x52, 0x1e, 0xf7, 0xc9, 0x40, 0x47, 0x5f, 0xe9, 0xe0, 0xf0, 0x1b, 0x02, 0xb8, 0xe8,
	0x1f, 0x3f, 0x82, 0x9e, 0x3e, 0x57, 0xd5, 0xad, 0xd0, 0x69, 0x1a, 0xff, 0x44, 0xa9, 0x88, 0x51,
	0xe3, 0xe3, 0xdf, 0x7b, 0xea, 0x78, 0x56, 0x38, 0x6a, 0x4a, 0x27, 0x35, 0xf1, 0x66, 0xe7, 0x4f,
	0xff, 0xd8, 0xb9, 0x15, 0xee, 0xff, 0xbd, 0x17, 0x33, 0xd1, 0xe6, 0x80, 0x9f, 0xda, 0xb0, 0x5b,
	0xb3, 0xfa, 0x48, 0x2d, 0x10, 0x7e, 0x01, 0x46, 0x16, 0xe9, 0x8d, 0xb2, 0x91, 0xea, 0xf7, 0x5e,
	0x53, 0x09, 0xad, 0xf2, 0x4f, 0xd4, 0xde, 0x3d, 0xd3, 0x6b, 0x47, 0xae, 0x68, 0x84, 0x41, 0x46,
	0x70, 0x63, 0x6d, 0x8e, 0x15, 0xb9, 0xbd, 0x05, 0xf9, 0xfa, 0x1a, 0xc9, 0x14, 0x78, 0x0d, 0x3b,
	0xc6, 0x10, 0xc5, 0x56, 0xf0, 0xce, 0x16, 0xf0, 0x6b, 0x2b, 0xca, 0x91, 0x86, 0x1c, 0x1c, 0x7f,
	0x59, 0x38, 0xe8, 0x7c, 0xe1, 0xa0, 0xef, 0x0b, 0x07, 0x7d, 0x5c, 0x3a, 0xad, 0xf3, 0xa5, 0xd3,
	0xfa, 0xba, 0x74, 0x5a, 0x6f, 0x1e, 0x4c, 0x32, 0x99, 0xce, 0xc6, 0x7e, 0xcc, 0xa7, 0x8d, 0xaf,
	0xdc, 0xe9, 0xa4, 0xf6, 0xd2, 0x8d, 0x7b, 0xea, 0x55, 0x78, 0xf8, 0x23, 0x00, 0x00, 0xff, 0xff,
	0xe8, 0x8a, 0x59, 0x8f, 0x15, 0x05, 0x00, 0x00,
}

func (m *SandboxRecord) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SandboxRecord) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SandboxRecord) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Initialized {
		i--
		if m.Initialized {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x30
	}
	if m.Ctx != nil {
		{
			size, err := m.Ctx.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSandbox(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.Scratchpad != nil {
		{
			size, err := m.Scratchpad.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSandbox(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.Data != nil {
		{
			size, err := m.Data.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSandbox(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintSandbox(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x12
	}
	if m.Handle != 0 {
		i = encodeVarintSandbox(dAtA, i, uint64(m.Handle))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *SandboxData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SandboxData) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SandboxData) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.PolicyTheorem) > 0 {
		i -= len(m.PolicyTheorem)
		copy(dAtA[i:], m.PolicyTheorem)
		i = encodeVarintSandbox(dAtA, i, uint64(len(m.PolicyTheorem)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Relationships) > 0 {
		i -= len(m.Relationships)
		copy(dAtA[i:], m.Relationships)
		i = encodeVarintSandbox(dAtA, i, uint64(len(m.Relationships)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.PolicyDefinition) > 0 {
		i -= len(m.PolicyDefinition)
		copy(dAtA[i:], m.PolicyDefinition)
		i = encodeVarintSandbox(dAtA, i, uint64(len(m.PolicyDefinition)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *SandboxCtx) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SandboxCtx) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SandboxCtx) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.PolicyTheorem != nil {
		{
			size, err := m.PolicyTheorem.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSandbox(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Relationships) > 0 {
		for iNdEx := len(m.Relationships) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Relationships[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSandbox(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Policy != nil {
		{
			size, err := m.Policy.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSandbox(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *SandboxDataErrors) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SandboxDataErrors) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SandboxDataErrors) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TheoremsErrrors) > 0 {
		for iNdEx := len(m.TheoremsErrrors) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.TheoremsErrrors[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSandbox(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.RelationshipsErrors) > 0 {
		for iNdEx := len(m.RelationshipsErrors) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.RelationshipsErrors[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSandbox(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.PolicyErrors) > 0 {
		for iNdEx := len(m.PolicyErrors) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.PolicyErrors[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSandbox(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintSandbox(dAtA []byte, offset int, v uint64) int {
	offset -= sovSandbox(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SandboxRecord) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Handle != 0 {
		n += 1 + sovSandbox(uint64(m.Handle))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovSandbox(uint64(l))
	}
	if m.Data != nil {
		l = m.Data.Size()
		n += 1 + l + sovSandbox(uint64(l))
	}
	if m.Scratchpad != nil {
		l = m.Scratchpad.Size()
		n += 1 + l + sovSandbox(uint64(l))
	}
	if m.Ctx != nil {
		l = m.Ctx.Size()
		n += 1 + l + sovSandbox(uint64(l))
	}
	if m.Initialized {
		n += 2
	}
	return n
}

func (m *SandboxData) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.PolicyDefinition)
	if l > 0 {
		n += 1 + l + sovSandbox(uint64(l))
	}
	l = len(m.Relationships)
	if l > 0 {
		n += 1 + l + sovSandbox(uint64(l))
	}
	l = len(m.PolicyTheorem)
	if l > 0 {
		n += 1 + l + sovSandbox(uint64(l))
	}
	return n
}

func (m *SandboxCtx) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Policy != nil {
		l = m.Policy.Size()
		n += 1 + l + sovSandbox(uint64(l))
	}
	if len(m.Relationships) > 0 {
		for _, e := range m.Relationships {
			l = e.Size()
			n += 1 + l + sovSandbox(uint64(l))
		}
	}
	if m.PolicyTheorem != nil {
		l = m.PolicyTheorem.Size()
		n += 1 + l + sovSandbox(uint64(l))
	}
	return n
}

func (m *SandboxDataErrors) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.PolicyErrors) > 0 {
		for _, e := range m.PolicyErrors {
			l = e.Size()
			n += 1 + l + sovSandbox(uint64(l))
		}
	}
	if len(m.RelationshipsErrors) > 0 {
		for _, e := range m.RelationshipsErrors {
			l = e.Size()
			n += 1 + l + sovSandbox(uint64(l))
		}
	}
	if len(m.TheoremsErrrors) > 0 {
		for _, e := range m.TheoremsErrrors {
			l = e.Size()
			n += 1 + l + sovSandbox(uint64(l))
		}
	}
	return n
}

func sovSandbox(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSandbox(x uint64) (n int) {
	return sovSandbox(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SandboxRecord) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSandbox
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SandboxRecord: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SandboxRecord: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Handle", wireType)
			}
			m.Handle = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSandbox
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Handle |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSandbox
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSandbox
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSandbox
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSandbox
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSandbox
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSandbox
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Data == nil {
				m.Data = &SandboxData{}
			}
			if err := m.Data.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Scratchpad", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSandbox
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSandbox
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSandbox
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Scratchpad == nil {
				m.Scratchpad = &SandboxData{}
			}
			if err := m.Scratchpad.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ctx", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSandbox
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSandbox
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSandbox
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Ctx == nil {
				m.Ctx = &SandboxCtx{}
			}
			if err := m.Ctx.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Initialized", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSandbox
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Initialized = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipSandbox(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSandbox
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SandboxData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSandbox
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SandboxData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SandboxData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PolicyDefinition", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSandbox
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSandbox
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSandbox
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PolicyDefinition = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Relationships", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSandbox
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSandbox
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSandbox
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Relationships = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PolicyTheorem", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSandbox
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSandbox
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSandbox
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PolicyTheorem = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSandbox(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSandbox
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SandboxCtx) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSandbox
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SandboxCtx: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SandboxCtx: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Policy", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSandbox
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSandbox
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSandbox
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Policy == nil {
				m.Policy = &types.Policy{}
			}
			if err := m.Policy.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Relationships", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSandbox
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSandbox
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSandbox
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Relationships = append(m.Relationships, &types.Relationship{})
			if err := m.Relationships[len(m.Relationships)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PolicyTheorem", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSandbox
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSandbox
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSandbox
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.PolicyTheorem == nil {
				m.PolicyTheorem = &types.PolicyTheorem{}
			}
			if err := m.PolicyTheorem.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSandbox(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSandbox
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SandboxDataErrors) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSandbox
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SandboxDataErrors: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SandboxDataErrors: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PolicyErrors", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSandbox
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSandbox
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSandbox
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PolicyErrors = append(m.PolicyErrors, &errors.ParserMessage{})
			if err := m.PolicyErrors[len(m.PolicyErrors)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RelationshipsErrors", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSandbox
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSandbox
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSandbox
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RelationshipsErrors = append(m.RelationshipsErrors, &errors.ParserMessage{})
			if err := m.RelationshipsErrors[len(m.RelationshipsErrors)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TheoremsErrrors", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSandbox
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSandbox
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSandbox
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TheoremsErrrors = append(m.TheoremsErrrors, &errors.ParserMessage{})
			if err := m.TheoremsErrrors[len(m.TheoremsErrrors)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSandbox(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSandbox
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipSandbox(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSandbox
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSandbox
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSandbox
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthSandbox
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSandbox
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSandbox
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSandbox        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSandbox          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSandbox = fmt.Errorf("proto: unexpected end of group")
)
