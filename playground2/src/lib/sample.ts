export const SAMPLE_POLICY = `name: folders
resources:
- name: folder
  relations:
  - name: parent
    types: [folder]
  - name: viewer
    types: [actor]
  permissions:
  - name: view
    expr: viewer + parent->view
- name: doc
  relations:
  - name: folder
    types: [folder]
  - name: viewer
    types: [actor]
  permissions:
  - name: view
    expr: viewer + folder->view
`;

export const SAMPLE_CONSTRAINTS = `constraints:
- acyclic: folder.parent
- at_most_one: doc.folder
- disjoint: [doc.owner, doc.viewer]
`;

export const SAMPLE_REQUIREMENTS = `requirements:
- name: doc-view
  statement: Only people who were given access to a document can see it.
  notes:
  - Assumed that folder viewers count as given access.
  contracts:
  - name: only-owners-and-viewers-view
    description: Only a doc's owner and its direct viewers can view it.
    status: proposed
    properties:
    - type: assert
      description: every actor who can view a doc is its owner or one of its viewers
      formula: |
        all d: Doc, a: Actor | a in d.doc_view implies a in d.doc_owner + d.doc_viewer
    - type: pred
      description: some actor can view a doc
      formula: |
        some d: Doc | some d.doc_view
  - name: viewers-can-view
    description: A doc's viewers can view it.
    status: accepted
    properties:
    - type: assert
      description: every viewer of a doc can view it
      formula: |
        all d: Doc | d.doc_viewer in d.doc_view
    - type: pred
      description: a doc can have a viewer
      formula: |
        some d: Doc | some d.doc_viewer
  - name: owners-can-view
    description: A doc's owners can view it.
    properties:
    - type: assert
      description: every owner of a doc can view it (doc_owners is a typo, the relation is doc_owner)
      formula: |
        all d: Doc | d.doc_owners in d.doc_view
  - name: banned-actors-cannot-view
    description: No banned actor can view a doc, whatever else they hold.
- name: folder-inheritance
  statement: Whoever can see a folder can see the docs in it.
  contracts:
  - name: folder-viewers-view-docs
    description: Anyone who can view a doc's folder can view the doc.
    properties:
    - type: assert
      description: the viewers of a doc's folder can view the doc
      formula: |
        all d: Doc | d.folder.folder_view in d.doc_view
`;
