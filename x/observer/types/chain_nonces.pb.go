// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: zetachain/zetacore/observer/chain_nonces.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

type ChainNonces struct {
	Creator string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	// deprecated(v19): index has been replaced by chain_id for unique identifier
	Index           string   `protobuf:"bytes,2,opt,name=index,proto3" json:"index,omitempty"` // Deprecated: Do not use.
	ChainId         int64    `protobuf:"varint,3,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	Nonce           uint64   `protobuf:"varint,4,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Signers         []string `protobuf:"bytes,5,rep,name=signers,proto3" json:"signers,omitempty"`
	FinalizedHeight uint64   `protobuf:"varint,6,opt,name=finalizedHeight,proto3" json:"finalizedHeight,omitempty"`
}

func (m *ChainNonces) Reset()         { *m = ChainNonces{} }
func (m *ChainNonces) String() string { return proto.CompactTextString(m) }
func (*ChainNonces) ProtoMessage()    {}
func (*ChainNonces) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8bc11608907ed56, []int{0}
}
func (m *ChainNonces) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ChainNonces) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ChainNonces.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ChainNonces) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChainNonces.Merge(m, src)
}
func (m *ChainNonces) XXX_Size() int {
	return m.Size()
}
func (m *ChainNonces) XXX_DiscardUnknown() {
	xxx_messageInfo_ChainNonces.DiscardUnknown(m)
}

var xxx_messageInfo_ChainNonces proto.InternalMessageInfo

func (m *ChainNonces) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

// Deprecated: Do not use.
func (m *ChainNonces) GetIndex() string {
	if m != nil {
		return m.Index
	}
	return ""
}

func (m *ChainNonces) GetChainId() int64 {
	if m != nil {
		return m.ChainId
	}
	return 0
}

func (m *ChainNonces) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *ChainNonces) GetSigners() []string {
	if m != nil {
		return m.Signers
	}
	return nil
}

func (m *ChainNonces) GetFinalizedHeight() uint64 {
	if m != nil {
		return m.FinalizedHeight
	}
	return 0
}

func init() {
	proto.RegisterType((*ChainNonces)(nil), "zetachain.zetacore.observer.ChainNonces")
}

func init() {
	proto.RegisterFile("zetachain/zetacore/observer/chain_nonces.proto", fileDescriptor_d8bc11608907ed56)
}

var fileDescriptor_d8bc11608907ed56 = []byte{
	// 272 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0xbd, 0x4e, 0xc3, 0x30,
	0x14, 0x85, 0xeb, 0xa6, 0x69, 0xc1, 0x0c, 0x48, 0x56, 0x07, 0x03, 0x92, 0x15, 0x31, 0x65, 0xc1,
	0x1e, 0x78, 0x83, 0xb2, 0xd0, 0x85, 0x21, 0x23, 0x0b, 0xca, 0xcf, 0xc5, 0xb1, 0x04, 0x76, 0x65,
	0x1b, 0x54, 0xfa, 0x14, 0x3c, 0x0d, 0xcf, 0xc0, 0xd8, 0x91, 0x11, 0x25, 0x2f, 0x82, 0xec, 0x10,
	0x90, 0xba, 0xdd, 0x73, 0x74, 0xce, 0xbd, 0x57, 0x1f, 0xe6, 0x3b, 0xf0, 0x65, 0xdd, 0x96, 0x4a,
	0x8b, 0x38, 0x19, 0x0b, 0xc2, 0x54, 0x0e, 0xec, 0x2b, 0x58, 0x11, 0xfd, 0x07, 0x6d, 0x74, 0x0d,
	0x8e, 0x6f, 0xac, 0xf1, 0x86, 0x5c, 0xfc, 0xe5, 0xf9, 0x98, 0xe7, 0x63, 0xfe, 0x7c, 0x29, 0x8d,
	0x34, 0x31, 0x27, 0xc2, 0x34, 0x54, 0x2e, 0x3f, 0x10, 0x3e, 0xb9, 0x09, 0x8d, 0xbb, 0xb8, 0x88,
	0x50, 0xbc, 0xa8, 0x2d, 0x94, 0xde, 0x58, 0x8a, 0x32, 0x94, 0x1f, 0x17, 0xa3, 0x24, 0x14, 0xa7,
	0x4a, 0x37, 0xb0, 0xa5, 0xd3, 0xe0, 0xaf, 0xa6, 0x14, 0x15, 0x83, 0x41, 0xce, 0xf0, 0xd1, 0xf0,
	0x8c, 0x6a, 0x68, 0x92, 0xa1, 0x3c, 0x29, 0x16, 0x51, 0xaf, 0x1b, 0xb2, 0xc4, 0x69, 0xfc, 0x90,
	0xce, 0x32, 0x94, 0xcf, 0x8a, 0x41, 0x84, 0x23, 0x4e, 0x49, 0x0d, 0xd6, 0xd1, 0x34, 0x4b, 0xc2,
	0x91, 0x5f, 0x49, 0x72, 0x7c, 0xfa, 0xa8, 0x74, 0xf9, 0xa4, 0x76, 0xd0, 0xdc, 0x82, 0x92, 0xad,
	0xa7, 0xf3, 0xd8, 0x3c, 0xb4, 0x57, 0xeb, 0xcf, 0x8e, 0xa1, 0x7d, 0xc7, 0xd0, 0x77, 0xc7, 0xd0,
	0x7b, 0xcf, 0x26, 0xfb, 0x9e, 0x4d, 0xbe, 0x7a, 0x36, 0xb9, 0x17, 0x52, 0xf9, 0xf6, 0xa5, 0xe2,
	0xb5, 0x79, 0x8e, 0xd8, 0xae, 0x0e, 0x08, 0x6e, 0xff, 0x19, 0xfa, 0xb7, 0x0d, 0xb8, 0x6a, 0x1e,
	0x51, 0x5c, 0xff, 0x04, 0x00, 0x00, 0xff, 0xff, 0x59, 0x6a, 0x82, 0x3e, 0x6f, 0x01, 0x00, 0x00,
}

func (m *ChainNonces) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ChainNonces) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ChainNonces) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.FinalizedHeight != 0 {
		i = encodeVarintChainNonces(dAtA, i, uint64(m.FinalizedHeight))
		i--
		dAtA[i] = 0x30
	}
	if len(m.Signers) > 0 {
		for iNdEx := len(m.Signers) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Signers[iNdEx])
			copy(dAtA[i:], m.Signers[iNdEx])
			i = encodeVarintChainNonces(dAtA, i, uint64(len(m.Signers[iNdEx])))
			i--
			dAtA[i] = 0x2a
		}
	}
	if m.Nonce != 0 {
		i = encodeVarintChainNonces(dAtA, i, uint64(m.Nonce))
		i--
		dAtA[i] = 0x20
	}
	if m.ChainId != 0 {
		i = encodeVarintChainNonces(dAtA, i, uint64(m.ChainId))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Index) > 0 {
		i -= len(m.Index)
		copy(dAtA[i:], m.Index)
		i = encodeVarintChainNonces(dAtA, i, uint64(len(m.Index)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintChainNonces(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintChainNonces(dAtA []byte, offset int, v uint64) int {
	offset -= sovChainNonces(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ChainNonces) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovChainNonces(uint64(l))
	}
	l = len(m.Index)
	if l > 0 {
		n += 1 + l + sovChainNonces(uint64(l))
	}
	if m.ChainId != 0 {
		n += 1 + sovChainNonces(uint64(m.ChainId))
	}
	if m.Nonce != 0 {
		n += 1 + sovChainNonces(uint64(m.Nonce))
	}
	if len(m.Signers) > 0 {
		for _, s := range m.Signers {
			l = len(s)
			n += 1 + l + sovChainNonces(uint64(l))
		}
	}
	if m.FinalizedHeight != 0 {
		n += 1 + sovChainNonces(uint64(m.FinalizedHeight))
	}
	return n
}

func sovChainNonces(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozChainNonces(x uint64) (n int) {
	return sovChainNonces(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ChainNonces) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowChainNonces
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
			return fmt.Errorf("proto: ChainNonces: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ChainNonces: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChainNonces
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
				return ErrInvalidLengthChainNonces
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthChainNonces
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChainNonces
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
				return ErrInvalidLengthChainNonces
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthChainNonces
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Index = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainId", wireType)
			}
			m.ChainId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChainNonces
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ChainId |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nonce", wireType)
			}
			m.Nonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChainNonces
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Nonce |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signers", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChainNonces
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
				return ErrInvalidLengthChainNonces
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthChainNonces
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signers = append(m.Signers, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FinalizedHeight", wireType)
			}
			m.FinalizedHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChainNonces
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FinalizedHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipChainNonces(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthChainNonces
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
func skipChainNonces(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowChainNonces
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
					return 0, ErrIntOverflowChainNonces
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
					return 0, ErrIntOverflowChainNonces
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
				return 0, ErrInvalidLengthChainNonces
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupChainNonces
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthChainNonces
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthChainNonces        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowChainNonces          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupChainNonces = fmt.Errorf("proto: unexpected end of group")
)
