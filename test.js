let req = {
    data: {
        policy_definition: `
name: test
resources:
  file:
    relations:
      owner:
        types:
          - actor
      reader:
        types:
          - actor
    permissions:
      read:
        expr: owner + reader
      write:
       expr: owner
`,
        relationships: `
file:readme#owner@did:example:bob
file:readme#reader@did:example:alice
        `,
        policy_theorem: `
				Authorizations {
				  file:readme#read@did:example:bob
				  file:readme#write@did:example:bob

				  !file:readme#write@did:example:alice
				  file:readme#read@did:example:alice
				}
				Delegations {}
				ImpliedRelations {}
        `,
    }
}