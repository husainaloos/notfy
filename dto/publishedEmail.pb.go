// Code generated by protoc-gen-go. DO NOT EDIT.
// source: publishedEmail.proto

package dto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type PublishedEmail struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	From                 string   `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	To                   []string `protobuf:"bytes,3,rep,name=to,proto3" json:"to,omitempty"`
	Cc                   []string `protobuf:"bytes,4,rep,name=cc,proto3" json:"cc,omitempty"`
	Bcc                  []string `protobuf:"bytes,5,rep,name=bcc,proto3" json:"bcc,omitempty"`
	Subject              string   `protobuf:"bytes,6,opt,name=subject,proto3" json:"subject,omitempty"`
	Body                 string   `protobuf:"bytes,7,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PublishedEmail) Reset()         { *m = PublishedEmail{} }
func (m *PublishedEmail) String() string { return proto.CompactTextString(m) }
func (*PublishedEmail) ProtoMessage()    {}
func (*PublishedEmail) Descriptor() ([]byte, []int) {
	return fileDescriptor_aafcfb180c438fef, []int{0}
}

func (m *PublishedEmail) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PublishedEmail.Unmarshal(m, b)
}
func (m *PublishedEmail) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PublishedEmail.Marshal(b, m, deterministic)
}
func (m *PublishedEmail) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PublishedEmail.Merge(m, src)
}
func (m *PublishedEmail) XXX_Size() int {
	return xxx_messageInfo_PublishedEmail.Size(m)
}
func (m *PublishedEmail) XXX_DiscardUnknown() {
	xxx_messageInfo_PublishedEmail.DiscardUnknown(m)
}

var xxx_messageInfo_PublishedEmail proto.InternalMessageInfo

func (m *PublishedEmail) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *PublishedEmail) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

func (m *PublishedEmail) GetTo() []string {
	if m != nil {
		return m.To
	}
	return nil
}

func (m *PublishedEmail) GetCc() []string {
	if m != nil {
		return m.Cc
	}
	return nil
}

func (m *PublishedEmail) GetBcc() []string {
	if m != nil {
		return m.Bcc
	}
	return nil
}

func (m *PublishedEmail) GetSubject() string {
	if m != nil {
		return m.Subject
	}
	return ""
}

func (m *PublishedEmail) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

func init() {
	proto.RegisterType((*PublishedEmail)(nil), "dto.PublishedEmail")
}

func init() { proto.RegisterFile("publishedEmail.proto", fileDescriptor_aafcfb180c438fef) }

var fileDescriptor_aafcfb180c438fef = []byte{
	// 162 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x29, 0x28, 0x4d, 0xca,
	0xc9, 0x2c, 0xce, 0x48, 0x4d, 0x71, 0xcd, 0x4d, 0xcc, 0xcc, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0x62, 0x4e, 0x29, 0xc9, 0x57, 0x9a, 0xc2, 0xc8, 0xc5, 0x17, 0x80, 0x22, 0x2b, 0xc4, 0xc7,
	0xc5, 0x94, 0x99, 0x22, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x1c, 0xc4, 0x94, 0x99, 0x22, 0x24, 0xc4,
	0xc5, 0x92, 0x56, 0x94, 0x9f, 0x2b, 0xc1, 0xa4, 0xc0, 0xa8, 0xc1, 0x19, 0x04, 0x66, 0x83, 0xd4,
	0x94, 0xe4, 0x4b, 0x30, 0x2b, 0x30, 0x6b, 0x70, 0x06, 0x31, 0x95, 0xe4, 0x83, 0xf8, 0xc9, 0xc9,
	0x12, 0x2c, 0x10, 0x7e, 0x72, 0xb2, 0x90, 0x00, 0x17, 0x73, 0x52, 0x72, 0xb2, 0x04, 0x2b, 0x58,
	0x00, 0xc4, 0x14, 0x92, 0xe0, 0x62, 0x2f, 0x2e, 0x4d, 0xca, 0x4a, 0x4d, 0x2e, 0x91, 0x60, 0x03,
	0x1b, 0x04, 0xe3, 0x82, 0xcc, 0x4f, 0xca, 0x4f, 0xa9, 0x94, 0x60, 0x87, 0x98, 0x0f, 0x62, 0x27,
	0xb1, 0x81, 0x9d, 0x68, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xa4, 0x1b, 0x88, 0xad, 0xba, 0x00,
	0x00, 0x00,
}
