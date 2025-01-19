// @generated by protoc-gen-es v2.2.3 with parameter "target=ts"
// @generated from file v1/test_service.proto (package v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import { file_google_api_annotations } from "../google/api/annotations_pb";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file v1/test_service.proto.
 */
export const file_v1_test_service: GenFile = /*@__PURE__*/
  fileDesc("ChV2MS90ZXN0X3NlcnZpY2UucHJvdG8SAnYxIg4KDFRlc3QxUmVxdWVzdCIgCg1UZXN0MVJlc3BvbnNlEg8KB21lc3NhZ2UYASABKAkiHAoMVGVzdDJSZXF1ZXN0EgwKBG5hbWUYASABKAkiIAoNVGVzdDJSZXNwb25zZRIPCgdtZXNzYWdlGAEgASgJIhwKDFRlc3QzUmVxdWVzdBIMCgRuYW1lGAEgASgJIiAKDVRlc3QzUmVzcG9uc2USDwoHbWVzc2FnZRgBIAEoCTLbAQoLVGVzdFNlcnZpY2USQgoFVGVzdDESEC52MS5UZXN0MVJlcXVlc3QaES52MS5UZXN0MVJlc3BvbnNlIhSC0+STAg46ASoiCS92MS90ZXN0MRJCCgVUZXN0MhIQLnYxLlRlc3QyUmVxdWVzdBoRLnYxLlRlc3QyUmVzcG9uc2UiFILT5JMCDjoBKiIJL3YxL3Rlc3QyEkQKBVRlc3QzEhAudjEuVGVzdDNSZXF1ZXN0GhEudjEuVGVzdDNSZXNwb25zZSIUgtPkkwIOOgEqIgkvdjEvdGVzdDMwAUJaCgZjb20udjFCEFRlc3RTZXJ2aWNlUHJvdG9QAVoWYXBwL2JhY2tlbmQvcGtnL2dlbi92MaICA1ZYWKoCAlYxygICVjHiAg5WMVxHUEJNZXRhZGF0YeoCAlYxYgZwcm90bzM", [file_google_api_annotations]);

/**
 * @generated from message v1.Test1Request
 */
export type Test1Request = Message<"v1.Test1Request"> & {
};

/**
 * Describes the message v1.Test1Request.
 * Use `create(Test1RequestSchema)` to create a new message.
 */
export const Test1RequestSchema: GenMessage<Test1Request> = /*@__PURE__*/
  messageDesc(file_v1_test_service, 0);

/**
 * @generated from message v1.Test1Response
 */
export type Test1Response = Message<"v1.Test1Response"> & {
  /**
   * @generated from field: string message = 1;
   */
  message: string;
};

/**
 * Describes the message v1.Test1Response.
 * Use `create(Test1ResponseSchema)` to create a new message.
 */
export const Test1ResponseSchema: GenMessage<Test1Response> = /*@__PURE__*/
  messageDesc(file_v1_test_service, 1);

/**
 * @generated from message v1.Test2Request
 */
export type Test2Request = Message<"v1.Test2Request"> & {
  /**
   * @generated from field: string name = 1;
   */
  name: string;
};

/**
 * Describes the message v1.Test2Request.
 * Use `create(Test2RequestSchema)` to create a new message.
 */
export const Test2RequestSchema: GenMessage<Test2Request> = /*@__PURE__*/
  messageDesc(file_v1_test_service, 2);

/**
 * @generated from message v1.Test2Response
 */
export type Test2Response = Message<"v1.Test2Response"> & {
  /**
   * @generated from field: string message = 1;
   */
  message: string;
};

/**
 * Describes the message v1.Test2Response.
 * Use `create(Test2ResponseSchema)` to create a new message.
 */
export const Test2ResponseSchema: GenMessage<Test2Response> = /*@__PURE__*/
  messageDesc(file_v1_test_service, 3);

/**
 * @generated from message v1.Test3Request
 */
export type Test3Request = Message<"v1.Test3Request"> & {
  /**
   * @generated from field: string name = 1;
   */
  name: string;
};

/**
 * Describes the message v1.Test3Request.
 * Use `create(Test3RequestSchema)` to create a new message.
 */
export const Test3RequestSchema: GenMessage<Test3Request> = /*@__PURE__*/
  messageDesc(file_v1_test_service, 4);

/**
 * @generated from message v1.Test3Response
 */
export type Test3Response = Message<"v1.Test3Response"> & {
  /**
   * @generated from field: string message = 1;
   */
  message: string;
};

/**
 * Describes the message v1.Test3Response.
 * Use `create(Test3ResponseSchema)` to create a new message.
 */
export const Test3ResponseSchema: GenMessage<Test3Response> = /*@__PURE__*/
  messageDesc(file_v1_test_service, 5);

/**
 * @generated from service v1.TestService
 */
export const TestService: GenService<{
  /**
   * @generated from rpc v1.TestService.Test1
   */
  test1: {
    methodKind: "unary";
    input: typeof Test1RequestSchema;
    output: typeof Test1ResponseSchema;
  },
  /**
   * @generated from rpc v1.TestService.Test2
   */
  test2: {
    methodKind: "unary";
    input: typeof Test2RequestSchema;
    output: typeof Test2ResponseSchema;
  },
  /**
   * 流式返回
   *
   * @generated from rpc v1.TestService.Test3
   */
  test3: {
    methodKind: "server_streaming";
    input: typeof Test3RequestSchema;
    output: typeof Test3ResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_v1_test_service, 0);

