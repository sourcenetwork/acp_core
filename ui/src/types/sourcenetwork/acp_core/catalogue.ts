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

function createBasePolicyCatalogue(): PolicyCatalogue {
  return { resourceCatalogue: {}, actorResourceName: "", actors: [] };
}

export const PolicyCatalogue: MessageFns<PolicyCatalogue> = {
  fromJSON(object: any): PolicyCatalogue {
    return {
      resourceCatalogue: isObject(object.resourceCatalogue)
        ? Object.entries(object.resourceCatalogue).reduce<{ [key: string]: ResourceCatalogue }>((acc, [key, value]) => {
          acc[key] = ResourceCatalogue.fromJSON(value);
          return acc;
        }, {})
        : {},
      actorResourceName: isSet(object.actorResourceName) ? globalThis.String(object.actorResourceName) : "",
      actors: globalThis.Array.isArray(object?.actors)
        ? object.actors.map((e: any) => globalThis.String(e))
        : [],
    };
  },

  toJSON(message: PolicyCatalogue): unknown {
    const obj: any = {};
    if (message.resourceCatalogue) {
      const entries = Object.entries(message.resourceCatalogue);
      if (entries.length > 0) {
        obj.resourceCatalogue = {};
        entries.forEach(([k, v]) => {
          obj.resourceCatalogue[k] = ResourceCatalogue.toJSON(v);
        });
      }
    }
    if (message.actorResourceName !== "") {
      obj.actorResourceName = message.actorResourceName;
    }
    if (message.actors?.length) {
      obj.actors = message.actors;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<PolicyCatalogue>, I>>(base?: I): PolicyCatalogue {
    return PolicyCatalogue.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<PolicyCatalogue>, I>>(object: I): PolicyCatalogue {
    const message = createBasePolicyCatalogue();
    message.resourceCatalogue = Object.entries(object.resourceCatalogue ?? {}).reduce<
      { [key: string]: ResourceCatalogue }
    >((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = ResourceCatalogue.fromPartial(value);
      }
      return acc;
    }, {});
    message.actorResourceName = object.actorResourceName ?? "";
    message.actors = object.actors?.map((e) => e) || [];
    return message;
  },
};

function createBasePolicyCatalogue_ResourceCatalogueEntry(): PolicyCatalogue_ResourceCatalogueEntry {
  return { key: "", value: undefined };
}

export const PolicyCatalogue_ResourceCatalogueEntry: MessageFns<PolicyCatalogue_ResourceCatalogueEntry> = {
  fromJSON(object: any): PolicyCatalogue_ResourceCatalogueEntry {
    return {
      key: isSet(object.key) ? globalThis.String(object.key) : "",
      value: isSet(object.value) ? ResourceCatalogue.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: PolicyCatalogue_ResourceCatalogueEntry): unknown {
    const obj: any = {};
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value !== undefined) {
      obj.value = ResourceCatalogue.toJSON(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<PolicyCatalogue_ResourceCatalogueEntry>, I>>(
    base?: I,
  ): PolicyCatalogue_ResourceCatalogueEntry {
    return PolicyCatalogue_ResourceCatalogueEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<PolicyCatalogue_ResourceCatalogueEntry>, I>>(
    object: I,
  ): PolicyCatalogue_ResourceCatalogueEntry {
    const message = createBasePolicyCatalogue_ResourceCatalogueEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? ResourceCatalogue.fromPartial(object.value)
      : undefined;
    return message;
  },
};

function createBaseResourceCatalogue(): ResourceCatalogue {
  return { permissions: [], relations: [], objectIds: [] };
}

export const ResourceCatalogue: MessageFns<ResourceCatalogue> = {
  fromJSON(object: any): ResourceCatalogue {
    return {
      permissions: globalThis.Array.isArray(object?.permissions)
        ? object.permissions.map((e: any) => globalThis.String(e))
        : [],
      relations: globalThis.Array.isArray(object?.relations)
        ? object.relations.map((e: any) => globalThis.String(e))
        : [],
      objectIds: globalThis.Array.isArray(object?.objectIds)
        ? object.objectIds.map((e: any) => globalThis.String(e))
        : [],
    };
  },

  toJSON(message: ResourceCatalogue): unknown {
    const obj: any = {};
    if (message.permissions?.length) {
      obj.permissions = message.permissions;
    }
    if (message.relations?.length) {
      obj.relations = message.relations;
    }
    if (message.objectIds?.length) {
      obj.objectIds = message.objectIds;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ResourceCatalogue>, I>>(base?: I): ResourceCatalogue {
    return ResourceCatalogue.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ResourceCatalogue>, I>>(object: I): ResourceCatalogue {
    const message = createBaseResourceCatalogue();
    message.permissions = object.permissions?.map((e) => e) || [];
    message.relations = object.relations?.map((e) => e) || [];
    message.objectIds = object.objectIds?.map((e) => e) || [];
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