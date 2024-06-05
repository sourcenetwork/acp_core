grammar TestSuite;

relationship_set: relationship+;

policy_tests: checks implied_relations delegation_assertions;

checks: 'Checks' '{' check* '}';
check: relationship | NEGATION relationship;

implied_relations: 'ImpliedRelations' '{' implied_relation* '}';
implied_relation: object_rel '=>' object_rel;
object_rel: object '#' relation;

delegation_assertions: 'DelegationAssertions' '{' delegation_assertion* '}';
delegation_assertion: actorid OPERATION relationship | NEGATION actorid OPERATION relationship;

relationship: object '#' relation '@' subject;

subject: object '#' relation #subj_uset
       | object #subj_obj
       | actorid #subj_actor
       ;
object: resource ':' object_id;
object_id: ID #ascii_id
         | STRING #utf_id
         ;
relation: ID;
resource: ID;
actorid: DID;

NEGATION: '!';

OPERATION: 'delete' | 'create';
ID: [a-zA-Z] [a-zA-Z0-9_]+;
STRING: '"' .*? '"';
DID: 'did:' [a-z0-9]+ ':' [a-z0-9A-Z._-]+;
HEX: '%' HEXDIG HEXDIG;
HEXDIG: [0-9a-fA-F];

COMMENT: '//' .*? '\r'? '\n' -> skip;
WS : [ \t]+ -> skip ; // skip spaces, tabs
NL: '\r'? '\n' -> skip;
