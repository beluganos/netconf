# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: rpcapi.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
from google.protobuf import descriptor_pb2
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='rpcapi.proto',
  package='cfgrpcapi',
  syntax='proto3',
  serialized_pb=_b('\n\x0crpcapi.proto\x12\tcfgrpcapi\".\n\x05Shell\x12\x0b\n\x03\x63md\x18\x01 \x01(\t\x12\x0c\n\x04\x61rgs\x18\x02 \x03(\t\x12\n\n\x02in\x18\x03 \x01(\x0c\"\x18\n\x06Result\x12\x0e\n\x06output\x18\x01 \x01(\x0c\"2\n\x0e\x45xecuteRequest\x12 \n\x06shells\x18\x01 \x03(\x0b\x32\x10.cfgrpcapi.Shell\"2\n\x0c\x45xecuteReply\x12\"\n\x07results\x18\x01 \x03(\x0b\x32\x11.cfgrpcapi.Result2I\n\x06RpcApi\x12?\n\x07\x45xecute\x12\x19.cfgrpcapi.ExecuteRequest\x1a\x17.cfgrpcapi.ExecuteReply\"\x00\x62\x06proto3')
)




_SHELL = _descriptor.Descriptor(
  name='Shell',
  full_name='cfgrpcapi.Shell',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='cmd', full_name='cfgrpcapi.Shell.cmd', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='args', full_name='cfgrpcapi.Shell.args', index=1,
      number=2, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='in', full_name='cfgrpcapi.Shell.in', index=2,
      number=3, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=_b(""),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=27,
  serialized_end=73,
)


_RESULT = _descriptor.Descriptor(
  name='Result',
  full_name='cfgrpcapi.Result',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='output', full_name='cfgrpcapi.Result.output', index=0,
      number=1, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=_b(""),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=75,
  serialized_end=99,
)


_EXECUTEREQUEST = _descriptor.Descriptor(
  name='ExecuteRequest',
  full_name='cfgrpcapi.ExecuteRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='shells', full_name='cfgrpcapi.ExecuteRequest.shells', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=101,
  serialized_end=151,
)


_EXECUTEREPLY = _descriptor.Descriptor(
  name='ExecuteReply',
  full_name='cfgrpcapi.ExecuteReply',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='results', full_name='cfgrpcapi.ExecuteReply.results', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=153,
  serialized_end=203,
)

_EXECUTEREQUEST.fields_by_name['shells'].message_type = _SHELL
_EXECUTEREPLY.fields_by_name['results'].message_type = _RESULT
DESCRIPTOR.message_types_by_name['Shell'] = _SHELL
DESCRIPTOR.message_types_by_name['Result'] = _RESULT
DESCRIPTOR.message_types_by_name['ExecuteRequest'] = _EXECUTEREQUEST
DESCRIPTOR.message_types_by_name['ExecuteReply'] = _EXECUTEREPLY
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

Shell = _reflection.GeneratedProtocolMessageType('Shell', (_message.Message,), dict(
  DESCRIPTOR = _SHELL,
  __module__ = 'rpcapi_pb2'
  # @@protoc_insertion_point(class_scope:cfgrpcapi.Shell)
  ))
_sym_db.RegisterMessage(Shell)

Result = _reflection.GeneratedProtocolMessageType('Result', (_message.Message,), dict(
  DESCRIPTOR = _RESULT,
  __module__ = 'rpcapi_pb2'
  # @@protoc_insertion_point(class_scope:cfgrpcapi.Result)
  ))
_sym_db.RegisterMessage(Result)

ExecuteRequest = _reflection.GeneratedProtocolMessageType('ExecuteRequest', (_message.Message,), dict(
  DESCRIPTOR = _EXECUTEREQUEST,
  __module__ = 'rpcapi_pb2'
  # @@protoc_insertion_point(class_scope:cfgrpcapi.ExecuteRequest)
  ))
_sym_db.RegisterMessage(ExecuteRequest)

ExecuteReply = _reflection.GeneratedProtocolMessageType('ExecuteReply', (_message.Message,), dict(
  DESCRIPTOR = _EXECUTEREPLY,
  __module__ = 'rpcapi_pb2'
  # @@protoc_insertion_point(class_scope:cfgrpcapi.ExecuteReply)
  ))
_sym_db.RegisterMessage(ExecuteReply)



_RPCAPI = _descriptor.ServiceDescriptor(
  name='RpcApi',
  full_name='cfgrpcapi.RpcApi',
  file=DESCRIPTOR,
  index=0,
  options=None,
  serialized_start=205,
  serialized_end=278,
  methods=[
  _descriptor.MethodDescriptor(
    name='Execute',
    full_name='cfgrpcapi.RpcApi.Execute',
    index=0,
    containing_service=None,
    input_type=_EXECUTEREQUEST,
    output_type=_EXECUTEREPLY,
    options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_RPCAPI)

DESCRIPTOR.services_by_name['RpcApi'] = _RPCAPI

# @@protoc_insertion_point(module_scope)
