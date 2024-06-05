// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sourcenetwork/acp_core/policy_record.proto

package types

import (
	fmt "fmt"
	proto "github.com/cosmos/gogoproto/proto"
	_ "github.com/cosmos/gogoproto/types"
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

// PolicyRecord represents a the Policy Document which will be persisted in the data layer
type PolicyRecord struct {
	Policy          *Policy          `protobuf:"bytes,1,opt,name=policy,proto3" json:"policy,omitempty"`
	ManagementGraph *ManagementGraph `protobuf:"bytes,2,opt,name=management_graph,json=managementGraph,proto3" json:"management_graph,omitempty"`
}

func (m *PolicyRecord) Reset()         { *m = PolicyRecord{} }
func (m *PolicyRecord) String() string { return proto.CompactTextString(m) }
func (*PolicyRecord) ProtoMessage()    {}
func (*PolicyRecord) Descriptor() ([]byte, []int) {
	return fileDescriptor_51a4aea7fe8278b3, []int{0}
}
func (m *PolicyRecord) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PolicyRecord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PolicyRecord.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PolicyRecord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PolicyRecord.Merge(m, src)
}
func (m *PolicyRecord) XXX_Size() int {
	return m.Size()
}
func (m *PolicyRecord) XXX_DiscardUnknown() {
	xxx_messageInfo_PolicyRecord.DiscardUnknown(m)
}

var xxx_messageInfo_PolicyRecord proto.InternalMessageInfo

func (m *PolicyRecord) GetPolicy() *Policy {
	if m != nil {
		return m.Policy
	}
	return nil
}

func (m *PolicyRecord) GetManagementGraph() *ManagementGraph {
	if m != nil {
		return m.ManagementGraph
	}
	return nil
}

// ManagementGraph represents a Policy's Relation Management Graph.
//
// The ManagementGraph is a directed graph which expresses the notion of Relation Management Authority.
// Relation Management Authority is the idea that a certain set of relationships with relation R will be managed by an actor with relation RM.
// Thus we can say RM manages R, meaning that if an actor A has a relationship 'actor {A} is a {RM} for {O}' where O is an object,
// then Actor A can create relationships 'actor {S} is a {R} for {O}' for any actor S.
//
// Nodes in the Graph are Relations in a Policy.
// Edges point from one Relation to another.
//
// NOTE: This proto definition should be treated as an *abstract data type*,
// meaning that the fields should not be manually editted.
type ManagementGraph struct {
	// map of node id to node definition
	Nodes map[string]*ManagerNode `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// stores all edges leaving a node
	ForwardEdges map[string]*ManagerEdges `protobuf:"bytes,2,rep,name=forward_edges,json=forwardEdges,proto3" json:"forward_edges,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// stores all edges pointing to a node
	BackwardEdges map[string]*ManagerEdges `protobuf:"bytes,3,rep,name=backward_edges,json=backwardEdges,proto3" json:"backward_edges,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *ManagementGraph) Reset()         { *m = ManagementGraph{} }
func (m *ManagementGraph) String() string { return proto.CompactTextString(m) }
func (*ManagementGraph) ProtoMessage()    {}
func (*ManagementGraph) Descriptor() ([]byte, []int) {
	return fileDescriptor_51a4aea7fe8278b3, []int{1}
}
func (m *ManagementGraph) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ManagementGraph) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ManagementGraph.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ManagementGraph) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ManagementGraph.Merge(m, src)
}
func (m *ManagementGraph) XXX_Size() int {
	return m.Size()
}
func (m *ManagementGraph) XXX_DiscardUnknown() {
	xxx_messageInfo_ManagementGraph.DiscardUnknown(m)
}

var xxx_messageInfo_ManagementGraph proto.InternalMessageInfo

func (m *ManagementGraph) GetNodes() map[string]*ManagerNode {
	if m != nil {
		return m.Nodes
	}
	return nil
}

func (m *ManagementGraph) GetForwardEdges() map[string]*ManagerEdges {
	if m != nil {
		return m.ForwardEdges
	}
	return nil
}

func (m *ManagementGraph) GetBackwardEdges() map[string]*ManagerEdges {
	if m != nil {
		return m.BackwardEdges
	}
	return nil
}

type ManagerNode struct {
	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Text string `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
}

func (m *ManagerNode) Reset()         { *m = ManagerNode{} }
func (m *ManagerNode) String() string { return proto.CompactTextString(m) }
func (*ManagerNode) ProtoMessage()    {}
func (*ManagerNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_51a4aea7fe8278b3, []int{2}
}
func (m *ManagerNode) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ManagerNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ManagerNode.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ManagerNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ManagerNode.Merge(m, src)
}
func (m *ManagerNode) XXX_Size() int {
	return m.Size()
}
func (m *ManagerNode) XXX_DiscardUnknown() {
	xxx_messageInfo_ManagerNode.DiscardUnknown(m)
}

var xxx_messageInfo_ManagerNode proto.InternalMessageInfo

func (m *ManagerNode) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ManagerNode) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

type ManagerEdges struct {
	Edges map[string]bool `protobuf:"bytes,1,rep,name=edges,proto3" json:"edges,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (m *ManagerEdges) Reset()         { *m = ManagerEdges{} }
func (m *ManagerEdges) String() string { return proto.CompactTextString(m) }
func (*ManagerEdges) ProtoMessage()    {}
func (*ManagerEdges) Descriptor() ([]byte, []int) {
	return fileDescriptor_51a4aea7fe8278b3, []int{3}
}
func (m *ManagerEdges) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ManagerEdges) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ManagerEdges.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ManagerEdges) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ManagerEdges.Merge(m, src)
}
func (m *ManagerEdges) XXX_Size() int {
	return m.Size()
}
func (m *ManagerEdges) XXX_DiscardUnknown() {
	xxx_messageInfo_ManagerEdges.DiscardUnknown(m)
}

var xxx_messageInfo_ManagerEdges proto.InternalMessageInfo

func (m *ManagerEdges) GetEdges() map[string]bool {
	if m != nil {
		return m.Edges
	}
	return nil
}

func init() {
	proto.RegisterType((*PolicyRecord)(nil), "sourcenetwork.acp_core.PolicyRecord")
	proto.RegisterType((*ManagementGraph)(nil), "sourcenetwork.acp_core.ManagementGraph")
	proto.RegisterMapType((map[string]*ManagerEdges)(nil), "sourcenetwork.acp_core.ManagementGraph.BackwardEdgesEntry")
	proto.RegisterMapType((map[string]*ManagerEdges)(nil), "sourcenetwork.acp_core.ManagementGraph.ForwardEdgesEntry")
	proto.RegisterMapType((map[string]*ManagerNode)(nil), "sourcenetwork.acp_core.ManagementGraph.NodesEntry")
	proto.RegisterType((*ManagerNode)(nil), "sourcenetwork.acp_core.ManagerNode")
	proto.RegisterType((*ManagerEdges)(nil), "sourcenetwork.acp_core.ManagerEdges")
	proto.RegisterMapType((map[string]bool)(nil), "sourcenetwork.acp_core.ManagerEdges.EdgesEntry")
}

func init() {
	proto.RegisterFile("sourcenetwork/acp_core/policy_record.proto", fileDescriptor_51a4aea7fe8278b3)
}

var fileDescriptor_51a4aea7fe8278b3 = []byte{
	// 461 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x93, 0x4f, 0x6f, 0xd3, 0x40,
	0x10, 0xc5, 0xb3, 0x4e, 0x53, 0xd1, 0x49, 0xfa, 0x87, 0x15, 0x42, 0x91, 0x0f, 0xa6, 0x4a, 0x91,
	0xa8, 0x40, 0xb2, 0x45, 0x90, 0x50, 0x9b, 0x63, 0xa5, 0x00, 0x17, 0x10, 0xf2, 0x11, 0x09, 0xa2,
	0xb5, 0x3d, 0x71, 0xad, 0xc4, 0x5e, 0x6b, 0xbd, 0xa6, 0xe4, 0x53, 0xc0, 0x99, 0x4f, 0xc4, 0x81,
	0x43, 0x8f, 0x1c, 0x51, 0xf2, 0x45, 0x90, 0x77, 0x0d, 0x6c, 0x68, 0x1b, 0x7c, 0xe8, 0x6d, 0x35,
	0x7a, 0xef, 0xfd, 0xfc, 0x34, 0x1e, 0x78, 0x5c, 0xf0, 0x52, 0x84, 0x98, 0xa1, 0xbc, 0xe0, 0x62,
	0xe6, 0xb1, 0x30, 0x9f, 0x84, 0x5c, 0xa0, 0x97, 0xf3, 0x79, 0x12, 0x2e, 0x26, 0x02, 0x43, 0x2e,
	0x22, 0x37, 0x17, 0x5c, 0x72, 0x7a, 0x7f, 0x4d, 0xeb, 0xfe, 0xd6, 0xda, 0x0f, 0x62, 0xce, 0xe3,
	0x39, 0x7a, 0x4a, 0x15, 0x94, 0x53, 0x4f, 0x26, 0x29, 0x16, 0x92, 0xa5, 0xb9, 0x36, 0xda, 0x47,
	0x1b, 0x21, 0x5a, 0x34, 0xf8, 0x4a, 0xa0, 0xf7, 0x56, 0x0d, 0x7c, 0x05, 0xa5, 0xcf, 0x61, 0x5b,
	0x0b, 0xfa, 0xe4, 0x90, 0x1c, 0x77, 0x87, 0x8e, 0x7b, 0x3d, 0xdf, 0xad, 0x5d, 0xb5, 0x9a, 0xfa,
	0x70, 0x90, 0xb2, 0x8c, 0xc5, 0x98, 0x62, 0x26, 0x27, 0xb1, 0x60, 0xf9, 0x79, 0xdf, 0x52, 0x09,
	0x8f, 0x6e, 0x4a, 0x78, 0xfd, 0x47, 0xff, 0xb2, 0x92, 0xfb, 0xfb, 0xe9, 0xfa, 0x60, 0xf0, 0x7d,
	0x0b, 0xf6, 0xff, 0x11, 0xd1, 0x57, 0xd0, 0xc9, 0x78, 0x84, 0x45, 0x9f, 0x1c, 0xb6, 0x8f, 0xbb,
	0xc3, 0x61, 0xc3, 0x70, 0xf7, 0x4d, 0x65, 0x1a, 0x67, 0x52, 0x2c, 0x7c, 0x1d, 0x40, 0x3f, 0xc0,
	0xee, 0x94, 0x8b, 0x0b, 0x26, 0xa2, 0x09, 0x46, 0x31, 0x16, 0x7d, 0x4b, 0x25, 0x9e, 0x36, 0x4d,
	0x7c, 0xa1, 0xcd, 0xe3, 0xca, 0xab, 0x83, 0x7b, 0x53, 0x63, 0x44, 0x19, 0xec, 0x05, 0x2c, 0x9c,
	0x19, 0x80, 0xb6, 0x02, 0x8c, 0x9a, 0x02, 0xce, 0x6a, 0xb7, 0x41, 0xd8, 0x0d, 0xcc, 0x99, 0xfd,
	0x1e, 0xe0, 0x6f, 0x2f, 0x7a, 0x00, 0xed, 0x19, 0xea, 0xbd, 0xed, 0xf8, 0xd5, 0x93, 0x9e, 0x42,
	0xe7, 0x23, 0x9b, 0x97, 0x58, 0x6f, 0xe2, 0x68, 0x33, 0x59, 0x54, 0x59, 0xbe, 0x76, 0x8c, 0xac,
	0x13, 0x62, 0x23, 0xdc, 0xbd, 0x52, 0xf2, 0x1a, 0xca, 0x68, 0x9d, 0xf2, 0xf0, 0x3f, 0x14, 0x95,
	0x65, 0x62, 0xa6, 0x40, 0xaf, 0x56, 0xbd, 0x7d, 0xce, 0xe0, 0x29, 0x74, 0x8d, 0xa2, 0x74, 0x0f,
	0xac, 0x24, 0xaa, 0xf3, 0xad, 0x24, 0xa2, 0x14, 0xb6, 0x24, 0x7e, 0x92, 0x2a, 0x7d, 0xc7, 0x57,
	0xef, 0xc1, 0x67, 0x02, 0x3d, 0x33, 0x8e, 0x8e, 0xa1, 0xa3, 0x77, 0xa9, 0x7f, 0x3f, 0xaf, 0xc9,
	0x37, 0xb8, 0xc6, 0x02, 0xb5, 0xdb, 0x3e, 0x01, 0xd8, 0x58, 0xf5, 0x9e, 0x59, 0xf5, 0x8e, 0x51,
	0xe2, 0x6c, 0xfc, 0x6d, 0xe9, 0x90, 0xcb, 0xa5, 0x43, 0x7e, 0x2e, 0x1d, 0xf2, 0x65, 0xe5, 0xb4,
	0x2e, 0x57, 0x4e, 0xeb, 0xc7, 0xca, 0x69, 0xbd, 0x7b, 0x12, 0x27, 0xf2, 0xbc, 0x0c, 0xdc, 0x90,
	0xa7, 0xde, 0x4d, 0xa7, 0x3f, 0x8b, 0x3d, 0xb9, 0xc8, 0xb1, 0x08, 0xb6, 0xd5, 0xf9, 0x3f, 0xfb,
	0x15, 0x00, 0x00, 0xff, 0xff, 0x91, 0x2a, 0xb4, 0x9c, 0x8a, 0x04, 0x00, 0x00,
}

func (m *PolicyRecord) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PolicyRecord) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PolicyRecord) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ManagementGraph != nil {
		{
			size, err := m.ManagementGraph.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPolicyRecord(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.Policy != nil {
		{
			size, err := m.Policy.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPolicyRecord(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ManagementGraph) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ManagementGraph) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ManagementGraph) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.BackwardEdges) > 0 {
		for k := range m.BackwardEdges {
			v := m.BackwardEdges[k]
			baseI := i
			if v != nil {
				{
					size, err := v.MarshalToSizedBuffer(dAtA[:i])
					if err != nil {
						return 0, err
					}
					i -= size
					i = encodeVarintPolicyRecord(dAtA, i, uint64(size))
				}
				i--
				dAtA[i] = 0x12
			}
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintPolicyRecord(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintPolicyRecord(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.ForwardEdges) > 0 {
		for k := range m.ForwardEdges {
			v := m.ForwardEdges[k]
			baseI := i
			if v != nil {
				{
					size, err := v.MarshalToSizedBuffer(dAtA[:i])
					if err != nil {
						return 0, err
					}
					i -= size
					i = encodeVarintPolicyRecord(dAtA, i, uint64(size))
				}
				i--
				dAtA[i] = 0x12
			}
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintPolicyRecord(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintPolicyRecord(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Nodes) > 0 {
		for k := range m.Nodes {
			v := m.Nodes[k]
			baseI := i
			if v != nil {
				{
					size, err := v.MarshalToSizedBuffer(dAtA[:i])
					if err != nil {
						return 0, err
					}
					i -= size
					i = encodeVarintPolicyRecord(dAtA, i, uint64(size))
				}
				i--
				dAtA[i] = 0x12
			}
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintPolicyRecord(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintPolicyRecord(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *ManagerNode) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ManagerNode) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ManagerNode) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Text) > 0 {
		i -= len(m.Text)
		copy(dAtA[i:], m.Text)
		i = encodeVarintPolicyRecord(dAtA, i, uint64(len(m.Text)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintPolicyRecord(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ManagerEdges) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ManagerEdges) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ManagerEdges) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Edges) > 0 {
		for k := range m.Edges {
			v := m.Edges[k]
			baseI := i
			i--
			if v {
				dAtA[i] = 1
			} else {
				dAtA[i] = 0
			}
			i--
			dAtA[i] = 0x10
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintPolicyRecord(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintPolicyRecord(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintPolicyRecord(dAtA []byte, offset int, v uint64) int {
	offset -= sovPolicyRecord(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *PolicyRecord) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Policy != nil {
		l = m.Policy.Size()
		n += 1 + l + sovPolicyRecord(uint64(l))
	}
	if m.ManagementGraph != nil {
		l = m.ManagementGraph.Size()
		n += 1 + l + sovPolicyRecord(uint64(l))
	}
	return n
}

func (m *ManagementGraph) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Nodes) > 0 {
		for k, v := range m.Nodes {
			_ = k
			_ = v
			l = 0
			if v != nil {
				l = v.Size()
				l += 1 + sovPolicyRecord(uint64(l))
			}
			mapEntrySize := 1 + len(k) + sovPolicyRecord(uint64(len(k))) + l
			n += mapEntrySize + 1 + sovPolicyRecord(uint64(mapEntrySize))
		}
	}
	if len(m.ForwardEdges) > 0 {
		for k, v := range m.ForwardEdges {
			_ = k
			_ = v
			l = 0
			if v != nil {
				l = v.Size()
				l += 1 + sovPolicyRecord(uint64(l))
			}
			mapEntrySize := 1 + len(k) + sovPolicyRecord(uint64(len(k))) + l
			n += mapEntrySize + 1 + sovPolicyRecord(uint64(mapEntrySize))
		}
	}
	if len(m.BackwardEdges) > 0 {
		for k, v := range m.BackwardEdges {
			_ = k
			_ = v
			l = 0
			if v != nil {
				l = v.Size()
				l += 1 + sovPolicyRecord(uint64(l))
			}
			mapEntrySize := 1 + len(k) + sovPolicyRecord(uint64(len(k))) + l
			n += mapEntrySize + 1 + sovPolicyRecord(uint64(mapEntrySize))
		}
	}
	return n
}

func (m *ManagerNode) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovPolicyRecord(uint64(l))
	}
	l = len(m.Text)
	if l > 0 {
		n += 1 + l + sovPolicyRecord(uint64(l))
	}
	return n
}

func (m *ManagerEdges) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Edges) > 0 {
		for k, v := range m.Edges {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovPolicyRecord(uint64(len(k))) + 1 + 1
			n += mapEntrySize + 1 + sovPolicyRecord(uint64(mapEntrySize))
		}
	}
	return n
}

func sovPolicyRecord(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPolicyRecord(x uint64) (n int) {
	return sovPolicyRecord(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *PolicyRecord) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPolicyRecord
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
			return fmt.Errorf("proto: PolicyRecord: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PolicyRecord: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Policy", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPolicyRecord
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
				return ErrInvalidLengthPolicyRecord
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPolicyRecord
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Policy == nil {
				m.Policy = &Policy{}
			}
			if err := m.Policy.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ManagementGraph", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPolicyRecord
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
				return ErrInvalidLengthPolicyRecord
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPolicyRecord
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ManagementGraph == nil {
				m.ManagementGraph = &ManagementGraph{}
			}
			if err := m.ManagementGraph.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPolicyRecord(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPolicyRecord
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
func (m *ManagementGraph) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPolicyRecord
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
			return fmt.Errorf("proto: ManagementGraph: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ManagementGraph: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nodes", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPolicyRecord
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
				return ErrInvalidLengthPolicyRecord
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPolicyRecord
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Nodes == nil {
				m.Nodes = make(map[string]*ManagerNode)
			}
			var mapkey string
			var mapvalue *ManagerNode
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowPolicyRecord
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
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowPolicyRecord
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthPolicyRecord
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthPolicyRecord
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var mapmsglen int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowPolicyRecord
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapmsglen |= int(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					if mapmsglen < 0 {
						return ErrInvalidLengthPolicyRecord
					}
					postmsgIndex := iNdEx + mapmsglen
					if postmsgIndex < 0 {
						return ErrInvalidLengthPolicyRecord
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &ManagerNode{}
					if err := mapvalue.Unmarshal(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipPolicyRecord(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthPolicyRecord
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Nodes[mapkey] = mapvalue
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ForwardEdges", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPolicyRecord
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
				return ErrInvalidLengthPolicyRecord
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPolicyRecord
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ForwardEdges == nil {
				m.ForwardEdges = make(map[string]*ManagerEdges)
			}
			var mapkey string
			var mapvalue *ManagerEdges
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowPolicyRecord
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
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowPolicyRecord
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthPolicyRecord
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthPolicyRecord
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var mapmsglen int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowPolicyRecord
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapmsglen |= int(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					if mapmsglen < 0 {
						return ErrInvalidLengthPolicyRecord
					}
					postmsgIndex := iNdEx + mapmsglen
					if postmsgIndex < 0 {
						return ErrInvalidLengthPolicyRecord
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &ManagerEdges{}
					if err := mapvalue.Unmarshal(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipPolicyRecord(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthPolicyRecord
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.ForwardEdges[mapkey] = mapvalue
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BackwardEdges", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPolicyRecord
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
				return ErrInvalidLengthPolicyRecord
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPolicyRecord
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.BackwardEdges == nil {
				m.BackwardEdges = make(map[string]*ManagerEdges)
			}
			var mapkey string
			var mapvalue *ManagerEdges
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowPolicyRecord
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
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowPolicyRecord
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthPolicyRecord
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthPolicyRecord
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var mapmsglen int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowPolicyRecord
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapmsglen |= int(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					if mapmsglen < 0 {
						return ErrInvalidLengthPolicyRecord
					}
					postmsgIndex := iNdEx + mapmsglen
					if postmsgIndex < 0 {
						return ErrInvalidLengthPolicyRecord
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &ManagerEdges{}
					if err := mapvalue.Unmarshal(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipPolicyRecord(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthPolicyRecord
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.BackwardEdges[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPolicyRecord(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPolicyRecord
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
func (m *ManagerNode) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPolicyRecord
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
			return fmt.Errorf("proto: ManagerNode: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ManagerNode: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPolicyRecord
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
				return ErrInvalidLengthPolicyRecord
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPolicyRecord
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Text", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPolicyRecord
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
				return ErrInvalidLengthPolicyRecord
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPolicyRecord
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Text = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPolicyRecord(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPolicyRecord
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
func (m *ManagerEdges) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPolicyRecord
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
			return fmt.Errorf("proto: ManagerEdges: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ManagerEdges: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Edges", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPolicyRecord
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
				return ErrInvalidLengthPolicyRecord
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPolicyRecord
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Edges == nil {
				m.Edges = make(map[string]bool)
			}
			var mapkey string
			var mapvalue bool
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowPolicyRecord
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
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowPolicyRecord
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthPolicyRecord
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthPolicyRecord
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var mapvaluetemp int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowPolicyRecord
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapvaluetemp |= int(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					mapvalue = bool(mapvaluetemp != 0)
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipPolicyRecord(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthPolicyRecord
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Edges[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPolicyRecord(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPolicyRecord
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
func skipPolicyRecord(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPolicyRecord
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
					return 0, ErrIntOverflowPolicyRecord
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
					return 0, ErrIntOverflowPolicyRecord
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
				return 0, ErrInvalidLengthPolicyRecord
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupPolicyRecord
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthPolicyRecord
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthPolicyRecord        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPolicyRecord          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupPolicyRecord = fmt.Errorf("proto: unexpected end of group")
)
