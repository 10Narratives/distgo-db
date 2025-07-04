// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: worker/database/v1/database_service.proto

package dbv1

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Database struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	DisplayName   string                 `protobuf:"bytes,2,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Database) Reset() {
	*x = Database{}
	mi := &file_worker_database_v1_database_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Database) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Database) ProtoMessage() {}

func (x *Database) ProtoReflect() protoreflect.Message {
	mi := &file_worker_database_v1_database_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Database.ProtoReflect.Descriptor instead.
func (*Database) Descriptor() ([]byte, []int) {
	return file_worker_database_v1_database_service_proto_rawDescGZIP(), []int{0}
}

func (x *Database) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Database) GetDisplayName() string {
	if x != nil {
		return x.DisplayName
	}
	return ""
}

func (x *Database) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Database) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type ListDatabasesRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PageSize      int32                  `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	PageToken     string                 `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListDatabasesRequest) Reset() {
	*x = ListDatabasesRequest{}
	mi := &file_worker_database_v1_database_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListDatabasesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListDatabasesRequest) ProtoMessage() {}

func (x *ListDatabasesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_worker_database_v1_database_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListDatabasesRequest.ProtoReflect.Descriptor instead.
func (*ListDatabasesRequest) Descriptor() ([]byte, []int) {
	return file_worker_database_v1_database_service_proto_rawDescGZIP(), []int{1}
}

func (x *ListDatabasesRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListDatabasesRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

type ListDatabasesResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Databases     []*Database            `protobuf:"bytes,1,rep,name=databases,proto3" json:"databases,omitempty"`
	NextPageToken string                 `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListDatabasesResponse) Reset() {
	*x = ListDatabasesResponse{}
	mi := &file_worker_database_v1_database_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListDatabasesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListDatabasesResponse) ProtoMessage() {}

func (x *ListDatabasesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_worker_database_v1_database_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListDatabasesResponse.ProtoReflect.Descriptor instead.
func (*ListDatabasesResponse) Descriptor() ([]byte, []int) {
	return file_worker_database_v1_database_service_proto_rawDescGZIP(), []int{2}
}

func (x *ListDatabasesResponse) GetDatabases() []*Database {
	if x != nil {
		return x.Databases
	}
	return nil
}

func (x *ListDatabasesResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

type GetDatabaseRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetDatabaseRequest) Reset() {
	*x = GetDatabaseRequest{}
	mi := &file_worker_database_v1_database_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetDatabaseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDatabaseRequest) ProtoMessage() {}

func (x *GetDatabaseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_worker_database_v1_database_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDatabaseRequest.ProtoReflect.Descriptor instead.
func (*GetDatabaseRequest) Descriptor() ([]byte, []int) {
	return file_worker_database_v1_database_service_proto_rawDescGZIP(), []int{3}
}

func (x *GetDatabaseRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type CreateDatabaseRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	DatabaseId    string                 `protobuf:"bytes,1,opt,name=database_id,json=databaseId,proto3" json:"database_id,omitempty"`
	Database      *Database              `protobuf:"bytes,2,opt,name=database,proto3" json:"database,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateDatabaseRequest) Reset() {
	*x = CreateDatabaseRequest{}
	mi := &file_worker_database_v1_database_service_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateDatabaseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateDatabaseRequest) ProtoMessage() {}

func (x *CreateDatabaseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_worker_database_v1_database_service_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateDatabaseRequest.ProtoReflect.Descriptor instead.
func (*CreateDatabaseRequest) Descriptor() ([]byte, []int) {
	return file_worker_database_v1_database_service_proto_rawDescGZIP(), []int{4}
}

func (x *CreateDatabaseRequest) GetDatabaseId() string {
	if x != nil {
		return x.DatabaseId
	}
	return ""
}

func (x *CreateDatabaseRequest) GetDatabase() *Database {
	if x != nil {
		return x.Database
	}
	return nil
}

type UpdateDatabaseRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Database      *Database              `protobuf:"bytes,1,opt,name=database,proto3" json:"database,omitempty"`
	UpdateMask    *fieldmaskpb.FieldMask `protobuf:"bytes,2,opt,name=update_mask,json=updateMask,proto3" json:"update_mask,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateDatabaseRequest) Reset() {
	*x = UpdateDatabaseRequest{}
	mi := &file_worker_database_v1_database_service_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateDatabaseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateDatabaseRequest) ProtoMessage() {}

func (x *UpdateDatabaseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_worker_database_v1_database_service_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateDatabaseRequest.ProtoReflect.Descriptor instead.
func (*UpdateDatabaseRequest) Descriptor() ([]byte, []int) {
	return file_worker_database_v1_database_service_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateDatabaseRequest) GetDatabase() *Database {
	if x != nil {
		return x.Database
	}
	return nil
}

func (x *UpdateDatabaseRequest) GetUpdateMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.UpdateMask
	}
	return nil
}

type DeleteDatabaseRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteDatabaseRequest) Reset() {
	*x = DeleteDatabaseRequest{}
	mi := &file_worker_database_v1_database_service_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteDatabaseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteDatabaseRequest) ProtoMessage() {}

func (x *DeleteDatabaseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_worker_database_v1_database_service_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteDatabaseRequest.ProtoReflect.Descriptor instead.
func (*DeleteDatabaseRequest) Descriptor() ([]byte, []int) {
	return file_worker_database_v1_database_service_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteDatabaseRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_worker_database_v1_database_service_proto protoreflect.FileDescriptor

const file_worker_database_v1_database_service_proto_rawDesc = "" +
	"\n" +
	")worker/database/v1/database_service.proto\x12\x12worker.database.v1\x1a\x1cgoogle/api/annotations.proto\x1a\x1fgoogle/api/field_behavior.proto\x1a\x19google/api/resource.proto\x1a\x1bgoogle/protobuf/empty.proto\x1a\x1fgoogle/protobuf/timestamp.proto\x1a google/protobuf/field_mask.proto\x1a\x17validate/validate.proto\"\x8d\x02\n" +
	"\bDatabase\x12\x17\n" +
	"\x04name\x18\x01 \x01(\tB\x03\xe0A\x03R\x04name\x120\n" +
	"\fdisplay_name\x18\x02 \x01(\tB\r\xe0A\x02\xfaB\ar\x05\x10\x01\x18\xff\x01R\vdisplayName\x12>\n" +
	"\n" +
	"created_at\x18\x03 \x01(\v2\x1a.google.protobuf.TimestampB\x03\xe0A\x03R\tcreatedAt\x12>\n" +
	"\n" +
	"updated_at\x18\x04 \x01(\v2\x1a.google.protobuf.TimestampB\x03\xe0A\x03R\tupdatedAt:6\xeaA3\n" +
	"\x1bworker.database.v1/Database\x12\x14databases/{database}\"^\n" +
	"\x14ListDatabasesRequest\x12'\n" +
	"\tpage_size\x18\x01 \x01(\x05B\n" +
	"\xfaB\a\x1a\x05\x18\xe8\a(\x00R\bpageSize\x12\x1d\n" +
	"\n" +
	"page_token\x18\x02 \x01(\tR\tpageToken\"{\n" +
	"\x15ListDatabasesResponse\x12:\n" +
	"\tdatabases\x18\x01 \x03(\v2\x1c.worker.database.v1.DatabaseR\tdatabases\x12&\n" +
	"\x0fnext_page_token\x18\x02 \x01(\tR\rnextPageToken\"G\n" +
	"\x12GetDatabaseRequest\x121\n" +
	"\x04name\x18\x01 \x01(\tB\x1d\xe0A\x02\xfaB\x17r\x152\x13^databases\\/[^\\/]+$R\x04name\"\xa1\x01\n" +
	"\x15CreateDatabaseRequest\x12A\n" +
	"\vdatabase_id\x18\x01 \x01(\tB \xe0A\x01\xfaB\x1ar\x18\x10\x01\x18@2\x0f^[a-z0-9\\-_.]*$\xd0\x01\x01R\n" +
	"databaseId\x12E\n" +
	"\bdatabase\x18\x02 \x01(\v2\x1c.worker.database.v1.DatabaseB\v\xe0A\x02\xfaB\x05\x8a\x01\x02\x10\x01R\bdatabase\"\x9b\x01\n" +
	"\x15UpdateDatabaseRequest\x12E\n" +
	"\bdatabase\x18\x01 \x01(\v2\x1c.worker.database.v1.DatabaseB\v\xe0A\x02\xfaB\x05\x8a\x01\x02\x10\x01R\bdatabase\x12;\n" +
	"\vupdate_mask\x18\x02 \x01(\v2\x1a.google.protobuf.FieldMaskR\n" +
	"updateMask\"J\n" +
	"\x15DeleteDatabaseRequest\x121\n" +
	"\x04name\x18\x01 \x01(\tB\x1d\xe0A\x02\xfaB\x17r\x152\x13^databases\\/[^\\/]+$R\x04name2\xa3\x05\n" +
	"\x0fDatabaseService\x12\x81\x01\n" +
	"\rListDatabases\x12(.worker.database.v1.ListDatabasesRequest\x1a).worker.database.v1.ListDatabasesResponse\"\x1b\x82\xd3\xe4\x93\x02\x15\x12\x13/v1alpha1/databases\x12y\n" +
	"\vGetDatabase\x12&.worker.database.v1.GetDatabaseRequest\x1a\x1c.worker.database.v1.Database\"$\x82\xd3\xe4\x93\x02\x1e\x12\x1c/v1alpha1/{name=databases/*}\x12\x80\x01\n" +
	"\x0eCreateDatabase\x12).worker.database.v1.CreateDatabaseRequest\x1a\x1c.worker.database.v1.Database\"%\x82\xd3\xe4\x93\x02\x1f:\bdatabase\"\x13/v1alpha1/databases\x12\x92\x01\n" +
	"\x0eUpdateDatabase\x12).worker.database.v1.UpdateDatabaseRequest\x1a\x1c.worker.database.v1.Database\"7\x82\xd3\xe4\x93\x021:\bdatabase2%/v1alpha1/{database.name=databases/*}\x12y\n" +
	"\x0eDeleteDatabase\x12).worker.database.v1.DeleteDatabaseRequest\x1a\x16.google.protobuf.Empty\"$\x82\xd3\xe4\x93\x02\x1e*\x1c/v1alpha1/{name=databases/*}Bf\n" +
	"\x1dcom.google.worker.database.v1P\x01ZCgithub.com/10Narratives/distgo-db/pkg/proto/worker/database/v1;dbv1b\x06proto3"

var (
	file_worker_database_v1_database_service_proto_rawDescOnce sync.Once
	file_worker_database_v1_database_service_proto_rawDescData []byte
)

func file_worker_database_v1_database_service_proto_rawDescGZIP() []byte {
	file_worker_database_v1_database_service_proto_rawDescOnce.Do(func() {
		file_worker_database_v1_database_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_worker_database_v1_database_service_proto_rawDesc), len(file_worker_database_v1_database_service_proto_rawDesc)))
	})
	return file_worker_database_v1_database_service_proto_rawDescData
}

var file_worker_database_v1_database_service_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_worker_database_v1_database_service_proto_goTypes = []any{
	(*Database)(nil),              // 0: worker.database.v1.Database
	(*ListDatabasesRequest)(nil),  // 1: worker.database.v1.ListDatabasesRequest
	(*ListDatabasesResponse)(nil), // 2: worker.database.v1.ListDatabasesResponse
	(*GetDatabaseRequest)(nil),    // 3: worker.database.v1.GetDatabaseRequest
	(*CreateDatabaseRequest)(nil), // 4: worker.database.v1.CreateDatabaseRequest
	(*UpdateDatabaseRequest)(nil), // 5: worker.database.v1.UpdateDatabaseRequest
	(*DeleteDatabaseRequest)(nil), // 6: worker.database.v1.DeleteDatabaseRequest
	(*timestamppb.Timestamp)(nil), // 7: google.protobuf.Timestamp
	(*fieldmaskpb.FieldMask)(nil), // 8: google.protobuf.FieldMask
	(*emptypb.Empty)(nil),         // 9: google.protobuf.Empty
}
var file_worker_database_v1_database_service_proto_depIdxs = []int32{
	7,  // 0: worker.database.v1.Database.created_at:type_name -> google.protobuf.Timestamp
	7,  // 1: worker.database.v1.Database.updated_at:type_name -> google.protobuf.Timestamp
	0,  // 2: worker.database.v1.ListDatabasesResponse.databases:type_name -> worker.database.v1.Database
	0,  // 3: worker.database.v1.CreateDatabaseRequest.database:type_name -> worker.database.v1.Database
	0,  // 4: worker.database.v1.UpdateDatabaseRequest.database:type_name -> worker.database.v1.Database
	8,  // 5: worker.database.v1.UpdateDatabaseRequest.update_mask:type_name -> google.protobuf.FieldMask
	1,  // 6: worker.database.v1.DatabaseService.ListDatabases:input_type -> worker.database.v1.ListDatabasesRequest
	3,  // 7: worker.database.v1.DatabaseService.GetDatabase:input_type -> worker.database.v1.GetDatabaseRequest
	4,  // 8: worker.database.v1.DatabaseService.CreateDatabase:input_type -> worker.database.v1.CreateDatabaseRequest
	5,  // 9: worker.database.v1.DatabaseService.UpdateDatabase:input_type -> worker.database.v1.UpdateDatabaseRequest
	6,  // 10: worker.database.v1.DatabaseService.DeleteDatabase:input_type -> worker.database.v1.DeleteDatabaseRequest
	2,  // 11: worker.database.v1.DatabaseService.ListDatabases:output_type -> worker.database.v1.ListDatabasesResponse
	0,  // 12: worker.database.v1.DatabaseService.GetDatabase:output_type -> worker.database.v1.Database
	0,  // 13: worker.database.v1.DatabaseService.CreateDatabase:output_type -> worker.database.v1.Database
	0,  // 14: worker.database.v1.DatabaseService.UpdateDatabase:output_type -> worker.database.v1.Database
	9,  // 15: worker.database.v1.DatabaseService.DeleteDatabase:output_type -> google.protobuf.Empty
	11, // [11:16] is the sub-list for method output_type
	6,  // [6:11] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_worker_database_v1_database_service_proto_init() }
func file_worker_database_v1_database_service_proto_init() {
	if File_worker_database_v1_database_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_worker_database_v1_database_service_proto_rawDesc), len(file_worker_database_v1_database_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_worker_database_v1_database_service_proto_goTypes,
		DependencyIndexes: file_worker_database_v1_database_service_proto_depIdxs,
		MessageInfos:      file_worker_database_v1_database_service_proto_msgTypes,
	}.Build()
	File_worker_database_v1_database_service_proto = out.File
	file_worker_database_v1_database_service_proto_goTypes = nil
	file_worker_database_v1_database_service_proto_depIdxs = nil
}
