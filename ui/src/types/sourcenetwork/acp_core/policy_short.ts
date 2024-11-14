// Code generated by protoc-gen-ts_proto. DO NOT EDIT.
// versions:
//   protoc-gen-ts_proto  v2.2.0
//   protoc               unknown
// source: sourcenetwork/acp_core/policy_short.proto

/* eslint-disable */
import { ActorResource } from "./policy";

export const protobufPackage = "sourcenetwork.acp_core";

/** PolicyEncodingType enumerates supported marshaling types for policies. */
export enum PolicyMarshalingType {
  /** UNKNOWN - Fallback value for a missing Marshaling Type */
  UNKNOWN = 0,
  /** SHORT_YAML - Policy Marshaled as a YAML Short Policy definition */
  SHORT_YAML = 1,
  /** SHORT_JSON - Policy Marshaled as a JSON Short Policy definition */
  SHORT_JSON = 2,
  UNRECOGNIZED = -1,
}

export function policyMarshalingTypeFromJSON(object: any): PolicyMarshalingType {
  switch (object) {
    case 0:
    case "UNKNOWN":
      return PolicyMarshalingType.UNKNOWN;
    case 1:
    case "SHORT_YAML":
      return PolicyMarshalingType.SHORT_YAML;
    case 2:
    case "SHORT_JSON":
      return PolicyMarshalingType.SHORT_JSON;
    case -1:
    case "UNRECOGNIZED":
    default:
      return PolicyMarshalingType.UNRECOGNIZED;
  }
}

export function policyMarshalingTypeToJSON(object: PolicyMarshalingType): string {
  switch (object) {
    case PolicyMarshalingType.UNKNOWN:
      return "UNKNOWN";
    case PolicyMarshalingType.SHORT_YAML:
      return "SHORT_YAML";
    case PolicyMarshalingType.SHORT_JSON:
      return "SHORT_JSON";
    case PolicyMarshalingType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

/**
 * PolicyShort is a compact Policy definition which is conveniently expressed
 * as JSON or YAML. The shorthand format is used created a Policy.
 */
export interface PolicyShort {
  name: string;
  description: string;
  /** meta field stores arbitrary key-values from users */
  meta: { [key: string]: string };
  /**
   * resources defines resources within a policy.
   * map keys define the name for a resource
   */
  resources: { [key: string]: ResourceShort };
  /**
   * actor resource defines the actor resource for the policy
   * optional.
   */
  actor:
    | ActorResource
    | undefined;
  /** specify the policy version */
  version: string;
}

export interface PolicyShort_MetaEntry {
  key: string;
  value: string;
}

export interface PolicyShort_ResourcesEntry {
  key: string;
  value: ResourceShort | undefined;
}

export interface ResourceShort {
  doc: string;
  permissions: { [key: string]: PermissionShort };
  relations: { [key: string]: RelationShort };
}

export interface ResourceShort_PermissionsEntry {
  key: string;
  value: PermissionShort | undefined;
}

export interface ResourceShort_RelationsEntry {
  key: string;
  value: RelationShort | undefined;
}

export interface RelationShort {
  doc: string;
  /** list of relations managed by the current relation */
  manages: string[];
  /**
   * types define a list of target types the current relation can point to.
   * Each type restriction points to a a resource's relation.
   * The syntax for a type restriction is "{resource}->{relation}", where relation is optional.
   * An empty relation means the relationship can only point to an object node, as opposed to an userset.
   */
  types: string[];
}

export interface PermissionShort {
  doc: string;
  expr: string;
}

function createBasePolicyShort(): PolicyShort {
  return { name: "", description: "", meta: {}, resources: {}, actor: undefined, version: "" };
}

export const PolicyShort: MessageFns<PolicyShort> = {
  fromJSON(object: any): PolicyShort {
    return {
      name: isSet(object.name) ? globalThis.String(object.name) : "",
      description: isSet(object.description) ? globalThis.String(object.description) : "",
      meta: isObject(object.meta)
        ? Object.entries(object.meta).reduce<{ [key: string]: string }>((acc, [key, value]) => {
          acc[key] = String(value);
          return acc;
        }, {})
        : {},
      resources: isObject(object.resources)
        ? Object.entries(object.resources).reduce<{ [key: string]: ResourceShort }>((acc, [key, value]) => {
          acc[key] = ResourceShort.fromJSON(value);
          return acc;
        }, {})
        : {},
      actor: isSet(object.actor) ? ActorResource.fromJSON(object.actor) : undefined,
      version: isSet(object.version) ? globalThis.String(object.version) : "",
    };
  },

  toJSON(message: PolicyShort): unknown {
    const obj: any = {};
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.description !== "") {
      obj.description = message.description;
    }
    if (message.meta) {
      const entries = Object.entries(message.meta);
      if (entries.length > 0) {
        obj.meta = {};
        entries.forEach(([k, v]) => {
          obj.meta[k] = v;
        });
      }
    }
    if (message.resources) {
      const entries = Object.entries(message.resources);
      if (entries.length > 0) {
        obj.resources = {};
        entries.forEach(([k, v]) => {
          obj.resources[k] = ResourceShort.toJSON(v);
        });
      }
    }
    if (message.actor !== undefined) {
      obj.actor = ActorResource.toJSON(message.actor);
    }
    if (message.version !== "") {
      obj.version = message.version;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<PolicyShort>, I>>(base?: I): PolicyShort {
    return PolicyShort.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<PolicyShort>, I>>(object: I): PolicyShort {
    const message = createBasePolicyShort();
    message.name = object.name ?? "";
    message.description = object.description ?? "";
    message.meta = Object.entries(object.meta ?? {}).reduce<{ [key: string]: string }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = globalThis.String(value);
      }
      return acc;
    }, {});
    message.resources = Object.entries(object.resources ?? {}).reduce<{ [key: string]: ResourceShort }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[key] = ResourceShort.fromPartial(value);
        }
        return acc;
      },
      {},
    );
    message.actor = (object.actor !== undefined && object.actor !== null)
      ? ActorResource.fromPartial(object.actor)
      : undefined;
    message.version = object.version ?? "";
    return message;
  },
};

function createBasePolicyShort_MetaEntry(): PolicyShort_MetaEntry {
  return { key: "", value: "" };
}

export const PolicyShort_MetaEntry: MessageFns<PolicyShort_MetaEntry> = {
  fromJSON(object: any): PolicyShort_MetaEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? globalThis.String(object.value) : "",
    };
  },

  toJSON(message: PolicyShort_MetaEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== "") {
      obj.value = message.value;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<PolicyShort_MetaEntry>, I>>(base?: I): PolicyShort_MetaEntry {
    return PolicyShort_MetaEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<PolicyShort_MetaEntry>, I>>(object: I): PolicyShort_MetaEntry {
    const message = createBasePolicyShort_MetaEntry();
    message.key = object.key ?? "";
    message.value = object.value ?? "";
    return message;
  },
};

function createBasePolicyShort_ResourcesEntry(): PolicyShort_ResourcesEntry {
  return { key: "", value: undefined };
}

export const PolicyShort_ResourcesEntry: MessageFns<PolicyShort_ResourcesEntry> = {
  fromJSON(object: any): PolicyShort_ResourcesEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? ResourceShort.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: PolicyShort_ResourcesEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== undefined) {
      obj.value = ResourceShort.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<PolicyShort_ResourcesEntry>, I>>(base?: I): PolicyShort_ResourcesEntry {
    return PolicyShort_ResourcesEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<PolicyShort_ResourcesEntry>, I>>(object: I): PolicyShort_ResourcesEntry {
    const message = createBasePolicyShort_ResourcesEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? ResourceShort.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseResourceShort(): ResourceShort {
  return { doc: "", permissions: {}, relations: {} };
}

export const ResourceShort: MessageFns<ResourceShort> = {
  fromJSON(object: any): ResourceShort {
    return {
      doc: isSet(object.doc) ? globalThis.String(object.doc) : "",
      permissions: isObject(object.permissions)
        ? Object.entries(object.permissions).reduce<{ [key: string]: PermissionShort }>((acc, [key, value]) => {
          acc[key] = PermissionShort.fromJSON(value);
          return acc;
        }, {})
        : {},
      relations: isObject(object.relations)
        ? Object.entries(object.relations).reduce<{ [key: string]: RelationShort }>((acc, [key, value]) => {
          acc[key] = RelationShort.fromJSON(value);
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: ResourceShort): unknown {
    const obj: any = {};
    if (message.doc !== "") {
      obj.doc = message.doc;
    }
    if (message.permissions) {
      const entries = Object.entries(message.permissions);
      if (entries.length > 0) {
        obj.permissions = {};
        entries.forEach(([k, v]) => {
          obj.permissions[k] = PermissionShort.toJSON(v);
        });
      }
    }
    if (message.relations) {
      const entries = Object.entries(message.relations);
      if (entries.length > 0) {
        obj.relations = {};
        entries.forEach(([k, v]) => {
          obj.relations[k] = RelationShort.toJSON(v);
        });
      }
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ResourceShort>, I>>(base?: I): ResourceShort {
    return ResourceShort.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ResourceShort>, I>>(object: I): ResourceShort {
    const message = createBaseResourceShort();
    message.doc = object.doc ?? "";
    message.permissions = Object.entries(object.permissions ?? {}).reduce<{ [key: string]: PermissionShort }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[key] = PermissionShort.fromPartial(value);
        }
        return acc;
      },
      {},
    );
    message.relations = Object.entries(object.relations ?? {}).reduce<{ [key: string]: RelationShort }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[key] = RelationShort.fromPartial(value);
        }
        return acc;
      },
      {},
    );
    return message;
  },
};

function createBaseResourceShort_PermissionsEntry(): ResourceShort_PermissionsEntry {
  return { key: "", value: undefined };
}

export const ResourceShort_PermissionsEntry: MessageFns<ResourceShort_PermissionsEntry> = {
  fromJSON(object: any): ResourceShort_PermissionsEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? PermissionShort.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: ResourceShort_PermissionsEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== undefined) {
      obj.value = PermissionShort.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ResourceShort_PermissionsEntry>, I>>(base?: I): ResourceShort_PermissionsEntry {
    return ResourceShort_PermissionsEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ResourceShort_PermissionsEntry>, I>>(
    object: I,
  ): ResourceShort_PermissionsEntry {
    const message = createBaseResourceShort_PermissionsEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? PermissionShort.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseResourceShort_RelationsEntry(): ResourceShort_RelationsEntry {
  return { key: "", value: undefined };
}

export const ResourceShort_RelationsEntry: MessageFns<ResourceShort_RelationsEntry> = {
  fromJSON(object: any): ResourceShort_RelationsEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? RelationShort.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: ResourceShort_RelationsEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== undefined) {
      obj.value = RelationShort.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ResourceShort_RelationsEntry>, I>>(base?: I): ResourceShort_RelationsEntry {
    return ResourceShort_RelationsEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ResourceShort_RelationsEntry>, I>>(object: I): ResourceShort_RelationsEntry {
    const message = createBaseResourceShort_RelationsEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? RelationShort.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseRelationShort(): RelationShort {
  return { doc: "", manages: [], types: [] };
}

export const RelationShort: MessageFns<RelationShort> = {
  fromJSON(object: any): RelationShort {
    return {
      doc: isSet(object.doc) ? globalThis.String(object.doc) : "",
      manages: globalThis.Array.isArray(object?.manages) ? object.manages.map((e: any) => globalThis.String(e)) : [],
      types: globalThis.Array.isArray(object?.types) ? object.types.map((e: any) => globalThis.String(e)) : [],
    };
  },

  toJSON(message: RelationShort): unknown {
    const obj: any = {};
    if (message.doc !== "") {
      obj.doc = message.doc;
    }
    if (message.manages?.length) {
      obj.manages = message.manages;
    }
    if (message.types?.length) {
      obj.types = message.types;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RelationShort>, I>>(base?: I): RelationShort {
    return RelationShort.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RelationShort>, I>>(object: I): RelationShort {
    const message = createBaseRelationShort();
    message.doc = object.doc ?? "";
    message.manages = object.manages?.map((e) => e) || [];
    message.types = object.types?.map((e) => e) || [];
    return message;
  },
};

function createBasePermissionShort(): PermissionShort {
  return { doc: "", expr: "" };
}

export const PermissionShort: MessageFns<PermissionShort> = {
  fromJSON(object: any): PermissionShort {
    return {
      doc: isSet(object.doc) ? globalThis.String(object.doc) : "",
      expr: isSet(object.expr) ? globalThis.String(object.expr) : "",
    };
  },

  toJSON(message: PermissionShort): unknown {
    const obj: any = {};
    if (message.doc !== "") {
      obj.doc = message.doc;
    }
    if (message.expr !== "") {
      obj.expr = message.expr;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<PermissionShort>, I>>(base?: I): PermissionShort {
    return PermissionShort.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<PermissionShort>, I>>(object: I): PermissionShort {
    const message = createBasePermissionShort();
    message.doc = object.doc ?? "";
    message.expr = object.expr ?? "";
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends globalThis.Array<infer U> ? globalThis.Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isObject(value: any): boolean {
  return typeof value === "object" && value !== null;
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

export interface MessageFns<T> {
  fromJSON(object: any): T;
  toJSON(message: T): unknown;
  create<I extends Exact<DeepPartial<T>, I>>(base?: I): T;
  fromPartial<I extends Exact<DeepPartial<T>, I>>(object: I): T;
}