// Code generated by protoc-gen-ts_proto. DO NOT EDIT.
// versions:
//   protoc-gen-ts_proto  v2.2.0
//   protoc               unknown
// source: sourcenetwork/acp_core/catalogue.proto

/* eslint-disable */

export const protobufPackage = "sourcenetwork.acp_core";

/** PolicyCatalogue represents a lookup table for entities definedw withing a Policy */
export interface PolicyCatalogue {
  resourceCatalogue: { [key: string]: ResourceCatalogue };
  actorResourceName: string;
  actors: string[];
}

export interface PolicyCatalogue_ResourceCatalogueEntry {
  key: string;
  value: ResourceCatalogue | undefined;
}

/** ResourceCatalogue models the set of known objects, permissions and relations for a Resource within a Policy */
export interface ResourceCatalogue {
  permissions: string[];
  relations: string[];
  objectIds: string[];
}
