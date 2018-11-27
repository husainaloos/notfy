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

//go:generate protoc -I=. --go_out=. *.proto
type PublishedEmail struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	From                 string   `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	To                   []string `protobuf:"bytes,3,rep,name=to,proto3" json:"to,omitempty"`
	Cc                   []string `protobuf:"bytes,4,rep,name=cc,proto3" json:"cc,omitempty"`
	Bcc                  []string `protobuf:"bytes,5,rep,name=bcc,proto3" json:"bcc,omitempty"`
	Subject              string   `protobuf:"bytes,6,opt,name=subject,proto3" json:"subject,omitempty"`
	Body                 string   `protobuf:"bytes,7,opt,name=body,proto3" json:"body,omitempty"`
	PublishedAt          int64    `protobuf:"varint,8,opt,name=publishedAt,proto3" json:"publishedAt,omitempty"`
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

func (m *PublishedEmail) GetPublishedAt() int64 {
	if m != nil {
		return m.PublishedAt
	}
	return 0
}

func init() {
	proto.RegisterType((*PublishedEmail)(nil), "dto.PublishedEmail")
}

func init() { proto.RegisterFile("publishedEmail.proto", fileDescriptor_aafcfb180c438fef) }

var fileDescriptor_aafcfb180c438fef = []byte{
	// 176 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x8f, 0xb1, 0x0a, 0xc2, 0x30,
	0x10, 0x86, 0x49, 0x52, 0x5b, 0x1b, 0xa1, 0xc8, 0xe1, 0x70, 0x63, 0x70, 0xea, 0xe4, 0xe2, 0x13,
	0x38, 0xb8, 0x4b, 0xdf, 0xc0, 0x5c, 0x2a, 0x46, 0x5a, 0xae, 0xb4, 0xe9, 0xe0, 0x8b, 0xf9, 0x7c,
	0x92, 0x14, 0x45, 0xb7, 0xff, 0xfb, 0x13, 0x3e, 0xfe, 0xd3, 0xbb, 0x61, 0xb6, 0x9d, 0x9f, 0xee,
	0xad, 0x3b, 0xf7, 0x57, 0xdf, 0x1d, 0x86, 0x91, 0x03, 0x83, 0x72, 0x81, 0xf7, 0x2f, 0xa1, 0xab,
	0xcb, 0xdf, 0x2b, 0x54, 0x5a, 0x7a, 0x87, 0xc2, 0x88, 0x5a, 0x35, 0xd2, 0x3b, 0x00, 0x9d, 0xdd,
	0x46, 0xee, 0x51, 0x1a, 0x51, 0x97, 0x4d, 0xca, 0xf1, 0x4f, 0x60, 0x54, 0x46, 0xd5, 0x65, 0x23,
	0x03, 0x47, 0x26, 0xc2, 0x6c, 0x61, 0x22, 0xd8, 0x6a, 0x65, 0x89, 0x70, 0x95, 0x8a, 0x18, 0x01,
	0x75, 0x31, 0xcd, 0xf6, 0xd1, 0x52, 0xc0, 0x3c, 0x89, 0x3e, 0x18, 0xfd, 0x96, 0xdd, 0x13, 0x8b,
	0xc5, 0x1f, 0x33, 0x18, 0xbd, 0xf9, 0x6e, 0x3e, 0x05, 0x5c, 0xa7, 0x31, 0xbf, 0x95, 0xcd, 0xd3,
	0x11, 0xc7, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x0b, 0xdd, 0xa6, 0xa2, 0xdc, 0x00, 0x00, 0x00,
}
