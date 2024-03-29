// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: internal/proto/goshort.proto

package proto

import (
	empty "github.com/golang/protobuf/ptypes/empty"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ShortenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Origin string `protobuf:"bytes,1,opt,name=origin,json=url,proto3" json:"origin,omitempty"`
}

func (x *ShortenRequest) Reset() {
	*x = ShortenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_goshort_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShortenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortenRequest) ProtoMessage() {}

func (x *ShortenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_goshort_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortenRequest.ProtoReflect.Descriptor instead.
func (*ShortenRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_goshort_proto_rawDescGZIP(), []int{0}
}

func (x *ShortenRequest) GetOrigin() string {
	if x != nil {
		return x.Origin
	}
	return ""
}

type ShortenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result  string `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	Existed bool   `protobuf:"varint,2,opt,name=existed,proto3" json:"existed,omitempty"`
}

func (x *ShortenResponse) Reset() {
	*x = ShortenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_goshort_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShortenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortenResponse) ProtoMessage() {}

func (x *ShortenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_goshort_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortenResponse.ProtoReflect.Descriptor instead.
func (*ShortenResponse) Descriptor() ([]byte, []int) {
	return file_internal_proto_goshort_proto_rawDescGZIP(), []int{1}
}

func (x *ShortenResponse) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

func (x *ShortenResponse) GetExisted() bool {
	if x != nil {
		return x.Existed
	}
	return false
}

type DeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UrlIDs []string `protobuf:"bytes,1,rep,name=UrlIDs,proto3" json:"UrlIDs,omitempty"`
}

func (x *DeleteRequest) Reset() {
	*x = DeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_goshort_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRequest) ProtoMessage() {}

func (x *DeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_goshort_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRequest.ProtoReflect.Descriptor instead.
func (*DeleteRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_goshort_proto_rawDescGZIP(), []int{2}
}

func (x *DeleteRequest) GetUrlIDs() []string {
	if x != nil {
		return x.UrlIDs
	}
	return nil
}

type Entry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result string `protobuf:"bytes,1,opt,name=result,json=short_url,proto3" json:"result,omitempty"`
	Origin string `protobuf:"bytes,2,opt,name=origin,json=original_url,proto3" json:"origin,omitempty"`
}

func (x *Entry) Reset() {
	*x = Entry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_goshort_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Entry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entry) ProtoMessage() {}

func (x *Entry) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_goshort_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entry.ProtoReflect.Descriptor instead.
func (*Entry) Descriptor() ([]byte, []int) {
	return file_internal_proto_goshort_proto_rawDescGZIP(), []int{3}
}

func (x *Entry) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

func (x *Entry) GetOrigin() string {
	if x != nil {
		return x.Origin
	}
	return ""
}

type Entries struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Entries []*Entry `protobuf:"bytes,1,rep,name=entries,proto3" json:"entries,omitempty"`
}

func (x *Entries) Reset() {
	*x = Entries{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_goshort_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Entries) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entries) ProtoMessage() {}

func (x *Entries) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_goshort_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entries.ProtoReflect.Descriptor instead.
func (*Entries) Descriptor() ([]byte, []int) {
	return file_internal_proto_goshort_proto_rawDescGZIP(), []int{4}
}

func (x *Entries) GetEntries() []*Entry {
	if x != nil {
		return x.Entries
	}
	return nil
}

type RestoreRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *RestoreRequest) Reset() {
	*x = RestoreRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_goshort_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RestoreRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RestoreRequest) ProtoMessage() {}

func (x *RestoreRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_goshort_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RestoreRequest.ProtoReflect.Descriptor instead.
func (*RestoreRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_goshort_proto_rawDescGZIP(), []int{5}
}

func (x *RestoreRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type RestoreResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Origin  string `protobuf:"bytes,1,opt,name=origin,proto3" json:"origin,omitempty"`
	Deleted bool   `protobuf:"varint,2,opt,name=deleted,proto3" json:"deleted,omitempty"`
}

func (x *RestoreResponse) Reset() {
	*x = RestoreResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_goshort_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RestoreResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RestoreResponse) ProtoMessage() {}

func (x *RestoreResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_goshort_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RestoreResponse.ProtoReflect.Descriptor instead.
func (*RestoreResponse) Descriptor() ([]byte, []int) {
	return file_internal_proto_goshort_proto_rawDescGZIP(), []int{6}
}

func (x *RestoreResponse) GetOrigin() string {
	if x != nil {
		return x.Origin
	}
	return ""
}

func (x *RestoreResponse) GetDeleted() bool {
	if x != nil {
		return x.Deleted
	}
	return false
}

type BatchShortenRequestEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string `protobuf:"bytes,1,opt,name=id,json=correlation_id,proto3" json:"id,omitempty"`
	Origin string `protobuf:"bytes,2,opt,name=origin,json=original_url,proto3" json:"origin,omitempty"`
}

func (x *BatchShortenRequestEntry) Reset() {
	*x = BatchShortenRequestEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_goshort_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BatchShortenRequestEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchShortenRequestEntry) ProtoMessage() {}

func (x *BatchShortenRequestEntry) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_goshort_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchShortenRequestEntry.ProtoReflect.Descriptor instead.
func (*BatchShortenRequestEntry) Descriptor() ([]byte, []int) {
	return file_internal_proto_goshort_proto_rawDescGZIP(), []int{7}
}

func (x *BatchShortenRequestEntry) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *BatchShortenRequestEntry) GetOrigin() string {
	if x != nil {
		return x.Origin
	}
	return ""
}

type BatchShortenResponseEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string `protobuf:"bytes,1,opt,name=id,json=correlation_id,proto3" json:"id,omitempty"`
	Result string `protobuf:"bytes,2,opt,name=result,json=short_url,proto3" json:"result,omitempty"`
}

func (x *BatchShortenResponseEntry) Reset() {
	*x = BatchShortenResponseEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_goshort_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BatchShortenResponseEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchShortenResponseEntry) ProtoMessage() {}

func (x *BatchShortenResponseEntry) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_goshort_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchShortenResponseEntry.ProtoReflect.Descriptor instead.
func (*BatchShortenResponseEntry) Descriptor() ([]byte, []int) {
	return file_internal_proto_goshort_proto_rawDescGZIP(), []int{8}
}

func (x *BatchShortenResponseEntry) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *BatchShortenResponseEntry) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

type BatchShortenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Entries []*BatchShortenRequestEntry `protobuf:"bytes,1,rep,name=Entries,proto3" json:"Entries,omitempty"`
}

func (x *BatchShortenRequest) Reset() {
	*x = BatchShortenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_goshort_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BatchShortenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchShortenRequest) ProtoMessage() {}

func (x *BatchShortenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_goshort_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchShortenRequest.ProtoReflect.Descriptor instead.
func (*BatchShortenRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_goshort_proto_rawDescGZIP(), []int{9}
}

func (x *BatchShortenRequest) GetEntries() []*BatchShortenRequestEntry {
	if x != nil {
		return x.Entries
	}
	return nil
}

type BatchShortenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Entries []*BatchShortenResponseEntry `protobuf:"bytes,1,rep,name=Entries,proto3" json:"Entries,omitempty"`
}

func (x *BatchShortenResponse) Reset() {
	*x = BatchShortenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_goshort_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BatchShortenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchShortenResponse) ProtoMessage() {}

func (x *BatchShortenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_goshort_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchShortenResponse.ProtoReflect.Descriptor instead.
func (*BatchShortenResponse) Descriptor() ([]byte, []int) {
	return file_internal_proto_goshort_proto_rawDescGZIP(), []int{10}
}

func (x *BatchShortenResponse) GetEntries() []*BatchShortenResponseEntry {
	if x != nil {
		return x.Entries
	}
	return nil
}

type Stats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Urls  int64 `protobuf:"varint,1,opt,name=urls,proto3" json:"urls,omitempty"`
	Users int64 `protobuf:"varint,2,opt,name=users,proto3" json:"users,omitempty"`
}

func (x *Stats) Reset() {
	*x = Stats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_goshort_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Stats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stats) ProtoMessage() {}

func (x *Stats) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_goshort_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Stats.ProtoReflect.Descriptor instead.
func (*Stats) Descriptor() ([]byte, []int) {
	return file_internal_proto_goshort_proto_rawDescGZIP(), []int{11}
}

func (x *Stats) GetUrls() int64 {
	if x != nil {
		return x.Urls
	}
	return 0
}

func (x *Stats) GetUsers() int64 {
	if x != nil {
		return x.Users
	}
	return 0
}

type RegisterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *RegisterResponse) Reset() {
	*x = RegisterResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_goshort_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterResponse) ProtoMessage() {}

func (x *RegisterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_goshort_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterResponse.ProtoReflect.Descriptor instead.
func (*RegisterResponse) Descriptor() ([]byte, []int) {
	return file_internal_proto_goshort_proto_rawDescGZIP(), []int{12}
}

func (x *RegisterResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

var File_internal_proto_goshort_proto protoreflect.FileDescriptor

var file_internal_proto_goshort_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x67, 0x6f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07,
	0x67, 0x6f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x25, 0x0a, 0x0e, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x13, 0x0a, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x43, 0x0a, 0x0f, 0x53,
	0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x78, 0x69, 0x73, 0x74, 0x65,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x65, 0x78, 0x69, 0x73, 0x74, 0x65, 0x64,
	0x22, 0x27, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x72, 0x6c, 0x49, 0x44, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x06, 0x55, 0x72, 0x6c, 0x49, 0x44, 0x73, 0x22, 0x40, 0x0a, 0x05, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x19, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x5f, 0x75, 0x72, 0x6c, 0x12, 0x1c, 0x0a,
	0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x6f,
	0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x5f, 0x75, 0x72, 0x6c, 0x22, 0x33, 0x0a, 0x07, 0x45,
	0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12, 0x28, 0x0a, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x67, 0x6f, 0x73, 0x68, 0x6f, 0x72,
	0x74, 0x2e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73,
	0x22, 0x20, 0x0a, 0x0e, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x22, 0x43, 0x0a, 0x0f, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x12, 0x18, 0x0a,
	0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07,
	0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x22, 0x54, 0x0a, 0x18, 0x42, 0x61, 0x74, 0x63, 0x68,
	0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x1a, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0e, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x12,
	0x1c, 0x0a, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x5f, 0x75, 0x72, 0x6c, 0x22, 0x52, 0x0a,
	0x19, 0x42, 0x61, 0x74, 0x63, 0x68, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x1a, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x12, 0x19, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x5f, 0x75, 0x72,
	0x6c, 0x22, 0x52, 0x0a, 0x13, 0x42, 0x61, 0x74, 0x63, 0x68, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3b, 0x0a, 0x07, 0x45, 0x6e, 0x74, 0x72,
	0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x67, 0x6f, 0x73, 0x68,
	0x6f, 0x72, 0x74, 0x2e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x45, 0x6e,
	0x74, 0x72, 0x69, 0x65, 0x73, 0x22, 0x54, 0x0a, 0x14, 0x42, 0x61, 0x74, 0x63, 0x68, 0x53, 0x68,
	0x6f, 0x72, 0x74, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3c, 0x0a,
	0x07, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22,
	0x2e, 0x67, 0x6f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x2e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x53, 0x68,
	0x6f, 0x72, 0x74, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x07, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x22, 0x31, 0x0a, 0x05, 0x53,
	0x74, 0x61, 0x74, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x72, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x04, 0x75, 0x72, 0x6c, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x22, 0x28,
	0x0a, 0x10, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x32, 0xb7, 0x03, 0x0a, 0x09, 0x53, 0x68, 0x6f,
	0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x12, 0x3c, 0x0a, 0x07, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65,
	0x6e, 0x12, 0x17, 0x2e, 0x67, 0x6f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x2e, 0x53, 0x68, 0x6f, 0x72,
	0x74, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x67, 0x6f, 0x73,
	0x68, 0x6f, 0x72, 0x74, 0x2e, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x16,
	0x2e, 0x67, 0x6f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x30,
	0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x10,
	0x2e, 0x67, 0x6f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x2e, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73,
	0x12, 0x3c, 0x0a, 0x07, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x17, 0x2e, 0x67, 0x6f,
	0x73, 0x68, 0x6f, 0x72, 0x74, 0x2e, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x67, 0x6f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x2e, 0x52,
	0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36,
	0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x4b, 0x0a, 0x0c, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65,
	0x6e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x73, 0x68, 0x6f, 0x72, 0x74,
	0x2e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x67, 0x6f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x2e, 0x42,
	0x61, 0x74, 0x63, 0x68, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x3d, 0x0a, 0x08, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x19, 0x2e, 0x67, 0x6f, 0x73, 0x68, 0x6f, 0x72,
	0x74, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x18, 0x5a, 0x16, 0x67, 0x6f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_proto_goshort_proto_rawDescOnce sync.Once
	file_internal_proto_goshort_proto_rawDescData = file_internal_proto_goshort_proto_rawDesc
)

func file_internal_proto_goshort_proto_rawDescGZIP() []byte {
	file_internal_proto_goshort_proto_rawDescOnce.Do(func() {
		file_internal_proto_goshort_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_proto_goshort_proto_rawDescData)
	})
	return file_internal_proto_goshort_proto_rawDescData
}

var file_internal_proto_goshort_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_internal_proto_goshort_proto_goTypes = []interface{}{
	(*ShortenRequest)(nil),            // 0: goshort.ShortenRequest
	(*ShortenResponse)(nil),           // 1: goshort.ShortenResponse
	(*DeleteRequest)(nil),             // 2: goshort.DeleteRequest
	(*Entry)(nil),                     // 3: goshort.Entry
	(*Entries)(nil),                   // 4: goshort.Entries
	(*RestoreRequest)(nil),            // 5: goshort.RestoreRequest
	(*RestoreResponse)(nil),           // 6: goshort.RestoreResponse
	(*BatchShortenRequestEntry)(nil),  // 7: goshort.BatchShortenRequestEntry
	(*BatchShortenResponseEntry)(nil), // 8: goshort.BatchShortenResponseEntry
	(*BatchShortenRequest)(nil),       // 9: goshort.BatchShortenRequest
	(*BatchShortenResponse)(nil),      // 10: goshort.BatchShortenResponse
	(*Stats)(nil),                     // 11: goshort.Stats
	(*RegisterResponse)(nil),          // 12: goshort.RegisterResponse
	(*empty.Empty)(nil),               // 13: google.protobuf.Empty
}
var file_internal_proto_goshort_proto_depIdxs = []int32{
	3,  // 0: goshort.Entries.entries:type_name -> goshort.Entry
	7,  // 1: goshort.BatchShortenRequest.Entries:type_name -> goshort.BatchShortenRequestEntry
	8,  // 2: goshort.BatchShortenResponse.Entries:type_name -> goshort.BatchShortenResponseEntry
	0,  // 3: goshort.Shortener.Shorten:input_type -> goshort.ShortenRequest
	2,  // 4: goshort.Shortener.Delete:input_type -> goshort.DeleteRequest
	13, // 5: goshort.Shortener.List:input_type -> google.protobuf.Empty
	5,  // 6: goshort.Shortener.Restore:input_type -> goshort.RestoreRequest
	13, // 7: goshort.Shortener.Ping:input_type -> google.protobuf.Empty
	9,  // 8: goshort.Shortener.ShortenBatch:input_type -> goshort.BatchShortenRequest
	13, // 9: goshort.Shortener.Register:input_type -> google.protobuf.Empty
	1,  // 10: goshort.Shortener.Shorten:output_type -> goshort.ShortenResponse
	13, // 11: goshort.Shortener.Delete:output_type -> google.protobuf.Empty
	4,  // 12: goshort.Shortener.List:output_type -> goshort.Entries
	6,  // 13: goshort.Shortener.Restore:output_type -> goshort.RestoreResponse
	13, // 14: goshort.Shortener.Ping:output_type -> google.protobuf.Empty
	10, // 15: goshort.Shortener.ShortenBatch:output_type -> goshort.BatchShortenResponse
	12, // 16: goshort.Shortener.Register:output_type -> goshort.RegisterResponse
	10, // [10:17] is the sub-list for method output_type
	3,  // [3:10] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_internal_proto_goshort_proto_init() }
func file_internal_proto_goshort_proto_init() {
	if File_internal_proto_goshort_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_proto_goshort_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShortenRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_proto_goshort_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShortenResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_proto_goshort_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_proto_goshort_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Entry); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_proto_goshort_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Entries); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_proto_goshort_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RestoreRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_proto_goshort_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RestoreResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_proto_goshort_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BatchShortenRequestEntry); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_proto_goshort_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BatchShortenResponseEntry); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_proto_goshort_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BatchShortenRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_proto_goshort_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BatchShortenResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_proto_goshort_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Stats); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_proto_goshort_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_proto_goshort_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_proto_goshort_proto_goTypes,
		DependencyIndexes: file_internal_proto_goshort_proto_depIdxs,
		MessageInfos:      file_internal_proto_goshort_proto_msgTypes,
	}.Build()
	File_internal_proto_goshort_proto = out.File
	file_internal_proto_goshort_proto_rawDesc = nil
	file_internal_proto_goshort_proto_goTypes = nil
	file_internal_proto_goshort_proto_depIdxs = nil
}
