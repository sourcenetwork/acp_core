package sandbox

import "github.com/sourcenetwork/acp_core/pkg/types"

var Samples []*types.SandboxTemplate = []*types.SandboxTemplate{
	{
		Name:        "Filesystem Example",
		Description: "Models a simple filesystem, with hiearchical files and user groups",
		Data: &types.SandboxData{
			PolicyDefinition: `
name: filesystem
resources:
  file:
    relations:
	  owner:
	    types:
		  - actor
	  reader:
	    types:
		  - actor
		  - groupe#participant
	  writer:
	    types:
		  - actor
		  - groupe#participant
	permissions:
	  read:
	    expr: owner + reader + writer
	  write:
	    expr: owner + writer
  group:
    relations:
	  owner:
	    types:
		  - actor
	  guest:
	    types:
		  - actor
	permissions:
	  participant:
	    expr: member + guest
`,
			Relationships: `
file:readme#owner@did:user:bob // bob owns file readme
file:readme#writer@did:user:alice // alice can read file readme
file:readme#reader@group:engineering#participant // participants of the engineering group can read file readme

group:engineering#owner@did:user:steve // steve owns the engineering group
group:engineering#guest@did:user:eve // eve is a guest in the engineering group
			`,
			PolicyTheorem: `
Authorizations {
  // bob can read and write to readme
  file:readme#read@did:user:bob
  file:readme#write@did:user:bob

  // alice can read and write to readme
  file:readme#read@did:user:alice
  file:readme#write@did:user:alice

  // steve and eve are participants of group engineering
  group:engineering#owner@did:user:steve 
  group:engineering#guest@did:user:eve 
  
  // steve and eve can read file readme but not write
  file:readme#read@did:user:steve
  file:readme#read@did:user:eve
  !file:readme#write@did:user:steve
  !file:readme#write@did:user:eve
}
`,
		},
	},
}
