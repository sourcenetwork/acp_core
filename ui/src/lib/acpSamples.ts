import { SandboxData } from "@/types/proto-js/sourcenetwork/acp_core/sandbox";

export type AcpSample = {
  id: string;
  title: string;
  description: string;
  contents: SandboxData;
};

export const allSamples: AcpSample[] = [
  {
    id: "basic-sample",
    title: "Simple Sample",
    description:
      "Sample description. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla sagittis, neque non condimentum vulputate.",
    contents: {
      policyDefinition: `name: simple\nresources:\n  file:\n    relations:\n	  owner:\n	    types:\n		  - actor\n	  reader:\n	    types:\n		  - actor\n	permissions:\n	  read:\n	    expr: owner + reader\n	  write:\n	    expr: owner`,
      policyTheorem: `Authorizations {\n  file:readme#read@did:user:bob\n  file:readme#read@did:user:alice\n}\nDelegations {   \n}`,
      relationships: `file:readme#owner@did:user:bob\n`,
    },
  },
  {
    id: "blank",
    title: "Blank Sample",
    description:
      "Sample description. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla sagittis, neque non condimentum vulputate.",
    contents: {
      policyDefinition: "name: blank",
      policyTheorem: "Authorizations {\n    \n}\n\nDelegations {\n    \n}",
      relationships: "",
    },
  },
  {
    id: "sample-3",
    title: "Sample 3",
    description:
      "Sample description. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla sagittis, neque non condimentum vulputate.",
    contents: {
      policyDefinition: "name: test3",
      policyTheorem: "Authorizations {\n    \n}\n\nDelegations {\n    \n}",
      relationships: "",
    },
  },
];

export const samples = new Map(allSamples.map((sample) => [sample.id, sample]));
