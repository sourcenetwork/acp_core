grammar PermissionExpr;

expr: term #atom 
    | expr operator term #nested
    ;

term: relation #cu_term
      |resource '->' relation #ttu_term
      | '(' expr ')' #expr_term
      ; 

relation: IDENTIFIER;
resource: IDENTIFIER;
operator: '+' #union
        | '-' #difference
        | '&' #intersection
        ;

IDENTIFIER: [a-zA-Z] [a-zA-Z0-9_]+;
WS : [ \t]+ -> skip ; // skip spaces, tabs