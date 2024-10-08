grammar Theorem;

relationship_document: relationship EOF;
relationship_set: relationship* EOF ;
policy_thorem: authorization_theorems delegation_theorems implied_relations?  EOF ;

authorization_theorems: 'Authorizations' '{' authorization_theorem* '}';
authorization_theorem: relationship | NEGATION relationship;

implied_relations: 'ImpliedRelations' '{' implied_relation* '}';
implied_relation: object_rel '=>' object_rel;
object_rel: object '#' relation;

delegation_theorems: 'Delegations' '{' delegation_theorem* '}';
delegation_theorem: actorid '>' operation  | NEGATION actorid '>' operation;

relationship: object '#' relation '@' subject;

operation: object '#' relation;

subject: object '#' relation #subj_uset
       | object #subj_obj
       | actorid #subj_actor
       ;
object: resource ':' object_id;
object_id: ID  #ascii_id
         | STRING  #utf_id
         ;
relation: ID;
resource: ID;
actorid: DID;

NEGATION: '!';

OPERATION: 'delete' | 'create';
ID: [a-zA-Z] [a-zA-Z0-9_]+;
STRING: '"' .*? '"';
DID: 'did:' [a-z0-9]+ ':' [a-z0-9A-Z._-]+;

COMMENT: '//' .*? '\r'? '\n' -> skip;
WS : [ \t]+ -> skip ; // skip spaces, tabs
NL: '\r'? '\n' -> skip;
