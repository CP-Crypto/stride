// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: stride/stakeibc/min_validator_requirements.proto

package types

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
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

type MinValidatorRequirements struct {
	CommissionRate int32 `protobuf:"varint,1,opt,name=commission_rate,json=commissionRate,proto3" json:"commission_rate,omitempty"`
	Uptime         int32 `protobuf:"varint,2,opt,name=uptime,proto3" json:"uptime,omitempty"`
}

func (m *MinValidatorRequirements) Reset()         { *m = MinValidatorRequirements{} }
func (m *MinValidatorRequirements) String() string { return proto.CompactTextString(m) }
func (*MinValidatorRequirements) ProtoMessage()    {}
func (*MinValidatorRequirements) Descriptor() ([]byte, []int) {
	return fileDescriptor_c53a743ee7ad112e, []int{0}
}
func (m *MinValidatorRequirements) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MinValidatorRequirements) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MinValidatorRequirements.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MinValidatorRequirements) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MinValidatorRequirements.Merge(m, src)
}
func (m *MinValidatorRequirements) XXX_Size() int {
	return m.Size()
}
func (m *MinValidatorRequirements) XXX_DiscardUnknown() {
	xxx_messageInfo_MinValidatorRequirements.DiscardUnknown(m)
}

var xxx_messageInfo_MinValidatorRequirements proto.InternalMessageInfo

func (m *MinValidatorRequirements) GetCommissionRate() int32 {
	if m != nil {
		return m.CommissionRate
	}
	return 0
}

func (m *MinValidatorRequirements) GetUptime() int32 {
	if m != nil {
		return m.Uptime
	}
	return 0
}

func init() {
	proto.RegisterType((*MinValidatorRequirements)(nil), "stride.stakeibc.MinValidatorRequirements")
}

func init() {
	proto.RegisterFile("stride/stakeibc/min_validator_requirements.proto", fileDescriptor_c53a743ee7ad112e)
}

var fileDescriptor_c53a743ee7ad112e = []byte{
	// 211 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0x28, 0x2e, 0x29, 0xca,
	0x4c, 0x49, 0xd5, 0x2f, 0x2e, 0x49, 0xcc, 0x4e, 0xcd, 0x4c, 0x4a, 0xd6, 0xcf, 0xcd, 0xcc, 0x8b,
	0x2f, 0x4b, 0xcc, 0xc9, 0x4c, 0x49, 0x2c, 0xc9, 0x2f, 0x8a, 0x2f, 0x4a, 0x2d, 0x2c, 0xcd, 0x2c,
	0x4a, 0xcd, 0x4d, 0xcd, 0x2b, 0x29, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x87, 0xe8,
	0xd0, 0x83, 0xe9, 0x50, 0x8a, 0xe6, 0x92, 0xf0, 0xcd, 0xcc, 0x0b, 0x83, 0xe9, 0x09, 0x42, 0xd2,
	0x22, 0xa4, 0xce, 0xc5, 0x9f, 0x9c, 0x9f, 0x9b, 0x9b, 0x59, 0x5c, 0x9c, 0x99, 0x9f, 0x17, 0x5f,
	0x94, 0x58, 0x92, 0x2a, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x1a, 0xc4, 0x87, 0x10, 0x0e, 0x4a, 0x2c,
	0x49, 0x15, 0x12, 0xe3, 0x62, 0x2b, 0x2d, 0x28, 0xc9, 0xcc, 0x4d, 0x95, 0x60, 0x02, 0xcb, 0x43,
	0x79, 0x4e, 0xde, 0x27, 0x1e, 0xc9, 0x31, 0x5e, 0x78, 0x24, 0xc7, 0xf8, 0xe0, 0x91, 0x1c, 0xe3,
	0x84, 0xc7, 0x72, 0x0c, 0x17, 0x1e, 0xcb, 0x31, 0xdc, 0x78, 0x2c, 0xc7, 0x10, 0x65, 0x98, 0x9e,
	0x59, 0x92, 0x51, 0x9a, 0xa4, 0x97, 0x9c, 0x9f, 0xab, 0x1f, 0x0c, 0x76, 0x92, 0xae, 0x4f, 0x62,
	0x52, 0xb1, 0x3e, 0xd4, 0x43, 0x65, 0x26, 0xfa, 0x15, 0x08, 0x5f, 0x95, 0x54, 0x16, 0xa4, 0x16,
	0x27, 0xb1, 0x81, 0x7d, 0x60, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x78, 0x1b, 0xa7, 0xe6, 0xf5,
	0x00, 0x00, 0x00,
}

func (m *MinValidatorRequirements) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MinValidatorRequirements) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MinValidatorRequirements) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Uptime != 0 {
		i = encodeVarintMinValidatorRequirements(dAtA, i, uint64(m.Uptime))
		i--
		dAtA[i] = 0x10
	}
	if m.CommissionRate != 0 {
		i = encodeVarintMinValidatorRequirements(dAtA, i, uint64(m.CommissionRate))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintMinValidatorRequirements(dAtA []byte, offset int, v uint64) int {
	offset -= sovMinValidatorRequirements(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MinValidatorRequirements) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.CommissionRate != 0 {
		n += 1 + sovMinValidatorRequirements(uint64(m.CommissionRate))
	}
	if m.Uptime != 0 {
		n += 1 + sovMinValidatorRequirements(uint64(m.Uptime))
	}
	return n
}

func sovMinValidatorRequirements(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMinValidatorRequirements(x uint64) (n int) {
	return sovMinValidatorRequirements(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MinValidatorRequirements) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMinValidatorRequirements
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
			return fmt.Errorf("proto: MinValidatorRequirements: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MinValidatorRequirements: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CommissionRate", wireType)
			}
			m.CommissionRate = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMinValidatorRequirements
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CommissionRate |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uptime", wireType)
			}
			m.Uptime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMinValidatorRequirements
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Uptime |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipMinValidatorRequirements(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMinValidatorRequirements
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
func skipMinValidatorRequirements(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMinValidatorRequirements
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
					return 0, ErrIntOverflowMinValidatorRequirements
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
					return 0, ErrIntOverflowMinValidatorRequirements
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
				return 0, ErrInvalidLengthMinValidatorRequirements
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMinValidatorRequirements
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMinValidatorRequirements
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMinValidatorRequirements        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMinValidatorRequirements          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMinValidatorRequirements = fmt.Errorf("proto: unexpected end of group")
)
