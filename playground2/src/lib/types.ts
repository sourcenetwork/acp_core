// Hand-written types for the Playground2 Connect JSON wire format.

export type SystemConstraint =
  | { disjoint: { resource: string; relations: string[] } }
  | { atMostOne: { resource: string; relation: string } }
  | { acyclic: { resource: string; relation: string } };

export type ContractStatus = "PROPOSED" | "ACCEPTED";
export type PropertyType = "ASSERT" | "PRED";

export interface WireProperty {
  type: PropertyType;
  description: string;
  formula: string;
}

export interface WireContract {
  name: string;
  description: string;
  status: ContractStatus;
  properties: WireProperty[];
}

export interface WireRequirement {
  name: string;
  statement: string;
  notes: string[];
  contracts: WireContract[];
}

export interface SystemModel {
  policy: string;
  constraints: SystemConstraint[];
  /** Only Verify uses requirements. */
  requirements?: WireRequirement[];
}

export interface GenerateConfig {
  count: number;
  scope: number;
  actors: number;
}

export interface GenerateRequest {
  systemModel: SystemModel;
  config?: Partial<GenerateConfig>;
}

export interface WireWorld {
  relationships?: string;
  assertions?: unknown;
  name?: string;
}

export interface GenerateResponse {
  worlds?: WireWorld[];
}

export interface Scenario {
  systemModel: SystemModel;
  world: WireWorld;
}

export interface CheckRequest {
  scenario: Scenario;
  object: { resource: string; id: string };
  /** Name of a permission or relation of the object's resource. */
  permission: string;
  actor: { id: string };
}

export interface CheckResponse {
  allowed?: boolean;
}

export type ExplainRequest = CheckRequest;

/** The wire format omits default values: no `nodeType` is UNION_NODE, no flag is false. */
export type ExplainNodeType = "UNION_NODE" | "DIFF_NODE" | "INTERSECTION_NODE" | "USERSET_NODE";

export interface ExplainNode {
  id: string;
  nodeType?: ExplainNodeType;
  text?: string;
  detail?: string;
  result?: { authorized?: boolean; explored?: boolean; exhausted?: boolean };
}

export interface ExplainEdge {
  sourceNodeId?: string;
  destNodeId?: string;
  message?: string;
}

export interface CheckExplainGraph {
  rootNodeId?: string;
  nodes?: ExplainNode[];
  edges?: ExplainEdge[];
}

export interface ExplainResponse {
  allowed?: boolean;
  graph?: CheckExplainGraph;
}

/** A world in the workspace. Generated worlds are ephemeral; saved ones are worlds/<name>.yaml in file-backed mode. */
export interface World {
  id: string;
  kind: "saved" | "ephemeral";
  /** The file's name without .yaml (saved worlds of a workspace directory). */
  name?: string;
  relationships: string;
  assertions?: string;
}

/** The playground's unit of state: one system model and its worlds. */
export interface Workspace {
  systemModel: { policy: string; constraints: string };
  worlds: World[];
}

/** The wire format may omit zero values (UNVERIFIED is the enum's zero), so every field is optional. */
export type VerifyStatus = "UNVERIFIED" | "SATISFIED" | "VIOLATED" | "DRAFT";

export interface VerifyRequest {
  systemModel: SystemModel;
  config?: { scope?: number };
}

export interface WireBinding {
  name?: string;
  values?: string[];
}

export interface PropertyResult {
  status?: VerifyStatus;
  /** Counterexample of a violated assert, or example of a satisfied pred. */
  world?: WireWorld | null;
  bindings?: WireBinding[];
  message?: string;
}

export interface ContractResult {
  requirement?: string;
  contract?: string;
  status?: VerifyStatus;
  properties?: PropertyResult[];
}

export interface VerifyResponse {
  results?: ContractResult[];
  scope?: number;
}
