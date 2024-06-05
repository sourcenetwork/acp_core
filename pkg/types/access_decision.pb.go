// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sourcenetwork/acp_core/access_decision.proto

package types

import (
	fmt "fmt"
	proto "github.com/cosmos/gogoproto/proto"
	types "github.com/cosmos/gogoproto/types"
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

// AccessDecision models the result of evaluating a set of AccessRequests for an Actor
type AccessDecision struct {
	Id                 string           `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	PolicyId           string           `protobuf:"bytes,2,opt,name=policy_id,json=policyId,proto3" json:"policy_id,omitempty"`
	Creator            string           `protobuf:"bytes,3,opt,name=creator,proto3" json:"creator,omitempty"`
	CreatorAccSequence uint64           `protobuf:"varint,4,opt,name=creator_acc_sequence,json=creatorAccSequence,proto3" json:"creator_acc_sequence,omitempty"`
	Operations         []*Operation     `protobuf:"bytes,5,rep,name=operations,proto3" json:"operations,omitempty"`
	Actor              string           `protobuf:"bytes,6,opt,name=actor,proto3" json:"actor,omitempty"`
	Params             *DecisionParams  `protobuf:"bytes,7,opt,name=params,proto3" json:"params,omitempty"`
	CreationTime       *types.Timestamp `protobuf:"bytes,8,opt,name=creation_time,json=creationTime,proto3" json:"creation_time,omitempty"`
	// issued_height stores the block height when the Decision was evaluated
	IssuedHeight uint64 `protobuf:"varint,9,opt,name=issued_height,json=issuedHeight,proto3" json:"issued_height,omitempty"`
}

func (m *AccessDecision) Reset()         { *m = AccessDecision{} }
func (m *AccessDecision) String() string { return proto.CompactTextString(m) }
func (*AccessDecision) ProtoMessage()    {}
func (*AccessDecision) Descriptor() ([]byte, []int) {
	return fileDescriptor_b2650781da68a114, []int{0}
}
func (m *AccessDecision) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AccessDecision) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AccessDecision.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AccessDecision) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccessDecision.Merge(m, src)
}
func (m *AccessDecision) XXX_Size() int {
	return m.Size()
}
func (m *AccessDecision) XXX_DiscardUnknown() {
	xxx_messageInfo_AccessDecision.DiscardUnknown(m)
}

var xxx_messageInfo_AccessDecision proto.InternalMessageInfo

func (m *AccessDecision) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *AccessDecision) GetPolicyId() string {
	if m != nil {
		return m.PolicyId
	}
	return ""
}

func (m *AccessDecision) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *AccessDecision) GetCreatorAccSequence() uint64 {
	if m != nil {
		return m.CreatorAccSequence
	}
	return 0
}

func (m *AccessDecision) GetOperations() []*Operation {
	if m != nil {
		return m.Operations
	}
	return nil
}

func (m *AccessDecision) GetActor() string {
	if m != nil {
		return m.Actor
	}
	return ""
}

func (m *AccessDecision) GetParams() *DecisionParams {
	if m != nil {
		return m.Params
	}
	return nil
}

func (m *AccessDecision) GetCreationTime() *types.Timestamp {
	if m != nil {
		return m.CreationTime
	}
	return nil
}

func (m *AccessDecision) GetIssuedHeight() uint64 {
	if m != nil {
		return m.IssuedHeight
	}
	return 0
}

// DecisionParams stores auxiliary information regarding the validity of a decision
type DecisionParams struct {
	// number of blocks a Decision is valid for
	DecisionExpirationDelta uint64 `protobuf:"varint,1,opt,name=decision_expiration_delta,json=decisionExpirationDelta,proto3" json:"decision_expiration_delta,omitempty"`
	// number of blocks a DecisionProof is valid for
	ProofExpirationDelta uint64 `protobuf:"varint,2,opt,name=proof_expiration_delta,json=proofExpirationDelta,proto3" json:"proof_expiration_delta,omitempty"`
	// number of blocks an AccessTicket is valid for
	TicketExpirationDelta uint64 `protobuf:"varint,3,opt,name=ticket_expiration_delta,json=ticketExpirationDelta,proto3" json:"ticket_expiration_delta,omitempty"`
}

func (m *DecisionParams) Reset()         { *m = DecisionParams{} }
func (m *DecisionParams) String() string { return proto.CompactTextString(m) }
func (*DecisionParams) ProtoMessage()    {}
func (*DecisionParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_b2650781da68a114, []int{1}
}
func (m *DecisionParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DecisionParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DecisionParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DecisionParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DecisionParams.Merge(m, src)
}
func (m *DecisionParams) XXX_Size() int {
	return m.Size()
}
func (m *DecisionParams) XXX_DiscardUnknown() {
	xxx_messageInfo_DecisionParams.DiscardUnknown(m)
}

var xxx_messageInfo_DecisionParams proto.InternalMessageInfo

func (m *DecisionParams) GetDecisionExpirationDelta() uint64 {
	if m != nil {
		return m.DecisionExpirationDelta
	}
	return 0
}

func (m *DecisionParams) GetProofExpirationDelta() uint64 {
	if m != nil {
		return m.ProofExpirationDelta
	}
	return 0
}

func (m *DecisionParams) GetTicketExpirationDelta() uint64 {
	if m != nil {
		return m.TicketExpirationDelta
	}
	return 0
}

// AccessRequest represents the wish to perform a set of operations by an actor
type AccessRequest struct {
	Operations []*Operation `protobuf:"bytes,1,rep,name=operations,proto3" json:"operations,omitempty"`
	// actor requesting operations
	Actor *Actor `protobuf:"bytes,2,opt,name=actor,proto3" json:"actor,omitempty"`
}

func (m *AccessRequest) Reset()         { *m = AccessRequest{} }
func (m *AccessRequest) String() string { return proto.CompactTextString(m) }
func (*AccessRequest) ProtoMessage()    {}
func (*AccessRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b2650781da68a114, []int{2}
}
func (m *AccessRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AccessRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AccessRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AccessRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccessRequest.Merge(m, src)
}
func (m *AccessRequest) XXX_Size() int {
	return m.Size()
}
func (m *AccessRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AccessRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AccessRequest proto.InternalMessageInfo

func (m *AccessRequest) GetOperations() []*Operation {
	if m != nil {
		return m.Operations
	}
	return nil
}

func (m *AccessRequest) GetActor() *Actor {
	if m != nil {
		return m.Actor
	}
	return nil
}

// Operation represents an action over an object.
type Operation struct {
	// target object for operation
	Object *Object `protobuf:"bytes,1,opt,name=object,proto3" json:"object,omitempty"`
	// permission required to perform operation
	Permission string `protobuf:"bytes,2,opt,name=permission,proto3" json:"permission,omitempty"`
}

func (m *Operation) Reset()         { *m = Operation{} }
func (m *Operation) String() string { return proto.CompactTextString(m) }
func (*Operation) ProtoMessage()    {}
func (*Operation) Descriptor() ([]byte, []int) {
	return fileDescriptor_b2650781da68a114, []int{3}
}
func (m *Operation) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Operation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Operation.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Operation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Operation.Merge(m, src)
}
func (m *Operation) XXX_Size() int {
	return m.Size()
}
func (m *Operation) XXX_DiscardUnknown() {
	xxx_messageInfo_Operation.DiscardUnknown(m)
}

var xxx_messageInfo_Operation proto.InternalMessageInfo

func (m *Operation) GetObject() *Object {
	if m != nil {
		return m.Object
	}
	return nil
}

func (m *Operation) GetPermission() string {
	if m != nil {
		return m.Permission
	}
	return ""
}

func init() {
	proto.RegisterType((*AccessDecision)(nil), "sourcenetwork.acp_core.AccessDecision")
	proto.RegisterType((*DecisionParams)(nil), "sourcenetwork.acp_core.DecisionParams")
	proto.RegisterType((*AccessRequest)(nil), "sourcenetwork.acp_core.AccessRequest")
	proto.RegisterType((*Operation)(nil), "sourcenetwork.acp_core.Operation")
}

func init() {
	proto.RegisterFile("sourcenetwork/acp_core/access_decision.proto", fileDescriptor_b2650781da68a114)
}

var fileDescriptor_b2650781da68a114 = []byte{
	// 530 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x53, 0xcd, 0x6e, 0x13, 0x31,
	0x18, 0xcc, 0x26, 0x69, 0xda, 0x7c, 0x69, 0x72, 0xb0, 0x42, 0xbb, 0x04, 0xb1, 0x84, 0x20, 0xa1,
	0x20, 0xd0, 0x2e, 0x4a, 0x51, 0x0f, 0x1c, 0x40, 0x41, 0xad, 0x04, 0x27, 0xd0, 0xc2, 0x89, 0xcb,
	0xca, 0xf1, 0x7e, 0x4d, 0x4c, 0x7e, 0x6c, 0x6c, 0x47, 0xd0, 0x27, 0x80, 0x23, 0xcf, 0xc3, 0x13,
	0x70, 0xec, 0x91, 0x23, 0x4a, 0x5e, 0x04, 0xad, 0x77, 0x5d, 0xb5, 0x29, 0xb9, 0x70, 0x8b, 0xbf,
	0x99, 0xb1, 0x27, 0x33, 0xdf, 0xc2, 0x13, 0x2d, 0x96, 0x8a, 0xe1, 0x02, 0xcd, 0x17, 0xa1, 0xa6,
	0x11, 0x65, 0x32, 0x61, 0x42, 0x61, 0x44, 0x19, 0x43, 0xad, 0x93, 0x14, 0x19, 0xd7, 0x5c, 0x2c,
	0x42, 0xa9, 0x84, 0x11, 0xe4, 0xe0, 0x1a, 0x3b, 0x74, 0xec, 0xce, 0xbd, 0xb1, 0x10, 0xe3, 0x19,
	0x46, 0x96, 0x35, 0x5a, 0x9e, 0x45, 0x86, 0xcf, 0x51, 0x1b, 0x3a, 0x97, 0xb9, 0xb0, 0xf3, 0x68,
	0xcb, 0x33, 0x0a, 0x67, 0xd4, 0x70, 0xb1, 0xd0, 0x13, 0x5e, 0x50, 0x7b, 0xdf, 0x2b, 0xd0, 0x1a,
	0xda, 0xd7, 0x4f, 0x8a, 0xc7, 0x49, 0x0b, 0xca, 0x3c, 0xf5, 0xbd, 0xae, 0xd7, 0xaf, 0xc7, 0x65,
	0x9e, 0x92, 0x3b, 0x50, 0x97, 0x62, 0xc6, 0xd9, 0x79, 0xc2, 0x53, 0xbf, 0x6c, 0xc7, 0x7b, 0xf9,
	0xe0, 0x4d, 0x4a, 0x7c, 0xd8, 0x65, 0x0a, 0xa9, 0x11, 0xca, 0xaf, 0x58, 0xc8, 0x1d, 0xc9, 0x53,
	0x68, 0x17, 0x3f, 0x13, 0xca, 0x58, 0xa2, 0xf1, 0xf3, 0x12, 0x17, 0x0c, 0xfd, 0x6a, 0xd7, 0xeb,
	0x57, 0x63, 0x52, 0x60, 0x43, 0xc6, 0xde, 0x17, 0x08, 0x19, 0x02, 0x08, 0x89, 0x2a, 0xb7, 0xe8,
	0xef, 0x74, 0x2b, 0xfd, 0xc6, 0xe0, 0x7e, 0xf8, 0xef, 0x10, 0xc2, 0xb7, 0x8e, 0x19, 0x5f, 0x11,
	0x91, 0x36, 0xec, 0x50, 0x96, 0x99, 0xa9, 0x59, 0x33, 0xf9, 0x81, 0xbc, 0x80, 0x9a, 0xa4, 0x8a,
	0xce, 0xb5, 0xbf, 0xdb, 0xf5, 0xfa, 0x8d, 0xc1, 0xc3, 0x6d, 0x97, 0xba, 0x0c, 0xde, 0x59, 0x76,
	0x5c, 0xa8, 0xc8, 0x4b, 0x68, 0x5a, 0xbb, 0x5c, 0x2c, 0x92, 0x2c, 0x6b, 0x7f, 0xcf, 0x5e, 0xd3,
	0x09, 0xf3, 0x22, 0x42, 0x57, 0x44, 0xf8, 0xc1, 0x15, 0x11, 0xef, 0x3b, 0x41, 0x36, 0x22, 0x0f,
	0xa0, 0xc9, 0xb5, 0x5e, 0x62, 0x9a, 0x4c, 0x90, 0x8f, 0x27, 0xc6, 0xaf, 0xdb, 0x10, 0xf6, 0xf3,
	0xe1, 0x6b, 0x3b, 0xeb, 0xfd, 0xf4, 0xa0, 0x75, 0xdd, 0x00, 0x79, 0x0e, 0xb7, 0xdd, 0x4e, 0x24,
	0xf8, 0x55, 0xf2, 0xfc, 0x6f, 0x26, 0x29, 0xce, 0x0c, 0xb5, 0x0d, 0x55, 0xe3, 0x43, 0x47, 0x38,
	0xbd, 0xc4, 0x4f, 0x32, 0x98, 0x3c, 0x83, 0x03, 0xa9, 0x84, 0x38, 0xbb, 0x29, 0x2c, 0x5b, 0x61,
	0xdb, 0xa2, 0x9b, 0xaa, 0x63, 0x38, 0x34, 0x9c, 0x4d, 0xd1, 0xdc, 0x94, 0x55, 0xac, 0xec, 0x56,
	0x0e, 0x6f, 0xe8, 0x7a, 0xdf, 0x3c, 0x68, 0xe6, 0x7b, 0x14, 0x67, 0x75, 0x6a, 0xb3, 0xd1, 0xa6,
	0xf7, 0x3f, 0x6d, 0x1e, 0xb9, 0x36, 0xcb, 0x36, 0xef, 0xbb, 0xdb, 0xd4, 0xc3, 0x8c, 0x54, 0x94,
	0xdd, 0x63, 0x50, 0xbf, 0xbc, 0x8d, 0x1c, 0x43, 0x4d, 0x8c, 0x3e, 0x21, 0x33, 0x36, 0xad, 0xc6,
	0x20, 0xd8, 0x6a, 0xc0, 0xb2, 0xe2, 0x82, 0x4d, 0x02, 0x00, 0x89, 0x6a, 0xce, 0x75, 0x96, 0x6c,
	0xb1, 0xf4, 0x57, 0x26, 0xaf, 0x4e, 0x7f, 0xad, 0x02, 0xef, 0x62, 0x15, 0x78, 0x7f, 0x56, 0x81,
	0xf7, 0x63, 0x1d, 0x94, 0x2e, 0xd6, 0x41, 0xe9, 0xf7, 0x3a, 0x28, 0x7d, 0x7c, 0x3c, 0xe6, 0x66,
	0xb2, 0x1c, 0x85, 0x4c, 0xcc, 0xa3, 0x2d, 0x9f, 0xa1, 0x9c, 0x8e, 0x23, 0x73, 0x2e, 0x51, 0x8f,
	0x6a, 0x76, 0x73, 0x8e, 0xfe, 0x06, 0x00, 0x00, 0xff, 0xff, 0xdb, 0x05, 0xa5, 0x54, 0x18, 0x04,
	0x00, 0x00,
}

func (m *AccessDecision) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AccessDecision) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AccessDecision) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.IssuedHeight != 0 {
		i = encodeVarintAccessDecision(dAtA, i, uint64(m.IssuedHeight))
		i--
		dAtA[i] = 0x48
	}
	if m.CreationTime != nil {
		{
			size, err := m.CreationTime.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintAccessDecision(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x42
	}
	if m.Params != nil {
		{
			size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintAccessDecision(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x3a
	}
	if len(m.Actor) > 0 {
		i -= len(m.Actor)
		copy(dAtA[i:], m.Actor)
		i = encodeVarintAccessDecision(dAtA, i, uint64(len(m.Actor)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.Operations) > 0 {
		for iNdEx := len(m.Operations) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Operations[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintAccessDecision(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if m.CreatorAccSequence != 0 {
		i = encodeVarintAccessDecision(dAtA, i, uint64(m.CreatorAccSequence))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintAccessDecision(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.PolicyId) > 0 {
		i -= len(m.PolicyId)
		copy(dAtA[i:], m.PolicyId)
		i = encodeVarintAccessDecision(dAtA, i, uint64(len(m.PolicyId)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintAccessDecision(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *DecisionParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DecisionParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DecisionParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.TicketExpirationDelta != 0 {
		i = encodeVarintAccessDecision(dAtA, i, uint64(m.TicketExpirationDelta))
		i--
		dAtA[i] = 0x18
	}
	if m.ProofExpirationDelta != 0 {
		i = encodeVarintAccessDecision(dAtA, i, uint64(m.ProofExpirationDelta))
		i--
		dAtA[i] = 0x10
	}
	if m.DecisionExpirationDelta != 0 {
		i = encodeVarintAccessDecision(dAtA, i, uint64(m.DecisionExpirationDelta))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *AccessRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AccessRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AccessRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Actor != nil {
		{
			size, err := m.Actor.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintAccessDecision(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Operations) > 0 {
		for iNdEx := len(m.Operations) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Operations[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintAccessDecision(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Operation) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Operation) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Operation) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Permission) > 0 {
		i -= len(m.Permission)
		copy(dAtA[i:], m.Permission)
		i = encodeVarintAccessDecision(dAtA, i, uint64(len(m.Permission)))
		i--
		dAtA[i] = 0x12
	}
	if m.Object != nil {
		{
			size, err := m.Object.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintAccessDecision(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintAccessDecision(dAtA []byte, offset int, v uint64) int {
	offset -= sovAccessDecision(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *AccessDecision) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovAccessDecision(uint64(l))
	}
	l = len(m.PolicyId)
	if l > 0 {
		n += 1 + l + sovAccessDecision(uint64(l))
	}
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovAccessDecision(uint64(l))
	}
	if m.CreatorAccSequence != 0 {
		n += 1 + sovAccessDecision(uint64(m.CreatorAccSequence))
	}
	if len(m.Operations) > 0 {
		for _, e := range m.Operations {
			l = e.Size()
			n += 1 + l + sovAccessDecision(uint64(l))
		}
	}
	l = len(m.Actor)
	if l > 0 {
		n += 1 + l + sovAccessDecision(uint64(l))
	}
	if m.Params != nil {
		l = m.Params.Size()
		n += 1 + l + sovAccessDecision(uint64(l))
	}
	if m.CreationTime != nil {
		l = m.CreationTime.Size()
		n += 1 + l + sovAccessDecision(uint64(l))
	}
	if m.IssuedHeight != 0 {
		n += 1 + sovAccessDecision(uint64(m.IssuedHeight))
	}
	return n
}

func (m *DecisionParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.DecisionExpirationDelta != 0 {
		n += 1 + sovAccessDecision(uint64(m.DecisionExpirationDelta))
	}
	if m.ProofExpirationDelta != 0 {
		n += 1 + sovAccessDecision(uint64(m.ProofExpirationDelta))
	}
	if m.TicketExpirationDelta != 0 {
		n += 1 + sovAccessDecision(uint64(m.TicketExpirationDelta))
	}
	return n
}

func (m *AccessRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Operations) > 0 {
		for _, e := range m.Operations {
			l = e.Size()
			n += 1 + l + sovAccessDecision(uint64(l))
		}
	}
	if m.Actor != nil {
		l = m.Actor.Size()
		n += 1 + l + sovAccessDecision(uint64(l))
	}
	return n
}

func (m *Operation) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Object != nil {
		l = m.Object.Size()
		n += 1 + l + sovAccessDecision(uint64(l))
	}
	l = len(m.Permission)
	if l > 0 {
		n += 1 + l + sovAccessDecision(uint64(l))
	}
	return n
}

func sovAccessDecision(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozAccessDecision(x uint64) (n int) {
	return sovAccessDecision(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *AccessDecision) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAccessDecision
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
			return fmt.Errorf("proto: AccessDecision: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AccessDecision: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccessDecision
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
				return ErrInvalidLengthAccessDecision
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAccessDecision
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PolicyId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccessDecision
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
				return ErrInvalidLengthAccessDecision
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAccessDecision
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PolicyId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccessDecision
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
				return ErrInvalidLengthAccessDecision
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAccessDecision
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreatorAccSequence", wireType)
			}
			m.CreatorAccSequence = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccessDecision
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CreatorAccSequence |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Operations", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccessDecision
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
				return ErrInvalidLengthAccessDecision
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAccessDecision
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Operations = append(m.Operations, &Operation{})
			if err := m.Operations[len(m.Operations)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Actor", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccessDecision
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
				return ErrInvalidLengthAccessDecision
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAccessDecision
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Actor = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccessDecision
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
				return ErrInvalidLengthAccessDecision
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAccessDecision
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Params == nil {
				m.Params = &DecisionParams{}
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreationTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccessDecision
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
				return ErrInvalidLengthAccessDecision
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAccessDecision
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.CreationTime == nil {
				m.CreationTime = &types.Timestamp{}
			}
			if err := m.CreationTime.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IssuedHeight", wireType)
			}
			m.IssuedHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccessDecision
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.IssuedHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipAccessDecision(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAccessDecision
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
func (m *DecisionParams) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAccessDecision
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
			return fmt.Errorf("proto: DecisionParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DecisionParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DecisionExpirationDelta", wireType)
			}
			m.DecisionExpirationDelta = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccessDecision
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.DecisionExpirationDelta |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProofExpirationDelta", wireType)
			}
			m.ProofExpirationDelta = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccessDecision
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ProofExpirationDelta |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TicketExpirationDelta", wireType)
			}
			m.TicketExpirationDelta = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccessDecision
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TicketExpirationDelta |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipAccessDecision(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAccessDecision
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
func (m *AccessRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAccessDecision
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
			return fmt.Errorf("proto: AccessRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AccessRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Operations", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccessDecision
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
				return ErrInvalidLengthAccessDecision
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAccessDecision
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Operations = append(m.Operations, &Operation{})
			if err := m.Operations[len(m.Operations)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Actor", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccessDecision
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
				return ErrInvalidLengthAccessDecision
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAccessDecision
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Actor == nil {
				m.Actor = &Actor{}
			}
			if err := m.Actor.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAccessDecision(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAccessDecision
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
func (m *Operation) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAccessDecision
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
			return fmt.Errorf("proto: Operation: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Operation: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Object", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccessDecision
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
				return ErrInvalidLengthAccessDecision
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAccessDecision
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Object == nil {
				m.Object = &Object{}
			}
			if err := m.Object.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Permission", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccessDecision
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
				return ErrInvalidLengthAccessDecision
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAccessDecision
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Permission = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAccessDecision(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAccessDecision
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
func skipAccessDecision(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAccessDecision
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
					return 0, ErrIntOverflowAccessDecision
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
					return 0, ErrIntOverflowAccessDecision
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
				return 0, ErrInvalidLengthAccessDecision
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupAccessDecision
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthAccessDecision
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthAccessDecision        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAccessDecision          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupAccessDecision = fmt.Errorf("proto: unexpected end of group")
)
